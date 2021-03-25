package md5_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	md5_ "github.com/kaydxh/golang/go/crypto/md5"
	"gotest.tools/assert"
)

func TestMd5File(t *testing.T) {
	file, err := ioutil.TempFile(".", "file")
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	defer os.RemoveAll(file.Name())

	strContext := "hello world"
	//	n, err := io.WriteString(file, "hello world")
	n, err := file.Write([]byte(strContext))
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	assert.Equal(t, len(strContext), n)
	fmt.Printf("fileName: %v, n: %v\n", file.Name(), n)

	sum, err := md5_.SumFile(file.Name())
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	assert.Equal(t, sum, md5_.SumString(strContext))
}
