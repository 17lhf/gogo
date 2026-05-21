package handler

import (
	"gogo/internal/i18n"
	"gogo/internal/pkg"
	"gogo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsSvc *service.StatsService
}

func NewStatsHandler(statsSvc *service.StatsService) *StatsHandler {
	return &StatsHandler{statsSvc: statsSvc}
}

// GetTerminals handles GET /api/v1/stats/terminals
func (h *StatsHandler) GetTerminals(c *gin.Context) {
	statsTerminals, err := h.statsSvc.GetTerminals(c.Request.Context())
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, i18n.Localize(c, i18n.MsgGetTerminalsStatsFailed)+err.Error())
		return
	}
	pkg.Success(c, statsTerminals)
}

// GetUsers handles GET /api/v1/stats/users
func (h *StatsHandler) GetUsers(c *gin.Context) {
	statsUsers, err := h.statsSvc.GetUsers(c.Request.Context())
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, pkg.CodeInternalError, i18n.Localize(c, i18n.MsgGetUsersStatsFailed)+err.Error())
		return
	}
	pkg.Success(c, statsUsers)
}
