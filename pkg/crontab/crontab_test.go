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
package crontab_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	crontab_ "github.com/kaydxh/golang/pkg/crontab"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
)

func TestCrontabSerivce(t *testing.T) {
	cfgFile := "./crontab.yaml"
	config := crontab_.NewConfig(crontab_.WithViper(viper_.GetViper(cfgFile, "crontab")))
	s, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	s.Run(context.Background())

	s.Register(func(ctx context.Context, logger *logrus.Entry) error {
		logger.Infof("doing...")
		filepath.Walk("./testdata", func(path string, info os.FileInfo, err error) error {

			if !info.IsDir() {
				now := time.Now()
				if now.Sub(info.ModTime()) > time.Minute {
					t.Logf("file[%v] expired 1 Minute, modify time: %v, now: %v", path, info.ModTime(), now)
				} else {
					t.Logf("file[%v] normal, modify time: %v, now: %v", path, info.ModTime(), now)

				}
			}

			return nil
		})

		return nil
	})

	time.Sleep(1 * time.Minute)

}
