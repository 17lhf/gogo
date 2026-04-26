// package middleware 提供HTTP中间件实现
// 演示Go中间件的标准写法，对比Java的Filter/Interceptor
package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// ── 中间件函数签名 ─────────────────────────────────────────────
// Java: implements Filter { doFilter(req, resp, chain) }
// Go:   func(next http.Handler) http.Handler — 返回新的Handler

// CORS 跨域中间件
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r) // 调用下一个handler（类似chain.doFilter）
	})
}

// contextKey 避免context key冲突（用私有类型）
type contextKey string

const requestIDKey contextKey = "request_id"

// RequestID 为每个请求生成唯一ID，存入context
// Java: ThreadLocal + Filter
// Go:   context + middleware
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := fmt.Sprintf("req-%08x", rand.Uint32())
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx)) // 传递带requestID的context
	})
}

// GetRequestID 从context中获取requestID
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return "unknown"
}

// Logger 请求日志中间件
// Java: HandlerInterceptor.preHandle + postHandle
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		elapsed := time.Since(start)
		fmt.Printf("[%s] %s %s %d %s\n",
			GetRequestID(r.Context()),
			r.Method,
			r.URL.Path,
			rw.statusCode,
			elapsed,
		)
	})
}

// responseWriter 包装ResponseWriter以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
