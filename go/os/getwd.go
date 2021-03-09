package os

import "os"

func Getwd() (string, error) {
	return os.Getwd()
}
