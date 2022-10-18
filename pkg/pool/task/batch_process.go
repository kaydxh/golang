package task

import "fmt"

func BatchProcess[T any](data []T, batchSize int, f func(d []T) error) error {
	if batchSize <= 0 {
	  return fmt.Errorf("invalid batchSize %d", batchSize )
	}

	if f == nil {
	  return fmt.Errorf("porcess func is nil")
	}

	for start, end := 0, 0; start < len(data); start = end {
		end = start + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[start:end]
		 err := f(batch)
		 if err != nil {
		   return err
		 }
	}

	return nil
}

