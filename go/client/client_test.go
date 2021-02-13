package client_test

import (
	"fmt"
	"testing"

	"github.com/kaydxh/golang/go/client"
)

func TestNew(t *testing.T) {
	conf := client.New(client.WithClientOptionPath("/data/conf"))
	fmt.Println(conf)
}
