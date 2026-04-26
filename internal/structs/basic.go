package structs

import (
	"net/http"

	"gogo/internal/lesson"
)

// BasicStructHandler GET /api/structs/basic
// 演示结构体的创建、方法调用、标签、匿名结构体
func BasicStructHandler(w http.ResponseWriter, r *http.Request) {
	// ── 创建结构体 ────────────────────────────────────────────
	p1 := Person{Name: "Alice", Age: 28, email: "alice@example.com"} // 字段名初始化（推荐）
	p2 := NewPerson("Bob", 32, "bob@example.com")                    // 构造函数
	var p3 Person                                                     // 零值初始化
	p3.Name = "Carol"
	p3.Age = 25

	// ── 方法调用 ──────────────────────────────────────────────
	greeting := p1.Greet()
	p2.SetAge(33) // 指针接收者，p2.Age变为33

	// ── 匿名结构体：临时数据、配置、测试 ─────────────────────
	config := struct {
		Host string
		Port int
		TLS  bool
	}{Host: "localhost", Port: 8080, TLS: false}

	// ── 结构体比较 ────────────────────────────────────────────
	// 所有字段可比较时，结构体可用==比较（Java需重写equals）
	pt1 := struct{ X, Y int }{1, 2}
	pt2 := struct{ X, Y int }{1, 2}
	equal := pt1 == pt2

	// ── iota模拟枚举 ─────────────────────────────────────────
	// Go没有enum关键字，用const+iota模拟
	type Status int
	const (
		Active   Status = iota + 1 // 1
		Inactive                   // 2
		Banned                     // 3
	)

	// ── 结构体标签（Tag）─────────────────────────────────────
	// Java: @JsonProperty("user_name")
	// Go: `json:"user_name"` 结构体标签
	type UserDTO struct {
		UserName string `json:"user_name"`
		Email    string `json:"email,omitempty"` // omitempty: 空值不序列化
		Password string `json:"-"`               // -: 永远不序列化
	}
	dto := UserDTO{UserName: "alice", Email: "alice@example.com", Password: "secret"}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go结构体 (Struct)",
		Java:  "Java class（没有继承，用组合/嵌入代替）",
		Summary: "Go用struct代替class，没有constructor关键字(用New*函数惯例)，" +
			"方法与类型分离定义，首字母大写控制访问权限。",
		Points: []string{
			"首字母大写=public(exported)，小写=package-private，无private/protected",
			"方法定义在类型外部: func (p Person) Method()，比Java更灵活",
			"值接收者(p Person)不修改原值，指针接收者(p *Person)可修改",
			"构造函数惯例: func NewXxx() *Xxx {}，Go无constructor关键字",
			"结构体标签: `json:\"name\"` 控制序列化，类似Java的@JsonProperty",
			"匿名结构体用于临时数据，避免为一次性使用定义类型",
			"const + iota 模拟枚举，Go无enum关键字",
		},
		Data: map[string]interface{}{
			"person_p1":        p1,
			"person_p2":        p2,
			"person_p3":        p3,
			"greeting":         greeting,
			"after_set_age":    p2.Age,
			"anonymous_struct": config,
			"struct_equal":     equal,
			"user_status":      Active,
			"dto_tags":         dto,
		},
	})
}
