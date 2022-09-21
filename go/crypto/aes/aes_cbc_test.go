/*
MIT License

Copyright (c) 2020 kay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package aes_test

import (
	"fmt"
	"testing"

	"flag"

	aes_ "github.com/kaydxh/golang/go/crypto/aes"
	io_ "github.com/kaydxh/golang/go/io"
)

func TestAesCbcEncryptDecrypt(t *testing.T) {
	plainText := []byte("Hello World")
	fmt.Println("plainText: ", string(plainText))

	key := []byte("daW3eDgPEa9TjknE")
	cryptText, err := aes_.AesCbcEncrypt(plainText, key)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println("cryptText: ", string(cryptText))

	newplainText, err := aes_.AesCbcDecrypt(cryptText, key)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println("newplainText: ", string(newplainText))
	if string(newplainText) != string(plainText) {
		t.Errorf("newplainText not equal plainText")
	}
}

// build: go test -c  aes_cbc_test.go -o decrypt_tool
// run: ./decrypt_tool -test.v  --test.run TestAesCbcDecrypt -srcePath=1.jpg -destPath=2.jpg
var (
	srcePath = flag.String("srcePath", "", "source encrypt image path")
	destPath = flag.String("destPath", "./decryptImage.jpg", "dest decrypt image path")
)

func TestAesCbcDecrypt(t *testing.T) {
	cryptText, err := io_.ReadFile(*srcePath)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
		return
	}

	key := []byte("daW3eDgPEa9TjknE")
	newplainText, err := aes_.AesCbcDecrypt(cryptText, key)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	io_.WriteFile(*destPath, newplainText, false)
}
