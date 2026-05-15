package repository

import (
	"context"
	"errors"

	"gogo/internal/model"

	"gorm.io/gorm"
)

// RoleRepository defines the data access interface for roles.
type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	GetByID(ctx context.Context, id int64) (*model.Role, error)
	GetByCode(ctx context.Context, code string) (*model.Role, error)
	List(ctx context.Context) ([]model.Role, error)
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id int64) error

	GetMenuIDs(ctx context.Context, roleID int64) ([]int64, error)
	SetMenusAndSyncPolicies(ctx context.Context, roleID int64, menuIDs []int64, roleCode string, policies [][2]string) error
	DeleteWithCleanup(ctx context.Context, roleID int64, roleCode string) error
}

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository.
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

func (r *roleRepository) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

func (r *roleRepository) List(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.WithContext(ctx).Order("id ASC").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

func (r *roleRepository) GetMenuIDs(ctx context.Context, roleID int64) ([]int64, error) {
	var menuIDs []int64
	err := r.db.WithContext(ctx).Table("role_menus").
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

type casbinRuleRow struct {
	Ptype string `gorm:"column:p_type"`
	V0    string `gorm:"column:v0"`
	V1    string `gorm:"column:v1"`
	V2    string `gorm:"column:v2"`
}

func (casbinRuleRow) TableName() string { return "casbin_rule" }

type roleMenuRow struct {
	RoleID int64 `gorm:"column:role_id"`
	MenuID int64 `gorm:"column:menu_id"`
}

func (roleMenuRow) TableName() string { return "role_menus" }

type userRoleRow struct {
	UserID int64 `gorm:"column:user_id"`
	RoleID int64 `gorm:"column:role_id"`
}

func (userRoleRow) TableName() string { return "user_roles" }

// SetMenusAndSyncPolicies updates role_menus and casbin_rule atomically.
//
// Within a single GORM transaction it:
//  1. Replaces role_menus rows for this role
//  2. Replaces casbin_rule rows for this role
func (r *roleRepository) SetMenusAndSyncPolicies(ctx context.Context, roleID int64, menuIDs []int64, roleCode string, policies [][2]string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&roleMenuRow{}).Error; err != nil {
			return err
		}
		for _, menuID := range menuIDs {
			if err := tx.Exec("INSERT INTO role_menus (role_id, menu_id) VALUES (?, ?)", roleID, menuID).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("p_type = ? AND v0 = ?", "p", roleCode).Delete(&casbinRuleRow{}).Error; err != nil {
			return err
		}
		for _, pol := range policies {
			if err := tx.Exec("INSERT INTO casbin_rule (p_type, v0, v1, v2) VALUES (?, ?, ?, ?)", "p", roleCode, pol[0], pol[1]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteWithCleanup removes a role and all its associated data in one transaction.
func (r *roleRepository) DeleteWithCleanup(ctx context.Context, roleID int64, roleCode string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&roleMenuRow{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", roleID).Delete(&userRoleRow{}).Error; err != nil {
			return err
		}
		if err := tx.Where("p_type = ? AND v0 = ?", "p", roleCode).Delete(&casbinRuleRow{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Role{}, roleID).Error
	})
}
