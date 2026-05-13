package service

import (
	"context"
	"fmt"

	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/repository"
)

// StoreService handles store management business logic.
type StoreService struct {
	storeRepo repository.StoreRepository
}

// NewStoreService creates a new StoreService.
func NewStoreService(storeRepo repository.StoreRepository) *StoreService {
	return &StoreService{storeRepo: storeRepo}
}

// Create creates a new store.
func (s *StoreService) Create(ctx context.Context, req dto.CreateStoreReq) (*model.Store, error) {
	store := &model.Store{
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
		Status:  1,
	}
	if err := s.storeRepo.Create(ctx, store); err != nil {
		return nil, fmt.Errorf("create store: %w", err)
	}
	return store, nil
}

// GetByID returns a store by ID.
func (s *StoreService) GetByID(ctx context.Context, id int64) (*model.Store, error) {
	store, err := s.storeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, fmt.Errorf("门店不存在")
	}
	return store, nil
}

// List returns paginated stores.
func (s *StoreService) List(ctx context.Context, req dto.StoreListReq) ([]model.Store, int64, error) {
	return s.storeRepo.List(ctx, req)
}

// Update updates a store.
func (s *StoreService) Update(ctx context.Context, id int64, req dto.UpdateStoreReq) error {
	store, err := s.storeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if store == nil {
		return fmt.Errorf("门店不存在")
	}
	if req.Name != "" {
		store.Name = req.Name
	}
	if req.Address != "" {
		store.Address = req.Address
	}
	return s.storeRepo.Update(ctx, store)
}

// Delete removes a store. Fails if it has terminals.
func (s *StoreService) Delete(ctx context.Context, id int64) error {
	store, err := s.storeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if store == nil {
		return fmt.Errorf("门店不存在")
	}

	hasTerminals, err := s.storeRepo.HasTerminals(ctx, id)
	if err != nil {
		return err
	}
	if hasTerminals {
		return fmt.Errorf("该门店下存在终端，无法删除")
	}

	return s.storeRepo.Delete(ctx, id)
}
