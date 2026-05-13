package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gogo/internal/dto"
	"gogo/internal/pkg"
	"gogo/internal/repository"
)

// LogHandler handles log query HTTP requests.
type LogHandler struct {
	logRepo repository.LogRepository
}

// NewLogHandler creates a new LogHandler.
func NewLogHandler(logRepo repository.LogRepository) *LogHandler {
	return &LogHandler{logRepo: logRepo}
}

// ListOperations handles GET /api/v1/logs/operations.
func (h *LogHandler) ListOperations(c *gin.Context) {
	var req dto.OperationLogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}

	logs, total, err := h.logRepo.ListOperations(c.Request.Context(), req)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeDBError, "查询操作日志失败")
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	pkg.Paginated(c, logs, total, req.Page, req.PageSize)
}

// ListTerminals handles GET /api/v1/logs/terminals.
func (h *LogHandler) ListTerminals(c *gin.Context) {
	var req dto.TerminalLogListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, "参数错误："+err.Error())
		return
	}

	logs, total, err := h.logRepo.ListTerminals(c.Request.Context(), req)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeDBError, "查询终端日志失败")
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	pkg.Paginated(c, logs, total, req.Page, req.PageSize)
}
