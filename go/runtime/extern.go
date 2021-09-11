package runtime

import (
	"fmt"
	"runtime"
)

// GetCaller returns the caller of the function that calls it.
//The argument skip is the number of stack frames to skip before recording in pc, with 0 identifying the frame for Callers itself and 1 identifying the caller of Callers
func GetCaller(skip int) string {
	var pc [1]uintptr
	runtime.Callers(skip, pc[:])
	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return fmt.Sprintf("Unable to find caller")
	}
	return f.Name()
}

func GetParentCaller() string {
	return GetCaller(1)
}
