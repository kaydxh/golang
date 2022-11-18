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
package errors

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
		//ok = true
		ok = s != nil
	}
	return s, ok
}

func ErrorToCodeString(err error) string {
	if err == nil {
		return codes.OK.String()
	}

	s, ok := FromError(err)
	if !ok {
		return err.Error()
	}
	codeString := s.Code().String()
	for _, detail := range s.Details() {
		switch detail := detail.(type) {
		case *errdetails.ErrorInfo:
			if detail.GetReason() != "" {
				codeString = detail.GetReason()
			}
		}
	}

	return codeString
}

// Error Message
func ErrorToString(err error) string {
	if err == nil {
		return codes.OK.String()
	}

	s, ok := FromError(err)
	if !ok {
		return err.Error()
	}
	return s.Message()
}

func ErrorToCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	s, ok := FromError(err)
	if !ok {
		return codes.Unknown
	}

	return s.Code()
}
