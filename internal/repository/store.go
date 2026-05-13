package repository

import (
	"context"
	"errors"

	"gogo/internal/dto"
	"gogo/internal/model"

	"gorm.io/gorm"
)

// StoreRepository defines the data access interface for stores.
type StoreRepository interface {
	Create(ctx context.Context, store *model.Store) error
	GetByID(ctx context.Context, id int64) (*model.Store, error)
	List(ctx context.Context, req dto.StoreListReq) ([]model.Store, int64, error)
	Update(ctx context.Context, store *model.Store) error
	Delete(ctx context.Context, id int64) error
	HasTerminals(ctx context.Context, storeID int64) (bool, error)
}

type storeRepository struct {
	db *gorm.DB
}

// NewStoreRepository creates a new StoreRepository.
func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(ctx context.Context, store *model.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *storeRepository) GetByID(ctx context.Context, id int64) (*model.Store, error) {
	var store model.Store
	err := r.db.WithContext(ctx).First(&store, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &store, err
}

func (r *storeRepository) List(ctx context.Context, req dto.StoreListReq) ([]model.Store, int64, error) {
	var stores []model.Store
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Store{})
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
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
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

func (r *storeRepository) Update(ctx context.Context, store *model.Store) error {
	return r.db.WithContext(ctx).Save(store).Error
}

func (r *storeRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Store{}, id).Error
}

func (r *storeRepository) HasTerminals(ctx context.Context, storeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Terminal{}).Where("store_id = ?", storeID).Count(&count).Error
	return count > 0, err
}
