package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// responseBodyWriter wraps gin.ResponseWriter to capture the response body.
type responseBodyWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// RequestLog logs every incoming request and its response at slog.Debug level.
// Set LOG_LEVEL=debug (or LOG_HTTP=true) in DEV to enable.
// Set LOG_HTTP_BODY=false to omit request/response bodies from logs.
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !slog.Default().Enabled(c.Request.Context(), slog.LevelDebug) {
			c.Next()
			return
		}

		start := time.Now()

		// Buffer request body
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Wrap response writer to capture response body
		w := &responseBodyWriter{ResponseWriter: c.Writer, buf: &bytes.Buffer{}}
		c.Writer = w

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", status),
			slog.Duration("duration", duration),
			slog.String("ip", c.ClientIP()),
		}

		if q := c.Request.URL.RawQuery; q != "" {
			attrs = append(attrs, slog.String("query", q))
		}

		if len(reqBody) > 0 {
			attrs = append(attrs, slog.String("req_body", truncate(string(reqBody), 4096)))
		}

		if w.buf.Len() > 0 {
			attrs = append(attrs, slog.String("resp_body", truncate(w.buf.String(), 4096)))
		}

		if errs := c.Errors; len(errs) > 0 {
			attrs = append(attrs, slog.String("errors", errs.String()))
		}

		slog.DebugContext(c.Request.Context(), "http request", attrs...)
	}
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}
