package md5

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	os_ "github.com/kaydxh/golang/go/os"
)

func SumBytes(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func SumString(s string) string {
	return SumBytes([]byte(s))
}

func SumReader(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func SumReaderN(r io.Reader, n int64) (string, error) {
	h := md5.New()
	if _, err := io.CopyN(h, r, n); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func SumReaderAt(r io.ReaderAt, offset, length int64) (string, error) {
	h := md5.New()
	buf := make([]byte, 1024)
	var total int64

	for total < length {
		n, err := r.ReadAt(buf, offset)
		if err == nil || err == io.EOF {
			offset += int64(n)
			total += int64(n)
			_, tmpErr := io.CopyN(h, bytes.NewReader(buf), int64(n))
			if tmpErr != nil {
				return "", tmpErr
			}

			if err == io.EOF {
				break
			}

		} else {
			return "", err
		}

	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func SumFile(fileName string) (string, error) {
	file, err := os_.OpenFile(fileName, true)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return SumReader(file)
}

func SumFileAt(fileName string, offset int64, length int64) (string, error) {
	file, err := os_.OpenAll(fileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return SumReaderAt(file, offset, length)
}
