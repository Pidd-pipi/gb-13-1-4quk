package dto

// SendCodeRequest 发送邮箱验证码请求。
type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// SendCodeResponse 发送验证码响应（开发环境直接返回 dev_code 便于联调）。
type SendCodeResponse struct {
	Email   string `json:"email"`
	Expire  int    `json:"expire_seconds"`
	DevCode string `json:"dev_code,omitempty"`
}

// RegisterRequest 注册请求。
type RegisterRequest struct {
	StudentNo string `json:"student_no" binding:"required,min=3,max=32"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=64"`
	Code      string `json:"code" binding:"required,len=6"`
}

// LoginRequest 登录请求（学号或邮箱 + 密码）。
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse 登录/注册成功返回的令牌。
type TokenResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}
