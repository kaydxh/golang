package rotatefile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
)

func TestRotateFileWithInterval(t *testing.T) {

	filename := "/Users/kayxhding/workspace/studyspace/git-kayxhding/github.com/kaydxh/golang/pkg/file-rotate/log/"
	rotateFiler, _ := rotate_.NewRotateFiler(
		filename,
		rotate_.WithRotateInterval(time.Second),
		rotate_.WithSuffixName(".log"),
		rotate_.WithPrefixName(filepath.Base(os.Args[0])),
	)

	for i := 0; i < 10; i++ {
		n, err := rotateFiler.Write([]byte("hello word"))
		if err != nil {
			t.Errorf("faild to write, err: %v", err)
		}
		time.Sleep(time.Second)
		t.Logf("successed to write %v bytes", n)
	}

	//t.Logf("successed to write %v bytes", n)
}
