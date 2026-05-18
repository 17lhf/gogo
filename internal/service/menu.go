package service

import (
	"context"
	"errors"
	"fmt"

	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/repository"
)

// Sentinel errors for menu service.
var (
	ErrMenuNotFound       = errors.New("menu not found")
	ErrParentMenuNotFound = errors.New("parent menu not found")
	ErrMenuHasChildren    = errors.New("menu has children")
)

// MenuService handles menu management business logic.
type MenuService struct {
	menuRepo repository.MenuRepository
	userRepo repository.UserRepository
}

// NewMenuService creates a new MenuService.
func NewMenuService(menuRepo repository.MenuRepository, userRepo repository.UserRepository) *MenuService {
	return &MenuService{menuRepo: menuRepo, userRepo: userRepo}
}

// Create creates a new menu item.
func (s *MenuService) Create(ctx context.Context, req dto.CreateMenuReq) (*model.Menu, error) {
	if req.ParentID != 0 {
		parent, err := s.menuRepo.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, ErrParentMenuNotFound
		}
	}

	menu := &model.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Type:      req.Type,
		Perms:     req.Perms,
		ApiPath:   req.ApiPath,
		ApiMethod: req.ApiMethod,
		SortOrder: req.SortOrder,
		Visible:   true,
		Status:    int16(model.UserStatusEnabled),
	}
	if err := s.menuRepo.Create(ctx, menu); err != nil {
		return nil, fmt.Errorf("create menu: %w", err)
	}
	return menu, nil
}

// GetByID returns a menu by ID.
func (s *MenuService) GetByID(ctx context.Context, id int64) (*model.Menu, error) {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, ErrMenuNotFound
	}
	return menu, nil
}

// Tree returns the full menu tree.
func (s *MenuService) Tree(ctx context.Context) ([]*model.Menu, error) {
	menus, err := s.menuRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return buildTree(menus, 0), nil
}

// TreeByRoleIDs returns the menu tree filtered by role menu assignments.
func (s *MenuService) TreeByRoleIDs(ctx context.Context, roleIDs []int64) ([]*model.Menu, error) {
	menuIDs := make(map[int64]bool)
	for _, roleID := range roleIDs {
		ids, err := s.menuRepo.GetMenusByRoleID(ctx, roleID)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			menuIDs[id] = true
		}
	}

	menus, err := s.menuRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Filter menus
	var filtered []model.Menu
	for _, m := range menus {
		if menuIDs[m.ID] {
			filtered = append(filtered, m)
		}
	}

	return buildTree(filtered, 0), nil
}

// TreeByUserID returns the menu tree for the given user, based on their roles.
func (s *MenuService) TreeByUserID(ctx context.Context, userID int64) ([]*model.Menu, error) {
	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	roleIDs := make([]int64, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}
	return s.TreeByRoleIDs(ctx, roleIDs)
}

// Update updates a menu.
func (s *MenuService) Update(ctx context.Context, id int64, req dto.UpdateMenuReq) error {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if menu == nil {
		return ErrMenuNotFound
	}

	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Perms != "" {
		menu.Perms = req.Perms
	}
	if req.ApiPath != "" {
		menu.ApiPath = req.ApiPath
	}
	if req.ApiMethod != "" {
		menu.ApiMethod = req.ApiMethod
	}
	if req.ParentID != nil {
		menu.ParentID = *req.ParentID
	}
	if req.SortOrder != nil {
		menu.SortOrder = *req.SortOrder
	}
	if req.Visible != nil {
		menu.Visible = *req.Visible
	}

	return s.menuRepo.Update(ctx, menu)
}

// Delete removes a menu. Fails if it has children.
func (s *MenuService) Delete(ctx context.Context, id int64) error {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if menu == nil {
		return ErrMenuNotFound
	}

	hasChildren, err := s.menuRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return ErrMenuHasChildren
	}

	return s.menuRepo.Delete(ctx, id)
}

func buildTree(menus []model.Menu, parentID int64) []*model.Menu {
	var tree []*model.Menu
	for i := range menus {
		if menus[i].ParentID == parentID {
			node := &menus[i]
			node.Children = buildTree(menus, node.ID)
			tree = append(tree, node)
		}
	}
	return tree
}
