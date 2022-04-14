package syscall_test

import (
	"testing"

	syscall_ "github.com/kaydxh/golang/go/syscall"
)

func TestDiskUsage(t *testing.T) {
	testCases := []struct {
		volumePath string
		expected   string
	}{
		{
			volumePath: "/",
			expected:   "",
		},
		{
			volumePath: "/dev",
			expected:   "",
		},
		{
			volumePath: "/data",
			expected:   "",
		},
		{
			volumePath: "/data/home/log",
			expected:   "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.volumePath, func(t *testing.T) {
			du, err := syscall_.NewDiskUsage(testCase.volumePath)
			if err != nil {
				t.Errorf("new disk for path[%v] err, got : %s", testCase.volumePath, err)
				return

			}
			t.Logf(
				"disk free[%v], avali[%v], size[%v], used[%v], usage: %v",
				du.Free(),
				du.Avail(),
				du.Size(),
				du.Used(),
				du.Usage(),
			)

		})
	}
}
