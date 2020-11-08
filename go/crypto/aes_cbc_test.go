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

package crypto_test

import (
	"fmt"
	"github.com/kaydxh/golang/go/crypto"
	"testing"
)

func TestAesCbcEncryptDecrypt(t *testing.T) {
	plainText := []byte("Hello World")
	fmt.Println("plainText: ", string(plainText))

	key := []byte("daW3eDgPEa9TjknE")
	cryptText, err := crypto.AesCbcEncrypt(plainText, key)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println("cryptText: ", string(cryptText))

	newplainText, err := crypto.AesCbcDecrypt(cryptText, key)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	fmt.Println("newplainText: ", string(newplainText))
	if string(newplainText) != string(plainText) {
		t.Errorf("newplainText not equal plainText")
	}
}
