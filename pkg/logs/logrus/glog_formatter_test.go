package logrus_test

import (
	"testing"

	logrus_ "github.com/kaydxh/golang/pkg/logs/logrus"
	"github.com/sirupsen/logrus"
)

func TestGlogFormatting(t *testing.T) {
	tf := &logrus_.GlogFormatter{
		DisableColors:     true,
		EnableGoroutineId: true,
	}

	testCases := []struct {
		value    string
		expected string
	}{
		{`foo`, "time=\"0001-01-01T00:00:00Z\" level=panic test=foo\n"},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(logrus.WithField("test", tc.value))

		if string(b) != tc.expected {
			t.Errorf("formatting expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}
