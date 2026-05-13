package service

import (
	"context"
	"fmt"

	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/pkg"
	"gogo/internal/repository"
)

// UserService handles user management business logic.
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Create creates a new user with validated password.
func (s *UserService) Create(ctx context.Context, req dto.CreateUserReq) (*model.User, error) {
	if err := pkg.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check username uniqueness
	existing, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("check username: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// Check email uniqueness
	existing, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("邮箱已存在")
	}

	hash, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hash,
		RealName: req.RealName,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

// GetByID returns a user by ID.
func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return user, nil
}

// List returns paginated users.
func (s *UserService) List(ctx context.Context, req dto.UserListReq) ([]model.User, int64, error) {
	return s.userRepo.List(ctx, req)
}

// Update updates a user's profile fields.
func (s *UserService) Update(ctx context.Context, id int64, req dto.UpdateUserReq) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}

	if req.Email != "" && req.Email != user.Email {
		existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if existing != nil && existing.ID != id {
			return fmt.Errorf("邮箱已存在")
		}
		user.Email = req.Email
	}
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	return s.userRepo.Update(ctx, user)
}

// Delete removes a user.
func (s *UserService) Delete(ctx context.Context, id int64) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}
	return s.userRepo.Delete(ctx, id)
}

// ResetPassword resets a user's password (admin operation).
func (s *UserService) ResetPassword(ctx context.Context, id int64, password string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}

	if err := pkg.ValidatePasswordStrength(password); err != nil {
		return err
	}

	hash, err := pkg.HashPassword(password)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, id, hash, true)
}

// AssignRoles assigns roles to a user.
func (s *UserService) AssignRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}
	return s.userRepo.SetRoles(ctx, userID, roleIDs)
}

// AssignStores assigns stores to a user.
func (s *UserService) AssignStores(ctx context.Context, userID int64, storeIDs []int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}
	return s.userRepo.SetStores(ctx, userID, storeIDs)
}
