# 健康检查增强功能设计文档

## 1. 概述

本文档描述了 golang 库中健康检查（Health Check）功能的设计思路和实现细节。该功能提供了 Kubernetes 兼容的健康检查端点，支持存活探针（Liveness）和就绪探针（Readiness），并具备良好的可扩展性。

## 2. 设计目标

- **Kubernetes 兼容**：支持标准的 `/healthz`、`/livez`、`/readyz` 端点
- **可扩展性**：通过接口抽象支持自定义健康检查器
- **组合能力**：支持聚合多个检查器，统一管理
- **优雅关闭**：关闭时先标记不就绪，允许负载均衡器排空连接
- **可观测性**：提供详细模式，返回每个检查器的状态和延迟

## 3. 架构设计

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      WebServer                              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              HealthzController                         │  │
│  │  ┌─────────────────┐    ┌─────────────────┐           │  │
│  │  │ livezCheckers   │    │ readyzCheckers  │           │  │
│  │  │ (Composite)     │    │ (Composite)     │           │  │
│  │  │  ┌───────────┐  │    │  ┌───────────┐  │           │  │
│  │  │  │   Ping    │  │    │  │   HTTP    │  │           │  │
│  │  │  ├───────────┤  │    │  ├───────────┤  │           │  │
│  │  │  │   TCP     │  │    │  │   TCP     │  │           │  │
│  │  │  ├───────────┤  │    │  ├───────────┤  │           │  │
│  │  │  │  Custom   │  │    │  │  Custom   │  │           │  │
│  │  │  └───────────┘  │    │  └───────────┘  │           │  │
│  │  └─────────────────┘    └─────────────────┘           │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  HTTP Endpoints:                                            │
│  - GET /healthz          综合健康检查                       │
│  - GET /livez            存活探针                           │
│  - GET /readyz           就绪探针                           │
│  - GET /healthz/verbose  详细健康检查结果                   │
│  - GET /livez/verbose    详细存活检查结果                   │
│  - GET /readyz/verbose   详细就绪检查结果                   │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 核心组件

| 组件 | 职责 |
|------|------|
| `HealthChecker` | 健康检查器接口，定义检查行为 |
| `CompositeHealthChecker` | 组合检查器，聚合多个检查器 |
| `Controller` | 健康检查控制器，管理端点和检查逻辑 |
| `WebServer` | 集成健康检查控制器，处理生命周期 |

## 4. 核心接口设计

### 4.1 HealthChecker 接口

```go
// HealthChecker 是健康检查的核心接口
type HealthChecker interface {
    // Name 返回检查器名称，用于日志和详细结果展示
    Name() string
    // Check 执行健康检查，返回 nil 表示健康，返回 error 表示不健康
    Check(ctx context.Context) error
}
```

### 4.2 内置检查器实现

| 检查器类型 | 用途 | 实现要点 |
|-----------|------|---------|
| `PingHealthChecker` | 基础检查，始终返回健康 | 用于默认存活检查 |
| `HTTPHealthChecker` | HTTP 端点健康检查 | 检查状态码 200-299 |
| `TCPHealthChecker` | TCP 端口连通性检查 | 使用 `net.Dialer` |
| `FuncHealthChecker` | 函数包装器 | 支持自定义检查逻辑 |
| `CompositeHealthChecker` | 组合检查器 | 聚合多个检查器 |

## 5. 关键实现

### 5.1 HTTP 健康检查器

```go
type HTTPHealthChecker struct {
    name    string
    url     string
    timeout time.Duration
    client  *http.Client
}

func (h *HTTPHealthChecker) Check(ctx context.Context) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    resp, err := h.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to connect: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("unhealthy status code: %d", resp.StatusCode)
    }
    return nil
}
```

### 5.2 TCP 健康检查器

```go
type TCPHealthChecker struct {
    name    string
    addr    string
    timeout time.Duration
}

func (t *TCPHealthChecker) Check(ctx context.Context) error {
    dialer := &net.Dialer{Timeout: t.timeout}
    conn, err := dialer.DialContext(ctx, "tcp", t.addr)
    if err != nil {
        return fmt.Errorf("failed to connect to %s: %w", t.addr, err)
    }
    conn.Close()
    return nil
}
```

### 5.3 组合健康检查器

```go
type CompositeHealthChecker struct {
    name     string
    mu       sync.RWMutex
    checkers []HealthChecker
}

func (c *CompositeHealthChecker) Check(ctx context.Context) error {
    c.mu.RLock()
    defer c.mu.RUnlock()

    for _, checker := range c.checkers {
        if err := checker.Check(ctx); err != nil {
            return fmt.Errorf("%s: %w", checker.Name(), err)
        }
    }
    return nil
}

// AddChecker 线程安全地添加检查器
func (c *CompositeHealthChecker) AddChecker(checker HealthChecker) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.checkers = append(c.checkers, checker)
}
```

### 5.4 健康检查控制器

