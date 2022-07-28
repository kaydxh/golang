package v1

import (
	apijsonpb_ "github.com/kaydxh/golang/pkg/middleware/api/jsonpb"
	jsonpb_ "github.com/kaydxh/golang/pkg/protobuf/jsonpb"
)

type JSONPb struct {
	*apijsonpb_.JSONPb
}

func NewDefaultJSONPb() *JSONPb {
	return &JSONPb{
		apijsonpb_.NewJSONPb(
			apijsonpb_.WithUseProtoNames(false),
			apijsonpb_.WithUseEnumNumbers(false),
			apijsonpb_.WithEmitUnpopulated(true),
			apijsonpb_.WithDiscardUnknown(true),
			apijsonpb_.WithIndent("\t"),
		),
	}
}

// Marshal marshals "v" into JSON.
func (j *JSONPb) Marshal(v interface{}) ([]byte, error) {
	respStruct, err := jsonpb_.MarshaToStructpb(v)
	if err != nil {
		return nil, err
	}

	return j.JSONPb.Marshal(respStruct)
}
