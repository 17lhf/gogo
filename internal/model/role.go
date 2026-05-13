package model

import "time"

// Role represents a role with associated permissions.
type Role struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Code        string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Description string    `gorm:"size:256" json:"description"`
	Status      int16     `gorm:"default:1;not null" json:"status"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`

	Menus []Menu `gorm:"many2many:role_menus;" json:"menus,omitempty"`
}

// TableName overrides the default table name.
func (Role) TableName() string { return "roles" }
