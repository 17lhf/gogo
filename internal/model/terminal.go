package model

import "time"

// Terminal represents an edge device registered in the system.
type Terminal struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	SN               string     `gorm:"size:64;uniqueIndex;not null" json:"sn"`
	Name             string     `gorm:"size:128;not null" json:"name"`
	Type             string     `gorm:"size:64" json:"type"`
	StoreID          int64      `gorm:"index" json:"store_id"`
	Status           string     `gorm:"size:16;default:offline;not null" json:"status"`
	IPAddress        string     `gorm:"size:45" json:"ip_address"`
	MACAddress       string     `gorm:"size:17" json:"mac_address"`
	DeviceToken      string     `gorm:"size:256;not null" json:"-"`
	LastHeartbeatAt  *time.Time `json:"last_heartbeat_at"`
	CreatedAt        time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"not null" json:"updated_at"`

	Store Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

// TableName overrides the default table name.
func (Terminal) TableName() string { return "terminals" }
