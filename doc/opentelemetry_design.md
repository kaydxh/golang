# OpenTelemetry 监控设计文档

## 1. 概述

本文档描述了 golang 库中 OpenTelemetry 监控功能的设计思路和实现细节。该功能提供统一的可观测性能力，支持 **Metrics（指标）**、**Traces（链路追踪）** 和 **Logs（日志）**，并支持多种后端导出器（Prometheus、Jaeger、OTLP 等）。

## 2. 设计目标

- **统一的可观测性**：通过 OpenTelemetry 标准实现 Metrics、Traces、Logs 的统一采集
- **多后端支持**：支持 Prometheus、Jaeger、OTLP、Stdout 等多种导出器
- **腾讯云兼容**：支持腾讯云 Prometheus Remote Write 接入
- **配置驱动**：通过 YAML 配置文件灵活配置各种导出器
- **低侵入性**：通过中间件自动采集 gRPC/HTTP 请求指标和链路追踪

## 3. 架构设计

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           应用层                                         │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                    gRPC/HTTP Server                                │ │
│  │  ┌─────────────────┐    ┌─────────────────┐    ┌───────────────┐  │ │
│  │  │ Metric 中间件    │    │ Trace 中间件    │    │ 业务代码       │  │ │
│  │  └────────┬────────┘    └────────┬────────┘    └───────────────┘  │ │
│  └───────────│──────────────────────│────────────────────────────────┘ │
└──────────────│──────────────────────│──────────────────────────────────┘
               │                      │
               ▼                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      OpenTelemetry SDK 层                                │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                 OpenTelemetryService                              │ │
│  │  ┌─────────────────┐    ┌─────────────────┐                       │ │
│  │  │   MeterProvider │    │  TracerProvider │                       │ │
│  │  │  ┌───────────┐  │    │  ┌───────────┐  │                       │ │
│  │  │  │  Reader   │  │    │  │  Exporter │  │                       │ │
│  │  │  └───────────┘  │    │  └───────────┘  │                       │ │
│  │  └─────────────────┘    └─────────────────┘                       │ │
│  └───────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
               │                      │
               ▼                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         导出器层 (Exporter)                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│  │ Prometheus  │  │    OTLP     │  │   Jaeger    │  │   Stdout    │   │
│  │   (Pull)    │  │   (Push)    │  │   (Push)    │  │   (Push)    │   │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘   │
└─────────│────────────────│────────────────│────────────────│──────────┘
          │                │                │                │
          ▼                ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         后端存储层                                       │
│  ┌─────────────┐  ┌─────────────────────┐  ┌─────────────┐             │
│  │ Prometheus  │  │ 腾讯云 Prometheus    │  │   Jaeger    │             │
│  │   Server    │  │   Remote Write      │  │   Server    │             │
│  └─────────────┘  └─────────────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.2 核心组件

| 组件 | 职责 |
|------|------|
| `OpenTelemetryService` | 统一管理 Meter 和 Tracer 的生命周期 |
| `MeterProvider` | 管理指标的采集和导出 |
| `TracerProvider` | 管理链路追踪的采集和导出 |
| `PullExporter` | Pull 模式导出器（如 Prometheus） |
| `PushExporter` | Push 模式导出器（如 OTLP、Jaeger、Stdout） |
| `ResourceStatsService` | 资源统计服务（CPU、内存等） |

## 4. 指标采集模式

### 4.1 Pull 模式 vs Push 模式

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Pull 模式 (Prometheus)                          │
│                                                                         │
│   ┌─────────────┐                           ┌─────────────────────┐     │
│   │   服务      │ ──────── /metrics ───────→ │   Prometheus       │     │
│   │ (暴露端口)  │ ←─────── 主动抓取 ────────│   Server           │     │
│   └─────────────┘                           └─────────────────────┘     │
│                                                                         │
│   特点：服务暴露端口，Prometheus 主动拉取                                │
│   适用：同 VPC、K8s 内部、ServiceMonitor                                │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                         Push 模式 (OTLP)                                │
│                                                                         │
│   ┌─────────────┐                           ┌─────────────────────┐     │
│   │   服务      │ ──────── 主动推送 ───────→ │   腾讯云 Prometheus │     │
│   │ (定时推送)  │                           │   Remote Write      │     │
│   └─────────────┘                           └─────────────────────┘     │
│                                                                         │
│   特点：服务主动推送，无需暴露端口                                       │
│   适用：跨 VPC、公网、防火墙限制场景                                     │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 模式对比

