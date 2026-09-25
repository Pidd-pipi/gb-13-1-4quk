package constants

// Role 用户角色枚举。
const (
	RoleStudent = "student"
	RoleAdmin   = "admin"
)

// 角色可选项，用于 handler 校验与 RBAC。
var RoleOptions = []string{RoleStudent, RoleAdmin}

// BookStatus 书籍状态机枚举：on_sale(在售) -> reserved(已预约) -> sold(已售出)。
const (
	BookStatusOnSale   = "on_sale"
	BookStatusReserved = "reserved"
	BookStatusSold     = "sold"
)

var BookStatusOptions = []string{BookStatusOnSale, BookStatusReserved, BookStatusSold}

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
