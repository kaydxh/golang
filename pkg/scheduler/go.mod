module github.com/kaydxh/golang/pkg/scheduler

go 1.25.3

replace github.com/kaydxh/golang/go => ../../go

require (
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	github.com/kaydxh/golang/go v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.47.0
	google.golang.org/protobuf v1.36.10
)

require github.com/stretchr/testify v1.11.1 // indirect