| 特性 | Pull 模式 (Prometheus) | Push 模式 (OTLP) |
|------|------------------------|------------------|
| **数据流向** | 后端主动拉取服务 | 服务主动推送到后端 |
| **端口暴露** | 需要暴露 `/metrics` 端点 | 无需暴露端口 |
| **网络要求** | Prometheus 需能访问服务 | 服务需能访问采集端 |
| **适用场景** | 同 VPC、K8s ServiceMonitor | 跨 VPC、无法暴露端口 |
| **配置类型** | `metric_prometheus` | `metric_otlp` |
| **Exporter** | `WithMeterPullExporter` | `WithMeterPushExporter` |

### 4.3 选择建议

| 场景 | 推荐模式 | 原因 |
|------|---------|------|
| K8s 集群内部 | Pull (Prometheus) | Prometheus Operator 原生支持 ServiceMonitor |
| 同 VPC 部署 | Pull (Prometheus) | 网络直连，配置简单 |
| 跨 VPC / 公网 | Push (OTLP) | 无需暴露端口，穿透防火墙 |
| 边缘节点 | Push (OTLP) | 边缘网络不稳定，主动推送更可靠 |
| 腾讯云托管 Prometheus | 两者皆可 | 根据网络拓扑选择 |

## 5. 配置设计

### 5.1 配置结构 (Proto 定义)

```protobuf
message OpenTelemetry {
  bool enabled = 1;
  google.protobuf.Duration metric_collect_duration = 2;  // 指标采集周期
  
  OtelTraceExporterType otel_trace_exporter_type = 3;    // 链路导出器类型
  OtelMetricExporterType otel_metric_exporter_type = 4;  // 指标导出器类型
  OtelLogExporterType otel_log_exporter_type = 5;        // 日志导出器类型
  
  OtelMetricExporter otel_metric_exporter = 6;           // 指标导出器配置
  OtelTraceExporter otel_trace_exporter = 7;             // 链路导出器配置
  
  Resource resource = 8;                                  // 资源标识
}

// 指标导出器配置
message OtelMetricExporter {
  Prometheus prometheus = 1;  // Prometheus Pull 模式
  Stdout stdout = 2;          // 标准输出
  OTLP otlp = 3;              // OTLP Push 模式
}

// OTLP 导出器配置
message OTLP {
  string endpoint = 1;              // 目标地址
  string token = 2;                 // 认证 Token
  string protocol = 3;              // 协议：http 或 grpc
  bool insecure = 4;                // 是否禁用 TLS
  map<string, string> headers = 5;  // 自定义请求头
  string url_path = 6;              // HTTP URL 路径
}
```

### 5.2 导出器类型枚举

```protobuf
enum OtelMetricExporterType {
  metric_none = 0;        // 不启用
  metric_prometheus = 1;  // Prometheus Pull 模式
  metric_stdout = 2;      // 标准输出
  metric_otlp = 3;        // OTLP Push 模式
}

enum OtelTraceExporterType {
  trace_none = 0;         // 不启用
  trace_jaeger = 1;       // Jaeger
  trace_stdout = 2;       // 标准输出
  trace_otlp = 3;         // OTLP Push 模式
}
```

### 5.3 配置示例

#### Prometheus Pull 模式（同 VPC）

```yaml
open_telemetry:
  enabled: true
  metric_collect_duration: 60s
  otel_metric_exporter_type: metric_prometheus
  otel_trace_exporter_type: trace_jaeger
  otel_metric_exporter:
    prometheus:
      url: /metrics  # 暴露的端点路径
  otel_trace_exporter:
    jaeger:
      url: http://jaeger:14268/api/traces
  resource:
    service_name: "my-service"
    attrs:
      env: "production"
```

#### OTLP Push 模式（腾讯云 Prometheus）

