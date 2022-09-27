/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
