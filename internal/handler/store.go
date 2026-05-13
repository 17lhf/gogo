package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gogo/internal/dto"
	"gogo/internal/middleware"
	"gogo/internal/pkg"
	"gogo/internal/service"
)

// StoreHandler handles store management HTTP requests.
type StoreHandler struct {
	storeSvc *service.StoreService
}

// NewStoreHandler creates a new StoreHandler.
func NewStoreHandler(storeSvc *service.StoreService) *StoreHandler {
	return &StoreHandler{storeSvc: storeSvc}
}

func (h *StoreHandler) Create(c *gin.Context) {
	var req dto.CreateStoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}
	store, err := h.storeSvc.Create(c.Request.Context(), req)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, err.Error())
		return
	}
	pkg.Success(c, store)
}

func (h *StoreHandler) GetByID(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}
	store, err := h.storeSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		pkg.Error(c, http.StatusNotFound, pkg.CodeParamError, err.Error())
		return
	}
	pkg.Success(c, store)
}

func (h *StoreHandler) List(c *gin.Context) {
	var req dto.StoreListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}
	stores, total, err := h.storeSvc.List(c.Request.Context(), req)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeDBError, "查询门店列表失败")
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	pkg.Paginated(c, stores, total, req.Page, req.PageSize)
}

func (h *StoreHandler) Update(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}
	var req dto.UpdateStoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}
	if err := h.storeSvc.Update(c.Request.Context(), id, req); err != nil {
		pkg.Error(c, http.StatusNotFound, pkg.CodeParamError, err.Error())
		return
	}
	pkg.Success(c, nil)
}

func (h *StoreHandler) Delete(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, "ID格式错误")
		return
	}
	if err := h.storeSvc.Delete(c.Request.Context(), id); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, err.Error())
		return
	}
	pkg.Success(c, nil)
}