```yaml
open_telemetry:
  enabled: true
  metric_collect_duration: 60s
  otel_metric_exporter_type: metric_otlp
  otel_trace_exporter_type: trace_stdout
  otel_metric_exporter:
    otlp:
      # 腾讯云 Prometheus Remote Write 地址（从控制台获取）
      endpoint: "your-instance.tencentcloudprom.com"
      # 认证 Token（从控制台获取）
      token: "your-prometheus-token"
      # 协议：http 或 grpc
      protocol: "http"
      # 内网访问可设为 true，公网访问设为 false
      insecure: false
      # OTLP 默认路径
      url_path: "/v1/metrics"
  resource:
    service_name: "my-service"
    attrs:
      env: "production"
      region: "ap-guangzhou"
```

## 6. 核心实现

### 6.1 OTLP Exporter 实现

```go
// OTLPExporterBuilder OTLP 导出器构建器
type OTLPExporterBuilder struct {
    opts struct {
        endpoint string              // 目标地址
        headers  map[string]string   // 请求头
        timeout  time.Duration       // 超时时间
        insecure bool                // 是否禁用 TLS
        protocol Protocol            // 协议类型
        urlPath  string              // URL 路径
    }
}

// Protocol 协议类型
type Protocol int

const (
    ProtocolHTTP Protocol = iota  // HTTP 协议
    ProtocolGRPC                  // gRPC 协议
)

// Build 构建 OTLP Exporter
func (b *OTLPExporterBuilder) Build(ctx context.Context) (metric.Exporter, error) {
    switch b.opts.protocol {
    case ProtocolHTTP:
        return b.buildHTTPExporter(ctx)
    case ProtocolGRPC:
        return b.buildGRPCExporter(ctx)
    default:
        return b.buildHTTPExporter(ctx)
    }
}

// buildHTTPExporter 构建 HTTP 导出器
func (b *OTLPExporterBuilder) buildHTTPExporter(ctx context.Context) (metric.Exporter, error) {
    opts := []otlpmetrichttp.Option{
        otlpmetrichttp.WithEndpoint(b.opts.endpoint),
    }
    
    if b.opts.insecure {
        opts = append(opts, otlpmetrichttp.WithInsecure())
    }
    if len(b.opts.headers) > 0 {
        opts = append(opts, otlpmetrichttp.WithHeaders(b.opts.headers))
    }
    if b.opts.timeout > 0 {
        opts = append(opts, otlpmetrichttp.WithTimeout(b.opts.timeout))
    }
    if b.opts.urlPath != "" {
        opts = append(opts, otlpmetrichttp.WithURLPath(b.opts.urlPath))
    }
    
    return otlpmetrichttp.New(ctx, opts...)
}
```

### 6.2 配置加载实现

```go
func (c *completedConfig) installMeter(ctx context.Context) ([]OpenTelemetryServiceOption, error) {
    var opts []OpenTelemetryServiceOption
    
    // 设置采集周期
    collectDuration := c.Proto.GetMetricCollectDuration().AsDuration()
    if collectDuration > 0 {
        opts = append(opts, WithMetricCollectDuration(collectDuration))
    }

    metricType := c.Proto.OtelMetricExporterType
    switch metricType {
    case OtelMetricExporterType_metric_prometheus:
        // Pull 模式：Prometheus 主动抓取
        builder := prometheus_.NewPrometheusExporterBuilder(
            prometheus_.WithMetricUrlPath(c.Proto.GetOtelMetricExporter().GetPrometheus().GetUrl()),
        )
        opts = append(opts, WithMeterPullExporter(builder))

    case OtelMetricExporterType_metric_otlp:
        // Push 模式：服务主动推送
        otlpConfig := c.Proto.GetOtelMetricExporter().GetOtlp()
        var otlpOpts []otlpmetric_.OTLPExporterBuilderOption
        
        // 配置 endpoint
        if otlpConfig.GetEndpoint() != "" {
            otlpOpts = append(otlpOpts, otlpmetric_.WithEndpoint(otlpConfig.GetEndpoint()))
        }
        
        // 配置协议
        if otlpConfig.GetProtocol() == "grpc" {
            otlpOpts = append(otlpOpts, otlpmetric_.WithProtocol(otlpmetric_.ProtocolGRPC))
        } else {
            otlpOpts = append(otlpOpts, otlpmetric_.WithProtocol(otlpmetric_.ProtocolHTTP))
        }
        
        // 配置认证头
        headers := make(map[string]string)
        for k, v := range otlpConfig.GetHeaders() {
            headers[k] = v
        }
        if otlpConfig.GetToken() != "" {
            headers["Authorization"] = "Bearer " + otlpConfig.GetToken()
        }
        if len(headers) > 0 {
            otlpOpts = append(otlpOpts, otlpmetric_.WithHeaders(headers))
        }
        
        builder := otlpmetric_.NewOTLPExporterBuilder(otlpOpts...)
        opts = append(opts, WithMeterPushExporter(builder))

    case OtelMetricExporterType_metric_stdout:
        // 标准输出（调试用）
        builder := stdoutmetric_.NewStdoutExporterBuilder(
            stdoutmetric_.WithPrettyPrint(c.Proto.GetOtelMetricExporter().GetStdout().GetPrettyPrint()),
        )
        opts = append(opts, WithMeterPushExporter(builder))
    }

    return opts, nil
}
```

