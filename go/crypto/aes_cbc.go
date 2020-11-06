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
