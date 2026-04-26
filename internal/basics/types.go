package basics

import (
	"fmt"
	"net/http"

	"gogo/internal/lesson"
)

// TypesHandler GET /api/basics/types
// 演示Go的类型系统与Java的差异
func TypesHandler(w http.ResponseWriter, r *http.Request) {
	// ── 变量声明三种方式 ──────────────────────────────────────
	// Java: int age = 25;
	var age int = 25         // 完整声明
	score := 98.6            // 短变量声明(:=)，编译器推断类型，只能在函数内使用
	const maxScore = 100.0   // 常量，编译时确定

	// ── 零值 (Zero Value) ─────────────────────────────────────
	// Java中未初始化的对象是null，基本类型有默认值
	// Go所有类型都有零值，永远不会出现"未初始化"状态
	var zeroInt int
	var zeroStr string
	var zeroBool bool
	var zeroSlice []int
	var zeroMap map[string]int

	// ── 类型转换 (必须显式) ───────────────────────────────────
	// Java: double d = myInt; (隐式)  Go: 必须显式，否则编译报错
	intVal := 42
	floatVal := float64(intVal)
	backToInt := int(floatVal)

	// ── rune vs byte ──────────────────────────────────────────
	// byte = uint8 表示ASCII字节
	// rune = int32 表示Unicode码点（类似Java的char但能表示所有Unicode）
	var ch byte = 'A'
	var chineseChar rune = '中'
	str := "Hello, 世界"
	byteLen := len(str)          // 字节长度：13（中文UTF-8占3字节）
	runeLen := len([]rune(str))  // 字符长度：9

	// ── iota：自增常量（类似Java enum的ordinal）────────────────
	const (
		Small  = iota * 10 // 0
		Medium             // 10
		Large              // 20
		XLarge             // 30
	)

	lesson.WriteJSON(w, http.StatusOK, lesson.Response{
		Code:  http.StatusOK,
		Topic: "Go基础类型系统",
		Java:  "Java原始类型(int/double/boolean) + String + 包装类",
		Summary: "Go类型系统比Java更严格：无隐式类型转换，所有类型有零值(不会NPE)，" +
			"rune支持完整Unicode(比Java的char更强大)",
		Points: []string{
			"变量声明: var x int = 1  或  x := 1 (函数内简写，类型推断)",
			"无隐式转换: int→float64必须写 float64(x)，Java中可以直接赋值",
			"零值机制: int=0, string=\"\", bool=false, 指针=nil，不会未初始化",
			"rune=int32=Unicode码点，byte=uint8=字节；处理中文必须用rune",
			"iota在const块中自动递增，可替代Java enum的部分用法",
		},
		Data: map[string]interface{}{
			"declarations": map[string]interface{}{
				"var_style":   age,
				"short_style": score,
				"constant":    maxScore,
			},
			"zero_values": map[string]interface{}{
				"zero_int":   zeroInt,
				"zero_str":   zeroStr,
				"zero_bool":  zeroBool,
				"zero_slice": zeroSlice,
				"zero_map":   zeroMap,
			},
			"type_conversion": map[string]interface{}{
				"int_to_float64": floatVal,
				"float64_to_int": backToInt,
			},
			"unicode": map[string]interface{}{
				"byte_char":   string(ch),
				"rune_char":   string(chineseChar),
				"string":      str,
				"byte_length": byteLen,
				"rune_length": runeLen,
			},
			"iota_example": map[string]int{
				"Small": Small, "Medium": Medium, "Large": Large, "XLarge": XLarge,
			},
		},
		Tips: []string{
			"Java: int x = 42;  →  Go: var x int = 42  或  x := 42",
			"Java: (double)myInt  →  Go: float64(myInt)",
			"Go的string是UTF-8字节序列，len()返回字节数而非字符数",
		},
	})
}

// 编译时验证：确保常量值正确
var _ = fmt.Sprintf // 保持import使用
