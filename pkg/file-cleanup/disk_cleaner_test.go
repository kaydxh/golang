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
package filecleanup_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	viper_ "github.com/kaydxh/golang/pkg/viper"

	filecleanup_ "github.com/kaydxh/golang/pkg/file-cleanup"
)

func diskUsageCallBack(diskPath string, diskUsage float32) {
	fmt.Printf("diskUsageCallBack diskPath: %v, diskUsage: %v\n", diskPath, diskUsage)
}

func TestDiskCleanerSerivce(t *testing.T) {
	cfgFile := "./diskcleaner.yaml"
	config := filecleanup_.NewConfig(
		filecleanup_.WithViper(viper_.GetViper(cfgFile, "diskcleaner")),
		filecleanup_.WithDiskUsageCallBack(diskUsageCallBack),
	)
	s, err := config.Complete().New(context.Background())
	if err != nil {
		t.Fatalf("failed to new config err: %v", err)
	}
	s.Run(context.Background())

	time.Sleep(1 * time.Minute)

}
