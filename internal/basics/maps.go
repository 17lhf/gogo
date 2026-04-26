package basics

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// MapsHandler GET /api/basics/maps
// 演示Map — Go内置的哈希表
func MapsHandler(w http.ResponseWriter, r *http.Request) {
	// ── 创建 ──────────────────────────────────────────────────
	// Java: new HashMap<String, Integer>()
	scores := map[string]int{"Alice": 95, "Bob": 87, "Carol": 92}

	phoneBook := make(map[string]string)
	phoneBook["Alice"] = "138-0000-0001"
	phoneBook["Bob"] = "139-0000-0002"

	// ── 安全读取（两值赋值）──────────────────────────────────
	// Java: map.get("key") 不存在返回null，可能NPE
	// Go:   v, ok := map[key]，ok=false时v是零值，安全！
	aliceScore, exists := scores["Alice"]
	nobody, notExists := scores["Nobody"] // nobody=0(零值), notExists=false

	// ── 删除 ──────────────────────────────────────────────────
	// Java: map.remove("key")
	delete(scores, "Bob")

	// ── 遍历（顺序随机！）────────────────────────────────────
	var entries []string
	for k, v := range scores {
		entries = append(entries, fmt.Sprintf("%s:%d", k, v))
	}

	// ── 嵌套Map ───────────────────────────────────────────────
	users := map[string]map[string]interface{}{
		"user_1": {"name": "Alice", "age": 28, "active": true},
		"user_2": {"name": "Bob", "age": 32, "active": false},
	}

	// ── 用map模拟Set ──────────────────────────────────────────
	// Go没有内置Set，用 map[T]struct{} 模拟（struct{}不占内存）
	seen := map[string]struct{}{}
	words := []string{"go", "is", "great", "go", "is", "fast"}
	for _, w := range words {
		seen[w] = struct{}{}
	}
	unique := make([]string, 0, len(seen))
	for k := range seen {
		unique = append(unique, k)
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go映射 (Map)",
		Java:  "Java HashMap<K, V>",
		Summary: "Go的map是内置哈希表。最重要的差异是用双值接收(value, ok)" +
			"判断key是否存在，避免Java的null陷阱。并发读写需要加锁！",
		Points: []string{
			"声明: map[KeyType]ValueType，key必须是可比较类型",
			"安全读取: v, ok := m[key]，ok=false时v是零值，不会panic",
			"delete(m, key) 删除键，删除不存在的key不报错",
			"遍历顺序随机，需要有序遍历要先对key排序",
			"map[T]struct{} 模拟Set，struct{}是空结构体，不占内存",
			"map是引用类型，函数传递map不需要指针",
			"并发读写map会panic，多goroutine需用sync.Map或加锁",
		},
		Data: map[string]interface{}{
			"scores_map":    scores,
			"phone_book":    phoneBook,
			"alice_score":   aliceScore,
			"key_exists":    exists,
			"nobody_score":  nobody,
			"key_missing":   notExists,
			"after_delete":  scores,
			"entries":       entries,
			"nested_map":    users,
			"set_unique":    unique,
		},
	})
}
