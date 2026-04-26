package concurrency

import (
	"net/http"

	"gogo/internal/lesson"
)

// ChannelsHandler GET /api/concurrency/channels
// 演示channel — goroutine间通信的管道
func ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	// ── 无缓冲channel ─────────────────────────────────────────
	// 发送方阻塞直到接收方ready，类似Java的SynchronousQueue
	unbuffered := make(chan int)
	go func() { unbuffered <- 42 }()
	val := <-unbuffered

	// ── 缓冲channel ───────────────────────────────────────────
	// 缓冲区未满时发送不阻塞，类似Java的ArrayBlockingQueue
	buffered := make(chan string, 3)
	buffered <- "消息1"
	buffered <- "消息2"
	buffered <- "消息3"

	msg1 := <-buffered
	msg2 := <-buffered
	msg3 := <-buffered

	// ── 单向channel ───────────────────────────────────────────
	// 限制channel方向，增强类型安全
	oneway := make(chan int, 1)
	sendOnly(oneway)
	received := recvOnly(oneway)

	// ── range遍历channel ──────────────────────────────────────
	// close(ch)后，range自动退出
	numbers := make(chan int, 5)
	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
		close(numbers) // 通知接收方没有更多数据
	}()

	var rangeResult []int
	for n := range numbers {
		rangeResult = append(rangeResult, n)
	}

	// ── 检查channel是否关闭 ───────────────────────────────────
	ch := make(chan int, 1)
	ch <- 100
	close(ch)
	v1, open1 := <-ch // 100, true
	v2, open2 := <-ch // 0, false（已关闭且空）

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Channel — Goroutine间通信",
		Java:  "Java BlockingQueue / SynchronousQueue / LinkedBlockingQueue",
		Summary: "Channel是Go并发的核心，实现CSP(通信顺序进程)模型。" +
			"通过channel传递数据而不是共享内存，天然避免数据竞争。",
		Points: []string{
			"make(chan T) 无缓冲channel，同步通信，发送接收必须同时ready",
			"make(chan T, n) 缓冲channel，异步通信，缓冲区未满/非空时不阻塞",
			"close(ch) 关闭channel；v, ok := <-ch 检测是否关闭",
			"for v := range ch 遍历channel，close后自动退出",
			"chan<- T 只写channel，<-chan T 只读channel",
			"向已关闭的channel发送会panic，从已关闭的channel接收返回零值",
		},
		Data: map[string]interface{}{
			"unbuffered_val":   val,
			"buffered_msgs":    []string{msg1, msg2, msg3},
			"unidirectional":   received,
			"range_results":    rangeResult,
			"closed_channel": map[string]interface{}{
				"first_val":  v1, "first_open":  open1,
				"second_val": v2, "second_open": open2,
			},
		},
	})
}

func sendOnly(ch chan<- int) { ch <- 999 } // 只能发送
func recvOnly(ch <-chan int) int { return <-ch } // 只能接收
