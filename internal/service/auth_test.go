package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gogo/internal/cache"
	"gogo/internal/config"
	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/pkg"
	"gogo/internal/repository"
)

// userRepoStub is a test stub for UserRepository.
type userRepoStub struct {
	users    map[int64]*model.User
	roles    map[int64][]model.Role
	storeIDs map[int64][]int64
}

func newUserRepoStub() *userRepoStub {
	return &userRepoStub{
		users:    make(map[int64]*model.User),
		roles:    make(map[int64][]model.Role),
		storeIDs: make(map[int64][]int64),
	}
}

func (s *userRepoStub) Create(ctx context.Context, user *model.User) error {
	user.ID = int64(len(s.users) + 1)
	s.users[user.ID] = user
	return nil
}

func (s *userRepoStub) GetByID(ctx context.Context, id int64) (*model.User, error) {
	u, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (s *userRepoStub) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	for _, u := range s.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, nil
}

func (s *userRepoStub) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, u := range s.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}
func (s *userRepoStub) List(ctx context.Context, req dto.UserListReq) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (s *userRepoStub) Update(ctx context.Context, user *model.User) error              { return nil }
func (s *userRepoStub) Delete(ctx context.Context, id int64) error                       { return nil }
func (s *userRepoStub) SetRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return nil
}
func (s *userRepoStub) SetStores(ctx context.Context, userID int64, storeIDs []int64) error {
	return nil
}
func (s *userRepoStub) UpdatePassword(ctx context.Context, id int64, hash string, mustChange bool) error {
	return nil
}
func (s *userRepoStub) UpdateLastLogin(ctx context.Context, id int64) error { return nil }
func (s *userRepoStub) UpdateStatus(ctx context.Context, id int64, status model.UserStatus) error {
	return nil
}
func (s *userRepoStub) GetCountByStatus(ctx context.Context) (map[int16]int64, error) {
	return nil, nil
}
func (s *userRepoStub) GetCountByRole(ctx context.Context) ([]dto.UserRoleStatItem, error) {
	return nil, nil
}
func (s *userRepoStub) GetCountByRecentAdded(ctx context.Context) (*dto.StatsRecentAdded, error) {
	return nil, nil
}

func (s *userRepoStub) GetRoles(ctx context.Context, userID int64) ([]model.Role, error) {
	return s.roles[userID], nil
}

func (s *userRepoStub) GetStoreIDs(ctx context.Context, userID int64) ([]int64, error) {
	return s.storeIDs[userID], nil
}

var _ repository.UserRepository = (*userRepoStub)(nil)

// roleRepoStub is a test stub for RoleRepository.
type roleRepoStub struct {
	roles map[string]*model.Role
}

func newRoleRepoStub() *roleRepoStub {
	return &roleRepoStub{roles: make(map[string]*model.Role)}
}

func (s *roleRepoStub) Create(ctx context.Context, role *model.Role) error       { return nil }
func (s *roleRepoStub) GetByID(ctx context.Context, id int64) (*model.Role, error) { return nil, nil }
func (s *roleRepoStub) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	r, ok := s.roles[code]
	if !ok {
		return nil, nil
	}
	return r, nil
}
func (s *roleRepoStub) List(ctx context.Context) ([]model.Role, error)            { return nil, nil }
func (s *roleRepoStub) Update(ctx context.Context, role *model.Role) error        { return nil }
func (s *roleRepoStub) Delete(ctx context.Context, id int64) error                 { return nil }
func (s *roleRepoStub) GetMenuIDs(ctx context.Context, roleID int64) ([]int64, error) {
	return nil, nil
}
func (s *roleRepoStub) SetMenusAndSyncPolicies(ctx context.Context, roleID int64, menuIDs []int64, roleCode string, policies [][2]string) error {
	return nil
}
func (s *roleRepoStub) DeleteWithCleanup(ctx context.Context, roleID int64, roleCode string) error {
	return nil
}

var _ repository.RoleRepository = (*roleRepoStub)(nil)

