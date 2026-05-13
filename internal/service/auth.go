package service

import (
	"context"
	"fmt"

	"gogo/internal/cache"
	"gogo/internal/config"
	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/pkg"
	"gogo/internal/repository"
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	sessionCache *cache.SessionCache
	lockoutCache *cache.LockoutCache
	cfg          config.AuthConfig
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	sessionCache *cache.SessionCache,
	lockoutCache *cache.LockoutCache,
	cfg config.AuthConfig,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		sessionCache: sessionCache,
		lockoutCache: lockoutCache,
		cfg:          cfg,
	}
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, req dto.LoginReq) (*dto.LoginResp, error) {
	// Check lockout
	locked, err := s.lockoutCache.IsLocked(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("check lockout: %w", err)
	}
	if locked {
		return nil, fmt.Errorf("%w: account locked", ErrAccountLocked)
	}

	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		s.lockoutCache.RecordFailure(ctx, req.Username)
		return nil, fmt.Errorf("%w: invalid credentials", ErrInvalidCredentials)
	}

	if user.Status != 1 {
		s.lockoutCache.RecordFailure(ctx, req.Username)
		switch user.Status {
		case 2:
			return nil, fmt.Errorf("%w: account disabled", ErrAccountDisabled)
		case 3:
			return nil, fmt.Errorf("%w: account locked", ErrAccountLocked)
		default:
			return nil, fmt.Errorf("%w: unknown status", ErrAccountDisabled)
		}
	}

	if err := pkg.CheckPassword(user.Password, req.Password); err != nil {
		locked, _ := s.lockoutCache.RecordFailure(ctx, req.Username)
		if locked {
			s.userRepo.UpdateStatus(ctx, user.ID, 3)
			return nil, fmt.Errorf("%w: account locked after %d failures", ErrAccountLocked, s.cfg.LockoutThreshold)
		}
		return nil, fmt.Errorf("%w: invalid credentials", ErrInvalidCredentials)
	}

	// Success - clear lockout
	s.lockoutCache.Reset(ctx, req.Username)

	// Gather role codes
	roles, err := s.userRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get roles: %w", err)
	}
	roleCodes := make([]string, len(roles))
	for i, role := range roles {
		roleCodes[i] = role.Code
	}

	// Generate token
	token, jti, err := pkg.GenerateToken(s.cfg.JWTSecret, user.ID, user.Username, roleCodes)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	// Store session in Redis
	claims, _ := pkg.ParseToken(s.cfg.JWTSecret, token)
	if err := s.sessionCache.Set(ctx, user.ID, jti, claims); err != nil {
		return nil, fmt.Errorf("store session: %w", err)
	}

	// Update last login
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &dto.LoginResp{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   8 * 3600,
	}, nil
}

// Logout removes the user's session from Redis.
func (s *AuthService) Logout(ctx context.Context, userID int64, jti string) error {
	return s.sessionCache.Delete(ctx, userID, jti)
}

// Me returns the current user's profile.
func (s *AuthService) Me(ctx context.Context, userID int64) (*dto.UserProfileResp, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	roleCodes := make([]string, len(roles))
	for i, role := range roles {
		roleCodes[i] = role.Code
	}

	storeIDs, err := s.userRepo.GetStoreIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	lastLogin := ""
	if user.LastLoginAt != nil {
		lastLogin = user.LastLoginAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return &dto.UserProfileResp{
		ID:                 user.ID,
		Username:           user.Username,
		Email:              user.Email,
		RealName:           user.RealName,
		Phone:              user.Phone,
		Status:             user.Status,
		MustChangePassword: user.MustChangePassword,
		PasswordUpdatedAt:  user.PasswordUpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastLoginAt:        lastLogin,
		Roles:              roleCodes,
		Stores:             storeIDs,
	}, nil
}

// ChangePassword validates the old password and updates to a new one.
func (s *AuthService) ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordReq) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if err := pkg.CheckPassword(user.Password, req.OldPassword); err != nil {
		return fmt.Errorf("%w: wrong password", ErrInvalidCredentials)
	}

	if err := pkg.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	hash, err := pkg.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, hash, false)
}

// Predefined service-level errors.
var (
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrAccountLocked      = fmt.Errorf("account locked")
	ErrAccountDisabled    = fmt.Errorf("account disabled")
)

// Ensure model is used for compilation.
var _ = (*model.User)(nil)
