package service

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
	"github.com/campusbooks/campusbooks/internal/util"
)

// ConversationService handles chat conversations and messages.
type ConversationService struct {
	convRepo *repository.ConversationRepository
	msgRepo  *repository.MessageRepository
	bookRepo *repository.BookRepository
	wishRepo *repository.WishRepository
	logger   *slog.Logger
}

// NewConversationService creates a ConversationService.
func NewConversationService(convRepo *repository.ConversationRepository, msgRepo *repository.MessageRepository,
	bookRepo *repository.BookRepository, wishRepo *repository.WishRepository, logger *slog.Logger) *ConversationService {
	return &ConversationService{convRepo: convRepo, msgRepo: msgRepo, bookRepo: bookRepo, wishRepo: wishRepo, logger: logger}
}

// CreateConversationFromBook opens a chat between a buyer and a book seller.
func (s *ConversationService) CreateConversationFromBook(userID, bookID uint, content string) (*dto.ConversationDTO, error) {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeBookNotFound, constants.MsgNotFound+": book id="+fmt.Sprint(bookID))
	}
	if book.SellerID == userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeConversationForbidden, "不能与自己发布的书籍发起会话")
	}
	return s.createConversation(bookID, 0, userID, book.SellerID, content)
}

// CreateConversationFromWish opens a chat between a seller and a wish owner.
func (s *ConversationService) CreateConversationFromWish(userID, wishID uint, content string) (*dto.ConversationDTO, error) {
	wish, err := s.wishRepo.FindByID(wishID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeWishNotFound, constants.MsgNotFound+": wish id="+fmt.Sprint(wishID))
	}
	if wish.UserID == userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeConversationForbidden, "不能联系自己发布的求购信息")
	}
	if !wish.IsOpen() {
		return nil, util.NewAppError(http.StatusConflict, constants.CodeWishAlreadyClosed, "该求购信息已关闭")
	}
	return s.createConversation(0, wishID, wish.UserID, userID, content)
}

// createConversation is the shared conversation factory reused by both
// book-based and wish-based flows (至少 2 个接口复用同一 service 方法).
func (s *ConversationService) createConversation(bookID, wishID, buyerID, sellerID uint, content string) (*dto.ConversationDTO, error) {
	existing, err := s.convRepo.FindExisting(bookID, wishID, buyerID, sellerID)
	if err == nil {
		if _, err := s.SendMessage(existing.ID, buyerID, dto.SendMessageRequest{Content: content}); err != nil {
			return nil, err
		}
		return s.convert(existing), nil
	}
	now := time.Now()
	conv := &model.Conversation{
		BookID:        bookID,
		WishID:        wishID,
		BuyerID:       buyerID,
		SellerID:      sellerID,
		LastMessage:   content,
		LastMessageAt: &now,
	}
	if err := s.convRepo.Create(conv); err != nil {
		s.logger.Error(constants.LogConversationCreated, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	if _, err := s.SendMessage(conv.ID, buyerID, dto.SendMessageRequest{Content: content}); err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogConversationCreated, "id", conv.ID, "buyer", buyerID, "seller", sellerID)
	return s.convert(conv), nil
}

// ListConversations returns the user's conversations.
func (s *ConversationService) ListConversations(userID uint, page, size int) (*dto.PageData, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = util.DefaultPageSize
	}
	items, total, err := s.convRepo.ListByUser(userID, (page-1)*size, size)
	if err != nil {
		s.logger.Error("conversation list failed", "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.ConversationDTO, 0, len(items))
	for i := range items {
		c := items[i]
		d := s.convert(&c)
		unread, _ := s.msgRepo.CountUnread(c.ID, userID)
		d.UnreadCount = unread
		list = append(list, *d)
	}
	return &dto.PageData{List: list, Total: total, Page: page, Size: size}, nil
}

// GetConversation returns a conversation if the user participates.
func (s *ConversationService) GetConversation(userID, convID uint) (*dto.ConversationDTO, error) {
	conv, err := s.convRepo.FindByID(convID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeConversationNotFound, constants.MsgNotFound+": conversation id="+fmt.Sprint(convID))
	}
	if conv.BuyerID != userID && conv.SellerID != userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeConversationForbidden, "无权查看该会话")
	}
	d := s.convert(conv)
	unread, _ := s.msgRepo.CountUnread(convID, userID)
	d.UnreadCount = unread
	return d, nil
}

