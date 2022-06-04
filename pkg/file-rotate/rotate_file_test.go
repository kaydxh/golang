package rotatefile_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
	"github.com/sirupsen/logrus"
)

func getWdOrDie() string {
	path, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("failed to get wd, err: %v", err)
	}

	return path
}

func TestRotateFileWithInterval(t *testing.T) {

	rotateFiler, _ := rotate_.NewRotateFiler(
		filepath.Join(getWdOrDie(), "log"),
		rotate_.WithRotateInterval(time.Second),
		rotate_.WithSuffixName(".log"),
		rotate_.WithPrefixName(filepath.Base(os.Args[0])),
	)

	for i := 0; i < 10; i++ {
		n, err := rotateFiler.Write([]byte("hello word"))
		if err != nil {
			t.Errorf("faild to write, err: %v", err)
		}
		time.Sleep(1 * time.Second)
		t.Logf("successed to write %v bytes", n)
	}

}

func TestRotateFileWithIntervalAndSize(t *testing.T) {

	rotateFiler, _ := rotate_.NewRotateFiler(
		filepath.Join(getWdOrDie(), "log"),
		rotate_.WithRotateInterval(time.Hour),
		rotate_.WithRotateSize(15),
		rotate_.WithSuffixName(".log"),
		rotate_.WithPrefixName(filepath.Base(os.Args[0])),
	)

	for i := 0; i < 0; i++ {
		n, err := rotateFiler.Write([]byte("hello word"))
		if err != nil {
			t.Errorf("faild to write, err: %v", err)
		}
		time.Sleep(time.Second)
		t.Logf("successed to write %v bytes", n)
	}

}

func TestRotateFileWithSize(t *testing.T) {

	rotateFiler, _ := rotate_.NewRotateFiler(
		filepath.Join(getWdOrDie(), "log"),
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

}

func TestRotateMaxCount(t *testing.T) {

	rotateFiler, _ := rotate_.NewRotateFiler(
		filepath.Join(getWdOrDie(), "log"),
		rotate_.WithRotateSize(15),
		rotate_.WithMaxCount(5),
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
