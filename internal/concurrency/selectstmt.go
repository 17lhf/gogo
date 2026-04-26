package concurrency

import (
	"net/http"
	"time"

	"gogo/internal/lesson"
)

// SelectHandler GET /api/concurrency/select
// 演示select — 多channel多路复用
func SelectHandler(w http.ResponseWriter, r *http.Request) {
	// ── 竞速：选第一个ready的channel ─────────────────────────
	// Java: CompletableFuture.anyOf()
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	go func() { time.Sleep(30 * time.Millisecond); ch1 <- "来自channel1" }()
	go func() { time.Sleep(10 * time.Millisecond); ch2 <- "来自channel2" }() // ch2更快

	var winner string
	select {
	case msg := <-ch1:
		winner = "ch1: " + msg
	case msg := <-ch2:
		winner = "ch2: " + msg // ch2先ready，被选中
	}

	// ── 非阻塞select（带default）────────────────────────────
	nonblockCh := make(chan int, 1)
	var nonblockResult string
	select {
	case v := <-nonblockCh:
		nonblockResult = "收到: " + string(rune(v+'0'))
	default:
		nonblockResult = "channel空，走default，不阻塞"
	}

	// ── 超时控制 ─────────────────────────────────────────────
	// Java: future.get(50, TimeUnit.MILLISECONDS)
	// Go:   select + time.After(d)
	slowCh := make(chan string, 1)
	go func() {
		time.Sleep(200 * time.Millisecond)
		slowCh <- "慢操作结果"
	}()

	var timeoutResult string
	select {
	case result := <-slowCh:
		timeoutResult = result
	case <-time.After(50 * time.Millisecond):
		timeoutResult = "超时！操作超过50ms未完成"
	}

	// ── for-select：持续监听多个channel ─────────────────────
	// 常见模式：worker监听任务channel和停止信号
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	stopCh := make(chan struct{})
	go func() { time.Sleep(55 * time.Millisecond); close(stopCh) }()

	tickCount := 0
loop:
	for {
		select {
		case <-ticker.C:
			tickCount++
		case <-stopCh:
			break loop // 退出for循环
		}
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Select — 多Channel多路复用",
		Java:  "Java NIO Selector / CompletableFuture.anyOf()",
		Summary: "select监听多个channel，选择第一个就绪的执行。" +
			"是Go处理超时、取消、多路复用的核心机制。",
		Points: []string{
			"select等待多个channel，选择第一个ready的case",
			"多个case同时ready时，随机选择一个（公平调度）",
			"default case让select非阻塞，用于轮询检查",
			"select + time.After(d) 实现超时，比Java的Future.get更优雅",
			"close(ch) 可广播信号给所有等待goroutine（关闭后所有<-ch立即返回零值）",
			"for { select { ... } } 是持续监听多个channel的标准模式",
		},
		Data: map[string]interface{}{
			"race_winner":       winner,
			"nonblock_result":   nonblockResult,
			"timeout_result":    timeoutResult,
			"ticker_fire_count": tickCount,
		},
	})
}
