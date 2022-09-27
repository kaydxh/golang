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
package webserver_test

import (
	"fmt"
	"testing"

	"context"

	prototext_ "github.com/kaydxh/golang/pkg/protobuf/prototext"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	webserver_ "github.com/kaydxh/golang/pkg/webserver"
)

func TestNew(t *testing.T) {
	/*
		viper.SetConfigName("webserver")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	*/
	/*
		viper.SetConfigFile("./webserver.yaml")
		err := viper.ReadInConfig()
		if err != nil {
			t.Errorf("failed to read config err: %v", err)
			return
		}
		subv := viper.Sub("web")
	*/

	cfgFile := "./webserver.yaml"
	config := webserver_.NewConfig(webserver_.WithViper(viper_.GetViper(cfgFile, "web")))

	s, err := config.Complete().New()
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	s.InstallWebHandlers()
	prepared, err := s.PrepareRun()
	if err != nil {
		t.Errorf("failed to PrepareRun err: %v", err)
	}

	prepared.Run(context.Background())
}

func TestPrintProto(t *testing.T) {
	config := &webserver_.Web{
		BindAddress: &webserver_.Web_Net{
			Host: "192.168.0.1",
		},
	}
	s := prototext_.FormatWithLength(config, 30)
	//s := prototext.MarshalOptions{}.Format(config)
	fmt.Printf("%s\n", s)
}
