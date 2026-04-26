package goerrors

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// DeferHandler GET /api/errors/defer
// 演示defer的用法和执行顺序
func DeferHandler(w http.ResponseWriter, r *http.Request) {
	// ── defer基础：LIFO顺序 ───────────────────────────────────
	// Java: try { ... } finally { cleanup(); }
	// Go:   defer cleanup()  更简洁，无需try块

	executionLog := demoLIFO()
	modifiedReturn := demoNamedReturn()

	// ── defer的实际使用场景 ───────────────────────────────────
	scenarios := []map[string]string{
		{
			"scenario": "文件关闭",
			"java":     "try (var f = new FileInputStream(...)) { ... }",
			"go":       "f, err := os.Open(...); if err != nil { return err }\ndefer f.Close()",
		},
		{
			"scenario": "数据库连接归还",
			"java":     "try (var conn = pool.getConnection()) { ... }",
			"go":       "conn := pool.Get(); defer pool.Put(conn)",
		},
		{
			"scenario": "互斥锁释放",
			"java":     "lock.lock(); try { ... } finally { lock.unlock(); }",
			"go":       "mu.Lock(); defer mu.Unlock()",
		},
		{
			"scenario": "HTTP响应体关闭",
			"java":     "try { ... } finally { resp.body().close(); }",
			"go":       "resp, err := http.Get(url); defer resp.Body.Close()",
		},
		{
			"scenario": "函数执行计时",
			"java":     "long start = System.currentTimeMillis(); // 手动在末尾计算",
			"go":       "defer func(start time.Time) { log(time.Since(start)) }(time.Now())",
		},
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "defer — 延迟执行",
		Java:  "Java try-with-resources / finally块",
		Summary: "defer注册的函数在外层函数返回前执行，遵循LIFO顺序。" +
			"比Java的finally更灵活，可在任意位置注册，紧跟资源获取代码。",
		Points: []string{
			"defer在函数返回前执行（无论正常return还是panic）",
			"多个defer遵循LIFO: 最后注册的最先执行",
			"defer参数立即求值，但函数体延迟执行",
			"defer闭包捕获变量引用，可看到变量的最新值",
			"defer可以修改命名返回值（高级用法）",
			"紧跟资源获取后写defer，逻辑紧凑不容易忘记释放",
		},
		Data: map[string]interface{}{
			"lifo_execution_log": executionLog,
			"named_return":       modifiedReturn,
			"use_scenarios":      scenarios,
		},
	})
}

func demoLIFO() []string {
	var log []string
	record := func(msg string) { log = append(log, msg) }

	// 用真实defer展示LIFO
	func() {
		record("步骤1: 函数开始")
		defer func() { record("步骤5: defer-C 执行（最先注册，最后执行）") }()
		record("步骤2: defer-C 已注册")
		defer func() { record("步骤4: defer-B 执行") }()
		record("步骤3: defer-B 已注册，defer-A 即将注册")
		defer func() { record("步骤3.5: defer-A 执行（最后注册，最先执行）") }()
		record("--- 函数正常逻辑结束，开始执行defer ---")
	}()

	return log
}

func demoNamedReturn() string {
	// defer可以修改命名返回值
	result := namedReturnDemo()
	return fmt.Sprintf("函数声明返回'原始值'，但被defer改为: %q", result)
}

func namedReturnDemo() (result string) {
	result = "原始值"
	defer func() {
		result = "被defer修改: " + result // 修改命名返回值
	}()
	return // 裸return，defer在此之后执行并修改result
}
