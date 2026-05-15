package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gogo/internal/cache"
	"gogo/internal/dto"
	"gogo/internal/i18n"
	"gogo/internal/middleware"
	"gogo/internal/pkg"
	"gogo/internal/service"
)

// TerminalHandler handles terminal management HTTP requests.
type TerminalHandler struct {
	terminalSvc    *service.TerminalService
	heartbeatCache *cache.HeartbeatCache
}

// NewTerminalHandler creates a new TerminalHandler.
func NewTerminalHandler(terminalSvc *service.TerminalService, heartbeatCache *cache.HeartbeatCache) *TerminalHandler {
	return &TerminalHandler{terminalSvc: terminalSvc, heartbeatCache: heartbeatCache}
}

func (h *TerminalHandler) Create(c *gin.Context) {
	var req dto.CreateTerminalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}
	t, err := h.terminalSvc.Create(c.Request.Context(), req)
	if err != nil {
		handleTerminalError(c, err)
		return
	}
	pkg.Success(c, t)
}

func (h *TerminalHandler) GetByID(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}
	t, err := h.terminalSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		handleTerminalError(c, err)
		return
	}
	pkg.Success(c, t)
}

func (h *TerminalHandler) List(c *gin.Context) {
	var req dto.TerminalListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}

	storeIDs := middleware.GetStoreIDs(c)
	terminals, total, err := h.terminalSvc.List(c.Request.Context(), req, storeIDs)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeDBError, i18n.Localize(c, i18n.MsgTerminalListFailed))
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	pkg.Paginated(c, terminals, total, req.Page, req.PageSize)
}

func (h *TerminalHandler) Update(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}
	var req dto.UpdateTerminalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeValidationError, i18n.Localize(c, i18n.MsgParamInvalid)+": "+err.Error())
		return
	}
	if err := h.terminalSvc.Update(c.Request.Context(), id, req); err != nil {
		handleTerminalError(c, err)
		return
	}
	pkg.Success(c, nil)
}

func (h *TerminalHandler) Delete(c *gin.Context) {
	id, err := middleware.GetInt64Param(c, "id")
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgIDFormat))
		return
	}
	if err := h.terminalSvc.Delete(c.Request.Context(), id); err != nil {
		handleTerminalError(c, err)
		return
	}
	pkg.Success(c, nil)
}

// Heartbeat handles POST /api/v1/terminals/:sn/heartbeat (device token auth).
func (h *TerminalHandler) Heartbeat(c *gin.Context) {
	sn := c.Param("sn")

	token := c.GetHeader("X-Device-Token")
	if token == "" {
		pkg.Error(c, http.StatusUnauthorized, pkg.CodeUnauthorized, i18n.Localize(c, i18n.MsgTerminalNoToken))
		return
	}

	storedSN, err := h.heartbeatCache.GetSNByDeviceToken(c.Request.Context(), token)
	if err != nil || storedSN == "" || storedSN != sn {
		t, err := h.terminalSvc.GetBySN(c.Request.Context(), sn)
		if err != nil || t == nil {
			pkg.Error(c, http.StatusNotFound, pkg.CodeTerminalNotFound, i18n.Localize(c, i18n.MsgTerminalNotFound))
			return
		}
		if t.DeviceToken != token {
			pkg.Error(c, http.StatusUnauthorized, pkg.CodeUnauthorized, i18n.Localize(c, i18n.MsgTerminalInvalidToken))
			return
		}
		h.heartbeatCache.SetDeviceToken(c.Request.Context(), token, sn)
	}

	var req dto.HeartbeatReq
	c.ShouldBindJSON(&req)

	if err := h.terminalSvc.Heartbeat(c.Request.Context(), sn, req.IPAddress, req.MACAddress); err != nil {
		if errors.Is(err, service.ErrTerminalDisabled) {
			pkg.Error(c, http.StatusForbidden, pkg.CodeTerminalDisabled, i18n.Localize(c, i18n.MsgTerminalDisabled))
			return
		}
		pkg.Error(c, http.StatusNotFound, pkg.CodeTerminalNotFound, i18n.Localize(c, i18n.MsgTerminalNotFound))
		return
	}

	pkg.Success(c, gin.H{"sn": sn, "timestamp": gin.H{}})
}

// RotateToken handles POST /api/v1/terminals/:sn/rotate-token.
func (h *TerminalHandler) RotateToken(c *gin.Context) {
	sn := c.Param("sn")
	token := c.GetHeader("X-Device-Token")
	if token == "" {
		pkg.Error(c, http.StatusUnauthorized, pkg.CodeUnauthorized, i18n.Localize(c, i18n.MsgTerminalNoToken))
		return
	}

	t, err := h.terminalSvc.GetBySN(c.Request.Context(), sn)
	if err != nil || t == nil {
		pkg.Error(c, http.StatusNotFound, pkg.CodeTerminalNotFound, i18n.Localize(c, i18n.MsgTerminalNotFound))
		return
	}
	if t.DeviceToken != token {
		pkg.Error(c, http.StatusUnauthorized, pkg.CodeUnauthorized, i18n.Localize(c, i18n.MsgTerminalInvalidToken))
		return
	}

	newToken, err := h.terminalSvc.RotateToken(c.Request.Context(), sn)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, i18n.Localize(c, i18n.MsgTerminalRotateFailed))
		return
	}

	pkg.Success(c, gin.H{"device_token": newToken})
}

func handleTerminalError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrTerminalNotFound):
		pkg.Error(c, http.StatusNotFound, pkg.CodeParamError, i18n.Localize(c, i18n.MsgTerminalNotFound))
	case errors.Is(err, service.ErrStoreNotFound):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgStoreNotFound))
	case errors.Is(err, service.ErrInvalidStatusChange):
		pkg.Error(c, http.StatusBadRequest, pkg.CodeParamError, i18n.Localize(c, i18n.MsgTerminalInvalidStatus))
	case errors.Is(err, service.ErrTerminalDisabled):
		pkg.Error(c, http.StatusForbidden, pkg.CodeTerminalDisabled, i18n.Localize(c, i18n.MsgTerminalDisabled))
	default:
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, err.Error())
	}
}
