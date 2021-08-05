package rotatefile_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func TestRotateFileWithIntervalAndSize(t *testing.T) {

	filename := "/Users/kayxhding/workspace/studyspace/git-kayxhding/github.com/kaydxh/golang/pkg/file-rotate/log/"
	rotateFiler, _ := rotate_.NewRotateFiler(
		filename,
		rotate_.WithRotateInterval(time.Hour),
		rotate_.WithRotateSize(15),
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

func TestRotateFileWithSize(t *testing.T) {

	filename := "/Users/kayxhding/workspace/studyspace/git-kayxhding/github.com/kaydxh/golang/pkg/file-rotate/log/"
	rotateFiler, _ := rotate_.NewRotateFiler(
		filename,
		rotate_.WithRotateSize(15),
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

func TestRotateMaxCount(t *testing.T) {

	filename := "/Users/kayxhding/workspace/studyspace/git-kayxhding/github.com/kaydxh/golang/pkg/file-rotate/log/"
	rotateFiler, _ := rotate_.NewRotateFiler(
		filename,
		rotate_.WithRotateSize(15),
		rotate_.WithMaxCount(2),
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

func TestRegex(t *testing.T) {
	var regexps = []*regexp.Regexp{
		regexp.MustCompile(`%[%+A-Za-z]`),
		regexp.MustCompile(`\*+`),
	}
	globPattern := "1%%%AA20160304"
	for _, re := range regexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
		fmt.Printf("re: %v , globPattern: %v\n", re, globPattern)
	}
	//	return globPattern + `*`
}
