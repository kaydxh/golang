package crypto

import "errors"

var (
	ErrKeyLength   = errors.New("err key lenght")
	ErrPaddingSize = errors.New("err padding size")
)
