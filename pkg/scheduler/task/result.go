package task

import (
	"fmt"
	"reflect"
	"strings"

	reflect_ "github.com/kaydxh/golang/go/reflect"
)

// TaskResult represents an actual return value of a processed task
type TaskResult struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// ReflectTaskResults ...
func ReflectTaskResults(taskResults []*TaskResult) ([]reflect.Value, error) {
	resultValues := make([]reflect.Value, len(taskResults))
	for i, taskResult := range taskResults {
		resultValue, err := reflect_.ReflectValue(taskResult.Type, taskResult.Value)
		if err != nil {
			return nil, err
		}
		resultValues[i] = resultValue
	}
	return resultValues, nil
}

// HumanReadableResults ...
func HumanReadableResults(results []reflect.Value) string {
	if len(results) == 1 {
		return fmt.Sprintf("%v", results[0].Interface())
	}

	readableResults := make([]string, len(results))
	for i := 0; i < len(results); i++ {
		readableResults[i] = fmt.Sprintf("%v", results[i].Interface())
	}

	return fmt.Sprintf("[%s]", strings.Join(readableResults, ", "))
}