package tcloud

import (
	marshaler_ "github.com/kaydxh/golang/go/runtime/marshaler"
	jsonpb_ "github.com/kaydxh/golang/pkg/protobuf/jsonpb"
)

type JSONPb struct {
	*marshaler_.JSONPb
}

func NewDefaultJSONPb() *JSONPb {
	return &JSONPb{
		marshaler_.NewJSONPb(
			marshaler_.WithUseProtoNames(false),
			marshaler_.WithUseEnumNumbers(false),
			marshaler_.WithEmitUnpopulated(true),
			marshaler_.WithDiscardUnknown(true),
			marshaler_.WithIndent("\t"),
		),
	}
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
