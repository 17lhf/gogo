package service

import (
	"context"
	"errors"
	"fmt"

	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/repository"
)

// Sentinel errors for role service.
var (
	ErrRoleNotFound   = errors.New("role not found")
	ErrRoleCodeExists = errors.New("role code already exists")
)

// RoleService handles role management business logic.
type RoleService struct {
	roleRepo repository.RoleRepository
	menuRepo repository.MenuRepository
}

// NewRoleService creates a new RoleService.
func NewRoleService(roleRepo repository.RoleRepository, menuRepo repository.MenuRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo, menuRepo: menuRepo}
}

// Create creates a new role.
func (s *RoleService) Create(ctx context.Context, req dto.CreateRoleReq) (*model.Role, error) {
	existing, err := s.roleRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrRoleCodeExists
	}

	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      int16(model.UserStatusEnabled),
	}
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}
	return role, nil
}

// GetByID returns a role by ID.
func (s *RoleService) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

// List returns all roles.
func (s *RoleService) List(ctx context.Context) ([]model.Role, error) {
	return s.roleRepo.List(ctx)
}

// Update updates a role.
func (s *RoleService) Update(ctx context.Context, id int64, req dto.UpdateRoleReq) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	return s.roleRepo.Update(ctx, role)
}

// Delete removes a role.
func (s *RoleService) Delete(ctx context.Context, id int64) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}
	return s.roleRepo.Delete(ctx, id)
}

// GetMenus returns menu IDs for a role.
func (s *RoleService) GetMenus(ctx context.Context, id int64) ([]int64, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return s.roleRepo.GetMenuIDs(ctx, id)
}

// AssignMenus assigns menus to a role and triggers Casbin policy reload.
func (s *RoleService) AssignMenus(ctx context.Context, id int64, menuIDs []int64) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}
	return s.roleRepo.SetMenus(ctx, id, menuIDs)
}
