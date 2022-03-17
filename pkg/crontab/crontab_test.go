package crontab_test

import (
	"context"
	"os"
	"path/filepath"
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
		filepath.Walk("./testdata", func(path string, info os.FileInfo, err error) error {

			if !info.IsDir() {
				now := time.Now()
				if now.Sub(info.ModTime()) > time.Minute {
					t.Logf("file[%v] expired 1 Minute, modify time: %v, now: %v", path, info.ModTime(), now)
				} else {
					t.Logf("file[%v] normal, modify time: %v, now: %v", path, info.ModTime(), now)

				}
			}

			return nil
		})

		return nil
	})

	time.Sleep(1 * time.Minute)

}
