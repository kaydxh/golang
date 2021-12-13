package protojson

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Unmarshaler struct {
	// UnmarshalOptions is a configurable JSON format parser.
	protojson.UnmarshalOptions
}

func Unmarshal(data []byte, pb proto.Message, options ...UnmarshalerOption) error {
	m := Unmarshaler{
		UnmarshalOptions: protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		},
	}
	m.ApplyOptions(options...)

	return m.Unmarshal(data, pb)
}
