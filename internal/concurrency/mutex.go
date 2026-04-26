package concurrency

import (
	"net/http"
	"sync"
	"sync/atomic"

	"gogo/internal/lesson"
)

// MutexHandler GET /api/concurrency/mutex
// 演示互斥锁、读写锁、原子操作、sync.Once
func MutexHandler(w http.ResponseWriter, r *http.Request) {
	// ── sync.Mutex ────────────────────────────────────────────
	// Java: synchronized(this) { ... }  或  ReentrantLock
	type SafeCounter struct {
		mu    sync.Mutex
		count int
	}

	counter := &SafeCounter{}
	var wg1 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			counter.mu.Lock()
			defer counter.mu.Unlock() // defer确保锁一定被释放
			counter.count++
		}()
	}
	wg1.Wait()

	// ── sync.RWMutex — 读写锁 ─────────────────────────────────
	// Java: ReadWriteLock；多读单写，读多写少时性能更好
	type SafeMap struct {
		mu   sync.RWMutex
		data map[string]int
	}
	sm := &SafeMap{data: map[string]int{"hits": 0}}

	sm.mu.Lock()
	sm.data["hits"]++ // 写锁，独占
	sm.mu.Unlock()

	sm.mu.RLock()
	hits := sm.data["hits"] // 读锁，可并发
	sm.mu.RUnlock()

	// ── sync/atomic — 无锁原子操作 ────────────────────────────
	// Java: AtomicInteger / AtomicLong
	var atomicCount int64
	var wg2 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			atomic.AddInt64(&atomicCount, 1) // 无锁，性能最佳
		}()
	}
	wg2.Wait()

	// ── sync.Once — 只执行一次 ────────────────────────────────
	// Java: 双重检查锁单例模式（DCL），Go的Once更简洁安全
	var once sync.Once
	initResult := "未初始化"
	for i := 0; i < 5; i++ {
		once.Do(func() {
			initResult = "初始化完成（只执行一次，无论调用多少次）"
		})
	}

	// ── sync.Map — 并发安全的Map ──────────────────────────────
	// Java: ConcurrentHashMap
	var syncMap sync.Map
	syncMap.Store("key1", "value1")
	syncMap.Store("key2", "value2")
	v, _ := syncMap.Load("key1")
	syncMap.Delete("key2")

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "并发同步 — Mutex / Atomic / Once / sync.Map",
		Java:  "Java synchronized / ReentrantLock / AtomicInteger / ConcurrentHashMap",
		Summary: "Go提供多种并发同步原语。优先用channel通信；需要共享状态时用Mutex/Atomic。" +
			"sync.Once解决单例初始化，sync.Map用于并发读写Map。",
		Points: []string{
			"sync.Mutex: Lock()/Unlock()；defer mu.Unlock()确保解锁，避免忘记",
			"sync.RWMutex: RLock()/RUnlock()读锁可并发，Lock()/Unlock()写锁独占",
			"sync/atomic: 无锁原子操作，性能最好，只支持int32/int64/pointer等简单类型",
			"sync.Once: 保证函数只执行一次，线程安全的单例初始化",
			"sync.Map: 并发安全Map，适合读多写少；不如加锁的普通map灵活",
			"go run -race 开启数据竞争检测，帮助发现并发bug",
			"优先级: channel > Mutex > atomic（越底层越难用正确）",
		},
		Data: map[string]interface{}{
			"mutex_counter":  counter.count,
			"rwmutex_hits":   hits,
			"atomic_counter": atomicCount,
			"once_result":    initResult,
			"sync_map_key1":  v,
		},
		Tips: []string{
			"Java: synchronized(obj) { ... }  →  Go: mu.Lock(); defer mu.Unlock()",
			"Java: AtomicInteger.incrementAndGet()  →  Go: atomic.AddInt64(&n, 1)",
			"Java: volatile  →  Go: sync/atomic或channel（Go无volatile关键字）",
		},
	})
}
