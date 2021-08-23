package tcloud

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsonpb_ "github.com/kaydxh/golang/pkg/protobuf/jsonpb"
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

// Marshal marshals "v" into JSON.
func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {
	respStruct, err := jsonpb_.MarshaToStructpb(v)
	if err != nil {
		return nil, err
	}
	body := &TCloudResponse{
		Response: respStruct,
	}

	return j.JSONPb.Marshal(body)
}
