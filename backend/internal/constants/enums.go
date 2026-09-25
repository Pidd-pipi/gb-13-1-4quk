package constants

// Role 用户角色枚举。
const (
	RoleStudent = "student"
	RoleAdmin   = "admin"
)

// 角色可选项，用于 handler 校验与 RBAC。
var RoleOptions = []string{RoleStudent, RoleAdmin}

// BookStatus 书籍状态机枚举：on_sale(在售) -> reserved(已预约) -> sold(已售出)；
// 短借流程额外使用 loaned(借出中)，归还确认后回到 on_sale。
const (
	BookStatusOnSale   = "on_sale"
	BookStatusReserved = "reserved"
	BookStatusSold     = "sold"
	BookStatusLoaned   = "loaned"
)

var BookStatusOptions = []string{BookStatusOnSale, BookStatusReserved, BookStatusSold, BookStatusLoaned}

// BorrowDuration 短借时长（天）。
const (
	BorrowDuration7  = 7
	BorrowDuration14 = 14
)

// BorrowDurationOptions 可借时长可选项，用于发布校验与 DTO。
var BorrowDurationOptions = []int{BorrowDuration7, BorrowDuration14}

// IsValidBorrowDuration 校验可借时长是否合法。
func IsValidBorrowDuration(days int) bool {
	for _, d := range BorrowDurationOptions {
		if d == days {
			return true
		}
	}
	return false
}

// BorrowStatus 短借申请状态机枚举：
// pending(待同意) -> approved(已同意/借出中) -> returning(待卖家确认归还) -> returned(已归还)；
// pending -> rejected(已拒绝，同意其他申请时其余申请自动拒绝)。
const (
	BorrowStatusPending   = "pending"
	BorrowStatusApproved  = "approved"
	BorrowStatusRejected  = "rejected"
	BorrowStatusReturning = "returning"
	BorrowStatusReturned  = "returned"
)

// BorrowStatusOptions 短借申请状态可选项。
var BorrowStatusOptions = []string{
	BorrowStatusPending,
	BorrowStatusApproved,
	BorrowStatusRejected,
	BorrowStatusReturning,
	BorrowStatusReturned,
}

// Condition 新旧程度枚举。
const (
	ConditionBrandNew = "brand_new" // 全新
	ConditionNineNew  = "nine_new"  // 九成新
	ConditionSevenNew = "seven_new" // 七成新
	ConditionFiveNew  = "five_new"  // 五成新
)

var ConditionOptions = []string{ConditionBrandNew, ConditionNineNew, ConditionSevenNew, ConditionFiveNew}

// SubjectCategory 学科分类枚举。
const (
	SubjectScience        = "science"         // 理工
	SubjectHumanities     = "humanities"      // 文史
	SubjectEconManagement = "econ_management" // 经管
	SubjectArt            = "art"             // 艺术
	SubjectOther          = "other"           // 其他
)

var SubjectCategoryOptions = []string{SubjectScience, SubjectHumanities, SubjectEconManagement, SubjectArt, SubjectOther}

// TradeType 交易方式枚举。
const (
	TradeTypeInPerson = "in_person" // 面交
	TradeTypeMail     = "mail"      // 邮寄
)

var TradeTypeOptions = []string{TradeTypeInPerson, TradeTypeMail}

// EvaluationType 交易评价类型枚举。
const (
	EvaluationGood    = "good"    // 好评
	EvaluationNeutral = "neutral" // 中评
	EvaluationBad     = "bad"     // 差评
)

var EvaluationTypeOptions = []string{EvaluationGood, EvaluationNeutral, EvaluationBad}

// WishStatus 求购信息状态枚举。
const (
	WishStatusOpen   = "open"
	WishStatusClosed = "closed"
)

var WishStatusOptions = []string{WishStatusOpen, WishStatusClosed}

// UserStatus 用户账号状态枚举。
const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)

var UserStatusOptions = []string{UserStatusActive, UserStatusDisabled}
