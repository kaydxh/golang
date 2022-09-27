/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package gocv

/*
#cgo pkg-config: protobuf opencv4 graphics-magick
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

	unsafe_ "github.com/kaydxh/golang/go/unsafe"
	gocvpb "github.com/kaydxh/golang/pkg/gocv/cgo/api/openapi-spec/gocv"
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

		C.sdk_gocv_magick_initialize_magick(
			unsafe_.BytesPointer(reqData),
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

func MagickImageDecode(req *gocvpb.MagickImageDecodeRequest) (*gocvpb.MagickImageDecodeResponse, error) {

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	var respData *C.char
	var respDataLen C.int
	defer func() {
		C.free(unsafe.Pointer(respData))
		respData = nil
	}()

	C.sdk_gocv_magick_image_decode(unsafe_.BytesPointer(reqData), C.int(len(reqData)), &respData, &respDataLen)

	var resp gocvpb.MagickImageDecodeResponse
	err = proto.Unmarshal(C.GoBytes(unsafe.Pointer(respData), C.int(respDataLen)), &resp)
	if err != nil {
		return nil, err
	}
	if resp.GetError() != nil {
		return nil, resp.GetError()
	}
	return &resp, nil
}
