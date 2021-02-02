package config_test

import (
	"fmt"
	"testing"

	"github.com/kaydxh/golang/go/config"
)

func TestNew(t *testing.T) {
	conf := config.New(config.WithConfigOptionPath("/data/conf"))
	fmt.Println(conf)
}
