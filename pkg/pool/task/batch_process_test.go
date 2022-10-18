package task_test

import (
	"fmt"
	"testing"

	task_ "github.com/kaydxh/golang/pkg/pool/task"
)

func processTask(data []int) error {
	fmt.Printf("process %v\n", data)
	return nil
}

func TestBatchProcessSync(t *testing.T) {
	testCases := []struct {
		data      []int
		batchSize int
		expected  string
	}{
		{
			data:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			batchSize: 4,
			expected:  "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d\n", i), func(t *testing.T) {
			err := task_.BatchProcess(testCase.data, testCase.batchSize, processTask)
			if err != nil {
				t.Fatalf("failed to batch process task %d, got : %v", i, err)

			}

		})
	}

}
