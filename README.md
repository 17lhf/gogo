# gogo — Go 语言学习后台系统

专为 Java 开发者设计，通过可调用的 RESTful API 端点学习 Go 核心特性，每个端点演示一个概念并提供 Java 对比说明。

## 快速启动

```bash
go run .
# 服务启动在 http://localhost:8080
```

访问 `GET /` 查看所有端点的 JSON 索引，访问 `GET /health` 检查服务状态。

---

## 项目结构

```
gogo/
├── main.go                        # 入口，启动 HTTP 服务（port 8080）
├── api/
│   └── router.go                  # Chi 路由注册 + 首页/健康检查 handler
└── internal/
    ├── lesson/
    │   └── response.go            # 统一响应结构 Response / WriteJSON / WriteError
    ├── middleware/
    │   └── middleware.go          # CORS / RequestID / Logger 中间件
    ├── basics/
    │   ├── types.go               # 基础类型：变量声明/零值/类型转换/iota
    │   ├── slices.go              # 切片：make/append/copy/range/二维切片
    │   ├── maps.go                # 映射：安全读取/delete/range/map-as-Set
    │   ├── functions.go           # 函数：多返回值/命名返回/可变参数/高阶函数
    │   └── pointers.go            # 指针：&取址/*解引用/值传递vs指针传递
    ├── structs/
    │   ├── types.go               # 共享类型定义（Person/Shape/Animal 等）
    │   ├── basic.go               # 结构体：创建方式/方法/匿名结构体/struct tag
    │   ├── interfaces.go          # 接口：隐式实现/类型断言/类型 switch/any
    │   ├── embedding.go           # 嵌入：组合代替继承/方法提升/方法覆盖
    │   └── methods.go             # 方法：值接收者 vs 指针接收者/Stringer 接口
    ├── concurrency/
    │   ├── goroutines.go          # Goroutine：go 关键字/channel 收集结果
    │   ├── channels.go            # Channel：无缓冲/有缓冲/方向类型/range+close
    │   ├── selectstmt.go          # Select：多路复用/非阻塞/超时/for-select
    │   ├── waitgroup.go           # WaitGroup：Add/Done/Wait vs CountDownLatch
    │   └── mutex.go               # Mutex/RWMutex/atomic/Once/sync.Map
    ├── goerrors/
    │   ├── basic.go               # 基础错误：errors.New/fmt.Errorf/早返回模式
    │   ├── custom.go              # 自定义错误：errors.Is/As/哨兵错误
    │   ├── panic_recover.go       # panic/recover：defer+recover 安全包装
    │   ├── defer.go               # defer：LIFO 顺序/命名返回值修改/资源清理
    │   └── wrapping.go            # 错误包装：%w/errors.Unwrap/错误链
    ├── closures/
    │   ├── basic.go               # 闭包基础：变量捕获/函数工厂/私有状态
    │   └── functional.go          # 函数式：Map/Filter/Reduce/柯里化/管道
    └── generics/
        ├── basic.go               # 泛型函数：类型参数/类型推断/comparable
        └── constraints.go         # 泛型约束：~运算符/泛型类型 Stack[T]/Pair[K,V]
```

---

## 学习路径

按以下顺序访问 API 端点，由浅入深逐步掌握 Go：

### 第一阶段：基础语法（对应 Java 基础）

| 端点 | 学习内容 | Java 对比 |
|------|----------|-----------|
| `GET /api/v1/basics/types` | 变量声明、零值、类型转换、iota 枚举 | `int`/`String` 声明，`enum` |
| `GET /api/v1/basics/slices` | 动态数组切片 | `ArrayList<T>` |
| `GET /api/v1/basics/maps` | 哈希映射 | `HashMap<K,V>` |
| `GET /api/v1/basics/functions` | 多返回值、可变参数、高阶函数 | 方法重载、lambda |
| `GET /api/v1/basics/pointers` | 指针取址与解引用 | Java 引用传递 |

### 第二阶段：面向对象（Go 的方式）

| 端点 | 学习内容 | Java 对比 |
|------|----------|-----------|
| `GET /api/v1/structs/basic` | 结构体定义、方法、struct tag | `class`，`@JsonProperty` |
| `GET /api/v1/structs/interfaces` | 隐式接口实现、类型断言 | `implements`，`instanceof` |
| `GET /api/v1/structs/embedding` | 嵌入组合代替继承 | `extends` |
| `GET /api/v1/structs/methods` | 值接收者 vs 指针接收者 | `this` 引用 |

