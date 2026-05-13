package dto

// LoginResp is the response for successful login.
type LoginResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// UserProfileResp is the response for GET /api/v1/auth/me.
type UserProfileResp struct {
	ID                 int64  `json:"id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	RealName           string `json:"real_name"`
	Phone              string `json:"phone"`
	Status             int16  `json:"status"`
	MustChangePassword bool   `json:"must_change_password"`
	PasswordUpdatedAt  string `json:"password_updated_at"`
	LastLoginAt        string `json:"last_login_at"`
	Roles              []string `json:"roles"`
	Stores             []int64  `json:"store_ids"`
	Menus              interface{} `json:"menus"`
}
