package goerrors

import (
	"errors"
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// ── 自定义错误类型 ────────────────────────────────────────────
// Java: public class ValidationException extends RuntimeException { ... }
// Go:   实现 error 接口即可（任何有 Error() string 方法的类型）

type ValidationError struct {
	Field   string
	Message string
	Code    int
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("验证失败[字段:%s, 代码:%d]: %s", e.Field, e.Code, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s(id=%v)不存在", e.Resource, e.ID)
}

// DatabaseError 带原始错误的包装类型
type DatabaseError struct {
	Operation string
	Table     string
	Cause     error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("DB操作失败[%s on %s]: %v", e.Operation, e.Table, e.Cause)
}

// Unwrap 实现此方法使 errors.Is/As 能穿透错误链
func (e *DatabaseError) Unwrap() error { return e.Cause }

// CustomErrorHandler GET /api/errors/custom
func CustomErrorHandler(w http.ResponseWriter, r *http.Request) {
	valErr := &ValidationError{Field: "email", Message: "格式不正确", Code: 1001}
	notFoundErr := &NotFoundError{Resource: "User", ID: 12345}
	dbErr := &DatabaseError{
		Operation: "SELECT",
		Table:     "users",
		Cause:     errors.New("connection timeout"),
	}

	// ── errors.As — 类型匹配 ─────────────────────────────────
	// Java: catch (ValidationException e) { e.getField() }
	// Go:   errors.As(err, &target) 在错误链中查找指定类型
	var ve *ValidationError
	var dbe *DatabaseError
	isValErr := errors.As(valErr, &ve)
	isDbErr := errors.As(dbErr, &dbe)

	// ── 哨兵错误（Sentinel Error）────────────────────────────
	// 预定义的固定错误值，用 errors.Is 比较（不可比较的类型用errors.As）
	var ErrPermissionDenied = errors.New("权限不足")
	wrapped := fmt.Errorf("调用API时: %w", ErrPermissionDenied) // %w 包装
	isPermDenied := errors.Is(wrapped, ErrPermissionDenied)    // true，穿透包装

	allErrors := []map[string]string{
		{"type": fmt.Sprintf("%T", valErr), "message": valErr.Error()},
		{"type": fmt.Sprintf("%T", notFoundErr), "message": notFoundErr.Error()},
		{"type": fmt.Sprintf("%T", dbErr), "message": dbErr.Error()},
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "自定义错误类型",
		Java:  "Java自定义Exception类（extends RuntimeException）",
		Summary: "Go通过实现error接口创建自定义错误，可携带结构化数据。" +
			"errors.Is/As能在包装的错误链中查找，替代Java的多catch块。",
		Points: []string{
			"自定义错误: 实现 Error() string 即可，无需继承任何类",
			"携带结构化数据: struct字段存储错误详情，比Java的getMessage()更丰富",
			"Unwrap() error: 实现此方法让errors.Is/As能穿透错误链",
			"errors.Is(err, target): 检查错误链中是否包含target（哨兵错误）",
			"errors.As(err, &target): 在错误链中找到target类型并赋值",
			"哨兵错误: var ErrXxx = errors.New(\"xxx\")，导出供调用方比较",
		},
		Data: map[string]interface{}{
			"all_errors": allErrors,
			"errors_as": map[string]interface{}{
				"validation_matched": isValErr,
				"validation_field":   ve.Field,
				"db_matched":         isDbErr,
				"db_table":           dbe.Table,
			},
			"sentinel_errors_is": isPermDenied,
		},
	})
}
