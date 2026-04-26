// package api 负责路由注册和统一响应
package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"gogo/internal/basics"
	"gogo/internal/closures"
	"gogo/internal/concurrency"
	"gogo/internal/generics"
	"gogo/internal/goerrors"
	"gogo/internal/middleware"
	"gogo/internal/structs"
)

// NewRouter 创建并返回配置好的HTTP路由器
// Java对比: @Configuration + @Bean RouterFunction 或 @RestController
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// ── 全局中间件（按顺序执行）──────────────────────────────
	r.Use(chimiddleware.Recoverer) // panic恢复，返回500
	r.Use(middleware.RequestID)    // 请求ID注入context
	r.Use(middleware.Logger)       // 请求日志
	r.Use(middleware.CORS)         // 跨域处理

	// ── 首页：API文档 ─────────────────────────────────────────
	r.Get("/", indexHandler)
	r.Get("/health", healthHandler)

	// ── 学习路由注册 ──────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {

		// 1. 基础类型与语法
		r.Route("/basics", func(r chi.Router) {
			r.Get("/types", basics.TypesHandler)
			r.Get("/slices", basics.SlicesHandler)
			r.Get("/maps", basics.MapsHandler)
			r.Get("/functions", basics.FunctionsHandler)
			r.Get("/pointers", basics.PointersHandler)
		})

		// 2. 结构体与面向对象
		r.Route("/structs", func(r chi.Router) {
			r.Get("/basic", structs.BasicStructHandler)
			r.Get("/interfaces", structs.InterfacesHandler)
			r.Get("/embedding", structs.EmbeddingHandler)
			r.Get("/methods", structs.MethodsHandler)
		})

		// 3. 并发编程
		r.Route("/concurrency", func(r chi.Router) {
			r.Get("/goroutines", concurrency.GoroutinesHandler)
			r.Get("/channels", concurrency.ChannelsHandler)
			r.Get("/select", concurrency.SelectHandler)
			r.Get("/waitgroup", concurrency.WaitGroupHandler)
			r.Get("/mutex", concurrency.MutexHandler)
		})

		// 4. 错误处理
		r.Route("/errors", func(r chi.Router) {
			r.Get("/basic", goerrors.BasicErrorHandler)
			r.Get("/custom", goerrors.CustomErrorHandler)
			r.Get("/panic-recover", goerrors.PanicRecoverHandler)
			r.Get("/defer", goerrors.DeferHandler)
			r.Get("/wrapping", goerrors.WrappingHandler)
		})

		// 5. 闭包与函数式
		r.Route("/closures", func(r chi.Router) {
			r.Get("/basic", closures.BasicClosureHandler)
			r.Get("/functional", closures.FunctionalHandler)
		})

		// 6. 泛型
		r.Route("/generics", func(r chi.Router) {
			r.Get("/basic", generics.BasicGenericsHandler)
			r.Get("/constraints", generics.ConstraintsHandler)
		})
	})

	return r
}

// indexHandler 返回所有可用端点的文档
func indexHandler(w http.ResponseWriter, r *http.Request) {
	type Endpoint struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Topic  string `json:"topic"`
	}
	type Group struct {
		Category  string     `json:"category"`
		Endpoints []Endpoint `json:"endpoints"`
	}

	groups := []Group{
		{
			Category: "基础类型与语法",
			Endpoints: []Endpoint{
				{"GET", "/api/v1/basics/types", "基础类型：变量/零值/类型转换/iota"},
				{"GET", "/api/v1/basics/slices", "切片：动态数组 vs Java ArrayList"},
				{"GET", "/api/v1/basics/maps", "映射：哈希表 vs Java HashMap"},
				{"GET", "/api/v1/basics/functions", "函数：多返回值/可变参数/高阶函数"},
				{"GET", "/api/v1/basics/pointers", "指针：&取址/*解引用 vs Java引用"},
			},
		},
		{
			Category: "结构体与面向对象",
			Endpoints: []Endpoint{
				{"GET", "/api/v1/structs/basic", "结构体：struct/标签/匿名结构体 vs Java class"},
				{"GET", "/api/v1/structs/interfaces", "接口：隐式实现/类型断言 vs Java interface"},
				{"GET", "/api/v1/structs/embedding", "嵌入：组合代替继承 vs Java extends"},
				{"GET", "/api/v1/structs/methods", "方法：值接收者vs指针接收者"},
			},
		},
		{
			Category: "并发编程",
			Endpoints: []Endpoint{
				{"GET", "/api/v1/concurrency/goroutines", "Goroutine：轻量级线程 vs Java Thread"},
				{"GET", "/api/v1/concurrency/channels", "Channel：通信管道 vs Java BlockingQueue"},
				{"GET", "/api/v1/concurrency/select", "Select：多路复用 vs Java Selector"},
				{"GET", "/api/v1/concurrency/waitgroup", "WaitGroup：等待多个goroutine vs CountDownLatch"},
				{"GET", "/api/v1/concurrency/mutex", "Mutex/Atomic/Once：同步原语 vs synchronized"},
			},
		},
		{
			Category: "错误处理",
			Endpoints: []Endpoint{
				{"GET", "/api/v1/errors/basic", "基础错误：error接口 vs Java Exception"},
				{"GET", "/api/v1/errors/custom", "自定义错误：errors.Is/As vs Java catch"},
				{"GET", "/api/v1/errors/panic-recover", "panic/recover：紧急停止/恢复"},
				{"GET", "/api/v1/errors/defer", "defer：延迟执行 vs Java finally"},
				{"GET", "/api/v1/errors/wrapping", "错误包装：%w和错误链"},
			},
		},
		{
			Category: "闭包与函数式编程",
			Endpoints: []Endpoint{
				{"GET", "/api/v1/closures/basic", "闭包基础：捕获变量/函数工厂/私有状态"},
				{"GET", "/api/v1/closures/functional", "函数式：Map/Filter/Reduce/柯里化"},
			},
		},
		{
			Category: "泛型 (Go 1.18+)",
			Endpoints: []Endpoint{
				{"GET", "/api/v1/generics/basic", "泛型函数：类型参数/类型推断"},
				{"GET", "/api/v1/generics/constraints", "泛型约束：自定义约束/泛型类型/~运算符"},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(map[string]interface{}{
		"name":        "Go语言学习后台系统",
		"description": "专为Java开发者设计的Go特性学习系统，每个端点演示一个Go核心概念并对比Java用法",
		"version":     "v1.0",
		"groups":      groups,
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`)) //nolint
}
