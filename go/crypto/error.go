package crypto

import "errors"

var (
	ErrKeyLength        = errors.New("err key lenght")
	ErrPaddingSize      = errors.New("err padding size")
	ErrCipherTextLength = errors.New("err cipherText lenght")
)
