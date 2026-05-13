package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gogo/internal/dto"
	"gogo/internal/middleware"
	"gogo/internal/pkg"
	"gogo/internal/service"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authSvc *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}

	resp, err := h.authSvc.Login(c.Request.Context(), req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	pkg.Success(c, resp)
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	jti := middleware.GetSessionID(c)
	if err := h.authSvc.Logout(c.Request.Context(), userID, jti); err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, "登出失败")
		return
	}
	pkg.Success(c, nil)
}

// Me handles GET /api/v1/auth/me.
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	profile, err := h.authSvc.Me(c.Request.Context(), userID)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, "获取用户信息失败")
		return
	}
	pkg.Success(c, profile)
}

// ChangePassword handles PUT /api/v1/auth/password.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.authSvc.ChangePassword(c.Request.Context(), userID, req); err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "原密码错误")
			return
		}
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, err.Error())
		return
	}

	pkg.Success(c, nil)
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrAccountLocked):
		pkg.Error(c, http.StatusForbidden, pkg.CodeAccountLocked, "账户已锁定，请30分钟后重试")
	case errors.Is(err, service.ErrAccountDisabled):
		pkg.Error(c, http.StatusForbidden, pkg.CodeAccountLocked, "账户已被禁用")
	case errors.Is(err, service.ErrInvalidCredentials):
		pkg.Error(c, http.StatusUnauthorized, pkg.CodeUnauthorized, "用户名或密码错误")
	default:
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, "登录失败")
	}
}
