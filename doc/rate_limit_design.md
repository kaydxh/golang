# 限流增强功能设计文档

## 1. 概述

本文档描述了 golang 库中限流（Rate Limit）功能的增强设计，支持**基于 QPS 的限流**、**并发数控制**和**不同接口不同限流策略**的能力。

## 2. 限流模式分析

### 2.1 并发控制限流器 (Limiter)

基于**Channel 信号量**的并发控制器：

```go
type Limiter struct {
  mu    sync.Mutex
  burst int             // 最大令牌数（最大并发数）
  sem   chan struct{}   // 信号量通道
}
```

**特点**：
- 控制**同时处理的请求数**（并发数）
- 请求完成后需要调用 `Put()` 归还令牌
- 适合控制资源占用，如数据库连接数、goroutine 数量
- 基于 Channel 实现，原生支持超时和 Context 取消

### 2.2 QPS 限流器 (QPSLimiter)

基于**令牌桶算法**的速率控制器：

```go
type QPSLimiter struct {
  qps        float64   // 每秒生成的令牌数
  burst      int       // 桶容量（允许的突发流量）
  tokens     float64   // 当前可用令牌
  lastUpdate time.Time // 上次更新时间
}
```

**特点**：
- 控制**每秒允许的请求数**（速率）
- 令牌按时间自动恢复，无需归还
- 适合流量控制和防止突发流量

### 2.3 两种模式对比

| 特性 | 并发限流器 (Limiter) | QPS 限流器 (QPSLimiter) |
|------|---------------------|------------------------|
| 控制目标 | 同时处理的请求数 | 每秒允许的请求数 |
| 令牌归还 | 需要 Put() | 不需要（时间自动恢复） |
| 适用场景 | 资源保护 | 流量控制 |
| 限制维度 | 资源维度 | 时间维度 |
| 实现方式 | Channel 信号量 | 令牌桶算法 |
| Context 支持 | 原生支持 | 原生支持 |

### 2.4 组合使用

**当 QPS 和并发数同时配置时，两种限流同时生效，取更严格的限制**：

```
请求进入
    ↓
┌─────────────────────┐
│  QPS限流检查        │ ← 令牌桶算法，按时间恢复
│  (default_qps=1000) │
└─────────┬───────────┘
          │ 通过
          ↓
┌─────────────────────┐
│  并发数检查         │ ← 信号量算法，请求完成归还
│  (max_concurrency=100) │
└─────────┬───────────┘
          │ 通过
          ↓
      处理请求
          ↓
      归还并发令牌
```

**实际吞吐量 = min(QPS限制, 并发数 × (1000/平均处理时间ms))**

## 3. 架构设计

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                         配置层 (YAML)                                │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │ web:                                                           │ │
│  │   qps_limit:                                                   │ │
│  │     grpc:                                                      │ │
│  │       default_qps: 1000                                        │ │
│  │       max_concurrency: 100                                     │ │
│  │     http:                                                      │ │
│  │       default_qps: 1000                                        │ │
│  │       max_concurrency: 100                                     │ │
│  └───────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       中间件层 (Interceptor)                         │
│  ┌─────────────────────────┐    ┌─────────────────────────────────┐│
│  │   gRPC 拦截器            │    │   HTTP 中间件                    ││
│  │   - QPS 限流拦截器       │    │   - QPS 限流中间件               ││
│  │   - 并发控制拦截器       │    │   - 并发控制中间件               ││
│  └─────────────────────────┘    └─────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       核心限流器层                                   │
│  ┌─────────────────────────┐    ┌─────────────────────────────────┐│
│  │   MethodQPSLimiter      │    │   MethodLimiter                 ││
│  │   (方法级 QPS 限流)      │    │   (方法级并发控制)               ││
│  └─────────────────────────┘    └─────────────────────────────────┘│
│                 │                              │                    │
│                 ▼                              ▼                    │
│  ┌─────────────────────────┐    ┌─────────────────────────────────┐│
│  │   QPSLimiter            │    │   Limiter                       ││
│  │   (令牌桶算法)           │    │   (信号量算法)                   ││
│  └─────────────────────────┘    └─────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────┘
```

### 3.2 方法级限流

支持为不同的 API 接口配置不同的限流策略：

```
┌─────────────────────────────────────────────────────────────┐
│                   MethodQPSLimiter                          │
│  ┌───────────────────────────────────────────────────────┐ │
│  │                   Method -> Limiter Map               │ │
│  │  ┌─────────────────────────────────────────────────┐  │ │
│  │  │ /api/v1/users    -> QPSLimiter(qps=100, burst=20)│ │ │
│  │  │ /api/v1/orders   -> QPSLimiter(qps=50,  burst=10)│ │ │
│  │  │ /api/v1/products -> QPSLimiter(qps=200, burst=50)│ │ │
│  │  │ *                -> QPSLimiter(qps=10,  burst=5) │ │ │
│  │  └─────────────────────────────────────────────────┘  │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 4. 配置设计

