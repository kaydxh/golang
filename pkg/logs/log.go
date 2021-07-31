package logs

/*

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type Log_Format int32

const (
	Log_json Log_Format = 0
	Log_text Log_Format = 1
)

type Log_Level int32

const (
	Log_panic Log_Level = iota //0
	Log_fatal                  // 1
	Log_error                  // 2
	Log_warn                   // 3
	Log_info                   // 4
	Log_debug                  // 5
	Log_trace                  // 6
)

type loggerContextKeyType int

const loggerContextKey loggerContextKeyType = 0

// Logger is a wrapper for logrus.Entry.
type Logger struct {
	*logrus.Entry
}

// ContextLogger interface for components which support
// logging with context, via setting a logger to an exisiting one,
// thereby inheriting its context.
type ContextLogger interface {
	UseLog(l *Logger)
}

// InitLog initializes logs
func InitLog(log_format Log_Format, log_level Log_Level) {
	log.SetPrefix(fmt.Sprintf("[%s] ", os.Args[0]))
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	updateLogger(log_format, log_level)
}

func updateLogger(log_format Log_Format, log_level Log_Level) {
	if log_format == Log_text {
		Log.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
		}
	}
	Log.Level = logrus.Level(log_level)
	Log.Hooks.Add(ContextHook{})
}

type ContextHook struct {
}

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	//'skip' = 6 is the default call stack skip, which
	//works ootb when Error(), Warn(), etc. are called
	//for Errorf(), Warnf(), etc. - we have to skip 1 lvl up
	for skip := 6; skip < 8; skip++ {
		if pc, file, line, ok := runtime.Caller(skip); ok {
			funcName := runtime.FuncForPC(pc).Name()

			//detect if we're still in logrus (formatting funcs)
			if !strings.Contains(funcName, "github.com/sirupsen/logrus") {
				entry.Data["file"] = path.Base(file)
				entry.Data["func"] = path.Base(funcName)
				entry.Data["line"] = line
				break
			}
		}
	}

	return nil
}

// WithContext adds logger to context `ctx` and returns the resulting context.
func WithContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, log)
}
*/
