package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/pkg"
)

func TestUserService_Create(t *testing.T) {
	repo := newUserRepoStub()
	svc := NewUserService(repo)

	req := dto.CreateUserReq{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "Test1234",
		RealName: "New User",
		Phone:    "1234567890",
	}

	user, err := svc.Create(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "newuser", user.Username)
	assert.Equal(t, "new@example.com", user.Email)
	assert.NotEmpty(t, user.Password)
}

func TestUserService_Create_DuplicateUsername(t *testing.T) {
	hash, _ := pkg.HashPassword("Test1234")
	repo := newUserRepoStub()
	repo.users[1] = &model.User{ID: 1, Username: "existing", Email: "old@example.com", Password: hash, Status: 1}

	svc := NewUserService(repo)

	req := dto.CreateUserReq{
		Username: "existing",
		Email:    "new@example.com",
		Password: "Test1234",
	}
	_, err := svc.Create(context.Background(), req)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrUsernameExists))
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	hash, _ := pkg.HashPassword("Test1234")
	repo := newUserRepoStub()
	repo.users[1] = &model.User{ID: 1, Username: "user1", Email: "existing@example.com", Password: hash, Status: 1}

	svc := NewUserService(repo)

	req := dto.CreateUserReq{
		Username: "user2",
		Email:    "existing@example.com",
		Password: "Test1234",
	}
	_, err := svc.Create(context.Background(), req)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmailExists))
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := newUserRepoStub()
	svc := NewUserService(repo)

	_, err := svc.GetByID(context.Background(), 999)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrUserNotFound))
}

func TestUserService_Delete_NotFound(t *testing.T) {
	repo := newUserRepoStub()
	svc := NewUserService(repo)

	err := svc.Delete(context.Background(), 999)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrUserNotFound))
}

func TestUserService_ResetPassword(t *testing.T) {
	hash, _ := pkg.HashPassword("Test1234")
	repo := newUserRepoStub()
	repo.users[1] = &model.User{ID: 1, Username: "user1", Email: "test@example.com", Password: hash, Status: 1}

	svc := NewUserService(repo)

	err := svc.ResetPassword(context.Background(), 1, "NewPass1")
	assert.NoError(t, err)
}

func TestUserService_AssignRoles(t *testing.T) {
	hash, _ := pkg.HashPassword("Test1234")
	repo := newUserRepoStub()
	repo.users[1] = &model.User{ID: 1, Username: "user1", Email: "test@example.com", Password: hash, Status: 1}

	svc := NewUserService(repo)

	err := svc.AssignRoles(context.Background(), 1, []int64{1, 2})
	assert.NoError(t, err)
}

func TestUserService_AssignStores(t *testing.T) {
	hash, _ := pkg.HashPassword("Test1234")
	repo := newUserRepoStub()
	repo.users[1] = &model.User{ID: 1, Username: "user1", Email: "test@example.com", Password: hash, Status: 1}

	svc := NewUserService(repo)

	err := svc.AssignStores(context.Background(), 1, []int64{1, 2})
	assert.NoError(t, err)
}

func TestUserService_Update_NotFound(t *testing.T) {
	repo := newUserRepoStub()
	svc := NewUserService(repo)

	err := svc.Update(context.Background(), 999, dto.UpdateUserReq{})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrUserNotFound))
}
