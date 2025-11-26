module github.com/kaydxh/golang/go

go 1.25.3

replace github.com/kaydxh/golang/pkg/viper => ../pkg/viper

replace github.com/kaydxh/golang/pkg/protobuf => ../pkg/protobuf

replace github.com/kaydxh/golang/pkg/grpc-gateway => ../pkg/grpc-gateway

replace github.com/kaydxh/golang/pkg/middleware => ../pkg/middleware

replace github.com/kaydxh/golang/pkg/webserver => ../pkg/webserver

replace github.com/kaydxh/golang/pkg/logs => ../pkg/logs

require (
	github.com/gin-gonic/gin v1.11.0
	github.com/go-playground/validator/v10 v10.28.0
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/kaydxh/golang/pkg/logs v0.0.0-20251125132602-62c26b682228
	github.com/kaydxh/golang/pkg/middleware v0.0.0-20251125160242-e06b25c89946
	github.com/moby/term v0.5.2
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.9.3
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/spf13/cobra v1.10.1
	github.com/spf13/pflag v1.0.10
	github.com/stretchr/testify v1.11.1
	golang.org/x/exp v0.0.0-20251113190631-e25ba8c21ef6
	golang.org/x/net v0.47.0
	golang.org/x/sys v0.38.0
	golang.org/x/text v0.31.0
	golang.org/x/time v0.14.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
	gotest.tools v2.2.0+incompatible
	gotest.tools/v3 v3.5.2
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.10 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kaydxh/golang/pkg/file-cleanup v0.0.0-20251125132602-62c26b682228 // indirect
	github.com/kaydxh/golang/pkg/file-rotate v0.0.0-20251125132602-62c26b682228 // indirect
	github.com/kaydxh/golang/pkg/protobuf v0.0.0-20251125160242-e06b25c89946 // indirect
	github.com/kaydxh/golang/pkg/viper v0.0.0-20251125160242-e06b25c89946 // indirect
	github.com/kaydxh/golang/pkg/webserver v0.0.0-20251125160242-e06b25c89946 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/arch v0.20.0 // indirect
	golang.org/x/crypto v0.44.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251022142026-3a174f9686a8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
