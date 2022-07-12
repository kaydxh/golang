package filecleanup_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	viper_ "github.com/kaydxh/golang/pkg/viper"

	filecleanup_ "github.com/kaydxh/golang/pkg/file-cleanup"
)

func diskUsageCallBack(diskUsage float32) {
	fmt.Printf("diskUsageCallBack diskUsage: %v\n", diskUsage)
}

func TestDiskCleanerSerivce(t *testing.T) {
	cfgFile := "./diskcleaner.yaml"
	config := filecleanup_.NewConfig(filecleanup_.WithViper(viper_.GetViper(cfgFile, "diskcleaner")), filecleanup_.WithDiskUsageCallBack(diskUsageCallBack))
	s, err := config.Complete().New(context.Background())
	if err != nil {
		t.Fatalf("failed to new config err: %v", err)
	}
	s.Run(context.Background())

	time.Sleep(1 * time.Minute)

}
