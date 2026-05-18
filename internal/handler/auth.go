package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gogo/internal/dto"
	"gogo/internal/i18n"
	"gogo/internal/middleware"
	"gogo/internal/pkg"
	"gogo/internal/service"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authSvc *service.AuthService
	menuSvc *service.MenuService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc *service.AuthService, menuSvc *service.MenuService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, menuSvc: menuSvc}
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
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
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, i18n.Localize(c, i18n.MsgAuthLogoutFailed))
		return
	}
	pkg.Success(c, nil)
}

// Me handles GET /api/v1/auth/me.
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)

	profile, err := h.authSvc.Me(c.Request.Context(), userID)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, i18n.Localize(c, i18n.MsgAuthGetProfileFail))
		return
	}

	profile.Menus, _ = h.menuSvc.TreeByUserID(c.Request.Context(), userID)

	pkg.Success(c, profile)
}

// ChangePassword handles PUT /api/v1/auth/password.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.authSvc.ChangePassword(c.Request.Context(), userID, req); err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgAuthWrongPassword))
			return
		}
		handleChangePasswordError(c, err)
		return
	}

	pkg.Success(c, nil)
}

func handleChangePasswordError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, pkg.ErrPasswordTooShort):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgPasswordLength))
	case errors.Is(err, pkg.ErrPasswordNoUpper):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgPasswordUpper))
	case errors.Is(err, pkg.ErrPasswordNoLower):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgPasswordLower))
	case errors.Is(err, pkg.ErrPasswordNoDigit):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgPasswordDigit))
	default:
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, err.Error())
	}
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrAccountLocked):
		pkg.Error(c, http.StatusForbidden, pkg.CodeAccountLocked, i18n.Localize(c, i18n.MsgAuthAccountLocked))
	case errors.Is(err, service.ErrAccountDisabled):
		pkg.Error(c, http.StatusForbidden, pkg.CodeAccountLocked, i18n.Localize(c, i18n.MsgAuthAccountDisabled))
	case errors.Is(err, service.ErrInvalidCredentials):
		pkg.Error(c, http.StatusUnauthorized, pkg.CodeUnauthorized, i18n.Localize(c, i18n.MsgAuthWrongCreds))
	default:
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, i18n.Localize(c, i18n.MsgAuthLoginFailed))
	}
}
