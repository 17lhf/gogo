package repository

import (
	"context"
	"errors"
	"log/slog"
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
	GetCountByStatus(ctx context.Context) (*dto.StatsStatusDistribution, error)
	GetCountByStore(ctx context.Context) ([]dto.StatsByStore, error)
	GetCountByRecentAdded(ctx context.Context) (*dto.StatsRecentAdded, error)
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

func (r *terminalRepository) GetCountByStatus(ctx context.Context) (*dto.StatsStatusDistribution, error) {
	var stats dto.StatsStatusDistribution
	err := r.db.WithContext(ctx).Raw(
		`select count(*) filter (where status = ?) as online,
			count(*) filter (where status = ?) as offline,
			count(*) filter (where status = ?) as enabled,
			count(*) filter (where status = ?) as disabled
		from terminals`,
		model.TerminalStatusOnline, model.TerminalStatusOffline,
		model.TerminalStatusEnabled, model.TerminalStatusDisabled,
	).Scan(&stats).Error
	if err != nil {
		slog.Error("failed to get terminal count by status", "error", err)
	}
	return &stats, err
}

func (r *terminalRepository) GetCountByStore(ctx context.Context) ([]dto.StatsByStore, error) {
	var stats []dto.StatsByStore
	err := r.db.WithContext(ctx).Raw(
		`select s.id as store_id, s.name as store_name,
            count(t.id) as total,
            count(t.id) filter (where t.status = ?) as online,
            count(t.id) filter (where t.status = ?) as offline,
            count(t.id) filter (where t.status = ?) as enabled,
            count(t.id) filter (where t.status = ?) as disabled
        from stores s
        left join terminals t on t.store_id = s.id
        group by s.id, s.name`,
		model.TerminalStatusOnline, model.TerminalStatusOffline,
		model.TerminalStatusEnabled, model.TerminalStatusDisabled,
	).Scan(&stats).Error
	if err != nil {
		slog.Error("failed to get terminal count by store", "error", err)
	}
	return stats, err
}

func (r *terminalRepository) GetCountByRecentAdded(ctx context.Context) (*dto.StatsRecentAdded, error) {
	var stats dto.StatsRecentAdded
	err := r.db.WithContext(ctx).Raw(
		`select count(*) filter (where created_at >= now() - interval '7 days') as last7_days,
			count(*) filter (where created_at >= now() - interval '30 days') as last30_days
		from terminals`,
	).Scan(&stats).Error
	if err != nil {
		slog.Error("failed to get terminal count by recent added", "error", err)
	}
	return &stats, err
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
		"status":            model.TerminalStatusOnline,
		"ip_address":        ip,
		"last_heartbeat_at": now,
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
