package goerrors

import (
	"errors"
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// WrappingHandler GET /api/errors/wrapping
// 演示错误包装与错误链（Go 1.13+）
func WrappingHandler(w http.ResponseWriter, r *http.Request) {
	// ── 逐层包装，添加上下文 ──────────────────────────────────
	// Java: new ServiceException("msg", cause)
	// Go:   fmt.Errorf("context: %w", err)  — %w 保留错误链
	dbErr := errors.New("connection refused")
	repoErr := fmt.Errorf("UserRepository.FindByID: %w", dbErr)
	serviceErr := fmt.Errorf("UserService.GetUser(id=42): %w", repoErr)
	handlerErr := fmt.Errorf("GET /users/42: %w", serviceErr)

	// ── errors.Is — 穿透整个错误链 ───────────────────────────
	// Java: 需要 while(cause != null) { if(cause instanceof X) ... cause = cause.getCause() }
	// Go:   errors.Is 自动遍历整个链
	isDbErr := errors.Is(handlerErr, dbErr) // true，穿透3层包装

	// ── errors.As — 提取链中的特定类型 ───────────────────────
	origValErr := &ValidationError{Field: "age", Message: "必须大于0", Code: 2001}
	wrapped := fmt.Errorf("处理用户数据: %w", origValErr)
	doubleWrapped := fmt.Errorf("API调用: %w", wrapped)

	var extracted *ValidationError
	asResult := errors.As(doubleWrapped, &extracted)

	// ── errors.Unwrap — 解一层包装 ───────────────────────────
	unwrapped := errors.Unwrap(wrapped)

	// ── 最佳实践：何时用%w vs 普通格式 ───────────────────────
	practices := []map[string]string{
		{
			"situation": "跨越抽象层，需要调用方能errors.Is/As",
			"use":       `fmt.Errorf("UserService.Create: %w", err)`,
			"note":      "用%w，保留错误链",
		},
		{
			"situation": "转换错误类型，不想暴露底层错误",
			"use":       `fmt.Errorf("操作失败")`,
			"note":      "不用%w，切断错误链，调用方无法穿透",
		},
		{
			"situation": "返回预定义哨兵错误",
			"use":       `return ErrNotFound`,
			"note":      "直接返回，调用方用errors.Is(err, ErrNotFound)",
		},
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "错误包装与错误链 (Error Wrapping)",
		Java:  "Java的 Exception.getCause() / 异常链",
		Summary: "Go 1.13+支持错误包装(%w)，形成错误链。" +
			"errors.Is在链中查找值，errors.As在链中查找类型，比Java的getCause()更强大。",
		Points: []string{
			"fmt.Errorf(\"context: %w\", err) 包装错误，添加上下文信息",
			"errors.Is(err, target) 在整个错误链中检查，支持多层包装",
			"errors.As(err, &target) 在整个错误链中找到target类型的错误",
			"errors.Unwrap(err) 解一层包装，返回被包装的原始错误",
			"实现 Unwrap() error 方法的自定义错误类型可参与错误链",
			"跨层传递建议包装: return fmt.Errorf(\"函数名: %w\", err)",
		},
		Data: map[string]interface{}{
			"error_chain": map[string]string{
				"1_db":      dbErr.Error(),
				"2_repo":    repoErr.Error(),
				"3_service": serviceErr.Error(),
				"4_handler": handlerErr.Error(),
			},
			"errors_is":  isDbErr,
			"errors_as": map[string]interface{}{
				"found":   asResult,
				"field":   extracted.Field,
				"message": extracted.Message,
			},
			"unwrap_result":  unwrapped.Error(),
			"best_practices": practices,
		},
	})
}
