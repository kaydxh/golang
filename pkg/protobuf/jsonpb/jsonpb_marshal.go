package jsonpb

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

func MarshaToStructpb(v interface{}) (*structpb.Struct, error) {
	var jb []byte
	switch v := v.(type) {
	case nil:
		return &structpb.Struct{}, nil

	case *structpb.Struct:
		return v, nil

	case proto.Message:
		m := jsonpb.Marshaler{EmitDefaults: true}
		data, err := m.MarshalToString(v)
		if err != nil {
			return nil, fmt.Errorf("failed to Marshal json: %v", err)
		}
		jb = []byte(data)

	case []byte:
		jb = v
	case *[]byte:
		jb = *v
	case string:
		jb = []byte(v)
	case *string:
		jb = []byte(*v)
	default:
		var err error
		jb, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to Marshal json: %v", err)
		}
	}

	var dataStructpb structpb.Struct
	if err := UnmarshalWithAllowUnknownFields(bytes.NewReader(jb), &dataStructpb); err != nil {
		return nil, err
	}
	return &dataStructpb, nil
}
