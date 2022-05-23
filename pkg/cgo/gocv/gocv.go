package gocv

import (
	"os"

	v1 "github.com/kaydxh/golang/pkg/cgo/api/openapi-spec/gocv/v1"
)

func NewMagickInitializeMagickRequest() *v1.MagickInitializeMagickRequest {
	return &v1.MagickInitializeMagickRequest{
		Path: os.Args[0],
	}
}
