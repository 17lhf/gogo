// package generics 演示Go 1.18+泛型特性
package generics

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// ── 泛型函数 ──────────────────────────────────────────────────
// Java: public <T extends Comparable<T>> T max(T a, T b)
// Go:   func Max[T constraints.Ordered](a, b T) T

// Ordered 约束：支持<>==的类型
// 等价于标准库中的 golang.org/x/exp/constraints.Ordered
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64 | ~string
}

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Contains[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

func Map[T, R any](s []T, fn func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func Filter[T any](s []T, fn func(T) bool) []T {
	var result []T
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T, R any](s []T, initial R, fn func(R, T) R) R {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// BasicGenericsHandler GET /api/generics/basic
func BasicGenericsHandler(w http.ResponseWriter, r *http.Request) {
	// ── 泛型函数调用 ──────────────────────────────────────────
	// Java: Integer.max(3, 7)，String需要另一套
	// Go:   Max(3, 7) 或 Max("apple", "banana")，同一个函数

	maxInt := Max(3, 7)
	maxStr := Max("apple", "banana")
	minFloat := Min(3.14, 2.71)

	// 类型推断：通常不需要显式指定类型参数
	// 显式写法: Max[int](3, 7)  推断写法: Max(3, 7)

	// ── 泛型切片操作 ──────────────────────────────────────────
	ints := []int{1, 2, 3, 4, 5}
	strs := []string{"Go", "Java", "Python", "Rust"}

	containsInt := Contains(ints, 3)
	containsStr := Contains(strs, "Java")
	notContains := Contains(strs, "Ruby")

	// ── 泛型Map/Filter/Reduce（通用版本）────────────────────
	// 在closures包中实现的只能处理int，泛型版本可以处理任何类型
	doubled := Map(ints, func(x int) int { return x * 2 })
	lengths := Map(strs, func(s string) int { return len(s) })
	longStrs := Filter(strs, func(s string) bool { return len(s) > 3 })
	sum := Reduce(ints, 0, func(acc, x int) int { return acc + x })
	joined := Reduce(strs, "", func(acc, s string) string {
		if acc == "" {
			return s
		}
		return acc + ", " + s
	})

	// ── 类型推断 ──────────────────────────────────────────────
	// 大多数情况下，Go能从参数推断类型参数
	typeInferDemo := map[string]string{
		"explicit": "Max[int](3, 7)  →  显式指定类型参数",
		"inferred": "Max(3, 7)  →  编译器自动推断为int",
		"note":     "Go泛型的类型推断比Java更强大，很少需要手动指定",
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go泛型基础 (Generics — Go 1.18+)",
		Java:  "Java泛型 <T extends Comparable<T>>",
		Summary: "Go 1.18引入泛型，用类型参数[T constraint]语法。" +
			"比Java泛型更简洁，支持~运算符约束底层类型，类型推断也更强。",
		Points: []string{
			"语法: func Func[T Constraint](args) ReturnType，类型参数在方括号中",
			"类型推断: 大多数情况编译器能自动推断，无需显式写Max[int](a, b)",
			"any = interface{}，comparable 是内置约束，支持==比较",
			"~ 运算符: ~int 表示底层类型是int的所有类型，包括自定义的 type MyInt int",
			"泛型函数比泛型类型更常用，Go推荐函数式风格",
			"Java的类型擦除导致运行时无类型信息；Go泛型是单态化，有完整类型信息",
		},
		Data: map[string]interface{}{
			"generic_max": map[string]interface{}{
				"max_int":    maxInt,
				"max_string": maxStr,
				"min_float":  fmt.Sprintf("%.2f", minFloat),
			},
			"contains": map[string]bool{
				"contains_3":    containsInt,
				"contains_java": containsStr,
				"contains_ruby": notContains,
			},
			"generic_collections": map[string]interface{}{
				"map_doubled":  doubled,
				"map_lengths":  lengths,
				"filter_long":  longStrs,
				"reduce_sum":   sum,
				"reduce_join":  joined,
			},
			"type_inference": typeInferDemo,
		},
	})
}
