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
package logs

import (
	"context"
	"fmt"
	"os"
	"time"

	http_ "github.com/kaydxh/golang/go/net/http"
	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
	"github.com/sirupsen/logrus"
)

type Rotate struct {
	maxAge         time.Duration
	maxCount       int64
	rotateSize     int64
	rotateInterval time.Duration
	prefixName     string
	suffixName     string
}

func WithRotate(log *logrus.Logger, filedir string, redirect Log_Redirct, options ...RotateOption) error {
	if log == nil {
		return fmt.Errorf("log is nil")
	}

	var rotate Rotate
	rotate.ApplyOptions(options...)

	rotateFiler, _ := rotate_.NewRotateFiler(
		filedir,
		rotate_.WithMaxAge(rotate.maxAge),
		rotate_.WithMaxCount(rotate.maxCount),
		rotate_.WithRotateInterval(rotate.rotateInterval),
		rotate_.WithRotateSize(rotate.rotateSize),
		rotate_.WithPrefixName(rotate.prefixName),
		rotate_.WithSuffixName(rotate.suffixName),
	)
	log.AddHook(HookHandler(func(entry *logrus.Entry) error {
		var (
			msg []byte
			err error
		)

		if log.Formatter == nil {
			msg, err = entry.Bytes()
		} else {
			msg, err = log.Formatter.Format(entry)
		}
		if err != nil {
			return err
		}
		//if opt.MuteDirectlyOutput && entry.Level <= logrus.WarnLevel {
		/*
			if entry.Level <= logrus.WarnLevel {
				if out != nil {
					_, _ = out.Write(msg)
				}
			}
		*/
		//https://eli.thegreenplace.net/2020/faking-stdin-and-stdout-in-go/
		//https://github.com/eliben/code-for-blog/blob/master/2020/go-fake-stdio/snippets/redirect-cgo-stdout.go
		file, _, err := rotateFiler.Write([]byte(msg))
		if err == nil {
			if redirect == Log_file {
				os.Stdout = file
				os.Stderr = file
			}
		}

		return err
	}))

	return nil
}

func GetLogger(ctx context.Context) *logrus.Entry {
	logger := logrus.WithField("request_id", http_.ExtractRequestIDFromContext(ctx))
	return logger
}

func GetLoggerOrFallback(ctx context.Context, defaultValue string) *logrus.Entry {
	requestId := http_.ExtractRequestIDFromContext(ctx)
	if requestId == "" {
		requestId = defaultValue
	}
	logger := logrus.WithField("request_id", requestId)
	return logger
}
