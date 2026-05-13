package model

import "time"

// User represents a system user.
type User struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username           string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Email              string     `gorm:"size:128;uniqueIndex;not null" json:"email"`
	Password           string     `gorm:"size:256;not null" json:"-"`
	RealName           string     `gorm:"size:64" json:"real_name"`
	Phone              string     `gorm:"size:20" json:"phone"`
	Status             int16      `gorm:"default:1;not null" json:"status"`
	MustChangePassword bool       `gorm:"default:false;not null" json:"must_change_password"`
	PasswordUpdatedAt  time.Time  `gorm:"not null" json:"password_updated_at"`
	LastLoginAt        *time.Time `json:"last_login_at"`
	CreatedAt          time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"not null" json:"updated_at"`

	Roles  []Role  `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Stores []Store `gorm:"many2many:user_stores;" json:"stores,omitempty"`
}

// TableName overrides the default table name.
func (User) TableName() string { return "users" }
