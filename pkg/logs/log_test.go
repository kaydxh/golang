package logs_test

import (
	"os"
	"testing"

	viper_ "github.com/kaydxh/golang/pkg/viper"

	logs_ "github.com/kaydxh/golang/pkg/logs"
	"github.com/sirupsen/logrus"
)

func TestInit(t *testing.T) {
	cfgFile := "./log.yaml"
	config := logs_.NewConfig(logs_.WithViper(viper_.GetViper(cfgFile, "log")))
	err := config.Complete().Apply()
	if err != nil {
		t.Errorf("failed to apply log config err: %v", err)
	}

	logrus.WithField(
		"module",
		os.Args,
	).WithField(
		"log_dir",
		config.Proto.GetFilepath(),
	).Infof(
		"successed to apply log config",
	)

}
