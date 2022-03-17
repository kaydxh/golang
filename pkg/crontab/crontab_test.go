package crontab_test

import (
	"context"
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
		return nil
	})

	time.Sleep(1 * time.Minute)

}
