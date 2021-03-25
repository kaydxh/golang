package md5

import (
	"crypto/md5"
	"io"

	os_ "github.com/kaydxh/golang/go/os"
)

func SumBytes(b []byte) string {
	h := md5.New()
	h.Write(b)
	return string(h.Sum(nil))
}

func SumString(s string) string {
	return SumBytes([]byte(s))
}

func SumReader(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}

func SumFile(fileName string) (string, error) {
	file, err := os_.OpenAllAt(fileName, true)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return SumReader(file)
}