```go
type Controller struct {
    livezCheckers  *CompositeHealthChecker  // 存活检查器
    readyzCheckers *CompositeHealthChecker  // 就绪检查器
    ready          atomic.Bool               // 就绪状态
    checkTimeout   time.Duration             // 检查超时（默认 10s）
}

// Healthz 综合健康检查端点
func (c *Controller) Healthz() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), c.checkTimeout)
        defer cancel()

        // 检查存活状态
        if err := c.livezCheckers.Check(checkCtx); err != nil {
            ctx.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "unhealthy", "type": "livez", "error": err.Error(),
            })
            return
        }

        // 检查就绪状态
        if !c.ready.Load() {
            ctx.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "not ready", "type": "readyz", "error": "server is shutting down",
            })
            return
        }

        if err := c.readyzCheckers.Check(checkCtx); err != nil {
            ctx.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "not ready", "type": "readyz", "error": err.Error(),
            })
            return
        }

        ctx.String(http.StatusOK, "ok")
    }
}
```

## 6. 优雅关闭设计

### 6.1 关闭流程

```
收到关闭信号 (ctx.Done())
       │
       ▼
标记服务为不就绪 (SetReady(false))
       │
       ▼
/readyz 开始返回 503
       │
       ▼
等待 ShutdownDelayDuration (允许 LB 排空连接)
       │
       ▼
关闭 HTTP/gRPC 服务器
```

### 6.2 实现代码

```go
func (s preparedGenericWebServer) Run(ctx context.Context) error {
    // ... 启动代码 ...

    <-ctx.Done()  // 等待关闭信号

    // 1. 标记服务器为不就绪
    if s.HealthzController != nil {
        s.HealthzController.SetReady(false)
    }

    // 2. 等待关闭延迟，允许负载均衡器排空连接
    if s.ShutdownDelayDuration > 0 {
        time.Sleep(s.ShutdownDelayDuration)
    }

    // 3. 关闭服务器
    // ...
}
```

## 7. 与 Consul 服务发现集成

```go
func (srv *ServiceRegistryServer) Register() error {
    checkUrl := fmt.Sprintf("http://%v:%v/api/%v/v1/health", 
        srv.Ip, srv.Port, srv.ServiceName)

    reg := &api.AgentServiceRegistration{
        ID:   srv.ServiceId,
        Name: srv.ServiceName,
        Check: &api.AgentServiceCheck{
            Interval:                       srv.CheckInterval.String(),
            HTTP:                           checkUrl,
            DeregisterCriticalServiceAfter: srv.TTL.String(),
        },
    }
    return agent.ServiceRegister(reg)
}
```

## 8. 使用示例

### 8.1 添加自定义检查器

```go
// 添加数据库健康检查
dbChecker := healthz.NewTCPHealthChecker("mysql", "localhost:3306", 5*time.Second)
server.HealthzController.AddReadyzChecker(dbChecker)

// 添加 Redis 健康检查
redisChecker := healthz.NewTCPHealthChecker("redis", "localhost:6379", 5*time.Second)
server.HealthzController.AddReadyzChecker(redisChecker)

// 添加自定义检查逻辑
customChecker := healthz.NewFuncHealthChecker("custom", func(ctx context.Context) error {
    // 自定义检查逻辑
    if someCondition {
        return errors.New("unhealthy")
    }
    return nil
})
server.HealthzController.AddLivezChecker(customChecker)
```

### 8.2 Kubernetes 部署配置

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    livenessProbe:
      httpGet:
        path: /livez
        port: 8080
      initialDelaySeconds: 10
      periodSeconds: 10
    readinessProbe:
      httpGet:
        path: /readyz
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5
```

## 9. API 响应格式

### 9.1 健康状态正常

```
GET /healthz
Response: 200 OK
Body: ok
```

### 9.2 健康状态异常

```
GET /healthz
Response: 503 Service Unavailable
Body: {"status": "unhealthy", "type": "livez", "error": "mysql: connection refused"}
```

### 9.3 详细检查结果

```
GET /healthz/verbose
Response: 200 OK
Body:
{
  "status": "ok",
  "livez": {
    "status": "ok",
    "checks": [
      {"name": "ping", "status": "ok", "latency": "0.1ms"}
    ]
  },
  "readyz": {
    "status": "ok",
    "checks": [
      {"name": "mysql", "status": "ok", "latency": "2.3ms"},
      {"name": "redis", "status": "ok", "latency": "1.1ms"}
    ]
  }
}
```

## 10. 设计特点总结

| 特点 | 说明 |
|------|------|
| Kubernetes 兼容 | 支持标准的存活/就绪探针端点 |
| 可扩展性 | 通过接口支持自定义检查器 |
| 组合模式 | 支持聚合多个检查器 |
| 详细模式 | 提供每个检查器的详细结果和延迟 |
| 优雅关闭 | 关闭时先标记不就绪 |
| 超时控制 | 默认 10 秒检查超时，可配置 |
| 线程安全 | 使用互斥锁和原子操作保证并发安全 |
| 服务发现集成 | 支持 Consul HTTP 健康检查 |

## 11. 文件结构

```
pkg/webserver/controller/healthz/
├── checker.go      # 健康检查器接口和实现
├── healthz.go      # 健康检查控制器
└── healthz_test.go # 单元测试

pkg/webserver/
├── webserver.go    # WebServer 集成健康检查
└── config.go       # 配置和初始化

pkg/discovery/consul/
└── discovery.go    # Consul 服务注册健康检查
```
