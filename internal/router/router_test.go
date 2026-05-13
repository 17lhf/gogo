package router

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gogo/internal/pkg"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupMinimalRouter creates a router with just the health and no-route endpoints.
func setupMinimalRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.NoRoute(func(c *gin.Context) {
		pkg.Error(c, 404, pkg.CodeParamError, "接口不存在")
	})
	return r
}

func TestRouter_HealthEndpoint(t *testing.T) {
	r := setupMinimalRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "ok", body["status"])
}

func TestRouter_NoRoute(t *testing.T) {
	r := setupMinimalRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestRouter_ResponseFormat(t *testing.T) {
	r := setupMinimalRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	r.ServeHTTP(w, req)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, float64(40001), body["code"])
	assert.Contains(t, body["msg"], "不存在")
	assert.Nil(t, body["data"])
}

func TestRouter_HealthResponseFormat(t *testing.T) {
	r := setupMinimalRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
