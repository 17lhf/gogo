package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gogo/internal/model"
	"gogo/internal/repository"
)

const (
	whitelistReadPath = "/api/v1/auth/me"
)

// Audit returns a middleware that logs write operations and whitelisted reads.
// It should be applied AFTER the handler to capture response info.
func Audit(logRepo repository.LogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read request body for logging
		var reqBody json.RawMessage
		if c.Request.Body != nil {
			data, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
			if len(data) > 0 {
				reqBody = json.RawMessage(data)
			}
		}

		c.Next()

		method := c.Request.Method
		path := c.FullPath()

		// Skip non-auditable requests
		if !isAuditable(method, path) {
			return
		}

		durationMs := int(time.Since(start).Milliseconds())
		status := model.LogStatusSuccess
		if c.Writer.Status() >= 400 {
			status = model.LogStatusFailure
		}

		userID := GetUserID(c)
		username := GetUsername(c)

		detail := gin.H{
			"method":  method,
			"path":    path,
			"query":   c.Request.URL.RawQuery,
			"body":    reqBody,
		}

		detailJSON, _ := json.Marshal(detail)

		var uid *int64
		if userID != 0 {
			uid = &userID
		}

		logRepo.CreateOperation(c.Request.Context(), &model.OperationLog{
			UserID:       uid,
			Username:     username,
			Action:       method + " " + path,
			ResourceType: extractResource(path),
			ResourceID:   c.Param("id"),
			Detail:       detailJSON,
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			Status:       status,
			DurationMs:   durationMs,
		})
	}
}

func isAuditable(method, path string) bool {
	// Always log write operations
	if method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
		return true
	}
	// Log whitelist read operations
	if method == "GET" && path == whitelistReadPath {
		return true
	}
	return false
}

func extractResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		return parts[2] // e.g. /api/v1/users → "users"
	}
	return ""
}
