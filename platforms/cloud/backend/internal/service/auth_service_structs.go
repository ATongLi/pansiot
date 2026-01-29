package service

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Phone            string `json:"phone"`
	TenantSerial     string `json:"tenant_serial"`
	TenantName       string `json:"tenant_name"`
	Industry         string `json:"industry"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"verification_code"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	RealName     string `json:"real_name"`
	TenantID     int64  `json:"tenant_id"`
	TenantName   string `json:"tenant_name"`
	SerialNumber string `json:"serial_number"`
	TenantType   string `json:"tenant_type"`
}
