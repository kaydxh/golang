package task

import (
	"errors"
	"reflect"
)

var (
	ErrTaskNil = errors.New("Task must be not nil")
	// ErrTaskMustBeFunc ...
	ErrTaskMustBeFunc = errors.New("Task must be a func type")
	// ErrTaskReturnsNoValue ...
	ErrTaskReturnsNoValue = errors.New("Task must return at least a single value")
	// ErrLastReturnValueMustBeError ..
	ErrLastReturnValueMustBeError = errors.New("Last return value of a task must be error")
)

// Validate validates task function using reflection and makes sure
// it has a proper signature. Functions used as tasks must return at least a
// single value and the last return type must be error
func Validate(task *Task) error {
	if task == nil {
		return ErrTaskNil
	}
	v := reflect.ValueOf(task.TaskFunc)
	t := v.Type()

	// TaskFunc must be a function
	if t.Kind() != reflect.Func {
		return ErrTaskMustBeFunc
	}

	// TaskFunc must return at least a single value
	if t.NumOut() < 1 {
		return ErrTaskReturnsNoValue
	}

	// Last return value must be error
	lastReturnType := t.Out(t.NumOut() - 1)
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !lastReturnType.Implements(errorInterface) {
		return ErrLastReturnValueMustBeError
	}

	return nil
}
