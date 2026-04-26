package basics

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// FunctionsHandler GET /api/basics/functions
// 演示Go函数的特性：多返回值、可变参数、函数作为值、defer
func FunctionsHandler(w http.ResponseWriter, r *http.Request) {
	// ── 多返回值 ──────────────────────────────────────────────
	// Java没有多返回值，通常用异常或包装对象
	// Go可以返回多个值，最常见的是 (result, error) 模式
	quotient, remainder := divide(10, 3)

	// ── 命名返回值 ────────────────────────────────────────────
	// Java没有此特性
	min, max := minMax([]int{5, 2, 8, 1, 9, 3})

	// ── 可变参数 ──────────────────────────────────────────────
	// Java: public int sum(int... nums)
	// Go:   func sum(nums ...int) int
	total := sum(1, 2, 3, 4, 5)

	// 将切片展开传入可变参数
	nums := []int{10, 20, 30}
	totalFromSlice := sum(nums...) // ... 展开切片

	// ── 函数作为值（一等公民）────────────────────────────────
	// Java: Function<Integer, Integer> f = x -> x * 2;
	// Go:   f := func(x int) int { return x * 2 }
	double := func(x int) int { return x * 2 }
	triple := func(x int) int { return x * 3 }
	applyResults := []int{apply(5, double), apply(5, triple)}

	// ── 立即调用函数（IIFE）──────────────────────────────────
	iife := func(name string) string {
		return fmt.Sprintf("Hello, %s!", name)
	}("Go") // 定义后立即调用

	// ── defer示意（详见 /api/errors/defer）────────────────────
	deferOrder := explainDefer()

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go函数特性",
		Java:  "Java方法 + java.util.function包 + varargs",
		Summary: "Go函数是一等公民，支持多返回值(最常用于返回error)、" +
			"可变参数、函数作为参数/返回值。比Java更简洁，无需函数式接口。",
		Points: []string{
			"多返回值: func divide(a, b int) (int, int)，避免Java的包装对象",
			"命名返回值: func f() (result int, err error)，可以使用裸return",
			"可变参数: func sum(nums ...int)，用 slice... 展开切片传入",
			"函数是一等公民: 赋给变量、作为参数传递、作为返回值",
			"匿名函数: func(x int) int { return x*2 }，可立即调用",
			"defer在函数退出时按LIFO顺序执行，多个defer后进先出",
		},
		Data: map[string]interface{}{
			"multi_return": map[string]int{
				"quotient": quotient, "remainder": remainder,
			},
			"named_return": map[string]int{
				"min": min, "max": max,
			},
			"variadic_sum":        total,
			"variadic_from_slice": totalFromSlice,
			"function_as_value":   applyResults,
			"iife_result":         iife,
			"defer_order":         deferOrder,
		},
		Tips: []string{
			"_ 丢弃不需要的返回值: quotient, _ := divide(10, 3)",
			"(result, error) 是Go最常见函数签名，调用后立即检查error",
			"defer常用: defer file.Close(), defer mu.Unlock()",
		},
	})
}

func divide(a, b int) (int, int) {
	return a / b, a % b
}

// 命名返回值：函数签名声明返回变量名，可直接裸return
func minMax(arr []int) (min, max int) {
	min, max = arr[0], arr[0]
	for _, v := range arr[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return // 裸return，自动返回命名的min和max
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// apply 演示高阶函数：接收函数作为参数
// Java: public int apply(int x, Function<Integer, Integer> fn)
func apply(x int, fn func(int) int) int {
	return fn(x)
}

func explainDefer() []string {
	return []string{
		"1. 函数开始执行",
		"2. defer '第三个执行' 注册（LIFO：最后注册最先执行）",
		"3. defer '第二个执行' 注册",
		"4. defer '第一个执行' 注册",
		"5. 函数正常逻辑结束",
		"--- 函数返回前，defer按LIFO顺序执行 ---",
		"6. '第一个执行' (最后注册)",
		"7. '第二个执行'",
		"8. '第三个执行' (最早注册，最后执行)",
	}
}
