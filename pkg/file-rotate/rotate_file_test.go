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
package rotatefile_test

import (
	"context"
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
		_, n, err := rotateFiler.Write([]byte("hello word"))
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
		_, n, err := rotateFiler.Write([]byte("hello word"))
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
		rotate_.WithRotateSize(10),
		rotate_.WithRotateInterval(15*time.Second),
		rotate_.WithSuffixName(".log"),
		rotate_.WithPrefixName(filepath.Base(os.Args[0])),
		rotate_.WithRotateCallback(func(ctx context.Context, path string) {
			t.Logf("=======callback path: %v", path)
		}),
	)

	for i := 0; i < 5; i++ {
		//_, n, err := rotateFiler.Write([]byte("hello word"))
		_, n, err := rotateFiler.WriteBytesLine([][]byte{[]byte("hello word")})
		if err != nil {
			t.Errorf("faild to write, err: %v", err)
		}
		time.Sleep(time.Second)
		t.Logf("successed to write %v bytes", n)
	}

	select {}

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
		_, n, err := rotateFiler.Write([]byte("hello word"))
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
