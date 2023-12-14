package gocv_test

import (
	"testing"

	io_ "github.com/kaydxh/golang/go/io"
	gocv_ "github.com/kaydxh/golang/pkg/gocv/cgo/swig/gocv"
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

	sdk := gocv_.NewMagicImage()
	initReq := gocv_.NewMagickInitializeMagickRequest()
	_, err = gocv_.MagickInitializeMagick(initReq)
	if err != nil {
		t.Error(err.Error())
		return
	}

	req := gocv_.NewMagickImageDecodeRequest()
	req.SetImage(string(data))
	resp, err := sdk.MagickImageDecode(req)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("resp row: %d, colunms: %v, magick: %v", resp.GetRows(), resp.GetColumns(), resp.GetMagick())
}
