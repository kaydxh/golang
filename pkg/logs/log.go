package logs

import (
	"context"
	"fmt"
	"time"

	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
	trace_ "github.com/kaydxh/golang/pkg/middleware/http-middleware/trace"
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

func WithRotate(log *logrus.Logger, filedir string, options ...RotateOption) error {
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
		_, err = rotateFiler.Write([]byte(msg))
		return err
	}))

	return nil
}

func GetLogger(ctx context.Context) *logrus.Entry {
	logger := logrus.WithField("request_id", trace_.ExtractRequestIDFromContext(ctx))
	return logger
}

func GetLoggerOrFallback(ctx context.Context, defaultValue string) *logrus.Entry {
	requestId := trace_.ExtractRequestIDFromContext(ctx)
	if requestId == "" {
		requestId = defaultValue
	}
	logger := logrus.WithField("request_id", requestId)
	return logger
}