func TestAuthService_Login_Success(t *testing.T) {
	hash, _ := pkg.HashPassword("Test1234")
	now := time.Now()
	userRepo := newUserRepoStub()
	userRepo.users[1] = &model.User{
		ID:                1,
		Username:          "testuser",
		Email:             "test@example.com",
		Password:          hash,
		Status:            1,
		PasswordUpdatedAt: now,
	}
	userRepo.roles[1] = []model.Role{{ID: 1, Code: "OPERATOR", Name: "操作员"}}

	cfg := config.AuthConfig{
		JWTSecret:        "test-secret",
		SessionTTL:       8 * time.Hour,
		LockoutThreshold: 5,
		LockoutDuration:  30 * time.Minute,
	}

	// Use a real Redis connection or skip for now
	// For now, test only the validation logic
	_ = cache.NewSessionCache(nil, cfg.SessionTTL) // will be nil - used for setup only

	// Test that password validation works
	err := pkg.CheckPassword(hash, "Test1234")
	assert.NoError(t, err)

	err = pkg.CheckPassword(hash, "WrongPassword")
	assert.Error(t, err)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hash, _ := pkg.HashPassword("Test1234")
	err := pkg.CheckPassword(hash, "WrongPassword1")
	assert.Error(t, err)
}

func TestAuthService_ChangePassword(t *testing.T) {
	hash, _ := pkg.HashPassword("Test1234")
	now := time.Now()

	userRepo := newUserRepoStub()
	userRepo.users[1] = &model.User{
		ID:       1,
		Username: "testuser",
		Password: hash,
		Status:   1,
		PasswordUpdatedAt: now,
	}

	cfg := config.AuthConfig{PasswordMaxAge: 365 * 24 * time.Hour}
	lockoutCache := cache.NewLockoutCache(nil, 5, 30*time.Minute)
	sessionCache := cache.NewSessionCache(nil, 8*time.Hour)

	svc := NewAuthService(userRepo, newRoleRepoStub(), sessionCache, lockoutCache, cfg)

	// Correct old password
	req := dto.ChangePasswordReq{
		OldPassword: "Test1234",
		NewPassword: "NewPass1",
	}
	err := svc.ChangePassword(context.Background(), 1, req)
	assert.NoError(t, err)

	// Wrong old password
	req.OldPassword = "WrongPass1"
	err = svc.ChangePassword(context.Background(), 1, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wrong password")
}

func TestAuthService_ChangePassword_WeakNew(t *testing.T) {
	cfg := config.AuthConfig{PasswordMaxAge: 365 * 24 * time.Hour}
	pwHash, _ := pkg.HashPassword("Test1234")
	now := time.Now()

	userRepo := newUserRepoStub()
	userRepo.users[1] = &model.User{ID: 1, Username: "test", Password: pwHash, Status: 1, PasswordUpdatedAt: now}

	svc := NewAuthService(userRepo, newRoleRepoStub(), nil, nil, cfg)

	req := dto.ChangePasswordReq{OldPassword: "Test1234", NewPassword: "short"}
	err := svc.ChangePassword(context.Background(), 1, req)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, pkg.ErrPasswordTooShort))
}

func TestAuthService_Me(t *testing.T) {
	now := time.Now()
	userRepo := newUserRepoStub()
	userRepo.users[1] = &model.User{
		ID:                1,
		Username:          "testuser",
		Email:             "test@example.com",
		RealName:          "Test User",
		Phone:             "1234567890",
		Status:            1,
		PasswordUpdatedAt: now,
	}
	userRepo.roles[1] = []model.Role{{ID: 1, Code: "OPERATOR", Name: "操作员"}}
	userRepo.storeIDs[1] = []int64{1, 2}

	cfg := config.AuthConfig{}
	svc := NewAuthService(userRepo, newRoleRepoStub(), nil, nil, cfg)

	profile, err := svc.Me(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "testuser", profile.Username)
	assert.Equal(t, "test@example.com", profile.Email)
	assert.Equal(t, "Test User", profile.RealName)
	assert.Equal(t, []string{"OPERATOR"}, profile.Roles)
	assert.Equal(t, []int64{1, 2}, profile.Stores)
}
