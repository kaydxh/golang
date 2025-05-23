package fsnotify_test

import (
	"context"
	"testing"

	fsnotify_ "github.com/kaydxh/golang/pkg/fsnotify"
	viper_ "github.com/kaydxh/golang/pkg/viper"
)

func TestFsnotify(t *testing.T) {

	cfgFile := "./fsnotify.yaml"
	config := fsnotify_.NewConfig(fsnotify_.WithViper(viper_.GetViper(cfgFile, "fsnotify")))

	fn, err := config.Complete().New(context.Background())
	if err != nil {
		t.Fatalf("failed to new config err: %v", err)
	}
	fn.ApplyOptions(fsnotify_.WithWriteCallbackFunc(func(ctx context.Context, path string) {
		t.Logf("%s write call back", path)
	}))

	select {}
}