### 4.1 配置结构

限流配置统一放在 `web.qps_limit` 下，支持 gRPC 和 HTTP 分别配置：

```yaml
web:
  bind_address:
    port: 10001
  grpc:
    timeout: 0s
  http:
    api_formatter: trivial_api_v20
  # QPS限流配置（统一在 web 配置内）
  qps_limit:
    # gRPC QPS限流配置
    grpc:
      default_qps: 1000       # 默认QPS，0表示不限制
      default_burst: 1500     # 默认突发容量，0表示使用default_qps值
      max_concurrency: 100    # 最大并发数限制，0表示不限制
      method_qps:             # 方法级配置（可选）
      - method: "/seadate.v1.SeaDateService/Now"
        qps: 500
        burst: 750
        max_concurrency: 50
    # HTTP QPS限流配置
    http:
      default_qps: 1000
      default_burst: 1500
      max_concurrency: 100
      method_qps:             # 路径级配置（可选）
      - method: "/v1/now"
        qps: 500
        burst: 750
        max_concurrency: 50
```

### 4.2 配置优先级

```
方法级配置 (method_qps) > 全局默认配置 (default_*)
```

| 请求路径 | method_qps 配置 | 实际使用 |
|---------|----------------|---------|
| `/v1/now` | 有配置 `qps:500` | 使用 500 QPS |
| `/v1/users` | 无配置 | 使用 `default_qps: 1000` |

### 4.3 配置参数说明

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `default_qps` | float64 | 默认每秒请求数限制，0表示不限制 | 0 |
| `default_burst` | int | 默认突发容量，0表示使用default_qps值 | 0 |
| `max_concurrency` | int | 最大并发数限制，0表示不限制 | 0 |
| `method_qps` | array | 方法级配置列表 | [] |

## 5. 核心实现

### 5.1 Limiter (并发控制)

基于 **Channel 信号量**实现，移除了对自定义 `sync.Cond` 的依赖：

```go
type Limiter struct {
  mu    sync.Mutex
  burst int             // 最大并发数
  sem   chan struct{}   // 信号量通道
}

// 创建并发限流器
func NewLimiter(b int) *Limiter

// 核心方法
func (lim *Limiter) Allow() bool                        // 立即获取，不等待
func (lim *Limiter) AllowN(n int) bool                  // 获取 n 个令牌
func (lim *Limiter) AllowFor(timeout time.Duration) bool // 带超时等待
func (lim *Limiter) AllowContext(ctx context.Context) error // Context 支持
func (lim *Limiter) Put()                               // 归还令牌
func (lim *Limiter) PutN(n int)                         // 归还 n 个令牌

// 等待方法
func (lim *Limiter) WaitFor(timeout time.Duration) error
func (lim *Limiter) WaitN(timeout time.Duration, n int) error
func (lim *Limiter) WaitContext(ctx context.Context) error
func (lim *Limiter) WaitNContext(ctx context.Context, n int) error

// 查询和动态调整
func (lim *Limiter) Burst() int     // 获取最大并发数
func (lim *Limiter) Tokens() int    // 获取当前可用令牌
func (lim *Limiter) Bursting() int  // 获取当前正在使用的令牌数
func (lim *Limiter) SetBurst(int)   // 动态调整最大并发数
```

**实现原理**：

