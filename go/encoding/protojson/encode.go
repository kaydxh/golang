package protojson

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Marshaler struct {
	// MarshalOptions is a configurable JSON format marshaler.
	protojson.MarshalOptions
}

func Marshal(pb proto.Message, options ...MarshalerOption) ([]byte, error) {
	m := Marshaler{
		MarshalOptions: protojson.MarshalOptions{
			AllowPartial: true,
		},
	}
	m.ApplyOptions(options...)

	return m.Marshal(pb)
}
