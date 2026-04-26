package basics

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// PointersHandler GET /api/basics/pointers
// 演示Go的指针 — 比C简单，比Java更明确
func PointersHandler(w http.ResponseWriter, r *http.Request) {
	// ── 基本指针操作 ──────────────────────────────────────────
	// Java: 对象都是引用，但无法直接操作地址
	// Go: 有真正的指针，但没有指针运算（更安全）
	x := 42
	ptr := &x    // & 取地址，ptr是 *int 类型
	*ptr = 100   // * 解引用，修改ptr指向的值
	// 此时 x == 100

	// ── 函数参数传递 ──────────────────────────────────────────
	// Java: 基本类型值传递，对象引用传递（实际上也是值传递，传的是引用的副本）
	// Go: 所有参数都是值传递！要修改原值必须传指针
	a, b := 5, 10
	swapByValue(a, b)      // a,b不变（拷贝）
	swapByPointer(&a, &b)  // a,b互换（指针）

	// ── new vs make ───────────────────────────────────────────
	// new(T): 分配T类型零值内存，返回 *T（很少用）
	// make(T): 只用于 slice/map/channel，返回 T（非指针），初始化内部结构
	numPtr := new(int) // *int，指向0
	*numPtr = 999

	// ── nil指针 ───────────────────────────────────────────────
	var nilPtr *int // nil，不指向任何地址，解引用会panic

	// ── 结构体指针 ────────────────────────────────────────────
	type Point struct{ X, Y int }
	p1 := Point{1, 2}   // 值类型
	p2 := &Point{3, 4}  // 指针类型 *Point
	p2.X = 10           // Go自动解引用，等价于 (*p2).X = 10

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go指针",
		Java:  "Java的对象引用（无法直接操作地址）",
		Summary: "Go有真正的指针(&取地址, *解引用)，但无指针运算。" +
			"关键差异：Go所有参数值传递，要修改原值必须传指针；Java对象是引用无需指针。",
		Points: []string{
			"& 取变量地址得到指针，* 解引用指针得到值",
			"Go所有参数值传递，修改原始值必须传指针(*T)",
			"结构体指针可用.直接访问字段，Go自动解引用: p.X 等价于 (*p).X",
			"new(T) 分配零值内存返回*T；make(T) 用于slice/map/channel初始化",
			"nil指针解引用会panic，使用前需判断 ptr != nil",
			"传递大结构体用指针避免拷贝；小struct直接传值即可",
		},
		Data: map[string]interface{}{
			"pointer_demo": map[string]interface{}{
				"original_x":   x,
				"after_deref_modify": fmt.Sprintf("x通过指针被修改为%d", x),
			},
			"value_vs_pointer_swap": map[string]interface{}{
				"before":            fmt.Sprintf("a=%d, b=%d", 5, 10),
				"after_value_swap":  fmt.Sprintf("a=%d, b=%d (未改变)", 5, 10),
				"after_pointer_swap": fmt.Sprintf("a=%d, b=%d (已互换)", a, b),
			},
			"new_example":  *numPtr,
			"nil_pointer":  nilPtr == nil,
			"struct_value":   fmt.Sprintf("Point{%d, %d}", p1.X, p1.Y),
			"struct_pointer": fmt.Sprintf("&Point{%d, %d}", p2.X, p2.Y),
		},
		Tips: []string{
			"何时使用指针: 1)需要修改参数 2)大结构体避免拷贝 3)方法需修改接收者字段",
			"方法接收者用(*T)可修改字段，用(T)只操作副本",
			"Go没有Java的null关键字，对应的是nil",
		},
	})
}

func swapByValue(a, b int) {
	a, b = b, a // 只交换局部变量
}

func swapByPointer(a, b *int) {
	*a, *b = *b, *a // 通过指针修改调用者变量
}
