package errors_test

import (
	"fmt"
	"testing"

	errors_ "github.com/kaydxh/golang/go/errors"
)

func TestError(t *testing.T) {
	var errs = []error{fmt.Errorf("error 1"), fmt.Errorf("error 2")}

	var err error
	//Aggregate  implemnet interface error 	Error() string
	err = errors_.NewAggregate(errs)
	//multiErrorStrings := err errors_.NewAggregate(errs).Errors()
	multiErrorStrings := err.Error()
	t.Logf("multiErrorStrings: %v", multiErrorStrings)

}
