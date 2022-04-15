package filecleanup_test

import (
	"context"
	"testing"
	"time"

	viper_ "github.com/kaydxh/golang/pkg/viper"

	filecleanup_ "github.com/kaydxh/golang/pkg/file-cleanup"
)

func TestDiskCleanerSerivce(t *testing.T) {
	cfgFile := "./diskcleaner.yaml"
	config := filecleanup_.NewConfig(filecleanup_.WithViper(viper_.GetViper(cfgFile, "diskcleaner")))
	s, err := config.Complete().New(context.Background())
	if err != nil {
		t.Fatalf("failed to new config err: %v", err)
	}
	s.Run(context.Background())

	time.Sleep(1 * time.Minute)

}
