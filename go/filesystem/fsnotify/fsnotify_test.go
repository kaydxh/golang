package fsnotify_test

import (
	"context"
	"testing"

	fsnotify_ "github.com/kaydxh/golang/go/filesystem/fsnotify"
	viper_ "github.com/kaydxh/golang/pkg/viper"
)

func TestFsnotify(t *testing.T) {

	cfgFile := "./fsnotify.yaml"
	config := fsnotify_.NewConfig(fsnotify_.WithViper(viper_.GetViper(cfgFile, "fsnotify")))

	fn, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	_ = fn
	select {}
}
