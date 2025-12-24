# golang

Golang lib is a comprehensive programming toolkit for building microservices in Go. It provides a rich set of interfaces and services to accelerate your development.

## Overview

This library offers a complete set of tools and utilities for developing Go microservices, including web servers, database connections, message queues, service discovery, monitoring, logging, and much more. It's designed to speed up development by providing production-ready components with sensible defaults.

## Features

### Core Components

- **Web Server** (`pkg/webserver`): High-performance web server built on Gin and gRPC-Gateway
  - HTTP and gRPC support
  - Graceful shutdown
  - Health checks and profiling endpoints
  - Service registry integration (Consul)
  - Request/response middleware

- **Configuration** (`pkg/viper`, `pkg/config`): Flexible configuration management
  - Viper integration for YAML/JSON/TOML configs
  - Environment variable support
  - Configuration validation

- **Logging** (`pkg/logs`): Advanced logging capabilities
  - Logrus-based logging
  - Log rotation (size/time-based)
  - Context-aware logging with request IDs
  - Multiple output formats (JSON, text)

- **Database** (`pkg/database`): Database connectivity
  - MySQL support with connection pooling
  - Redis client with advanced features
  - Transaction management

- **Message Queue** (`pkg/mq`): Message queue integration
  - Kafka producer and consumer
  - Batch processing
  - Error handling and retry mechanisms

- **Service Discovery** (`pkg/discovery`): Service registration and discovery
  - Consul integration
  - Etcd support
  - Health check integration

- **Monitoring** (`pkg/monitor`): Observability tools
  - OpenTelemetry integration
  - Prometheus metrics
  - Distributed tracing

- **Middleware** (`pkg/middleware`): HTTP and gRPC middleware
  - Authentication/Authorization
  - Rate limiting
  - Request logging
  - Error handling

- **Storage** (`pkg/storage`): Object storage support
  - S3-compatible storage
  - File mount utilities

- **Task Management** (`pkg/scheduler`, `pkg/crontab`, `pkg/pool`):
  - Task scheduler
  - Cron job support
  - Connection and task pools
  - Work queues

- **File Operations**:
  - File rotation (`pkg/file-rotate`)
  - File transfer (`pkg/file-transfer`)
  - File cleanup (`pkg/file-cleanup`)
  - File system monitoring (`pkg/fsnotify`)

- **gRPC Gateway** (`pkg/grpc-gateway`): RESTful API gateway for gRPC services

- **DNS Resolver** (`pkg/resolver`): Advanced DNS resolution with caching

- **Binary Log** (`pkg/binlog`): Binary log management for data replication

- **Other Utilities**:
  - Proxy support (`pkg/proxy`)
  - Performance profiling (`pkg/profile`)
  - Protobuf utilities (`pkg/protobuf`)
  - OpenCV bindings (`pkg/gocv`)

### Standard Library Enhancements (`go/`)

The `go/` directory contains enhanced versions of standard library packages with additional functionality:

- **Context**: Request ID extraction and propagation
- **Errors**: Structured error handling with codes
- **Time**: Exponential backoff, rate limiting, time utilities
- **Net**: gRPC/HTTP clients, DNS resolver, IP/MAC utilities
- **Crypto**: AES, MD5, SHA256 utilities
- **IO**: File operations, copy utilities
- **Reflect**: Advanced reflection utilities
- **Container**: Heap, Set, WorkQueue data structures
- **And more...**

### Tutorials (`tutorial/`)

Programming paradigm examples including:
- Function options pattern
- Decorator pattern
- Factory pattern
- Error handling patterns
- Pipeline pattern
- Visitor pattern
- And more...

## Quick Start

### Installation

```bash
go get github.com/kaydxh/golang
```

### Basic Web Server Example

```go
package main

import (
	"context"
	"testing"

	viper_ "github.com/kaydxh/golang/pkg/viper"
	webserver_ "github.com/kaydxh/golang/pkg/webserver"
)

func main() {
	cfgFile := "./webserver.yaml"
	config := webserver_.NewConfig(webserver_.WithViper(viper_.GetViper(cfgFile, "web")))

	s, err := config.Complete().New(context.Background())
	if err != nil {
		panic(err)
	}
	
	// Install your web handlers
	s.InstallWebHandlers(/* your handlers */)
	
	prepared, err := s.PrepareRun()
	if err != nil {
		panic(err)
	}

	prepared.Run(context.Background())
}
```

### Configuration Example

Create a `webserver.yaml` file:

```yaml
web:
  bind_address:
    host: 0.0.0.0
    port: 10000
  grpc:
    max_concurrency_unary: 0
    max_concurrency_stream: 0
    max_receive_message_size: 0
    max_send_message_size: 0
    timeout: 0s
  http:
    max_concurrency: 0
    timeout: 0s
  debug:
    enable_profiling: true
  monitor:
    open_telemetry:
      enabled: true
      metric_collect_duration: 60s
      otel_trace_exporter_type: trace_stdout
      otel_metric_exporter_type: metric_stdout
```

## Package Structure

```
golang/
├── go/              # Standard library enhancements
├── pkg/             # Main packages
│   ├── webserver/   # Web server framework
│   ├── logs/        # Logging
│   ├── database/    # Database clients
│   ├── mq/          # Message queues
│   ├── discovery/   # Service discovery
│   ├── monitor/     # Monitoring
│   └── ...          # Other packages
├── tutorial/        # Programming examples
└── script/          # Utility scripts
```

## Requirements

- Go 1.24.0 or higher
- See `go.mod` for specific dependency versions

## Documentation

- Each package contains detailed documentation
- Check `tutorial/` directory for usage examples
- Review package-level comments for API documentation

## Contributing

If you need support, start with your branch, and create a pull request for us. We appreciate your help!

### Development Guidelines

1. Follow Go coding standards
2. Add tests for new features
3. Update documentation
4. Ensure backward compatibility when possible

## License

See [LICENSE](LICENSE) file for details.

## Evolution

Golang started in Oct 8, 2020.

## Related Projects

This library is designed to work seamlessly with:
- Gin web framework
- gRPC
- Consul
- Etcd
- Kafka
- Prometheus
- OpenTelemetry

## Support

For issues, questions, or contributions, please open an issue or pull request on GitHub.
