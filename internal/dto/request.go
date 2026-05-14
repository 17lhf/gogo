package dto

import "gogo/internal/model"

// LoginReq is the request body for POST /api/v1/auth/login.
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordReq is the request body for PUT /api/v1/auth/password.
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// CreateUserReq is the request body for creating a user.
type CreateUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email,max=128"`
	Password string `json:"password" binding:"required,min=8"`
	RealName string `json:"real_name" binding:"max=64"`
	Phone    string `json:"phone" binding:"max=20"`
}

// UpdateUserReq is the request body for updating a user.
type UpdateUserReq struct {
	Email    string `json:"email" binding:"omitempty,email,max=128"`
	RealName string `json:"real_name" binding:"max=64"`
	Phone    string `json:"phone" binding:"max=20"`
	Status *model.UserStatus `json:"status" binding:"omitempty,userstatus"`
}

// ResetPasswordReq is the request body for admin password reset.
type ResetPasswordReq struct {
	Password string `json:"password" binding:"required,min=8"`
}

// AssignRolesReq is the request body for assigning roles to a user.
type AssignRolesReq struct {
	RoleIDs []int64 `json:"role_ids" binding:"required"`
}

// AssignStoresReq is the request body for assigning stores to a user.
type AssignStoresReq struct {
	StoreIDs []int64 `json:"store_ids" binding:"required"`
}

// UserListReq is the query parameters for listing users.
type UserListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Username string `form:"username"`
	Status   *model.UserStatus `form:"status"`
}

// CreateRoleReq is the request body for creating a role.
type CreateRoleReq struct {
	Name        string `json:"name" binding:"required,min=2,max=64"`
	Code        string `json:"code" binding:"required,min=2,max=64"`
	Description string `json:"description" binding:"max=256"`
}

// UpdateRoleReq is the request body for updating a role.
type UpdateRoleReq struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=64"`
	Description string `json:"description" binding:"max=256"`
}

// AssignMenusReq is the request body for assigning menus to a role.
type AssignMenusReq struct {
	MenuIDs []int64 `json:"menu_ids" binding:"required"`
}

// CreateMenuReq is the request body for creating a menu.
type CreateMenuReq struct {
	ParentID  int64  `json:"parent_id"`
	Name      string `json:"name" binding:"required,min=2,max=64"`
	Path      string `json:"path" binding:"max=256"`
	Component string `json:"component" binding:"max=256"`
	Icon      string `json:"icon" binding:"max=64"`
	Type      model.MenuType `json:"type" binding:"required,menutype"`
	Perms     string `json:"perms" binding:"max=128"`
	SortOrder int    `json:"sort_order"`
}

// UpdateMenuReq is the request body for updating a menu.
type UpdateMenuReq struct {
	ParentID  *int64 `json:"parent_id"`
	Name      string `json:"name" binding:"max=64"`
	Path      string `json:"path" binding:"max=256"`
	Component string `json:"component" binding:"max=256"`
	Icon      string `json:"icon" binding:"max=64"`
	Perms     string `json:"perms" binding:"max=128"`
	SortOrder *int   `json:"sort_order"`
	Visible   *bool  `json:"visible"`
}

// CreateStoreReq is the request body for creating a store.
type CreateStoreReq struct {
	Name    string `json:"name" binding:"required,min=2,max=128"`
	Code    string `json:"code" binding:"required,min=2,max=64"`
	Address string `json:"address" binding:"max=256"`
}

// UpdateStoreReq is the request body for updating a store.
type UpdateStoreReq struct {
	Name    string `json:"name" binding:"omitempty,min=2,max=128"`
	Address string `json:"address" binding:"max=256"`
}

// StoreListReq is the query parameters for listing stores.
type StoreListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Name     string `form:"name"`
}

// CreateTerminalReq is the request body for creating a terminal.
type CreateTerminalReq struct {
	SN      string `json:"sn" binding:"required,min=2,max=64"`
	Name    string `json:"name" binding:"required,min=2,max=128"`
	Type    string `json:"type" binding:"max=64"`
	StoreID int64  `json:"store_id" binding:"required"`
}

// UpdateTerminalReq is the request body for updating a terminal.
type UpdateTerminalReq struct {
	Name  string  `json:"name" binding:"omitempty,min=2,max=128"`
	Type  string  `json:"type" binding:"max=64"`
	StoreID *int64 `json:"store_id"`
	Status *model.TerminalStatus `json:"status" binding:"omitempty,terminalstatus"`
}

// TerminalListReq is the query parameters for listing terminals.
type TerminalListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	SN       string `form:"sn"`
	Status   model.TerminalStatus `form:"status"`
	StoreID  *int64 `form:"store_id"`
}

// HeartbeatReq is the request body for terminal heartbeat.
type HeartbeatReq struct {
	IPAddress  string `json:"ip_address"`
	MACAddress string `json:"mac_address"`
}

// RotateTokenReq is the request body for token rotation.
type RotateTokenReq struct {
	// No fields needed; auth is via X-Device-Token header
}

// OperationLogListReq is the query parameters for listing operation logs.
type OperationLogListReq struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	UserID    *int64 `form:"user_id"`
	Action    string `form:"action"`
	Status    *model.LogStatus `form:"status"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}

// TerminalLogListReq is the query parameters for listing terminal logs.
type TerminalLogListReq struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	SN        string `form:"sn"`
	EventType string `form:"event_type"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}
