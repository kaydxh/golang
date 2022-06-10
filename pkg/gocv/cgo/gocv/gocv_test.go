package gocv_test

import (
	"testing"

	io_ "github.com/kaydxh/golang/go/io"
	gocv_ "github.com/kaydxh/golang/pkg/gocv/cgo/gocv"
)

func TestMagickInitializeMagick(t *testing.T) {
}

func TestMagickImageDecode(t *testing.T) {
	filename := "testdata/test.jpg"
	data, err := io_.ReadFile(filename)
	if err != nil {
		t.Error("Invalid ReadFile in TestMagickImageDecode")
		return
	}
	t.Logf("data size: %v", len(data))
	req := gocv_.NewMagickImageDecodeRequest()
	req.Image = data
	resp, err := gocv_.MagickImageDecode(req)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("resp: %v", resp)
}
