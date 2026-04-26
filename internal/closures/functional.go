package closures

import (
	"net/http"
	"strings"

	"gogo/internal/lesson"
)

// FunctionalHandler GET /api/closures/functional
// 演示Go的函数式编程：map/filter/reduce/柯里化
func FunctionalHandler(w http.ResponseWriter, r *http.Request) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// ── Map：对每个元素应用函数 ───────────────────────────────
	// Java: list.stream().map(x -> x * 2).collect(toList())
	// Go:   mapInts(list, func(x int) int { return x * 2 })
	doubled := mapInts(nums, func(x int) int { return x * 2 })
	squared := mapInts(nums, func(x int) int { return x * x })

	// ── Filter：过滤满足条件的元素 ────────────────────────────
	// Java: list.stream().filter(x -> x % 2 == 0).collect(toList())
	evens := filterInts(nums, func(x int) bool { return x%2 == 0 })
	greaterThan5 := filterInts(nums, func(x int) bool { return x > 5 })

	// ── Reduce：聚合所有元素 ──────────────────────────────────
	// Java: list.stream().reduce(0, Integer::sum)
	sum := reduceInts(nums, 0, func(acc, x int) int { return acc + x })
	product := reduceInts([]int{1, 2, 3, 4, 5}, 1, func(acc, x int) int { return acc * x })

	// ── 函数组合（Pipeline）──────────────────────────────────
	// 先过滤偶数，再翻倍，再求和
	pipelineResult := reduceInts(
		mapInts(
			filterInts(nums, func(x int) bool { return x%2 == 0 }),
			func(x int) int { return x * 2 },
		),
		0,
		func(acc, x int) int { return acc + x },
	)

	// ── 字符串处理链 ─────────────────────────────────────────
	words := []string{"Go", "is", "awesome", "and", "fast"}
	upper := mapStrings(words, strings.ToUpper)
	longWords := filterStrings(words, func(s string) bool { return len(s) > 2 })

	// ── 柯里化（Currying）────────────────────────────────────
	// Java: 需要多层Function嵌套
	// Go:   返回函数的函数
	add := curry(func(a, b int) int { return a + b })
	add5 := add(5)     // 部分应用，固定第一个参数为5
	add10 := add(10)

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "函数式编程 — Map/Filter/Reduce/柯里化",
		Java:  "Java Stream API (map/filter/reduce)",
		Summary: "Go没有内置的Stream API，但通过高阶函数和闭包可以实现相同效果。" +
			"Go 1.23+的 range-over-func 和 iter 包开始原生支持函数式迭代。",
		Points: []string{
			"高阶函数: 接收或返回函数的函数，是函数式编程的基础",
			"Go没有内置map/filter/reduce，但很容易用泛型实现通用版本",
			"函数组合: 将多个函数串联，每个函数的输出作为下一个的输入",
			"柯里化: 多参数函数转化为一系列单参数函数，实现部分应用",
			"Go 1.18+泛型可以写出真正通用的map/filter/reduce，不需要为每个类型重写",
			"Java Stream是惰性求值；Go的函数式代码是立即求值的",
		},
		Data: map[string]interface{}{
			"map": map[string][]int{
				"original": nums,
				"doubled":  doubled,
				"squared":  squared,
			},
			"filter": map[string][]int{
				"evens":         evens,
				"greater_than5": greaterThan5,
			},
			"reduce": map[string]int{
				"sum":     sum,
				"product": product,
			},
			"pipeline": map[string]interface{}{
				"steps":  "filter(even) → map(*2) → reduce(+)",
				"result": pipelineResult,
			},
			"string_ops": map[string]interface{}{
				"upper":      upper,
				"long_words": longWords,
			},
			"currying": map[string]int{
				"add5_to_3":  add5(3),
				"add5_to_7":  add5(7),
				"add10_to_3": add10(3),
			},
		},
	})
}

func mapInts(s []int, fn func(int) int) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func filterInts(s []int, fn func(int) bool) []int {
	var result []int
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func reduceInts(s []int, initial int, fn func(int, int) int) int {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

func mapStrings(s []string, fn func(string) string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func filterStrings(s []string, fn func(string) bool) []string {
	var result []string
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// curry 将二元函数转为柯里化形式
func curry(fn func(int, int) int) func(int) func(int) int {
	return func(a int) func(int) int {
		return func(b int) int {
			return fn(a, b)
		}
	}
}