```
┌─────────────────────────────────────────────────────────────┐
│                    Limiter (Channel 信号量)                  │
│  ┌───────────────────────────────────────────────────────┐ │
│  │   sem: chan struct{} (buffered channel)               │ │
│  │   ┌─────────────────────────────────────────────────┐ │ │
│  │   │  [token] [token] [token] ... [token]            │ │ │
│  │   │  ← burst 个预填充令牌                            │ │ │
│  │   └─────────────────────────────────────────────────┘ │ │
│  │                                                       │ │
│  │   Allow():  select { case <-sem: return true }        │ │
│  │   Put():    select { case sem <- struct{}{}: }        │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

**相比旧实现的改进**：

| 改进点 | 旧实现 | 新实现 |
|-------|-------|-------|
| 依赖 | 自定义 `sync_.Cond` | 标准库 `channel` |
| 超时机制 | `cond.WaitForDo` | `select` + `time.After` |
| Context 支持 | 无 | 原生支持 |
| 动态调整 | 无 | `SetBurst()` |
| 代码复杂度 | 较高 | 简洁 |

### 5.2 QPSLimiter

基于**令牌桶算法**实现：

```go
type QPSLimiter struct {
  mu         sync.Mutex
  qps        float64   // 每秒生成的令牌数
  burst      int       // 桶容量（最大令牌数）
  tokens     float64   // 当前可用令牌
  lastUpdate time.Time // 上次更新时间
}

// 创建 QPS 限流器
func NewQPSLimiter(qps float64, burst int) *QPSLimiter

// 核心方法
func (l *QPSLimiter) Allow() bool                    // 立即判断是否允许
func (l *QPSLimiter) AllowN(n int) bool              // 请求 n 个令牌
func (l *QPSLimiter) Wait(ctx context.Context) error // 阻塞等待令牌
func (l *QPSLimiter) AllowFor(timeout time.Duration) bool // 带超时等待
```

### 5.3 MethodQPSLimiter

```go
type MethodQPSLimiter struct {
  mu       sync.RWMutex
  limiters map[string]*QPSLimiter // 方法 -> 限流器映射
  global   *QPSLimiter            // 默认全局限流器
}

// 方法 QPS 配置
type MethodQPSConfig struct {
  Method string  // API 方法名
  QPS    float64 // 每秒请求数
  Burst  int     // 突发容量
}

// 创建方法级限流器
func NewMethodQPSLimiter(defaultQPS float64, defaultBurst int) *MethodQPSLimiter
func NewMethodQPSLimiterWithConfigs(defaultQPS float64, defaultBurst int, 
    configs []MethodQPSConfig) (*MethodQPSLimiter, error)

// 动态管理方法限流
func (m *MethodQPSLimiter) AddMethod(method string, qps float64, burst int) error
func (m *MethodQPSLimiter) SetMethodQPS(method string, qps float64, burst int) error
func (m *MethodQPSLimiter) RemoveMethod(method string)

// 限流判断（自动选择对应的限流器）
func (m *MethodQPSLimiter) Allow(method string) bool
func (m *MethodQPSLimiter) AllowFor(method string, timeout time.Duration) bool
```

### 5.4 配置转换

```go
// QPSLimitConfig 限流配置结构
type QPSLimitConfig struct {
  DefaultQPS     float64
  DefaultBurst   int
  MaxConcurrency int
  MethodQPS      []MethodQPSConfigItem
}

// 转换为gRPC网关QPS限流配置
func (c *QPSLimitConfig) ToGRPCQPSLimitConfig() gw_.QPSLimitConfig

// 转换为HTTP QPS限流配置
func (c *QPSLimitConfig) ToHTTPQPSLimitConfig() gw_.HTTPQPSLimitConfig
```

## 6. 中间件集成

### 6.1 gRPC 拦截器

```go
// 创建方法级 QPS 限流器
limiter := ratelimit.NewMethodQPSLimiter(100, 20) // 默认 100 QPS
limiter.AddMethod("/service.UserService/GetUser", 500, 100)
limiter.AddMethod("/service.OrderService/CreateOrder", 50, 10)

// 应用拦截器
server := grpc.NewServer(
  grpc.UnaryInterceptor(ratelimit.UnaryServerInterceptorQPS(limiter)),
  grpc.StreamInterceptor(ratelimit.StreamServerInterceptorQPS(limiter)),
)
```

### 6.2 HTTP 中间件

```go
// 创建 QPS 限流器
limiter := ratelimit.NewQPSRateLimiter(100, 20) // 默认 100 QPS
limiter.AddPath("/api/v1/users", 500, 100)
limiter.AddPath("/api/v1/orders", 50, 10)

// 应用中间件
http.Handle("/", limiter.Handler(yourHandler))

// 暴露限流统计
http.Handle("/debug/ratelimit", limiter.StatsHandler())
```

### 6.3 通过配置自动集成

在 `webserver.Config` 中自动安装限流中间件：

```go
// installHttpMiddlewareChain 安装 HTTP 中间件链
func (c *Config) installHttpMiddlewareChain() []gw_.GRPCGatewayOption {
  // ... 其他中间件

  // QPS限流和并发控制（通过扩展配置）
  if c.opts.httpQPSLimit != nil {
    opts = append(opts, gw_.WithHttpHandlerInterceptorsQPSLimitOptions(
      c.opts.httpQPSLimit.ToHTTPQPSLimitConfig(),
    ))
  }
  
  return opts
}

