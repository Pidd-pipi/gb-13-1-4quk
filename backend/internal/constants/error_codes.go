package constants

// 统一业务错误码：0 成功，其余按模块分段。
const (
	CodeSuccess         = 0
	CodeBadRequest      = 40000
	CodeUnauthorized    = 40100
	CodeForbidden       = 40300
	CodeNotFound        = 40400
	CodeConflict        = 40900
	CodeRateLimited     = 42900
	CodeValidationError = 42200
	CodeInternalError   = 50000
)

// 用户/认证模块错误码。
const (
	CodeEmailCodeInvalid   = 40101
	CodeEmailCodeExpired   = 40102
	CodeStudentNoTaken     = 40901
	CodeEmailTaken         = 40902
	CodeInvalidCredentials = 40103
	CodeUserDisabled       = 40301
	CodeProfileNotComplete = 42201
)

// 书籍模块错误码。
const (
	CodeBookNotFound        = 40401
	CodeBookStatusInvalid   = 40903
	CodeBookNotOwned        = 40302
	CodeBookStatusConflict  = 40904
	CodeBookAlreadyFavorite = 40905
	CodeBookNotFavorite     = 40402
)

// 短借模块错误码。
const (
	CodeBorrowNotFound        = 40405
	CodeBorrowForbidden       = 40306
	CodeBorrowNotAllowed      = 42203
	CodeBorrowConflict        = 40908
	CodeBorrowDuplicate       = 40909
	CodeBorrowDurationInvalid = 42204
)

// 求购模块错误码。
const (
	CodeWishNotFound      = 40403
	CodeWishNotOwned      = 40303
	CodeWishAlreadyClosed = 40906
)

// 会话/消息模块错误码。
const (
	CodeConversationNotFound  = 40404
	CodeConversationForbidden = 40304
)

// 评价模块错误码。
const (
	CodeEvaluationDuplicate = 40907
	CodeEvaluationInvalid   = 42202
	CodeEvaluationForbidden = 40305
)

// 通用错误码提示文案。
const (
	MsgInvalidParam  = "invalid request parameter"
	MsgUnauthorized  = "authentication required"
	MsgForbidden     = "permission denied"
	MsgNotFound      = "resource not found"
	MsgConflict      = "resource conflict"
	MsgRateLimited   = "too many requests"
	MsgInternalError = "internal server error"
)
