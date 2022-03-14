package reslover_test

import (
	"context"
	"testing"
	"time"

	reslover_ "github.com/kaydxh/golang/pkg/reslover"
	viper_ "github.com/kaydxh/golang/pkg/viper"
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

	cfgFile := "./reslover.yaml"
	config := reslover_.NewConfig(reslover_.WithViper(viper_.GetViper(cfgFile, "reslover")))

	c := config.Complete()
	s, err := c.New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	s.Run(context.Background())
	domain := "www.baid.com"
	rq := reslover_.ResloverQuery{
		Domain: domain,
	}
	rq.SetDefault()
	s.AddService(rq)
	time.Sleep(2 * time.Second)
	node, has := s.PickNode(domain, "1")
	t.Logf("has: %v, node: %v", has, node)
}
