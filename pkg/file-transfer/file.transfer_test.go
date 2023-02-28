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
package filetransfer_test

import (
	"context"
	"fmt"
	"testing"

	filetransfer_ "github.com/kaydxh/golang/pkg/file-transfer"
	viper_ "github.com/kaydxh/golang/pkg/viper"
)

func TestDownload(t *testing.T) {
	cfgFile := "./ft.yaml"
	config := filetransfer_.NewConfig(filetransfer_.WithViper(viper_.GetViper(cfgFile, "filetransfer")))
	ft, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	testCases := []struct {
		Url string
	}{
		{
			Url: "https://ai-media-1256936300.cos.ap-guangzhou.myqcloud.com/find.sh?q-sign-algorithm=sha1&q-ak=AKIDCDyve81SJuISPkMq0IukLg7tupWyoqCg&q-sign-time=1659955959;8640000001659870000&q-key-time=1659955959;8640000001659870000&q-header-list=&q-url-param-list=&q-signature=5695f17ee30c3cd2d37197c773a445eda8a70c8c",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			data, err := ft.Download(context.Background(), testCase.Url)
			if err != nil {
				t.Fatalf("failed to download: %v, got : %s", testCase.Url, err)
			}
			t.Logf("data len: %v", len(data))
			t.Logf("data : %v", string(data))
		})
	}

}

func TestDownloadByProxy(t *testing.T) {
	cfgFile := "./ft.yaml"
	config := filetransfer_.NewConfig(filetransfer_.WithViper(viper_.GetViper(cfgFile, "filetransfer")))
	ft, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	testCases := []struct {
		Url string
	}{
		{
			Url: "http://quyujiaofu-new-1300074211.cos.ap-guangzhou.myqcloud.com/hk_test/480p.jpg",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			data, err := ft.Download(context.Background(), testCase.Url)
			if err != nil {
				t.Fatalf("failed to download: %v, got : %s", testCase.Url, err)
			}
			t.Logf("data len: %v", len(data))
			t.Logf("data : %v", string(data))
		})
	}

}
