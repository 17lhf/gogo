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

// RoleHandler handles role management HTTP requests.
type RoleHandler struct {
	roleSvc *service.RoleService
}

// NewRoleHandler creates a new RoleHandler.
func NewRoleHandler(roleSvc *service.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

// Create handles POST /api/v1/roles.
func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	role, err := h.roleSvc.Create(c.Request.Context(), req)
	if err != nil {
		handleRoleError(c, err)
		return
	}

	pkg.Success(c, role)
}

// List handles GET /api/v1/roles.
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roleSvc.List(c.Request.Context())
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeDBError, i18n.Localize(c, i18n.MsgRoleListFailed))
		return
	}

	pkg.Success(c, roles)
}

// GetByID handles GET /api/v1/roles/:id.
func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	role, err := h.roleSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		handleRoleError(c, err)
		return
	}

	pkg.Success(c, role)
}

// Update handles PUT /api/v1/roles/:id.
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	var req dto.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	if err := h.roleSvc.Update(c.Request.Context(), id, req); err != nil {
		handleRoleError(c, err)
		return
	}

	pkg.Success(c, nil)
}

// Delete handles DELETE /api/v1/roles/:id.
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	if err := h.roleSvc.Delete(c.Request.Context(), id); err != nil {
		handleRoleError(c, err)
		return
	}

	pkg.Success(c, nil)
}

// GetMenus handles GET /api/v1/roles/:id/menus.
func (h *RoleHandler) GetMenus(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	menuIDs, err := h.roleSvc.GetMenus(c.Request.Context(), id)
	if err != nil {
		handleRoleError(c, err)
		return
	}

	pkg.Success(c, menuIDs)
}

// AssignMenus handles PUT /api/v1/roles/:id/menus.
func (h *RoleHandler) AssignMenus(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}

	var req dto.AssignMenusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	if err := h.roleSvc.AssignMenus(c.Request.Context(), id, req.MenuIDs); err != nil {
		handleRoleError(c, err)
		return
	}

	pkg.Success(c, nil)
}

func handleRoleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrRoleNotFound):
		pkg.Error(c, http.StatusNotFound, pkg.CodeParamError, i18n.Localize(c, i18n.MsgRoleNotFound))
	case errors.Is(err, service.ErrRoleCodeExists):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgRoleCodeExists))
	default:
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, err.Error())
	}
}
