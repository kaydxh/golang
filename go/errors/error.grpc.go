package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromError(err error) (s *status.Status, ok bool) {
	s, ok = status.FromError(err)
	if ok {
		return s, ok
	}
	var gRPCStatus interface {
		GRPCStatus() *status.Status
	}
	if errors.As(err, &gRPCStatus) {
		s = gRPCStatus.GRPCStatus()
		ok = true
	}
	return s, ok
}

func ErrorToString(err error) string {
	if err == nil {
		return codes.OK.String()
	}

	//st, ok := status.FromError(err)
	st, ok := FromError(err)
	if !ok {
		return err.Error()
	}

	return st.Message()
}

func ErrorToCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	//st, ok := status.FromError(err)
	st, ok := FromError(err)
	if !ok {
		return codes.Unknown
	}

	return st.Code()
}
