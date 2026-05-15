package repository

import (
	"context"
	"errors"

	"gogo/internal/model"

	"gorm.io/gorm"
)

// MenuRepository defines the data access interface for menus.
type MenuRepository interface {
	Create(ctx context.Context, menu *model.Menu) error
	GetByID(ctx context.Context, id int64) (*model.Menu, error)
	List(ctx context.Context) ([]model.Menu, error)
	Update(ctx context.Context, menu *model.Menu) error
	Delete(ctx context.Context, id int64) error
	HasChildren(ctx context.Context, parentID int64) (bool, error)
	GetMenusByRoleID(ctx context.Context, roleID int64) ([]int64, error)
	GetButtonAPIsByIDs(ctx context.Context, ids []int64) (map[int64][2]string, error)
}

type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new MenuRepository.
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(ctx context.Context, menu *model.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *menuRepository) GetByID(ctx context.Context, id int64) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.WithContext(ctx).First(&menu, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &menu, err
}

func (r *menuRepository) List(ctx context.Context) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) Update(ctx context.Context, menu *model.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

func (r *menuRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Menu{}, id).Error
}

func (r *menuRepository) HasChildren(ctx context.Context, parentID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Menu{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count > 0, err
}

func (r *menuRepository) GetMenusByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var menuIDs []int64
	err := r.db.WithContext(ctx).Table("role_menus").
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

func (r *menuRepository) GetButtonAPIsByIDs(ctx context.Context, ids []int64) (map[int64][2]string, error) {
	if len(ids) == 0 {
		return map[int64][2]string{}, nil
	}
	type row struct {
		ID        int64
		ApiPath   string
		ApiMethod string
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&model.Menu{}).
		Select("id, api_path, api_method").
		Where("id IN ? AND type = ? AND api_path != '' AND api_method != ''", ids, model.MenuTypeButton).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64][2]string, len(rows))
	for _, r := range rows {
		result[r.ID] = [2]string{r.ApiPath, r.ApiMethod}
	}
	return result, nil
}
