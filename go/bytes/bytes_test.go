package bytes_test

import (
	"fmt"
	"testing"

	"unsafe"

	bytes_ "github.com/kaydxh/golang/go/bytes"
)

func TestEncodeAndDecode(t *testing.T) {
	// to bytes
	type MyStruct struct {
		X, Y, Z int
		Name    string
	}

	testCases := []struct {
		s        MyStruct
		expected string
	}{
		{
			s:        MyStruct{1, 20, 30, "foobar"},
			expected: "",
		},
	}
	const sz = int(unsafe.Sizeof(MyStruct{}))
	fmt.Printf("MyStruct size: %v\n", int(unsafe.Sizeof(MyStruct{})))

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {

			b := bytes_.Encode(testCase.s)
			fmt.Println(b)
			fmt.Printf("size of b: %v\n", len(b))

			/*
				var p []byte = (*(*[sz]byte)(unsafe.Pointer(&testCase.s)))[:]
				fmt.Println(p)
				fmt.Printf("size of p: %v\n", len(p))
			*/

			/*
				var m = *(*MyStruct)(unsafe.Pointer(&b[0]))
				fmt.Println(m)
			*/

			r := bytes_.Decode[MyStruct](b)
			fmt.Println(r)

		})
	}
}
