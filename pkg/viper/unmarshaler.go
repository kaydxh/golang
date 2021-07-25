package viper

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/mapstructure"
	"github.com/ory/viper"
)

func UnmarshalProtoMessageWithJsonPb(v *viper.Viper, msg proto.Message, options ...viper.DecoderConfigOption) error {

	var opts []viper.DecoderConfigOption
	opts = append(opts, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json" // trick of protobuf, which generates json tag only
		decoderConfig.WeaklyTypedInput = true
		decoderConfig.DecodeHook = UnmarshalProtoMessageWithJsonpbHookFunc(msg)
	})
	opts = append(opts, options...)
	return v.Unmarshal(msg, opts...)
}

func UnmarshalProtoMessageWithJsonpbHookFunc(v proto.Message) mapstructure.DecodeHookFunc {
	return func(src reflect.Type, dst reflect.Type, data interface{}) (interface{}, error) {
		// Convert it by parsing
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		err = jsonpb.Unmarshal(bytes.NewReader(dataBytes), v)
		if err != nil {
			return data, err
		}

		return v, nil
	}
}
