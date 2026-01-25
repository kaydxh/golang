module github.com/kaydxh/golang/pkg/gocv

go 1.24.0

replace github.com/kaydxh/golang/go => ../../go

require (
	github.com/kaydxh/golang/go v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20251113190631-e25ba8c21ef6
	google.golang.org/protobuf v1.36.10
)

require (
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.77.0 // indirect
)
