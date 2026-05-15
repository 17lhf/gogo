package service

import (
	"context"
	"errors"
	"fmt"

	casbinSDK "github.com/casbin/casbin/v2"

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
	enforcer *casbinSDK.Enforcer
}

// NewRoleService creates a new RoleService.
func NewRoleService(roleRepo repository.RoleRepository, menuRepo repository.MenuRepository, enforcer *casbinSDK.Enforcer) *RoleService {
	return &RoleService{roleRepo: roleRepo, menuRepo: menuRepo, enforcer: enforcer}
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

// Delete removes a role and all associated data (role_menus, user_roles, casbin_rule).
func (s *RoleService) Delete(ctx context.Context, id int64) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}
	if err := s.roleRepo.DeleteWithCleanup(ctx, id, role.Code); err != nil {
		return err
	}
	return s.enforcer.LoadPolicy()
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

// AssignMenus assigns menus to a role and syncs Casbin policies.
//
// It reads api_path and api_method directly from button-type menus and
// updates both role_menus and casbin_rule in a single transaction.
func (s *RoleService) AssignMenus(ctx context.Context, id int64, menuIDs []int64) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}

	oldIDs, err := s.roleRepo.GetMenuIDs(ctx, id)
	if err != nil {
		return fmt.Errorf("get old menu ids: %w", err)
	}

	idSet := make(map[int64]bool, len(oldIDs)+len(menuIDs))
	for _, mid := range oldIDs {
		idSet[mid] = true
	}
	for _, mid := range menuIDs {
		idSet[mid] = true
	}
	allIDs := make([]int64, 0, len(idSet))
	for mid := range idSet {
		allIDs = append(allIDs, mid)
	}

	buttonAPIs, err := s.menuRepo.GetButtonAPIsByIDs(ctx, allIDs)
	if err != nil {
		return fmt.Errorf("get button apis: %w", err)
	}

	policySet := make(map[[2]string]bool)
	for _, mid := range menuIDs {
		if api, ok := buttonAPIs[mid]; ok {
			policySet[api] = true
		}
	}

	policies := make([][2]string, 0, len(policySet))
	for pol := range policySet {
		policies = append(policies, pol)
	}

	if err := s.roleRepo.SetMenusAndSyncPolicies(ctx, id, menuIDs, role.Code, policies); err != nil {
		return fmt.Errorf("set menus and policies: %w", err)
	}

	return s.enforcer.LoadPolicy()
}
