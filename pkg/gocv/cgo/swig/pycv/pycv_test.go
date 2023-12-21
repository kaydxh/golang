package pycv_test

import (
	"testing"

	"github.com/kaydxh/golang/pkg/gocv/cgo/swig/pycv"
	pycv_ "github.com/kaydxh/golang/pkg/gocv/cgo/swig/pycv"
)

func TestDo(t *testing.T) {
	err := pycv_.GlobalInit("model_dir", -1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	defer pycv.GlobalRelease()

	sdk := pycv_.NewPyImage()
	initReq := pycv_.NewLocalInitRequest()
	_, err = sdk.LocalInit(initReq)
	if err != nil {
		t.Errorf("failed to local init, err: %v", err)
	}
	req := pycv_.NewDoRequest()
	req.SetArg1("arg1")
	req.SetArg2("arg2")
	_, err = sdk.Do(req)
	if err != nil {
		t.Error(err.Error())
		return
	}
}
