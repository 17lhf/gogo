package model

import "time"

// Store represents a physical store location.
type Store struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Code      string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Address   string    `gorm:"size:256" json:"address"`
	Status    int16     `gorm:"default:1;not null" json:"status"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

// TableName overrides the default table name.
func (Store) TableName() string { return "stores" }
