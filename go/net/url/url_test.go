package url_test

import (
	"fmt"
	"testing"

	url_ "github.com/kaydxh/golang/go/net/url"
	"golang.org/x/net/context"
)

func TestUrlEncode(t *testing.T) {
	type SimpleChild struct {
		Status bool
		Name   string
	}

	type SimpleData struct {
		Id           int
		Name         string
		Child        SimpleChild
		ParamsInt8   map[string]int8
		ParamsString map[string]string
		Array        [3]uint16
	}

	data := SimpleData{
		Id:   2,
		Name: "http://localhost/test.php?id=2",
		Child: SimpleChild{
			Status: true,
		},
		ParamsInt8: map[string]int8{
			"one": 1,
		},
		ParamsString: map[string]string{
			"two": "你好",
		},

		Array: [3]uint16{2, 3, 300},
	}

	c, _ := url_.New(context.Background())
	bytes, err := c.Encode(data)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	fmt.Println(string(bytes))
}
