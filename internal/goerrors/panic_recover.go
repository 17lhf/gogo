package goerrors

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// PanicRecoverHandler GET /api/errors/panic-recover
// 演示panic/recover — Go的紧急停止/恢复机制
func PanicRecoverHandler(w http.ResponseWriter, r *http.Request) {
	// ── panic vs error ────────────────────────────────────────
	// Go的错误处理哲学：
	// - 可预期的错误（IO失败、参数错误）→ 用 error 值传递
	// - 不可恢复的程序错误（bug、不变量违反）→ 用 panic
	//
	// Java几乎所有错误都用Exception，没有这个区分
	// Go的panic类似Java的Error（OutOfMemoryError等），而非Exception

	// ── 正常执行 ──────────────────────────────────────────────
	normalResult := safeRun(func() interface{} {
		return "正常执行完成"
	})

	// ── recover捕获字符串panic ────────────────────────────────
	// Java: try { ... } catch (Exception e) { ... }
	// Go:   recover() 只在 defer 函数中有效！
	panicCaught := safeRun(func() interface{} {
		panic("模拟程序panic！")
	})

	// ── recover捕获runtime panic ─────────────────────────────
	nilPanic := safeRun(func() interface{} {
		var s []int
		return s[10] // index out of range，runtime自动触发panic
	})

	// ── panic可以是任意类型 ───────────────────────────────────
	type PanicInfo struct{ Code int; Message string }
	structPanic := safeRun(func() interface{} {
		panic(PanicInfo{Code: 500, Message: "内部错误"})
	})

	// ── panic的合理使用场景 ───────────────────────────────────
	validUseCases := []map[string]string{
		{
			"scenario": "Must*初始化函数",
			"example":  `regexp.MustCompile("pattern") // 编译失败直接panic`,
			"reason":   "程序启动时的不可恢复错误，早失败比晚失败好",
		},
		{
			"scenario": "不变量违反",
			"example":  `if p == nil { panic("p不应该为nil") }`,
			"reason":   "调用者违反了函数前置条件，属于bug而非运行时错误",
		},
		{
			"scenario": "不该到达的分支",
			"example":  `default: panic(fmt.Sprintf("未知类型: %T", v))`,
			"reason":   "switch穷举了所有情况，default意味着程序逻辑有误",
		},
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "panic / recover",
		Java:  "Java Error / RuntimeException（但使用场景不同）",
		Summary: "panic用于不可恢复的程序错误，recover（只能在defer中）捕获panic。" +
			"普通业务错误用error，不要滥用panic作为流程控制。",
		Points: []string{
			"panic立即停止当前函数，沿调用栈向上传播，直到recover或程序崩溃",
			"recover()只在defer函数中有效，返回panic的值（nil=没有panic）",
			"runtime panic: 空指针解引用、数组越界、类型断言失败等自动触发",
			"panic可以是任意类型的值，不限于字符串或error",
			"库代码不应让panic传出包边界，应在边界recover并转为error",
			"HTTP框架的Recoverer中间件会捕获handler的panic返回500",
		},
		Data: map[string]interface{}{
			"normal":          normalResult,
			"panic_caught":    panicCaught,
			"nil_panic":       nilPanic,
			"struct_panic":    structPanic,
			"valid_use_cases": validUseCases,
		},
		Tips: []string{
			"Java: throw new RuntimeException()  →  Go: panic(\"msg\")，但别滥用",
			"Java: catch (Exception e)  →  Go: defer func() { if r := recover(); r != nil { ... } }()",
			"绝大多数情况用error，panic只用于真正的程序bug",
		},
	})
}

// safeRun 包装执行，捕获任何panic并返回结果
func safeRun(fn func() interface{}) map[string]interface{} {
	result := map[string]interface{}{"panicked": false}
	func() {
		defer func() {
			if r := recover(); r != nil {
				result["panicked"] = true
				result["panic_value"] = fmt.Sprintf("%v", r)
				result["panic_type"] = fmt.Sprintf("%T", r)
			}
		}()
		result["value"] = fn()
	}()
	return result
}
