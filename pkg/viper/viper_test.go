package viper_test

import (
	"fmt"
	"testing"

	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/kaydxh/golang/pkg/viper/code"
)

func TestViper(t *testing.T) {
	v := viper_.GetViper("./error.yaml", "")
	var ErrorTemplateConfig code.ErrorTemplate
	err := viper_.UnmarshalProtoMessageWithJsonPb(v, &ErrorTemplateConfig)
	if err != nil {
		return
	}

	fmt.Printf("ErrorTemplateConfig: %v\n", ErrorTemplateConfig.String())

}
