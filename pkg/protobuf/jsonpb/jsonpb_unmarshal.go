package jsonpb

import (
	"io"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
)

func UnmarshalWithAllowUnknownFields(r io.Reader, pb proto.Message) error {
	unmarshaler := &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	return unmarshaler.Unmarshal(r, pb)
}
