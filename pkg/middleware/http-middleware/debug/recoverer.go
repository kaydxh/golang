package interceptordebug

import (
	"net/http"
	"runtime/debug"

	runtime_ "github.com/kaydxh/golang/go/runtime"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recoverer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil && r != http.ErrAbortHandler {

				out, err := runtime_.FormatStack(r)
				if err != nil {
					logrus.WithError(status.Errorf(codes.Internal, "%s", r)).Errorf("%s", debug.Stack())
				} else {
					logrus.WithError(status.Errorf(codes.Internal, "%s", r)).Errorf("%s", out)
				}
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
