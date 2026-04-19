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

## 8. K8s Resource 属性

### 8.1 概述

自动从 K8s 环境变量中采集容器运行时信息，作为 OpenTelemetry Resource 属性注入到所有 Metrics 和 Traces 中。

### 8.2 支持的 K8s 属性

| 环境变量 | 属性名 | 说明 |
|----------|--------|------|
| `NODE_IP` | `k8s.node.ip` | 宿主机节点 IP |
| `POD_NAMESPACE` | `k8s.namespace.name` | Pod 命名空间 |
| `POD_NAME` | `k8s.pod.name` | Pod 名称 |
| `POD_IP` | `k8s.pod.ip` | Pod IP |
| `CONTAINER_NAME` | `k8s.container.name` | 容器名称 |
| `CONTAINER_PLATFORM` | `k8s.container.platform` | 容器平台 (STKE/TKE) |

### 8.3 K8s Deployment 配置示例

```yaml
spec:
  containers:
  - name: my-app
    env:
    - name: NODE_IP
      valueFrom:
        fieldRef:
          fieldPath: status.hostIP
    - name: POD_NAMESPACE
      valueFrom:
        fieldRef:
          fieldPath: metadata.namespace
    - name: POD_NAME
      valueFrom:
        fieldRef:
          fieldPath: metadata.name
    - name: POD_IP
      valueFrom:
        fieldRef:
          fieldPath: status.podIP
    - name: CONTAINER_NAME
      value: "my-app"
    - name: CONTAINER_PLATFORM
      value: "STKE"
```

### 8.4 YAML 配置

```yaml
open_telemetry:
  resource:
    service_name: "my-service"
    service_version: "1.0.0"
    attrs:
      env: "production"
    k8s:
      enabled: true  # 是否启用 K8s 属性自动检测（默认 true）
      # 以下属性会自动从环境变量读取，也可以手动覆盖
      # node_ip: ""
      # pod_namespace: ""
      # pod_name: ""
      # pod_ip: ""
      # container_name: ""
      # container_platform: ""
    apm:
      token: ""  # APM Token，与业务系统关联，从 APM 控制台获取
```

### 8.5 APM Token 配置

#### 概述

APM（Application Performance Management）Token 用于腾讯云 APM 服务的身份认证和数据关联。Token 会作为 Resource attribute 添加到所有 trace span 中。

#### 获取 APM Token

