package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorToString(err error) string {
	if err == nil {
		return codes.OK.String()
	}

	st, ok := status.FromError(err)
	if !ok {
		return err.Error()
	}

	return st.Code().String()
}

func ErrorToCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	st, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}

	return st.Code()
}
