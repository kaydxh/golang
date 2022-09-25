package marshaler

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type JSONPb struct {
	runtime.JSONPb
	opts struct {
		useProtoNames   bool
		useEnumNumbers  bool
		emitUnpopulated bool
		discardUnknown  bool
		indent          string
	}
}

func NewDefaultJSONPb() *JSONPb {
	return NewJSONPb(
		// use json name
		WithUseProtoNames(false),
		WithUseEnumNumbers(false),
		WithEmitUnpopulated(true),
		WithDiscardUnknown(true),
		WithIndent("\t"),
	)
}

func NewJSONPb(options ...JSONPbOption) *JSONPb {
	j := &JSONPb{}
	j.ApplyOptions(options...)
	return j
}
