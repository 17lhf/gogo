package structs

import (
	"net/http"

	"gogo/internal/lesson"
)

// EmbeddingHandler GET /api/structs/embedding
// 演示嵌入（Embedding）— Go的组合代替继承
func EmbeddingHandler(w http.ResponseWriter, r *http.Request) {
	dog := Dog{Animal: Animal{Name: "旺财"}, Breed: "柴犬"}
	cat := Cat{Animal: Animal{Name: "咪咪"}, Indoor: true}
	duck := Duck{Animal: Animal{Name: "唐老鸭"}}

	// Dog继承Animal的Eat，但覆盖了Breathe
	dogEat := dog.Eat("骨头")           // 来自Animal（提升的方法）
	dogBreathe := dog.Breathe()        // Dog自己覆盖的方法
	dogFetch := dog.Fetch()            // Dog特有的方法
	dogAnimalBreathe := dog.Animal.Breathe() // 访问被覆盖的原方法

	catPurr := cat.Purr()
	catEat := cat.Eat("鱼") // 来自Animal

	// Duck隐式实现了Swimmer接口
	var swimmer Swimmer = duck
	swimResult := swimmer.Swim()

	// ── 多嵌入 ────────────────────────────────────────────────
	svc := Service{
		Animal:  Animal{Name: "订单服务"},
		Logger:  Logger{Prefix: "INFO"},
		Version: "v1.0",
	}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go嵌入 (Embedding) — 组合优于继承",
		Java:  "Java extends（但Go用组合而非继承）",
		Summary: "Go不支持继承，通过嵌入(Embedding)实现代码复用。" +
			"嵌入类型的方法被提升到外层类型，可覆盖，也可通过类型名访问原方法。",
		Points: []string{
			"嵌入写法: type Dog struct { Animal }，不是 animal Animal（非字段名）",
			"方法提升: Dog自动拥有Animal的所有方法，无需委托代码",
			"方法覆盖: Dog定义同名方法即覆盖，不需要override关键字",
			"访问原方法: dog.Animal.Breathe() 仍可访问被覆盖的实现",
			"多嵌入: 可同时嵌入多个类型，解决Java单继承的限制",
			"嵌入接口: struct中嵌入接口类型，可实现装饰器模式",
			"组合优于继承: 比Java继承更灵活，避免继承链的强耦合",
		},
		Data: map[string]interface{}{
			"dog": map[string]string{
				"eat":              dogEat,
				"breathe_override": dogBreathe,
				"fetch":            dogFetch,
				"breathe_original": dogAnimalBreathe,
			},
			"cat": map[string]string{
				"purr": catPurr,
				"eat":  catEat,
			},
			"duck_as_swimmer": swimResult,
			"service_multi_embed": map[string]string{
				"log": svc.Log("服务已启动"),
				"eat": svc.Eat("数据"),
			},
		},
		Tips: []string{
			"Java: class Dog extends Animal  →  Go: type Dog struct { Animal }",
			"Java: @Override  →  Go: 直接定义同名方法即可覆盖",
			"Go推荐组合：Dog has-an Animal，而非Dog is-an Animal",
		},
	})
}
