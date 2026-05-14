package model

import (
	"encoding/json"
	"time"
)

// OperationLog records admin operations for audit purposes.
type OperationLog struct {
	ID           int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       *int64          `json:"user_id"`
	Username     string          `gorm:"size:64" json:"username"`
	Action       string          `gorm:"size:128;not null" json:"action"`
	ResourceType string          `gorm:"size:64" json:"resource_type"`
	ResourceID   string          `gorm:"size:64" json:"resource_id"`
	Detail       json.RawMessage `gorm:"type:jsonb" json:"detail"`
	IP           string          `gorm:"size:45" json:"ip"`
	UserAgent    string          `gorm:"size:512" json:"user_agent"`
	Status       LogStatus       `gorm:"not null" json:"status"`
	DurationMs   int             `json:"duration_ms"`
	CreatedAt    time.Time       `gorm:"not null;index" json:"created_at"`
}

// TableName overrides the default table name.
func (OperationLog) TableName() string { return "operation_logs" }