// SendMessage adds a message to a conversation.
func (s *ConversationService) SendMessage(convID, senderID uint, req dto.SendMessageRequest) (*dto.MessageDTO, error) {
	conv, err := s.convRepo.FindByID(convID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeConversationNotFound, constants.MsgNotFound+": conversation id="+fmt.Sprint(convID))
	}
	if conv.BuyerID != senderID && conv.SellerID != senderID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeConversationForbidden, "无权在该会话中发言")
	}
	if req.Content == "" && req.ImageURL == "" {
		return nil, util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "消息内容不能为空")
	}
	msg := &model.Message{
		ConversationID: convID,
		SenderID:       senderID,
		Content:        req.Content,
		ImageURL:       req.ImageURL,
		IsRead:         false,
	}
	if err := s.msgRepo.Create(msg); err != nil {
		s.logger.Error(constants.LogMessageSendSuccess, "conversation_id", convID, "sender", senderID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	if err := s.convRepo.Touch(convID, truncate(req.Content, 50)); err != nil {
		s.logger.Warn("conversation touch failed", "conversation_id", convID, "error", err)
	}
	s.logger.Info(constants.LogMessageSendSuccess, "conversation_id", convID, "sender", senderID)
	m := dto.FromMessage(msg)
	return &m, nil
}

// ListMessages returns messages of a conversation.
func (s *ConversationService) ListMessages(userID, convID uint) ([]dto.MessageDTO, error) {
	conv, err := s.convRepo.FindByID(convID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeConversationNotFound, constants.MsgNotFound+": conversation id="+fmt.Sprint(convID))
	}
	if conv.BuyerID != userID && conv.SellerID != userID {
		return nil, util.NewAppError(http.StatusForbidden, constants.CodeConversationForbidden, "无权查看该会话")
	}
	items, err := s.msgRepo.ListByConversation(convID)
	if err != nil {
		s.logger.Error("message list failed", "conversation_id", convID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.MessageDTO, 0, len(items))
	for i := range items {
		list = append(list, dto.FromMessage(&items[i]))
	}
	return list, nil
}

// MarkRead marks all messages in a conversation as read for the user.
func (s *ConversationService) MarkRead(userID, convID uint) error {
	conv, err := s.convRepo.FindByID(convID)
	if err != nil {
		return util.NewAppError(http.StatusNotFound, constants.CodeConversationNotFound, constants.MsgNotFound+": conversation id="+fmt.Sprint(convID))
	}
	if conv.BuyerID != userID && conv.SellerID != userID {
		return util.NewAppError(http.StatusForbidden, constants.CodeConversationForbidden, "无权操作该会话")
	}
	if err := s.msgRepo.MarkConversationRead(convID, userID); err != nil {
		s.logger.Error(constants.LogConversationRead, "id", convID, "user_id", userID, "error", err)
		return util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogConversationRead, "id", convID, "user_id", userID)
	return nil
}

func (s *ConversationService) convert(c *model.Conversation) *dto.ConversationDTO {
	d := dto.FromConversation(c)
	if c.Book != nil {
		b := dto.FromBook(c.Book)
		d.Book = &b
	}
	if c.Buyer != nil {
		u := dto.FromUser(c.Buyer)
		d.Buyer = &u
	}
	if c.Seller != nil {
		u := dto.FromUser(c.Seller)
		d.Seller = &u
	}
	return &d
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "..."
}
