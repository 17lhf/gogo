package concurrency

import (
	"fmt"
	"net/http"
	"time"

	"gogo/internal/lesson"
)

// GoroutinesHandler GET /api/concurrency/goroutines
// 演示goroutine — Go的轻量级线程
func GoroutinesHandler(w http.ResponseWriter, r *http.Request) {
	// ── Goroutine基础 ─────────────────────────────────────────
	// Java: new Thread(() -> task()).start()  或  executor.submit(task)
	// Go:   go funcName()  或  go func() { ... }()
	//
	// 关键差异:
	// Java线程: ~1MB内存，OS线程，创建/切换开销大
	// Goroutine: ~2KB起始栈，Go运行时调度，可轻松创建百万个

	results := make(chan string, 5) // 缓冲channel收集结果
	start := time.Now()

	for i := 1; i <= 5; i++ {
		i := i // 捕获循环变量（Go 1.22之前必须写这行！）
		go func() {
			time.Sleep(time.Duration(i*20) * time.Millisecond)
			results <- fmt.Sprintf("任务%d完成(耗时%dms)", i, i*20)
		}()
	}

	// 收集全部结果（每个 <-results 阻塞直到有数据）
	collected := make([]string, 0, 5)
	for range 5 {
		collected = append(collected, <-results)
	}

	elapsed := time.Since(start)
	// 顺序执行需要 20+40+60+80+100=300ms
	// 并发执行只需最长任务 ~100ms

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Goroutine — 轻量级并发",
		Java:  "Java Thread / ExecutorService / CompletableFuture",
		Summary: "Goroutine是Go并发的基本单元，比Java线程轻量100倍以上。" +
			"用 go 关键字启动，用channel通信(CSP模型)。",
		Points: []string{
			"go func() 启动goroutine，'go'关键字就是全部语法",
			"goroutine初始栈约2KB，按需增长；Java线程默认1MB",
			"Go运行时用M:N调度: M个goroutine映射到N个OS线程",
			"循环中启动goroutine必须捕获变量: i := i (Go 1.22+已修复)",
			"goroutine泄漏: 无接收方的goroutine永久阻塞，要注意生命周期",
			"CSP原则: 不要通过共享内存通信，要通过通信共享内存",
		},
		Data: map[string]interface{}{
			"goroutine_results":    collected,
			"concurrent_elapsed":  elapsed.String(),
			"sequential_would_be": "300ms",
			"speedup":             fmt.Sprintf("并发约快%.1fx", 300.0/float64(elapsed.Milliseconds())),
		},
		Tips: []string{
			"Java: new Thread(() -> task()).start()  →  Go: go task()",
			"Java: ExecutorService  →  Go: goroutine + channel worker pool",
			"Java: Future.get()  →  Go: v := <-resultChan",
		},
	})
}