// installGrpcMiddlewareChain 安装 gRPC 中间件链
func (c *Config) installGrpcMiddlewareChain() []gw_.GRPCGatewayOption {
  // ... 其他中间件

  // QPS限流和并发控制（通过扩展配置）
  if c.opts.grpcQPSLimit != nil {
    opts = append(opts, gw_.WithServerInterceptorsQPSLimitOptions(
      c.opts.grpcQPSLimit.ToGRPCQPSLimitConfig(),
    ))
  }

  return opts
}
```

## 7. 使用示例

### 7.1 并发限流器基本使用

```go
// 创建并发限流器：最大 100 并发
limiter := rate.NewLimiter(100)

// 立即获取令牌
if limiter.Allow() {
  defer limiter.Put()
  // 处理请求
}

// 带超时等待
if limiter.AllowFor(100 * time.Millisecond) {
  defer limiter.Put()
  // 处理请求
}

// 使用 Context（推荐）
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
if err := limiter.AllowContext(ctx); err == nil {
  defer limiter.Put()
  // 处理请求
}

// 动态调整最大并发数
limiter.SetBurst(200)
```

### 7.2 QPS 限流器基本使用

```go
// 创建 QPS 限流器：100 QPS，允许突发 20 个请求
limiter := rate.NewQPSLimiter(100, 20)

// 判断是否允许
if limiter.Allow() {
  // 处理请求
}

// 带超时等待
if limiter.AllowFor(100 * time.Millisecond) {
  // 处理请求
}

// 阻塞等待
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
if err := limiter.Wait(ctx); err == nil {
  // 处理请求
}
```

### 7.3 方法级限流

```go
// 配置不同接口的 QPS
configs := []rate.MethodQPSConfig{
  {Method: "/api/v1/users", QPS: 100, Burst: 20},
  {Method: "/api/v1/orders", QPS: 50, Burst: 10},
  {Method: "/api/v1/products", QPS: 200, Burst: 50},
}

// 创建限流器，默认 10 QPS
limiter, _ := rate.NewMethodQPSLimiterWithConfigs(10, 5, configs)

// 使用（自动选择对应的限流器）
limiter.Allow("/api/v1/users")    // 使用 100 QPS
limiter.Allow("/api/v1/orders")   // 使用 50 QPS
limiter.Allow("/api/v1/unknown")  // 使用默认 10 QPS
```

### 7.4 动态调整

```go
limiter := rate.NewMethodQPSLimiter(100, 20)

// 动态添加方法限流
limiter.AddMethod("/api/v1/hot-endpoint", 500, 100)

// 动态调整 QPS
limiter.SetMethodQPS("/api/v1/hot-endpoint", 1000, 200)

// 移除方法限流（使用全局限流）
limiter.RemoveMethod("/api/v1/hot-endpoint")

// 调整全局 QPS
limiter.SetGlobalQPS(200, 50)
```

### 7.5 查看统计

```go
limiter := rate.NewMethodQPSLimiter(100, 20)
limiter.AddMethod("/api/v1/users", 500, 100)

// 获取统计信息
stats := limiter.Stats()
for _, s := range stats {
  fmt.Printf("Method: %s, QPS: %.0f, Burst: %d, Available: %.1f\n",
    s.Method, s.QPS, s.Burst, s.Tokens)
}
// Output:
// Method: *, QPS: 100, Burst: 20, Available: 20.0
// Method: /api/v1/users, QPS: 500, Burst: 100, Available: 100.0
```

### 7.6 压测验证

```bash
# 配置 max_concurrency: 10 后进行压测
bombardier -n 30000 -c 3000 -t 10s -H Content-type:application/json -m POST http://127.0.0.1:10001/Now

# 预期结果：大量请求被限流拒绝，返回 429
# HTTP codes:
#   1xx - 0, 2xx - 2047, 3xx - 0, 4xx - 27869, 5xx - 0
# 4xx 表示被限流拒绝的请求（HTTP 429 Too Many Requests）
```

## 8. API 响应格式

### 8.1 HTTP 限流响应

**QPS 限流**：
```json
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 1

