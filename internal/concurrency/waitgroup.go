package concurrency

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gogo/internal/lesson"
)

// WaitGroupHandler GET /api/concurrency/waitgroup
// 演示sync.WaitGroup — 等待一组goroutine完成
func WaitGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Java: CountDownLatch / CyclicBarrier / CompletableFuture.allOf()
	// Go:   sync.WaitGroup

	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []string

	workers := []struct {
		name  string
		delay int
	}{
		{"数据库查询", 50},
		{"缓存获取", 20},
		{"外部API调用", 80},
		{"文件读取", 30},
		{"权限验证", 10},
	}

	start := time.Now()

	for _, worker := range workers {
		wg.Add(1)        // 计数+1，必须在goroutine启动前调用
		worker := worker // 捕获循环变量
		go func() {
			defer wg.Done() // goroutine结束时计数-1，类似Java的latch.countDown()

			time.Sleep(time.Duration(worker.delay) * time.Millisecond)

			mu.Lock()
			results = append(results, fmt.Sprintf("%s完成(耗时%dms)", worker.name, worker.delay))
			mu.Unlock()
		}()
	}

	wg.Wait() // 阻塞直到计数归零，类似Java的latch.await()
	elapsed := time.Since(start)

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "sync.WaitGroup — 等待多个Goroutine",
		Java:  "Java CountDownLatch / CompletableFuture.allOf()",
		Summary: "WaitGroup等待一组goroutine全部完成。" +
			"Add(n)增加计数，Done()减少计数，Wait()阻塞到计数为0。",
		Points: []string{
			"wg.Add(1) 必须在启动goroutine之前调用，否则可能在goroutine启动前Wait就返回",
			"defer wg.Done() 确保goroutine无论如何退出都会调用Done",
			"wg.Wait() 阻塞调用方，直到所有goroutine调用Done",
			"WaitGroup可复用，Wait()返回后可再次Add()",
			"WaitGroup不保护共享数据，并发写共享变量还需要Mutex",
			"收集goroutine结果：用Mutex保护共享slice，或用channel",
		},
		Data: map[string]interface{}{
			"workers_count":    len(workers),
			"results":          results,
			"concurrent_time":  elapsed.String(),
			"sequential_would": "190ms",
			"critical_path":    "≈80ms（最长任务）",
		},
	})
}
