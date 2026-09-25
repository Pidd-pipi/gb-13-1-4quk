package database

import (
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/util"
)

// Seed inserts baseline demo data when tables are empty.
func Seed(db *gorm.DB, logger *slog.Logger) error {
	var userCount int64
	if err := db.Model(&model.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		logger.Info("seed skipped: users already exist", "count", userCount)
		return nil
	}

	adminPwd, _ := util.HashPassword("admin123")
	stuPwd, _ := util.HashPassword("student123")
	users := []model.User{
		{StudentNo: "ADMIN001", Email: "admin@campusbooks.local", PasswordHash: adminPwd, Name: "系统管理员", Department: "教务处", Campus: "校本部", Contact: "010-00000000", Role: constants.RoleAdmin, EmailVerified: true, Status: constants.UserStatusActive},
		{StudentNo: "20230001", Email: "zhang@campusbooks.local", PasswordHash: stuPwd, Name: "张伟", Department: "计算机学院", Campus: "东校区", Contact: "13800000001", Role: constants.RoleStudent, EmailVerified: true, Status: constants.UserStatusActive},
		{StudentNo: "20230002", Email: "li@campusbooks.local", PasswordHash: stuPwd, Name: "李娜", Department: "计算机学院", Campus: "东校区", Contact: "13800000002", Role: constants.RoleStudent, EmailVerified: true, Status: constants.UserStatusActive},
		{StudentNo: "20230003", Email: "wang@campusbooks.local", PasswordHash: stuPwd, Name: "王芳", Department: "经济管理学院", Campus: "西校区", Contact: "13800000003", Role: constants.RoleStudent, EmailVerified: true, Status: constants.UserStatusActive},
		{StudentNo: "20230004", Email: "chen@campusbooks.local", PasswordHash: stuPwd, Name: "陈杰", Department: "外国语学院", Campus: "校本部", Contact: "13800000004", Role: constants.RoleStudent, EmailVerified: true, Status: constants.UserStatusActive},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}

	books := []model.Book{
		{SellerID: 2, Title: "高等数学（第七版）上册", Author: "同济大学数学系", ISBN: "9787040396638", CourseName: "高等数学", OriginalPrice: 48, Price: 20, Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Campus: "东校区", Description: "内页有少量笔记，不影响使用。", Images: []byte("[]"), Status: constants.BookStatusOnSale, ViewCount: 120, FavoriteCount: 5},
		{SellerID: 2, Title: "数据结构（C语言版）", Author: "严蔚敏", ISBN: "9787302147510", CourseName: "数据结构", OriginalPrice: 39, Price: 15, Condition: constants.ConditionSevenNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeMail, Campus: "东校区", Description: "教材封面轻微磨损，内容完整。", Images: []byte("[]"), Status: constants.BookStatusOnSale, ViewCount: 88, FavoriteCount: 3},
		{SellerID: 3, Title: "线性代数及其应用", Author: "David C. Lay", ISBN: "9787111476434", CourseName: "线性代数", OriginalPrice: 59, Price: 30, Condition: constants.ConditionBrandNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Campus: "东校区", Description: "全新未拆封，可小刀。", Images: []byte("[]"), Status: constants.BookStatusReserved, ReservedBy: 2, ViewCount: 210, FavoriteCount: 8},
		{SellerID: 3, Title: "微观经济学", Author: "曼昆", ISBN: "9787301257775", CourseName: "微观经济学", OriginalPrice: 62, Price: 25, Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectEconManagement, TradeType: constants.TradeTypeInPerson, Campus: "西校区", Description: "几乎全新，有荧光笔重点。", Images: []byte("[]"), Status: constants.BookStatusSold, ViewCount: 300, FavoriteCount: 12},
		{SellerID: 4, Title: "大学英语综合教程（第三册）", Author: "李荫华", ISBN: "9787544636010", CourseName: "大学英语", OriginalPrice: 45, Price: 18, Condition: constants.ConditionFiveNew, SubjectCategory: constants.SubjectHumanities, TradeType: constants.TradeTypeMail, Campus: "校本部", Description: "旧版教材，适合复习使用。", Images: []byte("[]"), Status: constants.BookStatusOnSale, ViewCount: 66, FavoriteCount: 2},
		{SellerID: 5, Title: "设计色彩基础", Author: "李立新", ISBN: "9787535637406", CourseName: "设计基础", OriginalPrice: 55, Price: 22, Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectArt, TradeType: constants.TradeTypeInPerson, Campus: "校本部", Description: "艺术设计专业用书。", Images: []byte("[]"), Status: constants.BookStatusOnSale, ViewCount: 45, FavoriteCount: 1},
		{SellerID: 3, Title: "计算机组成原理（第2版）", Author: "唐朔飞", ISBN: "9787040258424", CourseName: "计算机组成原理", OriginalPrice: 39, Price: 16, Condition: constants.ConditionSevenNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Campus: "东校区", Description: "同院系同学在售教材。", Images: []byte("[]"), Status: constants.BookStatusOnSale, ViewCount: 30, FavoriteCount: 1},
		{SellerID: 4, Title: "概率论与数理统计", Author: "盛骤", ISBN: "9787040238969", CourseName: "概率论", OriginalPrice: 36, Price: 12, Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Campus: "校本部", Description: "只用一个学期，支持短借。", Images: []byte("[]"), Status: constants.BookStatusOnSale, Lendable: true, LendDays: constants.LendDays7, ViewCount: 54, FavoriteCount: 2},
		{SellerID: 2, Title: "大学物理（上册）", Author: "张三慧", ISBN: "9787302193081", CourseName: "大学物理", OriginalPrice: 42, Price: 15, Condition: constants.ConditionSevenNew, SubjectCategory: constants.SubjectScience, TradeType: constants.TradeTypeInPerson, Campus: "东校区", Description: "可短借两周，到期请按时归还。", Images: []byte("[]"), Status: constants.BookStatusLentOut, Lendable: true, LendDays: constants.LendDays14, ViewCount: 73, FavoriteCount: 4},
		{SellerID: 3, Title: "宏观经济学", Author: "曼昆", ISBN: "9787301294345", CourseName: "宏观经济学", OriginalPrice: 58, Price: 20, Condition: constants.ConditionNineNew, SubjectCategory: constants.SubjectEconManagement, TradeType: constants.TradeTypeInPerson, Campus: "西校区", Description: "期末过渡用书，可借 7 天。", Images: []byte("[]"), Status: constants.BookStatusLentOut, Lendable: true, LendDays: constants.LendDays7, ViewCount: 41, FavoriteCount: 1},
	}
	if err := db.Create(&books).Error; err != nil {
		return err
	}

	now := time.Now()
	approvedAt := now.Add(-3 * 24 * time.Hour)
	dueAt := approvedAt.Add(14 * 24 * time.Hour)
	overdueApprovedAt := now.Add(-9 * 24 * time.Hour)
	overdueDueAt := overdueApprovedAt.Add(7 * 24 * time.Hour)
	borrows := []model.BorrowRequest{
		{BookID: 9, BorrowerID: 3, SellerID: 2, Status: constants.BorrowStatusApproved, ApprovedAt: &approvedAt, DueAt: &dueAt},
		{BookID: 10, BorrowerID: 5, SellerID: 3, Status: constants.BorrowStatusApproved, ApprovedAt: &overdueApprovedAt, DueAt: &overdueDueAt},
		{BookID: 8, BorrowerID: 2, SellerID: 4, Status: constants.BorrowStatusPending},
	}
	if err := db.Create(&borrows).Error; err != nil {
		return err
	}

	wishes := []model.Wish{
		{UserID: 3, BookTitle: "操作系统概念（第9版）", Author: "Abraham Silberschatz", ISBN: "9787111544935", ExpectedPrice: 30, ConditionRequirement: constants.ConditionNineNew, SubjectCategory: constants.SubjectScience, Description: "希望九成新以上，价格可谈。", Status: constants.WishStatusOpen},
		{UserID: 2, BookTitle: "管理学原理", Author: "周三多", ExpectedPrice: 15, ConditionRequirement: constants.ConditionSevenNew, SubjectCategory: constants.SubjectEconManagement, Description: "期末复习用。", Status: constants.WishStatusOpen},
		{UserID: 5, BookTitle: "中国文学史", Author: "袁行霈", ExpectedPrice: 25, ConditionRequirement: constants.ConditionFiveNew, SubjectCategory: constants.SubjectHumanities, Description: "任意版本均可。", Status: constants.WishStatusClosed},
	}
	if err := db.Create(&wishes).Error; err != nil {
		return err
	}

	conv := model.Conversation{BookID: 1, BuyerID: 3, SellerID: 2, LastMessage: "同学你好，这本书还在吗？", LastMessageAt: &now}
	if err := db.Create(&conv).Error; err != nil {
		return err
	}
	messages := []model.Message{
		{ConversationID: conv.ID, SenderID: 3, Content: "同学你好，这本书还在吗？", IsRead: true},
		{ConversationID: conv.ID, SenderID: 2, Content: "还在的，你在哪个校区？可以面交。", IsRead: false},
	}
	if err := db.Create(&messages).Error; err != nil {
		return err
	}

	evaluations := []model.Evaluation{
		{FromUserID: 3, ToUserID: 4, BookID: 4, Type: constants.EvaluationGood, Content: "书很新，卖家很耐心。"},
		{FromUserID: 4, ToUserID: 3, BookID: 4, Type: constants.EvaluationGood, Content: "买家准时面交，沟通顺畅。"},
	}
	if err := db.Create(&evaluations).Error; err != nil {
		return err
	}

	favorites := []model.Favorite{{UserID: 3, BookID: 1}, {UserID: 2, BookID: 3}}
	if err := db.Create(&favorites).Error; err != nil {
		return err
	}

	histories := []model.BrowseHistory{{UserID: 3, BookID: 1, ViewedAt: now}, {UserID: 3, BookID: 4, ViewedAt: now}}
	if err := db.Create(&histories).Error; err != nil {
		return err
	}

	logger.Info("seed completed", "users", len(users), "books", len(books), "wishes", len(wishes))
	return nil
}
