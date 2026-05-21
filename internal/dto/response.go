package dto

import "gogo/internal/model"

// LoginResp is the response for successful login.
type LoginResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// UserProfileResp is the response for GET /api/v1/auth/me.
type UserProfileResp struct {
	ID                 int64            `json:"id"`
	Username           string           `json:"username"`
	Email              string           `json:"email"`
	RealName           string           `json:"real_name"`
	Phone              string           `json:"phone"`
	Status             model.UserStatus `json:"status"`
	MustChangePassword bool             `json:"must_change_password"`
	PasswordUpdatedAt  string           `json:"password_updated_at"`
	LastLoginAt        string           `json:"last_login_at"`
	Roles              []string         `json:"roles"`
	Stores             []int64          `json:"store_ids"`
	Menus              interface{}      `json:"menus"`
}

type StatsStatusDistribution struct {
	Online   int64 `json:"online"`
	Offline  int64 `json:"offline"`
	Enabled  int64 `json:"enabled"`
	Disabled int64 `json:"disabled"`
}

type StatsByStore struct {
	StoreID   int64  `json:"store_id"`
	StoreName string `json:"store_name"`
	Total     int64  `json:"total"`
	Online    int64  `json:"online"`
	Offline   int64  `json:"offline"`
	Enabled   int64  `json:"enabled"`
	Disabled  int64  `json:"disabled"`
}

type StatsRecentAdded struct {
	Last7Days  int64 `json:"last_7_days"`
	Last30Days int64 `json:"last_30_days"`
}

type StatsTerminalsResp struct {
	StatusDistribution StatsStatusDistribution `json:"status_distribution"`
	ByStore            []StatsByStore          `json:"by_store"`
	RecentAdded        StatsRecentAdded        `json:"recent_added"`
}