1. 登录 [腾讯云 APM 控制台](https://console.cloud.tencent.com/apm/monitor/access)
2. 选择对应业务系统
3. 获取 Token

#### YAML 配置

```yaml
open_telemetry:
  enabled: true
  resource:
    service_name: "my-service"
    apm:
      token: "your-apm-token"  # 从 APM 控制台获取
```

#### Proto 定义

```protobuf
message Resource {
  string service_name = 1;
  string service_version = 2;
  map<string, string> attrs = 3;
  K8sResource k8s = 10;
  Apm apm = 11;      // 腾讯云 APM 配置
  ZhiYan zhiyan = 12; // 智研平台配置
}

message Apm {
  string token = 1;  // APM Token，与业务系统关联
}

message ZhiYan {
  string app_mark = 1;         // 业务指标上报 appMark
  string global_app_mark = 2;  // 全局指标上报 appMark
  string env = 3;              // 环境标识 (prod/test/dev)
  string instance_mark = 4;    // 实例标识
  string zhiyan_apm_token = 5; // 智研 APM Token
  string expand_key = 6;       // 是否扩展属性到维度
}
```

#### 核心实现

```go
// APM Token key
var ApmTokenKey = attribute.Key("token")

// WithApmToken 设置 APM Token
func WithApmToken(token string) ResourceOption {
    return func(o *ResourceOptions) {
        o.ApmToken = token
    }
}

// 在 NewResource 中添加 APM Token
if options.ApmToken != "" {
    attrs = append(attrs, ApmTokenKey.String(options.ApmToken))
}
```

### 8.6 智研 (ZhiYan) 平台配置

#### 概述

智研是腾讯内部的可观测平台，支持 Metrics、Traces、Logs 的统一上报和分析。通过 OpenTelemetry 协议接入。

#### 智研 vs 腾讯云 APM 对比

| 平台 | 访问范围 | 认证方式 | 上报地址 |
|------|---------|---------|---------|
| 腾讯云 APM | 公网 | `Authorization: Bearer <token>` | `xxx.tencentcloudapi.com` |
| 智研平台 | 内网 | `tps.tenant.id` Resource Attribute | `<智研内网trace地址>:4317` |

#### 获取智研租户

1. 登录 
2. 进入 监控宝 > 应用性能监控 > 接入管理
3. 获取租户信息，格式：`空间ID#日志租户#监控宝租户`
   - 示例：`<空间ID>#<日志租户>#<监控宝租户>`

#### 智研上报地址

| 网络环境 | 类型 | 服务地址 | 协议 | 备注 |
|----------|------|---------|------|------|
| **内网** | Metric | `<内网metric上报地址>:4318` | HTTP | 推荐，无需 TLS |
| **外网** | Metric | `<外网metric上报地址>:4318` | **HTTPS** | 必须启用 TLS |
| 内网 | Trace | `<内网trace上报地址>:4317` | gRPC | IDC 内网 |
| 内网 | Trace | `<内网trace上报地址>:4318` | HTTP | IDC 内网 |
| 公网 | Trace | `<外网trace上报地址>:4317` | gRPC | DevCloud/研发环境 |

> **⚠️ 重要**: 外网域名 **必须使用 HTTPS**（`insecure: false`），否则会连接超时。

#### YAML 配置示例

```yaml
open_telemetry:
  enabled: true
  
  # Trace 上报到智研
  otel_trace_exporter_type: trace_otlp
  otel_trace_exporter:
    otlp:
      endpoint: "<智研trace上报地址>:4317"
      protocol: "grpc"
      insecure: true
  
  # Metric 上报到智研
  otel_metric_exporter_type: metric_otlp
  otel_metric_exporter:
    otlp:
      endpoint: "<智研metric上报地址>:4317"
      protocol: "grpc"
      insecure: true
  
  resource:
    service_name: "my-service"
    zhiyan:
      # 业务指标上报 appMark（用于 App MeterProvider）
      app_mark: "<your_app_mark>"
      
      # 全局指标上报 appMark（用于 Global MeterProvider）
      global_app_mark: "<your_global_app_mark>"
      
      # 环境标识
      env: "prod"
      
      # 智研 APM Token（用于 Trace 上报）
      # 格式：空间ID#日志租户#监控宝租户
      zhiyan_apm_token: "<空间ID>#<日志租户>#<监控宝租户>"
      
      # 是否将 resource 属性扩展到指标维度
      expand_key: "no"
```

#### 智研 Resource Attribute Keys

| Attribute Key | 说明 | 是否必填 | 类型 |
|---------------|------|---------|------|
| `__zhiyan_app_mark__` | 上报应用标记 | **是** | string |
| `__zhiyan_env__` | 环境标识 (prod/test/dev) | **是** | string |
| `__zhiyan_instance_mark__` | 上报实例标识 | 否 | string |
| `__zhiyan_expand_tag_enable__` | 是否扩展属性到维度 | 否（默认 no） | string |
| `__zhiyan_data_grain__` | 数据粒度 (10/30/60) | 否（默认 60） | int |
| `__zhiyan_data_type__` | 数据类型，秒级填 "second" | 否 | string |
| `tps.tenant.id` | 智研 APM 租户 ID | Trace 上报必填 | string |

#### 智研上报必须配置

根据智研对接文档，以下配置项**必须正确设置**：

| 配置项 | 要求 | 说明 |
|--------|------|------|
| Delta Temporality | `temporality_delta: true` | 智研只支持 Delta 时间性 |
| Gzip 压缩 | `compression: true` | 建议启用以减少带宽 |
| 采集周期 | 与数据粒度一致 | 如分钟级数据设置 60s |
| HTTP 端口 | 4318 | OTLP HTTP 标准端口 |

#### 核心实现

```go
// 智研平台 attribute keys
var (
    ZhiYanAppMarkKey      = attribute.Key("__zhiyan_app_mark__")
    ZhiYanInstanceMarkKey = attribute.Key("__zhiyan_instance_mark__")
    ZhiYanEnvKey          = attribute.Key("__zhiyan_env__")
    ZhiYanExpandKey       = attribute.Key("__zhiyan_expand_tag_enable__")
    ZhiYanDataGrainKey    = attribute.Key("__zhiyan_data_grain__")
    ZhiYanDataTypeKey     = attribute.Key("__zhiyan_data_type__")
    ZhiYanTpsTenantIDKey  = attribute.Key("tps.tenant.id")
)

// 在 NewResource 中添加智研属性
if options.ZhiYanAppMark != "" {
    attrs = append(attrs, ZhiYanAppMarkKey.String(options.ZhiYanAppMark))
    attrs = append(attrs, ZhiYanEnvKey.String(options.ZhiYanEnv))
    attrs = append(attrs, ZhiYanExpandKey.String("no"))
}

if options.ZhiYanApmToken != "" {
    attrs = append(attrs, ZhiYanTpsTenantIDKey.String(options.ZhiYanApmToken))
}
```

### 8.7 核心实现

```go
// K8s 环境变量映射
var k8sEnvToAttribute = map[string]attribute.Key{
    "NODE_IP":            semconv.K8SNodeNameKey,
    "POD_NAMESPACE":      semconv.K8SNamespaceNameKey,
    "POD_NAME":           semconv.K8SPodNameKey,
    "POD_IP":             attribute.Key("k8s.pod.ip"),
    "CONTAINER_NAME":     semconv.K8SContainerNameKey,
    "CONTAINER_PLATFORM": attribute.Key("k8s.container.platform"),
}

// NewResource 创建包含 K8s 属性的 Resource
func NewResource(opts ...ResourceOption) (*resource.Resource, error) {
    cfg := defaultResourceConfig()
    for _, opt := range opts {
        opt(cfg)
    }
    
    var attrs []attribute.KeyValue
    
    // 添加服务信息
    if cfg.serviceName != "" {
        attrs = append(attrs, semconv.ServiceName(cfg.serviceName))
    }
    if cfg.serviceVersion != "" {
        attrs = append(attrs, semconv.ServiceVersion(cfg.serviceVersion))
    }
    
    // 添加 K8s 属性
    if cfg.enableK8s {
        attrs = append(attrs, getK8sAttributes()...)
    }
    
    // 添加自定义属性
    for k, v := range cfg.attrs {
        attrs = append(attrs, attribute.String(k, v))
    }
    
    return resource.NewWithAttributes(semconv.SchemaURL, attrs...)
}
```

## 9. 模调上报（主/被调）

### 9.1 概述

模调上报是一种标准化的服务调用监控方案，支持：
- **被调上报 (P)**：服务端记录被调用的指标
- **主调上报 (A)**：客户端记录调用其他服务的指标

### 9.2 上报指标

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `requests` | Counter | 请求数 |
| `success` | Counter | 成功数 |
| `timeout` | Counter | 超时数 |
| `abnormal` | Counter | 异常数 |
| `cost` | Counter | 累计耗时(ms) |
| `duration` | Histogram | 耗时分布 |

### 9.3 维度属性

#### 被调方 (P) 维度

| 属性 | 说明 |
|------|------|
| `ret_code` | 返回码 |
| `p_ip` | 被调 IP |
| `p_app` | 被调应用名 |
| `p_server` | 被调服务名 |
| `p_service` | 被调 Service |
| `p_interface` | 被调接口 |

#### 主调方 (A) 维度

| 属性 | 说明 |
|------|------|
| `a_ip` | 主调 IP |
| `a_app` | 主调应用名 |
| `a_server` | 主调服务名 |
| `a_service` | 主调 Service |
| `a_interface` | 主调接口 |

### 9.4 gRPC 拦截器使用示例

#### Server 端（被调上报）

```go
import interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/opentelemetry"

server := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        interceptoropentelemetry.UnaryServerModularInterceptor(
            interceptoropentelemetry.ModularServerConfig{
                AppName:    "my-app",
                ServerName: "my-server",
            },
        ),
    ),
    grpc.ChainStreamInterceptor(
        interceptoropentelemetry.StreamServerModularInterceptor(
            interceptoropentelemetry.ModularServerConfig{
                AppName:    "my-app",
                ServerName: "my-server",
            },
        ),
    ),
)
```

#### Client 端（主调上报）

```go
conn, err := grpc.Dial(target,
    grpc.WithChainUnaryInterceptor(
        interceptoropentelemetry.UnaryClientModularInterceptor(
            interceptoropentelemetry.ModularClientConfig{
                AppName:     "my-app",
                ServerName:  "my-server",
                ServiceName: "my-service",
            },
        ),
    ),
    grpc.WithChainStreamInterceptor(
        interceptoropentelemetry.StreamClientModularInterceptor(
            interceptoropentelemetry.ModularClientConfig{
                AppName:     "my-app",
                ServerName:  "my-server",
                ServiceName: "my-service",
            },
        ),
    ),
)
```

### 9.5 核心实现

```go
// MetricReporter 模调上报器
type MetricReporter struct {
    meterProvider otelmetric.MeterProvider
    counters      map[string]otelmetric.Int64Counter
    histograms    map[string]otelmetric.Float64Histogram
    mu            sync.RWMutex
}

// ReportServerMetric 被调上报
func (r *MetricReporter) ReportServerMetric(ctx context.Context, dim *ServerDimension, costMs float64) {
    attrs := dim.ToAttributes()
    
    // 请求数
    r.getCounter(ServerReportMeterName, RequestsMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
    
    // 成功/超时/异常
    if dim.IsSuccess() {
        r.getCounter(ServerReportMeterName, SuccessMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
    } else if dim.IsTimeout() {
        r.getCounter(ServerReportMeterName, TimeoutMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
    } else {
        r.getCounter(ServerReportMeterName, AbnormalMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
    }
    
    // 耗时
    r.getCounter(ServerReportMeterName, CostMetricName).Add(ctx, int64(costMs), otelmetric.WithAttributes(attrs...))
    r.getHistogram(ServerReportMeterName, DurationMetricName, DefaultDurationBounds).Record(ctx, costMs, otelmetric.WithAttributes(attrs...))
}
```

## 10. 自定义 Metric API

### 10.1 概述

提供便捷的 Metric API，简化指标上报操作，支持：
- 函数式 API（Context 属性传递）
- 函数式 API（显式属性传递）
- 面向对象 API（Instrument 封装）
- Global/App 双 Provider 支持

### 10.2 函数式 API（Context 属性）

```go
import "github.com/kaydxh/golang/pkg/opentelemetry/metric/api"

// 通过 context 传递属性
ctx = api.WithAttribute(ctx, "meter_name", "user_id", "12345")
ctx = api.WithAttributes(ctx, "meter_name", map[string]any{
    "region":     "us-west-2",
    "version":    "1.2.3",
})

// Counter
api.AddCounter(ctx, "business", "orders_count", 10)
api.IncrCounter(ctx, "business", "user_signups")

// Histogram
api.RecordHistogram(ctx, "latency", "api_duration", 125.5, nil)
api.RecordDuration(ctx, "latency", "db_query", 45.2)  // 使用默认 bounds
```

### 10.3 函数式 API（显式属性）

```go
import "go.opentelemetry.io/otel/attribute"

// 直接传入属性，不依赖 context
api.IncrCounterWithAttrs(ctx, "http", "requests",
    attribute.String("endpoint", "/api/v1/users"),
    attribute.Int("status", 200),
)

api.RecordDurationWithAttrs(ctx, "http", "latency", 45.2,
    attribute.String("method", "GET"),
    attribute.String("path", "/users"),
)
```

### 10.4 面向对象 API

```go
// Counter 对象
counter := api.NewCounter("business", "api_calls")
counter.SetAttribute("service", "user-service")
counter.Incr(ctx)

// 链式调用
counter.With("method", "GET").With("path", "/users").Incr(ctx)

// Histogram 对象
histogram := api.NewHistogram("latency", "processing_time",
    api.WithBounds([]float64{10, 50, 100, 500, 1000}),
)
histogram.SetAttribute("processor", "image")
histogram.Record(ctx, 234.5)

// Timer 对象
timer := api.NewTimer("performance", "operation_duration")
timer.SetAttribute("operation", "data_processing")
timer.Record(ctx, float64(elapsed.Milliseconds()))
```

### 10.5 Global vs App Provider

```go
// App Provider（应用级指标）- 默认
api.IncrCounter(ctx, "business", "orders")
counter := api.NewCounter("app", "signups")

// Global Provider（基础设施指标）
api.GlobalIncrCounter(ctx, "infra", "gc_cycles")
globalCounter := api.NewGlobalCounter("system", "allocations")

// 设置自定义 App Provider（可选）
api.SetAppMeterProvider(myCustomProvider)
```

### 10.6 配置驱动的双 MeterProvider

支持通过配置文件自动初始化 Global 和 App 两个独立的 MeterProvider：

#### Proto 配置定义

```protobuf
message OpenTelemetry {
  // ... 其他字段 ...

  // App MeterProvider configuration (separate from global)
  AppMeterProvider app_meter_provider = 11;
}

message AppMeterProvider {
  bool enabled = 1;  // Enable separate App MeterProvider

  // Exporter type for app metrics (if different from global)
  OtelMetricExporterType exporter_type = 2;

  // Exporter configuration (if different from global)
  OtelMetricExporter exporter = 3;

  // Collect duration for app metrics
  google.protobuf.Duration collect_duration = 4;

  // Resource attributes specific to app metrics
  Resource resource = 5;
}
```

#### YAML 配置示例

```yaml
open_telemetry:
  enabled: true

  # Global MeterProvider (基础设施指标)
  otel_metric_exporter_type: metric_prometheus
  otel_metric_exporter:
    prometheus:
      url: "/metrics"
  metric_collect_duration: "60s"

  # App MeterProvider (应用级指标，独立配置)
  app_meter_provider:
    enabled: true
    exporter_type: metric_otlp  # 可以使用不同的导出器
    exporter:
      otlp:
        endpoint: "prometheus.tencentcloudapi.com:4317"
        protocol: "grpc"
        token: "your-app-token"
    collect_duration: "30s"  # 可以使用不同的采集周期
    resource:
      service_name: "my-app-metrics"
      attrs:
        metric_type: "business"
```

#### 使用场景

| 场景 | Global Provider | App Provider |
|------|-----------------|--------------|
| **用途** | 基础设施/运维指标 | 业务/应用指标 |
| **指标类型** | CPU、内存、GC、网络 | 订单数、用户注册、业务延迟 |
| **采集周期** | 较长（60s） | 较短（15-30s） |
| **导出目标** | 本地 Prometheus | 腾讯云/远程 OTLP |
| **资源属性** | 通用服务信息 | 业务相关属性 |

#### 编程方式配置

```go
// 方式1：使用配置文件自动初始化
cfg := opentelemetry.NewConfig(
    opentelemetry.WithViper(v),
)
cfg.Complete().New(ctx)

// 方式2：编程方式手动配置双 Provider
ot := opentelemetry.NewOpenTelemetryService(
    // Global Provider
    opentelemetry.WithMeterPullExporter(prometheusBuilder),
    opentelemetry.WithMetricCollectDuration(time.Minute),

    // App Provider
    opentelemetry.WithAppMeterPushExporter(otlpBuilder),
    opentelemetry.WithAppMetricCollectDuration(30 * time.Second),
    opentelemetry.WithAppMeterResource(appResource),
)
ot.Install(ctx)
```

### 10.7 支持的属性类型

| 类型 | 示例 |
|------|------|
| `string` | `"value"` |
| `int`, `int8`, `int16`, `int32`, `int64` | `123` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `456` |
| `float32`, `float64` | `3.14` |
| `bool` | `true` |

### 10.8 默认 Histogram Bounds

```go
// 延迟类（毫秒）
DefaultDurationBounds = []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}

// 大小类（字节）
DefaultSizeBounds = []float64{100, 1000, 10000, 100000, 1000000, 10000000}
```

## 11. 资源统计

### 11.1 ResourceStatsService

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

### 11.2 采集的指标

| 指标 | 说明 | 类型 |
|------|------|------|
| `process_cpu_usage` | 进程 CPU 使用率 | Gauge |
| `process_memory_usage` | 进程内存使用率 | Gauge |
| `go_goroutines` | Goroutine 数量 | Gauge |
| `go_gc_duration_seconds` | GC 暂停时间 | Histogram |
| `go_memstats_alloc_bytes` | 已分配内存 | Gauge |

## 12. 使用示例

### 12.1 编程方式配置

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

### 12.2 配置文件方式

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

## 13. 文件结构

```
pkg/opentelemetry/
├── opentelemetry.go              # OpenTelemetryService 主入口
├── opentelemetry.proto           # Proto 配置定义
├── opentelemetry.pb.go           # 生成的 Proto 代码
├── opentelemetry.yaml            # 配置示例
├── opentelemetry.option.go       # OpenTelemetry 选项
├── config.go                     # 配置加载和安装
├── config.option.go              # 配置选项
├── metric/
│   ├── meter.go                  # MeterProvider 管理
│   ├── meter.option.go           # Meter 选项
│   ├── meter.pull.exporter.go    # Pull 模式导出器接口
│   ├── meter.push.exporter.go    # Push 模式导出器接口
│   ├── api/                      # 自定义 Metric API
│   │   ├── api.go                # 便捷 API 函数
│   │   └── instrument.go         # Instrument 对象封装
│   ├── report/                   # 模调上报
│   │   ├── report.go             # 上报核心实现
│   │   └── dimension.go          # 维度定义
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
│   ├── tracer.option.go          # Tracer 选项
│   ├── jaeger/                   # Jaeger 导出器
│   └── stdout/                   # Stdout 导出器
└── resource/
    └── resource.go               # K8s Resource 属性

pkg/middleware/grpc-middleware/opentelemetry/
├── metric_server.interceptor.go      # gRPC 指标拦截器
├── trace_server.interceptor.go       # gRPC 链路拦截器（Server）
├── trace_client.interceptor.go       # gRPC 链路拦截器（Client）
├── modular_server.interceptor.go     # 模调上报拦截器（Server/被调）
└── modular_client.interceptor.go     # 模调上报拦截器（Client/主调）
```

## 14. 设计特点总结

| 特点 | 说明 |
|------|------|
| OpenTelemetry 标准 | 基于 CNCF OpenTelemetry 标准实现 |
| 多后端支持 | 支持 Prometheus、OTLP、Jaeger、Stdout 等 |
| Pull/Push 双模式 | 同时支持 Pull（Prometheus）和 Push（OTLP）模式 |
| 腾讯云兼容 | 原生支持腾讯云 Prometheus Remote Write |
| APM Token 支持 | 支持腾讯云 APM Token 认证，作为 Resource attribute 注入 |
| 智研平台支持 | 支持腾讯内部智研平台，通过 tps.tenant.id 认证，支持 Metric/Trace 上报 |
| 配置驱动 | 通过 YAML/Proto 配置灵活切换后端 |
| 中间件集成 | gRPC/HTTP 请求自动采集指标和链路 |
| 资源监控 | 自动采集 CPU、内存、GC 等系统指标 |
| 低侵入性 | 通过中间件透明接入，业务代码无感知 |
| K8s 属性 | 自动采集 K8s 容器运行时信息 |
| 模调上报 | 支持主调/被调标准化监控 |
| 自定义 Metric API | 提供便捷的指标上报 API |
| Global/App 双 Provider | 支持基础设施和应用级指标分离 |

## 15. 常见问题排查

### 15.1 智研上报超时

**现象**：
```
[ERRO] metric export failed: context deadline exceeded: retry-able request failure
```

**原因与解决方案**：

| 错误原因 | 解决方案 |
|----------|----------|
| 使用内网地址但不在内网 | 切换到外网地址 + HTTPS |
| 使用外网地址但未启用 TLS | 设置 `insecure: false` |
| 网络不通 | 检查防火墙和网络连通性 |

**配置对照表**：

```yaml
# 内网环境（推荐）
endpoint: "<内网metric上报地址>:4318"
insecure: true   # HTTP

# 外网环境
endpoint: "<外网metric上报地址>:4318"
insecure: false  # 必须 HTTPS
```

### 15.2 指标未上报到智研

**检查清单**：

1. **必填字段是否配置**：
   - `__zhiyan_app_mark__`：必须配置 appMark
   - `__zhiyan_env__`：必须配置环境标识

2. **Temporality 是否正确**：
   - 智研只支持 `Delta Temporality`
   - 配置：`temporality_delta: true`

3. **采集周期是否匹配**：
   - 数据粒度为 60 时，采集周期应为 60s
   - 配置：`metric_collect_duration: 60s`

4. **压缩是否启用**：
   - 建议启用 gzip：`compression: true`

### 15.3 指标组/Scope Name 配置

智研模调分析需要正确的指标组：

| 指标组 | 用途 | 配置 |
|--------|------|------|
| `server_report` | 被调上报 | 服务端拦截器 |
| `client_report` | 主调上报 | 客户端拦截器 |
| `attr_report` | 属性监控 | 自定义业务指标 |
| `default` | 默认 | 通用指标 |

```yaml
resource:
  zhiyan:
    metric_group: "server_report"  # 指标组名称
```

### 15.4 完整的智研配置示例

```yaml
open_telemetry:
  enabled: true
  metric_collect_duration: 60s           # 采集周期与数据粒度一致
  otel_metric_exporter_type: metric_otlp
  otel_metric_exporter:
    otlp:
      # 根据网络环境选择内网或外网地址
      endpoint: "<智研metric上报地址>:4318"
      protocol: "http"
      insecure: true                      # 内网 HTTP / 外网必须 false
      url_path: "/v1/metrics"
      compression: true                   # 启用 gzip 压缩
      temporality_delta: true             # 必须使用 Delta

  resource:
    service_name: "my-service"
    zhiyan:
      global_app_mark: "<your_global_app_mark>" # 必填：appMark
      env: "prod"                         # 必填：环境
      instance_mark: ""                   # 选填：实例标识
      expand_key: "no"                    # 选填：是否扩展属性
      metric_group: "server_report"       # 选填：指标组
      # data_grain: 60                    # 选填：数据粒度 (10/30/60)
      # data_type: ""                     # 选填：秒级填 "second"
```

### 15.5 Trace 没有上报日志

**现象**：
- Meter 指标正常上报，但 Trace 没有任何 export 日志
- 调用 HTTP 接口后没有看到 span 创建

**原因分析**：

| 原因 | 说明 |
|------|------|
| TracerProvider 初始化时机错误 | otelgrpc interceptor 在创建时会缓存全局 TracerProvider，如果此时 TracerProvider 还没设置，会拿到 noop TracerProvider |
| HTTP 请求没有 trace interceptor | gRPC trace interceptor 只对 gRPC 请求生效，HTTP 请求需要单独的 HTTP trace middleware |
| Sampler 配置问题 | 默认 Sampler 可能不是 AlwaysSample，导致部分 span 被丢弃 |

**解决方案**：

1. **确保 TracerProvider 在 WebServer.Run() 之前初始化**：

```go
// options.go 中的正确顺序
s.installLogsOrDie()
s.installConfigOrDie()

ws, err := s.webServerConfig.Complete().New(ctx)  // ← 先创建 WebServer

// TracerProvider 可以在 WebServer.New() 之后初始化
// 但必须在 WebServer.Run() 之前完成
s.installOpenTelemetryTracerOrDie(ctx)
s.installOpenTelemetryMeterOrDie(ctx, ws)

ws.PrepareRun().Run(ctx)  // ← 启动前 Provider 已设置好
```

2. **HTTP trace middleware 动态获取 TracerProvider**：

```go
// HTTP middleware 在每次请求时动态获取 TracerProvider
func Trace(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tracer := otel.Tracer(tracerName)  // ← 请求时获取，此时已设置好
        ctx, span := tracer.Start(ctx, spanName, ...)
        // ...
    })
}
```

3. **gRPC trace interceptor 创建时缓存（注意）**：

```go
func UnaryServerTraceInterceptor(opts ...otelgrpc.Option) grpc.UnaryServerInterceptor {
    // gRPC interceptor 在创建时缓存 TracerProvider
    tp := otel.GetTracerProvider()
    defaultOpts := []otelgrpc.Option{
        otelgrpc.WithTracerProvider(tp),
    }
    defaultOpts = append(defaultOpts, opts...)
    return otelgrpc.UnaryServerInterceptor(defaultOpts...)
}
```

3. **添加 HTTP trace middleware**：

HTTP 请求（如 `GET /Now`）不经过 gRPC interceptor，需要单独添加 HTTP trace middleware：

```go
// http trace interceptor
func Trace(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 提取 trace context
        propagator := otel.GetTextMapPropagator()
        ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))
        
        // 创建 span
        tracer := otel.Tracer("http-trace")
        spanName := r.Method + " " + r.URL.Path
        ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
        defer span.End()
        
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}
```

4. **确保使用 AlwaysSample**：

```go
tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exp, batchOpts...),
    sdktrace.WithResource(res),
    sdktrace.WithSampler(sdktrace.AlwaysSample()),  // 确保所有 span 都被记录
)
```

**验证日志**：

正确配置后，启动日志应显示：
```
[INFO] TracerProvider created with BatchSpanProcessor (batch_timeout=5s, sampler=AlwaysSample)
[INFO] UnaryServerTraceInterceptor: creating trace interceptor, TracerProvider type=*trace.TracerProvider
```

如果显示 `TracerProvider type=*trace.noopTracerProvider`，说明 TracerProvider 初始化顺序有问题。

### 15.6 Trace 上报连接失败 (EOF)

**现象**：
```
[ERRO] [OTLP] trace export failed: traces export: Post "http://xxx:4318/v1/traces": EOF
```

Span 创建成功，但导出到远程服务器失败。

**原因分析**：

| 原因 | 说明 |
|------|------|
| Endpoint 地址不可达 | 网络不通、DNS 解析失败、防火墙阻断 |
| 协议不匹配 | 外网地址需要 HTTPS，但配置了 `insecure: true` |
| 端口错误 | gRPC 用 4317，HTTP 用 4318 |

**解决方案**：

根据网络环境选择正确的配置：

```yaml
# 方案1：内网 gRPC（推荐）
otel_trace_exporter:
  otlp:
    endpoint: "<内网trace上报地址>:4317"
    protocol: "grpc"
    insecure: true

# 方案2：内网 HTTP
otel_trace_exporter:
  otlp:
    endpoint: "<内网trace上报地址>:4318"
    protocol: "http"
    insecure: true

# 方案3：外网 HTTP（必须 HTTPS）
otel_trace_exporter:
  otlp:
    endpoint: "<外网trace上报地址>:4318"
    protocol: "http"
    insecure: false  # 外网必须启用 TLS

# 方案4：本地调试（stdout）
otel_trace_exporter_type: trace_stdout
```

**智研 Trace 上报地址汇总**：

| 网络环境 | 协议 | 地址 | 端口 | TLS |
|----------|------|------|------|-----|
| 内网 IDC | gRPC | `<内网trace上报地址>` | 4317 | 否 |
| 内网 IDC | HTTP | `<内网trace上报地址>` | 4318 | 否 |
| 外网/DevCloud | gRPC | `<外网trace上报地址>` | 4317 | 否 |
| 外网 | HTTP | `<外网metric上报地址>` | 4318 | **是** |

### 15.7 Tracer 和 Meter 初始化顺序问题

**问题**：Tracer 和 Meter 的初始化顺序有什么要求？

**原因分析**：

| 组件 | 必须在 WebServer 之前？ | 原因 |
|------|------------------------|------|
| **Tracer** | **否** | HTTP trace middleware 在请求时动态调用 `otel.Tracer()` 获取 TracerProvider |
| **Meter** | 否（但有例外） | 指标也是在请求处理时动态调用 `otel.Meter()` 获取 |

**HTTP Trace Middleware（动态获取）**：

```go
// HTTP trace middleware 在每次请求时动态获取 TracerProvider
func Trace(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 请求时动态获取，此时 TracerProvider 已设置
        tracer := otel.Tracer(tracerName)  // ← 每次请求时获取
        ctx, span := tracer.Start(ctx, spanName, ...)
        defer span.End()
        next.ServeHTTP(w, r)
    })
}
```

**gRPC Trace Interceptor（创建时缓存）**：

```go
// gRPC interceptor 在创建时获取并缓存 TracerProvider
func UnaryServerTraceInterceptor(opts ...otelgrpc.Option) grpc.UnaryServerInterceptor {
    // 创建时获取全局 TracerProvider 并传入
    tp := otel.GetTracerProvider()
    return otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp)...)
}
```

> **注意**：如果 gRPC interceptor 在 TracerProvider 设置之前创建，会缓存 noop TracerProvider。但当前实现中，HTTP middleware 是动态获取的，所以 Tracer 可以在 WebServer 之后初始化。

**Meter 的例外情况**：

当使用 **Prometheus Pull 模式** 时，需要在 Gin Router 上注册 `/metrics` 端点，此时需要先有 WebServer：

```go
// 需要 ws.GetGinEngine() 来注册 /metrics
s.opentelemetryConfig.ApplyOptions(opentelemetry_.WithGinRouter(ws.GetGinEngine()))
```

**推荐的初始化顺序**：

```
┌─────────────────────────────────────────────────────────────┐
│                     启动顺序                                  │
├─────────────────────────────────────────────────────────────┤
│  1. installLogs                                              │
│  2. installConfig                                            │
│  3. WebServer.New()                                          │
│  4. installMysql/Redis                                       │
│  5. installOpenTelemetryTracer  ←── 可以在 WebServer 之后     │
│  6. installOpenTelemetryMeter   ←── 可以在 WebServer 之后     │
│  7. WebServer.Run()             ←── 必须在 Run() 之前完成     │
└─────────────────────────────────────────────────────────────┘
```

> **关键点**：只要在 `WebServer.Run()` 之前完成 TracerProvider 和 MeterProvider 的初始化即可，因为 HTTP middleware 是在请求时动态获取 Provider 的。

### 15.8 调试技巧

**1. 启用 stdout exporter 验证 span 创建**：

```yaml
otel_trace_exporter_type: trace_stdout  # 临时切换到 stdout
otel_trace_exporter:
  stdout:
    pretty_print: true
```

**2. 启用 metric stdout exporter 查看上报数据**：

```yaml
otel_metric_exporter_type: metric_stdout
otel_metric_exporter:
  stdout:
    pretty_print: true
```

**3. 日志中确认 Resource Attribute**：
```
[DEBU] [OTLP] resource attribute: __zhiyan_app_mark__=<your_app_mark>
```

**4. 确认 TracerProvider 类型**：
```
[INFO] UnaryServerTraceInterceptor: creating trace interceptor, TracerProvider type=*trace.TracerProvider
```

如果显示 `*trace.noopTracerProvider`，说明初始化顺序有问题。

**5. 检查 span export 日志**：
```
[DEBU] [OTLP] exporting 1 spans to <trace上报地址>:4317
[INFO] [OTLP] trace export success: endpoint=xxx, spans=1, duration=xxx
```

**6. 网络连通性测试**：
```bash
# 测试 gRPC 端口
nc -zv <trace上报地址> 4317

# 测试 HTTP 端口
curl -v http://<trace上报地址>:4318/v1/traces
```
