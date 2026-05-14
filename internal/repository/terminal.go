package repository

import (
	"context"
	"errors"
	"time"

	"gogo/internal/dto"
	"gogo/internal/model"

	"gorm.io/gorm"
)

// TerminalRepository defines the data access interface for terminals.
type TerminalRepository interface {
	Create(ctx context.Context, terminal *model.Terminal) error
	GetByID(ctx context.Context, id int64) (*model.Terminal, error)
	GetBySN(ctx context.Context, sn string) (*model.Terminal, error)
	List(ctx context.Context, req dto.TerminalListReq, storeIDs []int64) ([]model.Terminal, int64, error)
	Update(ctx context.Context, terminal *model.Terminal) error
	Delete(ctx context.Context, id int64) error
	UpdateHeartbeat(ctx context.Context, id int64, ip string) error
	UpdateStatus(ctx context.Context, id int64, status model.TerminalStatus) error
	UpdateStatusBySN(ctx context.Context, sn string, status model.TerminalStatus) error
	UpdateDeviceToken(ctx context.Context, id int64, token string) error
}

type terminalRepository struct {
	db *gorm.DB
}

// NewTerminalRepository creates a new TerminalRepository.
func NewTerminalRepository(db *gorm.DB) TerminalRepository {
	return &terminalRepository{db: db}
}

func (r *terminalRepository) Create(ctx context.Context, terminal *model.Terminal) error {
	return r.db.WithContext(ctx).Create(terminal).Error
}

func (r *terminalRepository) GetByID(ctx context.Context, id int64) (*model.Terminal, error) {
	var t model.Terminal
	err := r.db.WithContext(ctx).Preload("Store").First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (r *terminalRepository) GetBySN(ctx context.Context, sn string) (*model.Terminal, error) {
	var t model.Terminal
	err := r.db.WithContext(ctx).Preload("Store").Where("sn = ?", sn).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (r *terminalRepository) List(ctx context.Context, req dto.TerminalListReq, storeIDs []int64) ([]model.Terminal, int64, error) {
	var terminals []model.Terminal
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Terminal{})
	if req.SN != "" {
		query = query.Where("sn ILIKE ?", "%"+req.SN+"%")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StoreID != nil {
		query = query.Where("store_id = ?", *req.StoreID)
	}
	if len(storeIDs) > 0 {
		query = query.Where("store_id IN ?", storeIDs)
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
	if err := query.Preload("Store").Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&terminals).Error; err != nil {
		return nil, 0, err
	}

	return terminals, total, nil
}

func (r *terminalRepository) Update(ctx context.Context, terminal *model.Terminal) error {
	return r.db.WithContext(ctx).Save(terminal).Error
}

func (r *terminalRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Terminal{}, id).Error
}

func (r *terminalRepository) UpdateHeartbeat(ctx context.Context, id int64, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.Terminal{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":              model.TerminalStatusOnline,
		"ip_address":          ip,
		"last_heartbeat_at":   now,
	}).Error
}

func (r *terminalRepository) UpdateStatus(ctx context.Context, id int64, status model.TerminalStatus) error {
	return r.db.WithContext(ctx).Model(&model.Terminal{}).Where("id = ?", id).Update("status", status).Error
}

func (r *terminalRepository) UpdateStatusBySN(ctx context.Context, sn string, status model.TerminalStatus) error {
	return r.db.WithContext(ctx).Model(&model.Terminal{}).Where("sn = ?", sn).Update("status", status).Error
}

func (r *terminalRepository) UpdateDeviceToken(ctx context.Context, id int64, token string) error {
	return r.db.WithContext(ctx).Model(&model.Terminal{}).Where("id = ?", id).Update("device_token", token).Error
}
