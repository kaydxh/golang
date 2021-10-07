package logs_test

import (
	"os"
	"testing"
	"time"

	viper_ "github.com/kaydxh/golang/pkg/viper"

	logs_ "github.com/kaydxh/golang/pkg/logs"
	"github.com/sirupsen/logrus"
)

func TestInit(t *testing.T) {
	cfgFile := "./log.yaml"
	config := logs_.NewConfig(logs_.WithViper(viper_.GetViper(cfgFile, "log")))
	err := config.Complete().Apply()
	if err != nil {
		t.Fatalf("failed to apply log config err: %v", err)
	}
	logrus.WithField(
		"module",
		os.Args,
	).WithField(
		"log_dir",
		config.Proto.GetFilepath(),
	).Infof(
		"successed to apply log config: %#v", config.Proto.String(),
	)

	for i := 1; i <= 10; i++ {
		logrus.Infof("test time: %v", i)
		time.Sleep(time.Second)
	}

}
