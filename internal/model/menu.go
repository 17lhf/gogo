package model

import "time"

// Menu represents a menu item (directory, page, or button).
type Menu struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID  int64     `gorm:"default:0;index;not null" json:"parent_id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Path      string    `gorm:"size:256" json:"path"`
	Component string    `gorm:"size:256" json:"component"`
	Icon      string    `gorm:"size:64" json:"icon"`
	Type      MenuType  `gorm:"not null" json:"type"`
	Perms     string    `gorm:"size:128" json:"perms"`
	SortOrder int       `gorm:"default:0;not null" json:"sort_order"`
	Visible   bool      `gorm:"default:true;not null" json:"visible"`
	Status    int16     `gorm:"default:1;not null" json:"status"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Children []*Menu `gorm:"-" json:"children,omitempty"`
}

// TableName overrides the default table name.
func (Menu) TableName() string { return "menus" }