### 第三阶段：错误处理（Go 风格）

| 端点 | 学习内容 | Java 对比 |
|------|----------|-----------|
| `GET /api/v1/errors/basic` | error 接口、早返回模式 | `Exception`，`try-catch` |
| `GET /api/v1/errors/custom` | 自定义错误、errors.Is/As | `catch (MyException e)` |
| `GET /api/v1/errors/defer` | defer 延迟执行、LIFO 顺序 | `finally` |
| `GET /api/v1/errors/panic-recover` | panic/recover 机制 | `throw`/`catch` |
| `GET /api/v1/errors/wrapping` | 错误包装与错误链 | `getCause()` |

### 第四阶段：并发编程（Go 的核心优势）

| 端点 | 学习内容 | Java 对比 |
|------|----------|-----------|
| `GET /api/v1/concurrency/goroutines` | 轻量级线程，go 关键字 | `new Thread()`，`ExecutorService` |
| `GET /api/v1/concurrency/channels` | 通信管道，CSP 模型 | `BlockingQueue` |
| `GET /api/v1/concurrency/select` | 多路复用，超时控制 | `Selector`，`CompletableFuture` |
| `GET /api/v1/concurrency/waitgroup` | 等待多个 goroutine | `CountDownLatch` |
| `GET /api/v1/concurrency/mutex` | 互斥锁、原子操作、单例 | `synchronized`，`AtomicLong`，DCL |

### 第五阶段：函数式编程

| 端点 | 学习内容 | Java 对比 |
|------|----------|-----------|
| `GET /api/v1/closures/basic` | 闭包、变量捕获、函数工厂 | Lambda 捕获 effectively-final 变量 |
| `GET /api/v1/closures/functional` | Map/Filter/Reduce、柯里化 | `Stream.map/filter/reduce` |

### 第六阶段：泛型（Go 1.18+）

| 端点 | 学习内容 | Java 对比 |
|------|----------|-----------|
| `GET /api/v1/generics/basic` | 类型参数、类型推断 | `<T extends Comparable<T>>` |
| `GET /api/v1/generics/constraints` | `~` 底层类型约束、泛型结构体 | 有界通配符 `? extends Number` |

---

## 统一响应格式

所有端点返回相同的 JSON 结构：

```json
{
  "code": 200,
  "topic": "端点主题",
  "java_equivalent": "对应的 Java 概念",
  "summary": "一句话概述",
  "key_points": ["要点 1", "要点 2"],
  "data": { },
  "tips": ["提示 1"]
}
```

---

## 技术栈

| 组件 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 需支持泛型（1.18+） |
| chi | v5 | 轻量级 HTTP 路由器 |

---

## 模块文件说明

### go.mod

声明模块名称、Go 版本和直接依赖，类似 Maven 的 `pom.xml`：

```
module gogo               ← 模块名，import 路径的前缀（如 gogo/api）
go 1.26.1                 ← 最低 Go 版本要求
require
  github.com/go-chi/chi/v5 v5.2.5  ← 依赖包及其精确版本
```

### go.sum

依赖的**安全校验文件**，由 `go` 命令自动维护，**不可手动编辑，不支持注释**。

```
github.com/go-chi/chi/v5 v5.2.5        h1:<hash>  ← chi 源码包的 SHA-256 指纹
github.com/go-chi/chi/v5 v5.2.5/go.mod h1:<hash>  ← chi 的 go.mod 文件的 SHA-256 指纹
```

每次执行 `go build` / `go mod download` 时，Go 会重新计算下载内容的哈希并与此文件对比，不一致则构建失败，防止依赖被篡改（供应链攻击）。

**与 Java Maven 的对比：**

| Go | Maven |
|----|-------|
| `go.mod` | `pom.xml`（依赖声明） |
| `go.sum` | `.sha1` / `.md5` 文件（完整性校验） |

区别：Maven 的校验文件分散在各 jar 旁，Go 将所有依赖的哈希集中在项目根目录，**需要提交到 git**，保证团队使用相同的校验基准。