### 6.3 中间件集成

#### gRPC 指标拦截器

```go
// UnaryServerInterceptorOfMetric gRPC 一元调用指标拦截器
func UnaryServerInterceptorOfMetric() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, 
        info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        
        start := time.Now()
        resp, err := handler(ctx, req)
        duration := time.Since(start)
        
        // 记录请求指标
        meter_.RecordRequest(ctx, info.FullMethod, duration, err)
        
        return resp, err
    }
}
```

#### HTTP 指标中间件

```go
// MetricMiddleware HTTP 指标中间件
func MetricMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        
        // 记录请求指标
        meter_.RecordHTTPRequest(c.Request.Context(), 
            c.Request.Method, c.Request.URL.Path, 
            c.Writer.Status(), duration)
    }
}
```

## 7. 腾讯云 Prometheus 接入指南

### 7.1 获取配置信息

1. 登录 [腾讯云 Prometheus 监控控制台](https://console.cloud.tencent.com/monitor/prometheus)
2. 选择对应实例，进入**基本信息** > **服务地址**
3. 获取以下信息：
   - **Remote Write 地址**：填入 `endpoint`
   - **Token**：填入 `token`

### 7.2 配置参数说明

| 参数 | 说明 | 示例 |
|------|------|------|
| `endpoint` | Remote Write 地址（不含 https://） | `your-instance.tencentcloudprom.com` |
| `token` | 认证 Token | `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| `protocol` | 传输协议 | `http`（推荐）或 `grpc` |
| `insecure` | 是否禁用 TLS | 同 VPC 可用 `true`，公网用 `false` |
| `url_path` | OTLP URL 路径 | `/v1/metrics` |

### 7.3 完整配置示例

```yaml
web:
  monitor:
    open_telemetry:
      enabled: true
      metric_collect_duration: 60s
      otel_metric_exporter_type: metric_otlp
      otel_trace_exporter_type: trace_stdout
      otel_metric_exporter:
        otlp:
          endpoint: "your-instance.tencentcloudprom.com"
          token: "your-prometheus-token"
          protocol: "http"
          insecure: false
          url_path: "/v1/metrics"
      resource:
        service_name: "sea-date"
        attrs:
          env: "production"
          region: "ap-guangzhou"
```

## 8. 资源统计

### 8.1 ResourceStatsService

自动采集系统资源指标：

```go
type ResourceStatsService struct {
    checkInterval time.Duration
    // 采集的指标
    cpuUsage      float64   // CPU 使用率
    memoryUsage   float64   // 内存使用率
    goroutines    int       // Goroutine 数量
    gcPauseNs     uint64    // GC 暂停时间
}
```

### 8.2 采集的指标

| 指标 | 说明 | 类型 |
|------|------|------|
| `process_cpu_usage` | 进程 CPU 使用率 | Gauge |
| `process_memory_usage` | 进程内存使用率 | Gauge |
| `go_goroutines` | Goroutine 数量 | Gauge |
| `go_gc_duration_seconds` | GC 暂停时间 | Histogram |
| `go_memstats_alloc_bytes` | 已分配内存 | Gauge |

## 9. 使用示例

### 9.1 编程方式配置

```go
import (
    "github.com/kaydxh/golang/pkg/opentelemetry"
    "github.com/kaydxh/golang/pkg/opentelemetry/metric/otlp"
)

// 方式一：使用腾讯云 Prometheus 快捷方法
otelService := opentelemetry.NewOpenTelemetryService(
    opentelemetry.WithMeterOptions(
        otlp.WithTencentCloudPrometheus(
            "your-instance.tencentcloudprom.com",
            "your-token",
            false,  // 使用 TLS
        ),
    ),
)
otelService.Install(ctx)

// 方式二：自定义 OTLP 配置
otelService := opentelemetry.NewOpenTelemetryService(
    opentelemetry.WithMeterOptions(
        otlp.WithOTLPPushExporter(
            "prometheus.example.com:4318",
            otlp.WithProtocol(otlp.ProtocolHTTP),
            otlp.WithHeaders(map[string]string{
                "Authorization": "Bearer xxx",
            }),
            otlp.WithInsecure(true),
        ),
    ),
)
```

### 9.2 配置文件方式

```yaml
# sea-date.yaml
web:
  monitor:
    open_telemetry:
      enabled: true
      metric_collect_duration: 60s
      otel_metric_exporter_type: metric_otlp
      otel_metric_exporter:
        otlp:
          endpoint: "your-instance.tencentcloudprom.com"
          token: "your-prometheus-token"
          protocol: "http"
          insecure: false
```

## 10. 文件结构

```
pkg/opentelemetry/
├── opentelemetry.go              # OpenTelemetryService 主入口
├── opentelemetry.proto           # Proto 配置定义
├── opentelemetry.pb.go           # 生成的 Proto 代码
├── opentelemetry.yaml            # 配置示例
├── config.go                     # 配置加载和安装
├── config.option.go              # 配置选项
├── metric/
│   ├── meter.go                  # MeterProvider 管理
│   ├── meter.option.go           # Meter 选项
│   ├── meter.pull.exporter.go    # Pull 模式导出器接口
│   ├── meter.push.exporter.go    # Push 模式导出器接口
│   ├── prometheus/               # Prometheus Pull 导出器
│   │   ├── prometheus.metric.go
│   │   └── prometheus.metric_option.go
│   ├── otlp/                     # OTLP Push 导出器
│   │   ├── otlp.metric.go
│   │   ├── otlp.metric.option.go
│   │   └── otlp.metric_option.go
│   └── stdout/                   # Stdout 导出器
│       └── stdout.metric.go
├── tracer/
│   ├── tracer.go                 # TracerProvider 管理
│   ├── jaeger/                   # Jaeger 导出器
│   └── stdout/                   # Stdout 导出器
└── resource/
    └── resource.go               # 资源统计服务

pkg/middleware/grpc-middleware/opentelemetry/
├── metric_server.interceptor.go  # gRPC 指标拦截器
├── trace_server.interceptor.go   # gRPC 链路拦截器
└── trace_client.interceptor.go   # gRPC 客户端链路拦截器
```

## 11. 设计特点总结

| 特点 | 说明 |
|------|------|
| OpenTelemetry 标准 | 基于 CNCF OpenTelemetry 标准实现 |
| 多后端支持 | 支持 Prometheus、OTLP、Jaeger、Stdout 等 |
| Pull/Push 双模式 | 同时支持 Pull（Prometheus）和 Push（OTLP）模式 |
| 腾讯云兼容 | 原生支持腾讯云 Prometheus Remote Write |
| 配置驱动 | 通过 YAML/Proto 配置灵活切换后端 |
| 中间件集成 | gRPC/HTTP 请求自动采集指标和链路 |
| 资源监控 | 自动采集 CPU、内存、GC 等系统指标 |
| 低侵入性 | 通过中间件透明接入，业务代码无感知 |
