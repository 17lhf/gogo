package generics

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// ── 泛型类型（Generic Types）────────────────────────────────

// Stack[T] 泛型栈
// Java: class Stack<T> { ... }
// Go:   type Stack[T any] struct { ... }
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int { return len(s.items) }

// Pair[K, V] 多类型参数
// Java: class Pair<K, V> { ... }
type Pair[K, V any] struct {
	Key   K
	Value V
}

func NewPair[K, V any](k K, v V) Pair[K, V] {
	return Pair[K, V]{Key: k, Value: v}
}

// Number 数字约束：只允许整数和浮点数
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

func Average[T Number](nums []T) float64 {
	if len(nums) == 0 {
		return 0
	}
	return float64(Sum(nums)) / float64(len(nums))
}

// ── ~运算符：约束底层类型 ─────────────────────────────────

type Celsius float64
type Fahrenheit float64

// 由于Celsius的底层类型是float64，满足Number约束中的~float64
// 但注意：Celsius + Fahrenheit不能直接相加（类型安全！）

// ConstraintsHandler GET /api/generics/constraints
func ConstraintsHandler(w http.ResponseWriter, r *http.Request) {
	// ── 泛型Stack使用 ─────────────────────────────────────────
	intStack := &Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)

	top, _ := intStack.Peek()
	popped, _ := intStack.Pop()

	strStack := &Stack[string]{}
	strStack.Push("Go")
	strStack.Push("Generics")
	strPeek, _ := strStack.Peek()

	// ── 多类型参数 ────────────────────────────────────────────
	pair1 := NewPair("user_id", 12345)
	pair2 := NewPair("name", "Alice")
	pair3 := NewPair[string, []int]("scores", []int{95, 87, 92})

	// ── 数字约束 ──────────────────────────────────────────────
	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.1, 2.2, 3.3}

	intSum := Sum(ints)
	floatSum := Sum(floats)
	intAvg := Average(ints)

	// ── ~运算符：自定义类型满足约束 ──────────────────────────
	temps := []Celsius{20, 25, 30, 35}
	celsiusSum := Sum(temps) // Celsius底层是float64，满足Number约束

	// ── 接口作为约束 vs 作为类型 ─────────────────────────────
	// Go 1.18之前：interface{}作为参数类型（运行时类型断言）
	// Go 1.18之后：[T any]作为类型参数（编译时类型安全）
	comparison := []map[string]string{
		{
			"approach":  "Go 1.18之前（interface{}）",
			"code":      `func Print(v interface{}) { fmt.Println(v) }`,
			"downside":  "运行时类型断言，可能panic，无类型安全",
		},
		{
			"approach":  "Go 1.18泛型（类型参数）",
			"code":      `func Print[T any](v T) { fmt.Println(v) }`,
			"benefit":   "编译时类型检查，无需类型断言，性能更好",
		},
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "泛型约束与泛型类型",
		Java:  "Java泛型类 class Stack<T> / 泛型约束 <T extends Number>",
		Summary: "Go泛型支持自定义约束接口，~运算符约束底层类型，" +
			"多类型参数，以及泛型类型（struct/interface）。",
		Points: []string{
			"泛型类型: type Stack[T any] struct{}，方法也需要类型参数",
			"多类型参数: func NewPair[K, V any](k K, v V) Pair[K, V]",
			"~ 运算符: ~int 匹配所有底层类型为int的类型（包括type MyInt int）",
			"接口作为约束: Number interface { ~int | ~float64 }",
			"泛型类型需要实例化: Stack[int]、Stack[string]",
			"Go泛型是单态化实现，与Java的类型擦除不同，运行时有完整类型信息",
		},
		Data: map[string]interface{}{
			"generic_stack": map[string]interface{}{
				"int_stack_peek":   top,
				"int_stack_pop":    popped,
				"int_stack_size":   intStack.Size(),
				"str_stack_peek":   strPeek,
			},
			"pairs": map[string]interface{}{
				"pair1": fmt.Sprintf("(%v, %v)", pair1.Key, pair1.Value),
				"pair2": fmt.Sprintf("(%v, %v)", pair2.Key, pair2.Value),
				"pair3": fmt.Sprintf("(%v, %v)", pair3.Key, pair3.Value),
			},
			"number_constraint": map[string]interface{}{
				"int_sum":      intSum,
				"float_sum":    fmt.Sprintf("%.1f", floatSum),
				"int_average":  intAvg,
				"celsius_sum":  celsiusSum,
			},
			"interface_vs_generics": comparison,
		},
	})
}
