package structs

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// MethodsHandler GET /api/structs/methods
// 演示值接收者vs指针接收者的深层区别
func MethodsHandler(w http.ResponseWriter, r *http.Request) {
	// ── 值接收者 vs 指针接收者 ────────────────────────────────
	type Counter struct{ count int }

	// 值接收者版：每次返回新Counter，不修改原值
	addValue := func(c Counter, n int) Counter {
		c.count += n
		return c
	}

	// 指针接收者版：直接修改原值
	addPointer := func(c *Counter, n int) {
		c.count += n
	}

	c1 := Counter{count: 0}
	r1 := addValue(c1, 5)   // c1.count仍为0
	r2 := addValue(c1, 10)  // c1.count仍为0

	c2 := &Counter{count: 0}
	addPointer(c2, 5)   // c2.count = 5
	addPointer(c2, 10)  // c2.count = 15

	// ── fmt.Stringer 接口 ─────────────────────────────────────
	// Java: @Override public String toString()
	// Go: 实现 String() string 方法
	p := Person{Name: "Alice", Age: 28}
	formatted := fmt.Sprintf("%v", p) // fmt自动调用String()

	// ── 方法集规则 ────────────────────────────────────────────
	// T的方法集: 只包含值接收者方法
	// *T的方法集: 包含值接收者方法 + 指针接收者方法
	// 所以: 指针*T可调用值接收者方法，但值T不能调用指针接收者方法
	// （这影响接口实现：如果方法用指针接收者，只有*T满足接口）

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go方法与接收者",
		Java:  "Java实例方法（this引用）",
		Summary: "Go方法通过接收者关联到类型。值接收者获得副本不影响原值，" +
			"指针接收者可修改原始数据。接口实现时，方法集规则决定是否满足接口。",
		Points: []string{
			"值接收者(p Person): 操作副本，原值不变，适合只读操作",
			"指针接收者(p *Person): 操作原值，适合修改操作或大struct",
			"一致性原则: 同一类型的方法，要么全用值接收者，要么全用指针接收者",
			"*T拥有T的所有方法，T不包含*T的方法（影响接口实现！）",
			"fmt.Stringer: 实现String() string，控制%v输出，类似Java的toString()",
			"方法可定义在任何包内的命名类型上，不限于struct",
		},
		Data: map[string]interface{}{
			"value_receiver": map[string]interface{}{
				"c1_original": c1.count,
				"result1":     r1.count,
				"result2":     r2.count,
				"note":        "c1.count仍为0，值接收者不修改原值",
			},
			"pointer_receiver": map[string]interface{}{
				"final_count": c2.count,
				"note":        "c2.count=15，指针接收者累加",
			},
			"stringer": map[string]string{
				"fmt_output": formatted,
				"note":       "fmt自动调用Person.String()方法",
			},
		},
		Tips: []string{
			"如果某个方法需要修改接收者，就用*T；否则根据大小和语义选择",
			"实现接口时：方法用*T接收者，则只有*T满足接口，T不满足",
			"大struct用*T避免拷贝；小struct(<=3字段)直接用T也可以",
		},
	})
}
