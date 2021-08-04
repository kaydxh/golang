package logs

import (
	"fmt"
	"time"

	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
	"github.com/sirupsen/logrus"
)

type Rotate struct {
	maxAge         time.Duration
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
		rotate_.WithRotateSize(rotate.rotateSize),
		rotate_.WithRotateInterval(rotate.rotateInterval),
		rotate_.WithPrefixName(rotate.prefixName),
		rotate_.WithSuffixName(rotate.suffixName),
	)
	log.AddHook(HookHandler(func(entry *logrus.Entry) error {
		//	var msg []byte
		//	var err error

		msg, err := entry.String()
		/*
			if log.Formatter == nil {
				msg_, err_ := entry.String()
				msg, err = []byte(msg_), err_
			} else {
				switch f := log.Formatter.(type) {
				case *logrus.TextFormatter:
					var disableColors = f.DisableColors
					// disable colors in log file
					f.DisableColors = true
					msg, err = log.Formatter.Format(entry)
					f.DisableColors = disableColors
				default:
					msg, err = log.Formatter.Format(entry)
				}
			}
		*/

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
