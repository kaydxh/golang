package strconv

import (
	"fmt"
	"strconv"
)

func ToFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}
func ToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func ParseNumOrDefault[T any](str string, _default T, convert func(string) (T, error)) T {
	if str == "" {
		return _default
	}
	n, err := convert(str)
	if err != nil {
		//logurs.Infof("Invalid number value: %s", err)
		return _default
	}
	return n
}

func ParseNum[T any](str string, convert func(string) (T, error)) (T, error) {
    var t T
	if str == "" {
		return t,  fmt.Errorf("string is empty")
	}
	return convert(str)
}
