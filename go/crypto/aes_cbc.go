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

package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func AesCbcEncrypt(plainText, key []byte) ([]byte, error) {
	if len(key) != 16 && len(key) != 32 {
		return nil, ErrKeyLength
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	paddingText := pad(plainText, aes.BlockSize)
	cipherText := make([]byte, aes.BlockSize+len(paddingText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherText[aes.BlockSize:], paddingText)
	return cipherText, nil
}

func AesCbcDecrypt(cipherText, key []byte) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, ErrKeyLength
	}

	if len(cipherText) < aes.BlockSize {
		return nil, ErrCipherTextLength
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := cipherText[:aes.BlockSize]
	cb := cipherText[aes.BlockSize:]

	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(cb, cb)
	plainText, err := unpad(cb)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func pad(plainText []byte, blockSize int) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	newText := append(plainText, padText...)
	return newText
}

func unpad(plainText []byte) ([]byte, error) {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	if unpadding >= length {
		return nil, ErrPaddingSize
	}

	return plainText[:length-unpadding], nil
}
