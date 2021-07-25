package webserver_test

import (
	"fmt"
	"testing"

	"context"

	webserver_ "github.com/kaydxh/golang/pkg/webserver"
	"github.com/ory/viper"
)

func TestNew(t *testing.T) {
	/*
		viper.SetConfigName("webserver")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	*/
	viper.SetConfigFile("./webserver.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("failed to read config err: %v", err)
		return
	}
	subv := viper.Sub("web")
	config := webserver_.NewConfig(webserver_.WithGetViper(func() *viper.Viper {
		return subv //viper.GetViper()
	}))

	//fmt.Println(viper.GetViper())
	fmt.Println(subv)

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
