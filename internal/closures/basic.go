// package closures 演示Go的闭包与函数式编程特性
package closures

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// BasicClosureHandler GET /api/closures/basic
// 演示闭包的基本概念和常见用法
func BasicClosureHandler(w http.ResponseWriter, r *http.Request) {
	// ── 闭包基础 ──────────────────────────────────────────────
	// 闭包：函数 + 捕获的外部变量（形成一个封闭的作用域）
	// Java: lambda表达式只能捕获effectively final变量
	// Go:   可以捕获并修改外部变量（更强大！）

	// 计数器闭包
	// Java: 需要AtomicInteger或用类封装状态
	// Go:   直接用闭包捕获变量
	counter := makeCounter()
	c1 := counter() // 1
	c2 := counter() // 2
	c3 := counter() // 3

	// 独立的计数器，互不影响
	counterA := makeCounter()
	counterB := makeCounter()
	a1, a2 := counterA(), counterA()
	b1 := counterB()

	// ── 闭包捕获变量引用 ──────────────────────────────────────
	// 注意：捕获的是变量本身，不是值的拷贝
	x := 10
	addX := func(n int) int { return n + x } // 捕获x的引用
	result1 := addX(5) // 15
	x = 20             // 修改x
	result2 := addX(5) // 25（看到了x的新值！）

	// ── 函数工厂（返回函数的函数）────────────────────────────
	// Java: 函数式接口 + lambda
	double := multiplier(2)
	triple := multiplier(3)

	// ── 闭包实现私有状态 ──────────────────────────────────────
	// Java: 需要class封装private字段
	// Go:   闭包自然封装，外部无法访问内部变量
	deposit, withdraw, balance := makeBankAccount(100)
	deposit(50)
	withdraw(30)
	currentBalance := balance()

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go闭包 (Closure)",
		Java:  "Java Lambda表达式（但Go可修改捕获的变量）",
		Summary: "闭包是函数+捕获的外部变量。Go的闭包比Java的lambda更强大：" +
			"可以捕获并修改外部变量，是实现私有状态、函数工厂的核心工具。",
		Points: []string{
			"闭包捕获变量引用，不是值的拷贝，外部修改变量，闭包能看到新值",
			"Java lambda只能捕获effectively final变量；Go可以修改捕获的变量",
			"函数工厂: 返回函数的函数，每次调用返回携带不同状态的函数",
			"每次调用工厂函数返回独立的闭包，各自维护独立状态",
			"闭包可封装私有状态，外部无法访问内部变量（类似Java的private）",
			"for循环中使用goroutine+闭包时注意变量捕获: i := i",
		},
		Data: map[string]interface{}{
			"single_counter": map[string]int{
				"call1": c1, "call2": c2, "call3": c3,
			},
			"independent_counters": map[string]int{
				"counterA_call1": a1, "counterA_call2": a2,
				"counterB_call1": b1,
			},
			"capture_by_reference": map[string]interface{}{
				"result_when_x_10": result1,
				"result_when_x_20": result2,
				"note":             "闭包看到的是x的当前值，不是创建时的值",
			},
			"function_factory": map[string]int{
				"double_5": double(5),
				"triple_5": triple(5),
			},
			"bank_account": map[string]interface{}{
				"initial":  100,
				"deposit":  50,
				"withdraw": 30,
				"balance":  currentBalance,
			},
		},
	})
}

func makeCounter() func() int {
	count := 0 // 被闭包捕获的局部变量
	return func() int {
		count++
		return count
	}
}

func multiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor // 捕获factor
	}
}

func makeBankAccount(initial int) (deposit func(int), withdraw func(int), balance func() int) {
	amount := initial // 私有状态，外部无法直接访问
	deposit = func(n int) { amount += n }
	withdraw = func(n int) {
		if n <= amount {
			amount -= n
		}
	}
	balance = func() int { return amount }
	return
}

// 确保fmt被使用
var _ = fmt.Sprintf
