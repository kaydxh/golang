package prototext

import (
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// Format formats the message as a multiline string and only return
// first n characters.
// This function is only intended for human consumption and ignores errors.
// Do not depend on the output being stable. It may change over time across
// different versions of the program.
func FormatWithLength(m proto.Message, n int) string {
	s := prototext.Format(m)
	if n >= len(s) || n == 0 {
		return s
	}
	return s[0:n]
}

func Marshal(m proto.Message) ([]byte, error) {
	return prototext.MarshalOptions{
		AllowPartial: true,
		EmitUnknown:  true,
	}.Marshal(m)
}
