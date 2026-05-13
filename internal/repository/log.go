package repository

import (
	"context"
	"time"

	"gogo/internal/dto"
	"gogo/internal/model"

	"gorm.io/gorm"
)

// LogRepository defines the data access interface for logs.
type LogRepository interface {
	CreateOperation(ctx context.Context, log *model.OperationLog) error
	ListOperations(ctx context.Context, req dto.OperationLogListReq) ([]model.OperationLog, int64, error)
	DeleteOperationBefore(ctx context.Context, before time.Time) error

	CreateTerminal(ctx context.Context, log *model.TerminalLog) error
	ListTerminals(ctx context.Context, req dto.TerminalLogListReq) ([]model.TerminalLog, int64, error)
	DeleteTerminalBefore(ctx context.Context, before time.Time) error
}

type logRepository struct {
	db *gorm.DB
}

// NewLogRepository creates a new LogRepository.
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

func (r *logRepository) CreateOperation(ctx context.Context, log *model.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *logRepository) ListOperations(ctx context.Context, req dto.OperationLogListReq) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.OperationLog{})
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *logRepository) DeleteOperationBefore(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&model.OperationLog{}).Error
}

func (r *logRepository) CreateTerminal(ctx context.Context, log *model.TerminalLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *logRepository) ListTerminals(ctx context.Context, req dto.TerminalLogListReq) ([]model.TerminalLog, int64, error) {
	var logs []model.TerminalLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TerminalLog{})
	if req.SN != "" {
		query = query.Where("sn = ?", req.SN)
	}
	if req.EventType != "" {
		query = query.Where("event_type = ?", req.EventType)
	}
	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *logRepository) DeleteTerminalBefore(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&model.TerminalLog{}).Error
}
