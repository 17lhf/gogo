package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess = 0

	// Param errors 40001-40099
	CodeParamError      = 40001
	CodeValidationError = 40002

	// Auth errors 40101-40199
	CodeUnauthorized    = 40101
	CodeTokenExpired    = 40102
	CodeAccountLocked   = 40103
	CodeSessionNotFound = 40104

	// Permission errors 40301-40399
	CodeForbidden             = 40301
	CodePasswordExpired       = 40302
	CodeMustChangePassword    = 40303
	CodeStoreHasTerminals     = 40304
	CodeTerminalDisabled      = 40305
	CodeTerminalNotFound      = 40306

	// Server errors 50001-50099
	CodeInternalError = 50001
	CodeDBError       = 50002
)

// Success writes a standard success response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": CodeSuccess, "msg": "success", "data": data})
}

// Paginated writes a paginated success response.
func Paginated(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  "success",
		"data": gin.H{"list": list, "total": total, "page": page, "page_size": pageSize},
	})
}

// Error writes a standard error response and aborts the request.
func Error(c *gin.Context, httpStatus int, code int, msg string) {
	c.AbortWithStatusJSON(httpStatus, gin.H{"code": code, "msg": msg, "data": nil})
}
