package logger_test

import (
	"testing"

	"github.com/kaydxh/golang/go/logger"
)

/*
func TestNew(t *testing.T) {
	l := logger.New(Ctx{"foo": "bar"})
	if nil == l {
		t.Error("expect not nil, got nil")
	}
}
*/

func TestInit(t *testing.T) {
	logger.InitLog(logger.Log_text, logger.Log_info)
	l := logger.Log.WithField("func", "init")
	l.Infof("info test: %v", "test info")
	l.Warnf("warn test: %v", "test info")
	l.Errorf("error test: %v", "test info")
	/*
		baselog.Level = logrus.PanicLevel
		baselog.Out = ioutil.Discard

		//	l := logger.NewFromLogger(baselog, Ctx{})
		assert.NotNil(t, baselog)
		assert.Equal(t, l.Logger.Level, logrus.PanicLevel)
		assert.Equal(t, l.Logger.Out, ioutil.Discard)
	*/
}
