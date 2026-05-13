package model

import (
	"encoding/json"
	"time"
)

// TerminalLog records terminal events (online/offline/heartbeat_timeout/disabled/enabled).
type TerminalLog struct {
	ID         int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	TerminalID *int64          `json:"terminal_id"`
	SN         string          `gorm:"size:64;not null" json:"sn"`
	EventType  string          `gorm:"size:32;not null" json:"event_type"`
	Detail     json.RawMessage `gorm:"type:jsonb" json:"detail"`
	CreatedAt  time.Time       `gorm:"not null;index" json:"created_at"`
}

// TableName overrides the default table name.
func (TerminalLog) TableName() string { return "terminal_logs" }