{
  "error": "rate_limit_exceeded",
  "message": "GET /api/v1/users is rejected by http_qps middleware, QPS limit exceeded",
  "code": 429
}
```

**并发数限流**：
```json
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
  "error": "concurrency_limit_exceeded",
  "message": "POST /Now is rejected by http_concurrency middleware, max concurrency exceeded",
  "code": 429
}
```

### 8.2 gRPC 限流响应

```
code: RESOURCE_EXHAUSTED
message: "/service.UserService/GetUser is rejected, QPS limit exceeded"
```

## 9. 文件结构

```
go/time/rate/
├── rate.go              # 并发控制限流器（信号量）
├── rate_method.go       # 方法级并发限流器
├── rate_qps.go          # QPS 限流器（令牌桶）
├── rate_qps_method.go   # 方法级 QPS 限流器
├── rate_qps_test.go     # QPS 限流器测试
└── rate_test.go         # 并发限流器测试

pkg/webserver/
├── config.go            # Web 服务配置（安装限流中间件）
├── config.option.go     # 配置选项
├── webserver.proto      # Proto 定义
├── webserver.yaml       # 示例配置
└── webserver_qps_limit.go # QPS 限流配置转换

pkg/grpc-gateway/
├── grpc_gateway_grpc.option.go  # gRPC QPS 限流选项
└── grpc_gateway_http.option.go  # HTTP QPS 限流选项

pkg/middleware/grpc-middleware/ratelimit/
├── ratelimit_server.interceptor.go      # 并发限流拦截器
└── ratelimit_qps_server.interceptor.go  # QPS 限流拦截器

pkg/middleware/http-middleware/ratelimiter/
├── ratelimiter.go       # 并发限流中间件
└── ratelimiter_qps.go   # QPS 限流中间件
```

## 10. 设计特点总结

| 特点 | 说明 |
|------|------|
| 双重限流 | 同时支持 QPS 限流（令牌桶）和并发控制（Channel 信号量） |
| 配置统一 | 限流配置统一放在 `web.qps_limit` 下，结构清晰 |
| 方法级配置 | 支持为不同 API 接口配置不同的限流策略 |
| 动态调整 | 运行时动态添加、修改、删除方法限流配置，支持 `SetBurst()` |
| 等待模式 | 支持立即返回、带超时等待、阻塞等待三种模式 |
| Context 支持 | 原生支持 `context.Context`，符合 Go 惯例 |
| 统计信息 | 提供限流统计接口，便于监控和调试 |
| 线程安全 | 使用 Channel 和互斥锁保证并发安全 |
| 无外部依赖 | 仅依赖标准库，移除了对自定义 `sync.Cond` 的依赖 |
| 向后兼容 | 旧的 `max_concurrency_unary` 等配置已废弃，统一使用新配置 |

## 11. 废弃配置说明

以下配置已废弃，统一使用 `web.qps_limit` 配置：

| 废弃配置 | 新配置 |
|---------|-------|
| `web.grpc.max_concurrency_unary` | `web.qps_limit.grpc.max_concurrency` |
| `web.grpc.max_concurrency_stream` | `web.qps_limit.grpc.max_concurrency` |
| `web.http.max_concurrency` | `web.qps_limit.http.max_concurrency` |

## 12. 实现演进

### 12.1 Limiter 重构（v2.0）

**变更原因**：
- 旧实现依赖自定义 `sync_.Cond`，增加了维护成本
- 自定义 Cond 的超时实现复杂，不如 Channel 原生支持

**变更内容**：

| 项目 | 旧实现 | 新实现 |
|------|-------|-------|
| 核心数据结构 | `tokens int` + `sync_.Cond` | `sem chan struct{}` |
| 超时等待 | `cond.WaitForDo(timeout, pred, do)` | `select { case <-sem: case <-time.After(timeout): }` |
| Context 支持 | 无 | `AllowContext()`, `WaitContext()` |
| 动态调整 | 无 | `SetBurst()` |
| 依赖 | `github.com/kaydxh/golang/go/sync` | 仅标准库 |

**API 兼容性**：
- ✅ `Allow()`, `AllowFor()`, `AllowWaitUntil()` - 签名不变
- ✅ `Put()`, `PutN()` - 签名不变
- ✅ `WaitFor()`, `WaitN()` - 签名不变
- ✅ `Burst()`, `Bursting()` - 签名不变
- ➕ 新增 `Tokens()`, `SetBurst()`, `AllowN()`, `AllowContext()`, `WaitContext()`, `WaitNContext()`
