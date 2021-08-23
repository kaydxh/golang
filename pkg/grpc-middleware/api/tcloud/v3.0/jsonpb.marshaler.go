package tcloud

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsonpb_ "github.com/kaydxh/golang/pkg/protobuf/jsonpb"
)

type JSONPb struct {
	runtime.JSONPb
	/*
		opts struct {
			useProtoNames   bool
			useEnumNumbers  bool
			emitUnpopulated bool
		}
	*/
}

func NewJSONPb() *JSONPb {
	j := &JSONPb{}
	j.UseProtoNames = false
	j.EmitUnpopulated = true
	j.DiscardUnknown = true
	j.Indent = "\t"
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

	fmt.Println("jsonpb: ", j.JSONPb)
	errResponse := &ErrorResponse{
		Error: &TCloudError{},
	}
	data, _ := j.JSONPb.Marshal(errResponse)
	fmt.Println("data: ", string(data))

	/*
		mash := &jsonpb.Marshaler{}

		data, err := mash.MarshalToString(body)
		fmt.Println("data ", data)
		return []byte(data), err
	*/

	return j.JSONPb.Marshal(body)
}
