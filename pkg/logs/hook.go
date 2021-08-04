package logs

import "github.com/sirupsen/logrus"

type HookHandler func(entry *logrus.Entry) error

func (h HookHandler) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h HookHandler) Fire(entry *logrus.Entry) error {
	return h(entry)
}

/*
type ContextHook struct {
}

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	data, err := entry.String()
	if err != nil {
		return err
	}

}
*/
