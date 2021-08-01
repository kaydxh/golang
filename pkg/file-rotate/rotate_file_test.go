package rotatefile_test

import (
	"testing"
	"time"

	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
)

func TestRotateFileWithInterval(t *testing.T) {

	filename := "/Users/kayxhding/workspace/studyspace/git-kayxhding/github.com/kaydxh/golang/pkg/file-rotate/log/"
	rotateFiler, _ := rotate_.NewRotateFiler(filename, rotate_.WithRotateInterval(time.Minute))
	n, err := rotateFiler.Write([]byte("hello word"))
	if err != nil {
		t.Errorf("faild to write, err: %v", err)
	}

	t.Logf("successed to write %v bytes", n)
}
