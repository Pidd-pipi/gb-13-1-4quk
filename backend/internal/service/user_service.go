package service

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/repository"
	"github.com/campusbooks/campusbooks/internal/util"
)

// UserService handles profile and user queries.
type UserService struct {
	userRepo *repository.UserRepository
	evalRepo *repository.EvaluationRepository
	logger   *slog.Logger
}

// NewUserService creates a UserService.
func NewUserService(userRepo *repository.UserRepository, evalRepo *repository.EvaluationRepository, logger *slog.Logger) *UserService {
	return &UserService{userRepo: userRepo, evalRepo: evalRepo, logger: logger}
}

// GetProfile returns a user by id.
func (s *UserService) GetProfile(userID uint) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound+": user id="+strconv.FormatUint(uint64(userID), 10))
	}
	d := dto.FromUser(user)
	return &d, nil
}

// UpdateProfile updates name/department/campus/contact.
func (s *UserService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound+": user id="+strconv.FormatUint(uint64(userID), 10))
	}
	user.Name = req.Name
	user.Department = req.Department
	user.Campus = req.Campus
	user.Contact = req.Contact
	if err := s.userRepo.Update(user); err != nil {
		s.logger.Error(constants.LogUserProfileUpdated, "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogUserProfileUpdated, "user_id", userID)
	d := dto.FromUser(user)
	return &d, nil
}

// UpdateAvatar updates the avatar url.
func (s *UserService) UpdateAvatar(userID uint, avatarURL string) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound+": user id="+strconv.FormatUint(uint64(userID), 10))
	}
	user.AvatarURL = avatarURL
	if err := s.userRepo.Update(user); err != nil {
		s.logger.Error(constants.LogUserAvatarUpdated, "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(constants.LogUserAvatarUpdated, "user_id", userID)
	d := dto.FromUser(user)
	return &d, nil
}

// ListUsers returns users for admin management.
func (s *UserService) ListUsers(role string, page, size int) (*dto.PageData, error) {
	items, total, err := s.userRepo.List(role, (page-1)*size, size)
	if err != nil {
		s.logger.Error(constants.LogUserListSuccess, "page", page, "size", size, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	list := make([]dto.UserDTO, 0, len(items))
	for _, u := range items {
		list = append(list, dto.FromUser(&u))
	}
	s.logger.Info(constants.LogUserListSuccess, "page", page, "size", size)
	return &dto.PageData{List: list, Total: total, Page: page, Size: size}, nil
}

// GetStats computes user homepage stats and risk flag.
func (s *UserService) GetStats(userID uint) (*dto.UserStatsDTO, error) {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return nil, util.NewAppError(http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound+": user id="+strconv.FormatUint(uint64(userID), 10))
	}
	total, onSale, sold, err := s.userRepo.CountBooksBySeller(userID)
	if err != nil {
		s.logger.Error("count seller books failed", "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	good, err := s.evalRepo.CountByType(strconv.FormatUint(uint64(userID), 10), constants.EvaluationGood)
	if err != nil {
		s.logger.Error("count good evaluation failed", "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	neutral, err := s.evalRepo.CountByType(strconv.FormatUint(uint64(userID), 10), constants.EvaluationNeutral)
	if err != nil {
		s.logger.Error("count neutral evaluation failed", "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	bad, err := s.evalRepo.CountByType(strconv.FormatUint(uint64(userID), 10), constants.EvaluationBad)
	if err != nil {
		s.logger.Error("count bad evaluation failed", "user_id", userID, "error", err)
		return nil, util.NewAppError(http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
	totalEval := good + neutral + bad
	rate := "100.0%"
	if totalEval > 0 {
		rate = fmt.Sprintf("%.1f%%", float64(good)/float64(totalEval)*100)
	}
	risk := totalEval >= 3 && float64(bad)/float64(totalEval) >= 0.4
	stats := &dto.UserStatsDTO{
		UserID:           userID,
		TotalBooks:       total,
		OnSaleBooks:      onSale,
		SoldBooks:        sold,
		TotalEvaluations: totalEval,
		GoodCount:        good,
		NeutralCount:     neutral,
		BadCount:         bad,
		GoodRate:         rate,
		RiskFlagged:      risk,
	}
	s.logger.Info(constants.LogUserStatsCalculated, "user_id", userID, "good_rate", rate)
	return stats, nil
}
