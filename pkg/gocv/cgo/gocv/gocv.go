package gocv

import (
	"os"

	"github.com/kaydxh/golang/pkg/gocv/cgo/api/openapi-spec/gocv"
)

func NewMagickInitializeMagickRequest() *gocv.MagickInitializeMagickRequest {
	return &gocv.MagickInitializeMagickRequest{
		Path: os.Args[0],
	}
}
