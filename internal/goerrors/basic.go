// package goerrors 演示Go的错误处理机制
// 这是Java开发者转Go时最大的思维转变之一
package goerrors

import (
	"errors"
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// BasicErrorHandler GET /api/errors/basic
// 演示Go基础错误处理 vs Java异常
func BasicErrorHandler(w http.ResponseWriter, r *http.Request) {
	// ── Go的错误哲学 ──────────────────────────────────────────
	// Java: 用Exception表示错误，通过throw/catch传递
	// Go:   错误是普通值，通过多返回值(result, error)传递
	//
	// Java: try { int r = divide(10, 0); } catch(ArithmeticException e) { ... }
	// Go:   r, err := divide(10, 3); if err != nil { ... }

	result, err := safeDivide(10, 3)
	var successCase map[string]interface{}
	if err != nil {
		successCase = map[string]interface{}{"error": err.Error()}
	} else {
		successCase = map[string]interface{}{"result": result, "error": nil}
	}

	_, errCase := safeDivide(10, 0)
	failCase := map[string]interface{}{
		"error":      errCase.Error(),
		"error_type": fmt.Sprintf("%T", errCase),
	}

	// ── 创建错误 ──────────────────────────────────────────────
	// Java: new IllegalArgumentException("msg")
	errSimple := errors.New("用户不存在")
	errFmt := fmt.Errorf("无效的用户ID: %d", -1)

	// ── 多步骤操作的错误处理风格 ──────────────────────────────
	// Go风格: 早返回，避免嵌套（"happy path"在左侧）
	steps := []map[string]string{
		{"step": "1. 验证参数", "go_style": "if param == \"\" { return errors.New(\"参数为空\") }"},
		{"step": "2. 查询数据库", "go_style": "user, err := db.Find(id); if err != nil { return err }"},
		{"step": "3. 业务逻辑", "go_style": "if !user.Active { return ErrUserInactive }"},
		{"step": "4. 返回结果", "go_style": "return user, nil"},
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go错误处理基础",
		Java:  "Java Exception / try-catch-finally",
		Summary: "Go中错误是普通值（实现error接口的类型），通过多返回值传递。" +
			"没有try-catch，通过 if err != nil 检查，强制调用者处理错误。",
		Points: []string{
			"error是内置接口: type error interface { Error() string }",
			"函数签名惯例: func f() (Result, error)，error是最后一个返回值",
			"立即检查: if err != nil { return ..., err } — 早返回避免嵌套",
			"errors.New(\"msg\") 创建简单错误，fmt.Errorf(\"格式%w\", err) 包装错误",
			"_ 忽略错误是危险的，除非确定不会失败",
			"Go没有checked/unchecked exception之分，所有错误都需要显式处理",
		},
		Data: map[string]interface{}{
			"success_case":  successCase,
			"fail_case":     failCase,
			"errors_new":    errSimple.Error(),
			"fmt_errorf":    errFmt.Error(),
			"early_return_pattern": steps,
		},
		Tips: []string{
			"Java: throw new IllegalArgumentException(\"msg\")  →  Go: return errors.New(\"msg\")",
			"Java: catch (Exception e)  →  Go: if err != nil { ... }",
			"Java: throws IOException  →  Go: 函数返回error，无需声明",
		},
	})
}

func safeDivide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}
