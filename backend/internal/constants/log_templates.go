package constants

// 日志模板集中管理（>= 25 条）：所有 handler/service/middleware 统一引用。
const (
	LogServerStarted        = "campusbooks backend started: addr=%s"
	LogServerStopped        = "campusbooks backend stopped: error=%v"
	LogServerShutdown       = "campusbooks backend shutting down"
	LogServerShutdownFailed = "campusbooks backend shutdown failed: error=%v"
	LogDBConnectFailed      = "database connect failed: error=%v"
	LogDBMigrateFailed      = "database migrate failed: error=%v"
	LogDBSeedFailed         = "database seed failed: error=%v"
	LogMinIOConnectFailed   = "minio connect failed: error=%v"
	LogRedisConnectFailed   = "redis connect failed, fallback to memory: error=%v"

	LogAuthCodeSent        = "email verify code sent: email=%s"
	LogAuthRegisterSuccess = "user register success: student_no=%s email=%s"
	LogAuthRegisterFailed  = "user register failed: student_no=%s error=%v"
	LogAuthLoginSuccess    = "user login success: email=%s role=%s"
	LogAuthLoginFailed     = "user login failed: email=%s error=%v"

	LogUserProfileUpdated  = "user profile updated: user_id=%d"
	LogUserAvatarUpdated   = "user avatar updated: user_id=%d"
	LogUserListSuccess     = "user list success: page=%d size=%d"
	LogUserStatsCalculated = "user stats calculated: user_id=%d good_rate=%s"

	LogBookCreateSuccess    = "book created: id=%d title=%s seller_id=%d"
	LogBookCreateFailed     = "book create failed: seller_id=%d error=%v"
	LogBookUpdateSuccess    = "book updated: id=%d title=%s"
	LogBookDeleteSuccess    = "book deleted: id=%d seller_id=%d"
	LogBookStatusChanged    = "book status changed: id=%d from=%s to=%s operator=%d"
	LogBookStatusChangeFail = "book status change failed: id=%d from=%s to=%s error=%v"
	LogBookViewed           = "book viewed: id=%d user_id=%d"
	LogBookFavoriteAdded    = "book favorite added: user_id=%d book_id=%d"
	LogBookFavoriteRemoved  = "book favorite removed: user_id=%d book_id=%d"
	LogBookRecommendList    = "book recommendation list success: user_id=%d"

	LogWishCreateSuccess  = "wish created: id=%d book_title=%s user_id=%d"
	LogWishCloseSuccess   = "wish closed: id=%d user_id=%d"
	LogWishContactSuccess = "wish contact created: wish_id=%d seller_id=%d"

	LogConversationCreated = "conversation created: id=%d buyer=%d seller=%d"
	LogMessageSendSuccess  = "message sent: conversation_id=%d sender=%d"
	LogConversationRead    = "conversation marked read: id=%d user_id=%d"

	LogBorrowApplied         = "borrow request created: id=%d book_id=%d borrower=%d"
	LogBorrowApproved        = "borrow request approved: id=%d book_id=%d due_at=%s"
	LogBorrowRejected        = "borrow request rejected: id=%d book_id=%d seller=%d"
	LogBorrowReturned        = "borrow request returned: id=%d book_id=%d borrower=%d"
	LogBorrowConfirmReturned = "borrow request confirm returned: id=%d book_id=%d seller=%d"
	LogBorrowReminded        = "borrow reminder sent: id=%d book_id=%d seller=%d"

	LogEvaluationCreate = "evaluation created: from=%d to=%d book_id=%d type=%s"
	LogAuditRecorded    = "audit log recorded: user_id=%d action=%s resource=%s/%d"
	LogUploadSuccess    = "file upload success: object=%s bucket=%s"
	LogUploadFailed     = "file upload failed: filename=%s error=%v"
	LogRateLimited      = "request rate limited: path=%s ip=%s"
	LogRequestHandled   = "request handled: request_id=%s method=%s path=%s status=%d latency_ms=%d"
	LogPanicRecovered   = "panic recovered: request_id=%s error=%v"
)
