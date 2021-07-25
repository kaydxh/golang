package webserver_test

import (
	"testing"

	"context"

	viper_ "github.com/kaydxh/golang/pkg/viper"
	webserver_ "github.com/kaydxh/golang/pkg/webserver"
)

func TestNew(t *testing.T) {
	/*
		viper.SetConfigName("webserver")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	*/
	/*
		viper.SetConfigFile("./webserver.yaml")
		err := viper.ReadInConfig()
		if err != nil {
			t.Errorf("failed to read config err: %v", err)
			return
		}
		subv := viper.Sub("web")
	*/

	cfgFile := "./webserver.yaml"
	config := webserver_.NewConfig(webserver_.WithViper(viper_.GetViper(cfgFile, "web")))

	s, err := config.Complete().New()
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	s.InstallWebHandlers()
	prepared, err := s.PrepareRun()
	if err != nil {
		t.Errorf("failed to PrepareRun err: %v", err)
	}

	prepared.Run(context.Background())
}
