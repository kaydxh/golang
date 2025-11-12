# OpenTelemetry gRPC Middleware

This package provides OpenTelemetry instrumentation for gRPC using the official `otelgrpc` package from OpenTelemetry.

## Overview

The package implements both server and client interceptors for distributed tracing using OpenTelemetry standards.

## Dependencies

- `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc` v0.53.0
- `go.opentelemetry.io/otel/trace` v1.30.0
- `go.opentelemetry.io/otel` v1.30.0

## Features

### Server Interceptors
- **UnaryServerTraceInterceptor**: Adds OpenTelemetry tracing to unary server RPCs
- **StreamServerTraceInterceptor**: Adds OpenTelemetry tracing to streaming server RPCs
- **UnaryServerTraceInterceptorWithTracer**: Unary server interceptor with custom tracer provider
- **StreamServerTraceInterceptorWithTracer**: Stream server interceptor with custom tracer provider

### Client Interceptors
- **UnaryClientTraceInterceptor**: Adds OpenTelemetry tracing to unary client RPCs
- **StreamClientTraceInterceptor**: Adds OpenTelemetry tracing to streaming client RPCs
- **UnaryClientTraceInterceptorWithTracer**: Unary client interceptor with custom tracer provider
- **StreamClientTraceInterceptorWithTracer**: Stream client interceptor with custom tracer provider

### Utility Functions
- **GetSpanFromContext**: Extracts the current span from context
- **GetTracerProvider**: Returns the global tracer provider
- **GetTracer**: Returns a tracer with the given name

## Usage

### Server Side

```go
import (
    interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
    "github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
    "github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer/jaeger"
    "google.golang.org/grpc"
)

func createGrpcServer() (*grpc.Server, error) {
    // Initialize tracer (optional, uses global tracer if not set)
    jaegerBuilder, err := jaeger.NewJaegerExporertBuilder("http://localhost:14268/api/traces")
    if err != nil {
        return nil, err
    }

    t := tracer.NewTracer(
        tracer.WithExporterBuilder(jaegerBuilder),
        tracer.WithServiceName("my-grpc-service"),
        tracer.WithServiceVersion("1.0.0"),
    )
    if err := t.Install(context.Background()); err != nil {
        return nil, err
    }

    // Create gRPC server with OpenTelemetry interceptors
    server := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            interceptoropentelemetry.UnaryServerTraceInterceptor(),
            // Add other interceptors as needed
        ),
        grpc.ChainStreamInterceptor(
            interceptoropentelemetry.StreamServerTraceInterceptor(),
        ),
    )

    return server, nil
}
```

### Client Side

```go
import (
    interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
    "google.golang.org/grpc"
)

func createGrpcClient(addr string) (*grpc.ClientConn, error) {
    // Create gRPC client with OpenTelemetry interceptors
    conn, err := grpc.NewClient(addr,
        grpc.WithChainUnaryInterceptor(
            interceptoropentelemetry.UnaryClientTraceInterceptor(),
            // Add other interceptors as needed
        ),
        grpc.WithChainStreamInterceptor(
            interceptoropentelemetry.StreamClientTraceInterceptor(),
        ),
    )
    if err != nil {
        return nil, err
    }

    return conn, nil
}
```

### Using Custom Tracer Provider

```go
import (
    interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
    "github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
    "google.golang.org/grpc"
)

func createGrpcServerWithCustomTracer() (*grpc.Server, error) {
    // Create and install tracer
    t := tracer.NewTracer(
        tracer.WithExporterBuilder(myExporterBuilder),
        tracer.WithServiceName("my-service"),
    )
    if err := t.Install(context.Background()); err != nil {
        return nil, err
    }

    // Create server with custom tracer provider
    server := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            interceptoropentelemetry.UnaryServerTraceInterceptorWithTracer(
                t.TracerProvider(),
                otelgrpc.WithMessageEvents(otelgrpc.SentEvents, otelgrpc.ReceivedEvents),
            ),
        ),
    )

    return server, nil
}
```

### Working with Spans

```go
import (
    interceptoropentelemetry "github.com/kaydxh/golang/pkg/middleware/grpc-middleware/monitor/opentelemetry"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
)

func myGrpcHandler(ctx context.Context, req *MyRequest) (*MyResponse, error) {
    // Get the current span from context
    span := interceptoropentelemetry.GetSpanFromContext(ctx)

    // Add custom attributes
    span.SetAttributes(
        attribute.String("custom.key", "value"),
        attribute.Int("user.id", 123),
    )

    // Add events
    span.AddEvent("processing request")

    // Set status
    span.SetStatus(codes.Ok, "success")

    // ... your handler logic ...

    return &MyResponse{}, nil
}
```

## Tracer Configuration

The tracer supports the following configuration options:

- **WithExporterBuilder**: Set the exporter builder (Jaeger, Stdout, etc.)
- **WithServiceName**: Set the service name for spans
- **WithServiceVersion**: Set the service version
- **WithServiceNamespace**: Set the service namespace

### Supported Exporters

1. **Jaeger**: For sending traces to Jaeger
2. **Stdout**: For debugging and development

## Integration with Existing Code

The interceptors are designed to work seamlessly with your existing gRPC setup. Simply add them to your interceptor chain alongside other interceptors like:

- Timer interceptors
- Debug interceptors
- Rate limiting interceptors
- Metric interceptors

## Best Practices

1. **Initialize tracer once**: Initialize the tracer provider at application startup
2. **Use context propagation**: Always pass context through your call chain
3. **Add meaningful attributes**: Use span attributes to add context-specific information
4. **Handle errors properly**: Set span status when errors occur
5. **Shutdown gracefully**: Call `tracer.Shutdown(ctx)` on application shutdown

## Related Packages

- `/pkg/monitor/opentelemetry/tracer` - Tracer initialization and configuration
- `/pkg/monitor/opentelemetry/metric` - Metric collection with OpenTelemetry
- `/pkg/middleware/grpc-middleware/monitor/opentelemetry` - This package

