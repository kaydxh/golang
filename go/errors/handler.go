package errors

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// ErrorHandlers is a list of functions which will be invoked when a nonreturnable
// error occurs.
// should be packaged up into a testable and reusable object.
var ErrorHandlers = []func(error){
	func() func(err error) {
		limiter := rate.NewLimiter(rate.Every(time.Millisecond), 1)
		return func(err error) {
			// 1ms was the number folks were able to stomach as a global rate limit.
			// If you need to log errors more than 1000 times a second you
			// should probably consider fixing your code instead. :)
			_ = limiter.Wait(context.Background())
		}
	}(),
}

// HandlerError is a method to invoke when a non-user facing piece of code cannot
// return an error and needs to indicate it has been ignored. Invoking this method
// is preferable to logging the error - the default behavior is to log but the
// errors may be sent to a remote server for analysis.
func HandleError(err error) {
	// this is sometimes called with a nil error.  We probably shouldn't fail and should do nothing instead
	if err == nil {
		return
	}

	for _, fn := range ErrorHandlers {
		fn(err)
	}
}
