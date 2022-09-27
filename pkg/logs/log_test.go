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
package logs_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	viper_ "github.com/kaydxh/golang/pkg/viper"

	logs_ "github.com/kaydxh/golang/pkg/logs"
	"github.com/sirupsen/logrus"
)

func TestInit(t *testing.T) {
	cfgFile := "./log.yaml"
	config := logs_.NewConfig(logs_.WithViper(viper_.GetViper(cfgFile, "log")))
	err := config.Complete().Apply()
	if err != nil {
		t.Fatalf("failed to apply log config err: %v", err)
	}
	logrus.WithField(
		"module",
		os.Args,
	).WithField(
		"log_dir",
		config.Proto.GetFilepath(),
	).Infof(
		"successed to apply log config: %#v", config.Proto.String(),
	)

	for i := 1; i <= 10; i++ {
		logrus.Infof("test time: %v", i)
		fmt.Println("---------write to stdout----------")
		time.Sleep(time.Second)
	}

}
