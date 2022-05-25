package gocv

/*
#cgo pkg-config: protobuf graphics-magick
#cgo CXXFLAGS: -std=c++11 -I${SRCDIR}/..
#cgo LDFLAGS: -L${SRCDIR}/../api/openapi-spec/gocv/ -lproto-gocv
#cgo LDFLAGS: -L${SRCDIR}/../api/openapi-spec/types/ -lproto-types
#include <stdlib.h>
#include "magick.h"
*/
import "C"

import (
	"sync"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

var once sync.Once

func init() {
	err := MagickInitializeMagick(nil)
	if err != nil {
		panic(err)
	}
}

func MagickInitializeMagick(req *gocvpb.MagickInitializeMagickRequest) error {
	var errOnce error
	once.Do(func() {
		reqData, err := proto.Marshal(req)
		if err != nil {
			errOnce = err
			return
		}

		var respData *C.char
		var respDataLen C.int
		defer func() {
			C.free(unsafe.Pointer(respData))
			respData = nil
		}()

		C.sdk_go_youtu_gocv_magick_initialize_magick(
			runtime_.BytesPointer(reqData),
			C.int(len(reqData)),
			&respData,
			&respDataLen,
		)

		var resp gocvpb.MagickInitializeMagickResponse
		err = proto.Unmarshal(C.GoBytes(unsafe.Pointer(respData), C.int(respDataLen)), &resp)
		if err != nil {
			errOnce = err
			return
		}
		if resp.GetError() != nil {
			errOnce = resp.GetError()
			return
		}
		errOnce = nil
		return
	})

	return errOnce
}
