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

// UserRepository defines the data access interface for users.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, req dto.UserListReq) ([]model.User, int64, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error

	GetRoles(ctx context.Context, userID int64) ([]model.Role, error)
	SetRoles(ctx context.Context, userID int64, roleIDs []int64) error
	GetStoreIDs(ctx context.Context, userID int64) ([]int64, error)
	SetStores(ctx context.Context, userID int64, storeIDs []int64) error

	UpdatePassword(ctx context.Context, id int64, hash string, mustChange bool) error
	UpdateLastLogin(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status model.UserStatus) error

	GetCountByStatus(ctx context.Context) (map[int16]int64, error)
	GetCountByRole(ctx context.Context) ([]dto.UserRoleStatItem, error)
	GetCountByRecentAdded(ctx context.Context) (*dto.StatsRecentAdded, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Stores").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Roles").Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) List(ctx context.Context, req dto.UserListReq) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})
	if req.Username != "" {
		query = query.Where("username ILIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", int16(*req.Status))
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
	if err := query.Preload("Roles").Preload("Stores").Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) GetRoles(ctx context.Context, userID int64) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *userRepository) SetRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&struct {
			UserID int64 `gorm:"column:user_id"`
			RoleID int64 `gorm:"column:role_id"`
		}{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepository) GetStoreIDs(ctx context.Context, userID int64) ([]int64, error) {
	var storeIDs []int64
	err := r.db.WithContext(ctx).Table("user_stores").
		Where("user_id = ?", userID).
		Pluck("store_id", &storeIDs).Error
	return storeIDs, err
}

func (r *userRepository) SetStores(ctx context.Context, userID int64, storeIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&struct {
			UserID  int64 `gorm:"column:user_id"`
			StoreID int64 `gorm:"column:store_id"`
		}{}).Error; err != nil {
			return err
		}
		for _, storeID := range storeIDs {
			if err := tx.Exec("INSERT INTO user_stores (user_id, store_id) VALUES (?, ?)", userID, storeID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int64, hash string, mustChange bool) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password":             hash,
		"must_change_password": mustChange,
		"password_updated_at":  time.Now(),
	}).Error
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}

func (r *userRepository) UpdateStatus(ctx context.Context, id int64, status model.UserStatus) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *userRepository) GetCountByStatus(ctx context.Context) (map[int16]int64, error) {
	type statusRow struct {
		Status int16 `gorm:"column:status"`
		Count  int64 `gorm:"column:count"`
	}
	var rows []statusRow
	err := r.db.WithContext(ctx).Raw(
		"SELECT status, COUNT(*) as count FROM users GROUP BY status",
	).Scan(&rows).Error
	if err != nil {
		slog.Error("failed to get user count by status", "error", err)
		return nil, err
	}
	result := make(map[int16]int64, len(rows))
	for _, row := range rows {
		result[row.Status] = row.Count
	}
	return result, nil
}

func (r *userRepository) GetCountByRole(ctx context.Context) ([]dto.UserRoleStatItem, error) {
	var stats []dto.UserRoleStatItem
	err := r.db.WithContext(ctx).Raw(
		`SELECT r.id as role_id, r.name as role_name, COUNT(ur.user_id) as count
		FROM roles r
		LEFT JOIN user_roles ur ON ur.role_id = r.id
		GROUP BY r.id, r.name
		ORDER BY r.id`,
	).Scan(&stats).Error
	if err != nil {
		slog.Error("failed to get user count by role", "error", err)
	}
	return stats, err
}

func (r *userRepository) GetCountByRecentAdded(ctx context.Context) (*dto.StatsRecentAdded, error) {
	var stats dto.StatsRecentAdded
	err := r.db.WithContext(ctx).Raw(
		`SELECT count(*) filter (where created_at >= now() - interval '7 days') as last7_days,
			count(*) filter (where created_at >= now() - interval '30 days') as last30_days
		FROM users`,
	).Scan(&stats).Error
	if err != nil {
		slog.Error("failed to get user count by recent added", "error", err)
	}
	return &stats, err
}
