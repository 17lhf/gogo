package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gogo/internal/dto"
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
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
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
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
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
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}

	users, total, err := h.userSvc.List(c.Request.Context(), req)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeDBError, "查询用户列表失败")
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
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}

	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
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
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
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
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}

	var req dto.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
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
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}

	var req dto.AssignRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
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
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}

	var req dto.AssignStoresReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}

	if err := h.userSvc.AssignStores(c.Request.Context(), id, req.StoreIDs); err != nil {
		handleUserError(c, err)
		return
	}

	pkg.Success(c, nil)
}

func handleUserError(c *gin.Context, err error) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "不存在"):
		pkg.Error(c, http.StatusNotFound, pkg.CodeParamError, msg)
	case strings.Contains(msg, "已存在"):
		pkg.Error(c, http.StatusConflict, pkg.CodeParamError, msg)
	case strings.Contains(msg, "密码"):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, msg)
	default:
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, msg)
	}
}
