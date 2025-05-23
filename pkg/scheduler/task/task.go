/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package task

import (
	"reflect"

	reflect_ "github.com/kaydxh/golang/go/reflect"
)

type Task struct {
	TaskId   string
	TaskFunc interface{}
	Args     []TaskArgument
	//	TaskStatus
	PreTaskHandler  func(*Task) error
	PostTaskHandler func(*Task) error
}

// Types which can be used: bool, string, int int8 int16 int32 int64, uint uint8 uint16 uint32 uint64, float32 float64
type TaskArgument struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (t *Task) Run() ([]*TaskResult, error) {

	reflectTaskFunc := reflect.ValueOf(t.TaskFunc)
	reflectTaskArgs, err := reflectTaskArgs(t.Args)
	if err != nil {
		return nil, err
	}

	results := reflectTaskFunc.Call(reflectTaskArgs)
	return reflectTaskResults(results)
}

func reflectTaskArgs(args []TaskArgument) ([]reflect.Value, error) {
	argValues := make([]reflect.Value, len(args))

	for i, arg := range args {
		argValue, err := reflect_.ReflectValue(arg.Type, arg.Value)
		if err != nil {
			return nil, err
		}
		argValues[i] = argValue
	}

	return argValues, nil
}

func reflectTaskResults(results []reflect.Value) ([]*TaskResult, error) {
	// Task must return at least a value
	if len(results) == 0 {
		return nil, ErrTaskReturnsNoValue
	}

	// Last returned value
	lastResult := results[len(results)-1]
	if !lastResult.IsNil() {
		// check that the result implements the standard error interface,
		// if not, return ErrLastReturnValueMustBeError error
		errorInterface := reflect.TypeOf((*error)(nil)).Elem()
		if !lastResult.Type().Implements(errorInterface) {
			return nil, ErrLastReturnValueMustBeError
		}

		// Return the standard error
		return nil, lastResult.Interface().(error)
	}

	// Convert reflect values to task results
	taskResults := make([]*TaskResult, len(results)-1)
	for i := 0; i < len(results)-1; i++ {
		val := results[i].Interface()
		typeStr := reflect.TypeOf(val).String()
		taskResults[i] = &TaskResult{
			Type:  typeStr,
			Value: val,
		}
	}

	return taskResults, nil
}
