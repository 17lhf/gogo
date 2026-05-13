package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gogo/internal/dto"
	"gogo/internal/middleware"
	"gogo/internal/pkg"
	"gogo/internal/service"
)

// MenuHandler handles menu management HTTP requests.
type MenuHandler struct {
	menuSvc *service.MenuService
}

// NewMenuHandler creates a new MenuHandler.
func NewMenuHandler(menuSvc *service.MenuService) *MenuHandler {
	return &MenuHandler{menuSvc: menuSvc}
}

// Create handles POST /api/v1/menus.
func (h *MenuHandler) Create(c *gin.Context) {
	var req dto.CreateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}

	menu, err := h.menuSvc.Create(c.Request.Context(), req)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, err.Error())
		return
	}

	pkg.Success(c, menu)
}

// Tree handles GET /api/v1/menus.
func (h *MenuHandler) Tree(c *gin.Context) {
	tree, err := h.menuSvc.Tree(c.Request.Context())
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeDBError, "查询菜单树失败")
		return
	}

	pkg.Success(c, tree)
}

// GetByID handles GET /api/v1/menus/:id.
func (h *MenuHandler) GetByID(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}

	menu, err := h.menuSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		pkg.Error(c, http.StatusNotFound, pkg.CodeParamError, err.Error())
		return
	}

	pkg.Success(c, menu)
}

// Update handles PUT /api/v1/menus/:id.
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}

	var req dto.UpdateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}

	if err := h.menuSvc.Update(c.Request.Context(), id, req); err != nil {
		pkg.Error(c, http.StatusNotFound, pkg.CodeParamError, err.Error())
		return
	}

	pkg.Success(c, nil)
}

// Delete handles DELETE /api/v1/menus/:id.
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}

	if err := h.menuSvc.Delete(c.Request.Context(), id); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, err.Error())
		return
	}

	pkg.Success(c, nil)
}
