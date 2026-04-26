package structs

import (
	"fmt"
	"math"
	"net/http"

	"gogo/internal/lesson"
)

// ── Shape接口的具体实现 ───────────────────────────────────────
// 以下类型在types.go中声明了Shape接口，此处只实现方法

type Circle struct{ Radius float64 }

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func (c Circle) Describe() string   { return fmt.Sprintf("Circle(r=%.2f)", c.Radius) }

type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }
func (r Rectangle) Describe() string {
	return fmt.Sprintf("Rectangle(%.2f×%.2f)", r.Width, r.Height)
}

type Triangle struct{ A, B, C float64 }

func (t Triangle) Area() float64 {
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}
func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }
func (t Triangle) Describe() string {
	return fmt.Sprintf("Triangle(%.2f,%.2f,%.2f)", t.A, t.B, t.C)
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func describeShape(s Shape) string {
	switch v := s.(type) { // 类型Switch
	case Circle:
		return fmt.Sprintf("这是圆形，半径=%.2f", v.Radius)
	case Rectangle:
		return fmt.Sprintf("这是矩形，%.2f×%.2f", v.Width, v.Height)
	default:
		return fmt.Sprintf("未知形状: %T", v)
	}
}

// 编译时验证Shape接口实现
var _ Shape = Circle{}
var _ Shape = Rectangle{}
var _ Shape = Triangle{}

// InterfacesHandler GET /api/structs/interfaces
func InterfacesHandler(w http.ResponseWriter, r *http.Request) {
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{A: 3, B: 4, C: 5},
	}

	// ── 多态：通过接口调用 ─────────────────────────────────────
	var shapeInfo []map[string]interface{}
	for _, s := range shapes {
		shapeInfo = append(shapeInfo, map[string]interface{}{
			"type":      s.Describe(),
			"area":      fmt.Sprintf("%.4f", s.Area()),
			"perimeter": fmt.Sprintf("%.4f", s.Perimeter()),
		})
	}

	// ── 类型断言（安全形式）──────────────────────────────────
	// Java: if (s instanceof Circle c) { ... }  // Java 16+
	// Go:   c, ok := s.(Circle)
	var s Shape = Circle{Radius: 3}
	circle, ok := s.(Circle)
	_, notRect := s.(Rectangle)

	// ── 类型Switch ────────────────────────────────────────────
	typeDesc1 := describeShape(Circle{Radius: 2})
	typeDesc2 := describeShape(Rectangle{Width: 3, Height: 4})

	// ── 空接口 interface{} / any ──────────────────────────────
	// Java: Object
	// Go 1.18+: any 是 interface{} 的别名
	var anything any = []int{1, 2, 3}

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go接口 (Interface)",
		Java:  "Java interface（但Go是隐式实现，无需implements）",
		Summary: "Go接口的最大特点是隐式实现：只要类型有接口要求的所有方法，" +
			"就自动实现该接口，无需implements声明。这使得解耦更彻底。",
		Points: []string{
			"隐式实现: 有方法就实现接口，无需implements，实现者无需知道接口存在",
			"类型断言: v, ok := i.(T)，ok=false时安全，不panic",
			"类型Switch: switch v := i.(type) { case int: ... }",
			"空接口 any(interface{}) 可接收任何类型，类似Java的Object",
			"接口组合: 一个接口可嵌入其他接口",
			"小接口原则: io.Reader只有1个方法，Go标准库偏好小而精的接口",
			"var _ Shape = Circle{} 编译时检查接口实现（推荐）",
		},
		Data: map[string]interface{}{
			"shapes":              shapeInfo,
			"total_area":          fmt.Sprintf("%.4f", totalArea(shapes)),
			"type_assert_ok":      circle.Radius,
			"assert_succeeded":    ok,
			"assert_failed":       notRect,
			"type_switch_circle":  typeDesc1,
			"type_switch_rect":    typeDesc2,
			"empty_interface":     fmt.Sprintf("%T: %v", anything, anything),
		},
	})
}
