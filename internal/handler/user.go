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

// UserHandler handles user management HTTP requests.
type UserHandler struct {
	userSvc *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// Create handles POST /api/v1/users.
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	user, err := h.userSvc.Create(c.Request.Context(), req)
	if err != nil {
		handleUserError(c, err)
		return
	}

	pkg.Success(c, user)
}

// GetByID handles GET /api/v1/users/:id.
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	user, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		handleUserError(c, err)
		return
	}

	pkg.Success(c, user)
}

// List handles GET /api/v1/users.
func (h *UserHandler) List(c *gin.Context) {
	var req dto.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	users, total, err := h.userSvc.List(c.Request.Context(), req)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeDBError, i18n.Localize(c, i18n.MsgUserListFailed))
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	pkg.Paginated(c, users, total, req.Page, req.PageSize)
}

// Update handles PUT /api/v1/users/:id.
func (h *UserHandler) Update(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	if err := h.userSvc.Update(c.Request.Context(), id, req); err != nil {
		handleUserError(c, err)
		return
	}

	pkg.Success(c, nil)
}

// Delete handles DELETE /api/v1/users/:id.
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	if err := h.userSvc.Delete(c.Request.Context(), id); err != nil {
		handleUserError(c, err)
		return
	}

	pkg.Success(c, nil)
}

// ResetPassword handles PUT /api/v1/users/:id/password.
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	var req dto.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	if err := h.userSvc.ResetPassword(c.Request.Context(), id, req.Password); err != nil {
		handleUserError(c, err)
		return
	}

	pkg.Success(c, nil)
}

// AssignRoles handles PUT /api/v1/users/:id/roles.
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	var req dto.AssignRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	if err := h.userSvc.AssignRoles(c.Request.Context(), id, req.RoleIDs); err != nil {
		handleUserError(c, err)
		return
	}

	pkg.Success(c, nil)
}

// AssignStores handles PUT /api/v1/users/:id/stores.
func (h *UserHandler) AssignStores(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	var req dto.AssignStoresReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	if err := h.userSvc.AssignStores(c.Request.Context(), id, req.StoreIDs); err != nil {
		handleUserError(c, err)
		return
	}

	pkg.Success(c, nil)
}

func handleUserError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		pkg.Error(c, http.StatusNotFound, pkg.CodeParamError, i18n.Localize(c, i18n.MsgUserNotFound))
	case errors.Is(err, service.ErrUsernameExists):
		pkg.Error(c, http.StatusConflict, pkg.CodeParamError, i18n.Localize(c, i18n.MsgUsernameExists))
	case errors.Is(err, service.ErrEmailExists):
		pkg.Error(c, http.StatusConflict, pkg.CodeParamError, i18n.Localize(c, i18n.MsgEmailExists))
	case errors.Is(err, pkg.ErrPasswordTooShort):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgPasswordLength))
	case errors.Is(err, pkg.ErrPasswordNoUpper):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgPasswordUpper))
	case errors.Is(err, pkg.ErrPasswordNoLower):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgPasswordLower))
	case errors.Is(err, pkg.ErrPasswordNoDigit):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgPasswordDigit))
	default:
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, err.Error())
	}
}
