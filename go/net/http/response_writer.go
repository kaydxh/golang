package http

import (
	"bytes"
	"fmt"
	"net/http"
)

// ResponseWriterWrapper
type ResponseWriterWrapper struct {
	w          http.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		w:          w,
		statusCode: http.StatusOK,
	}
}

func (rww *ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return rww.w.Write(buf)
}

// Header function overwrites the http.ResponseWriter Header() function
func (rww *ResponseWriterWrapper) Header() http.Header {
	return rww.w.Header()
}

// WriteHeader function overwrites the http.ResponseWriter WriteHeader() function
func (rww *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rww.statusCode = statusCode
	rww.w.WriteHeader(statusCode)
}

//  String function pack respose header, http status code and body
func (rww *ResponseWriterWrapper) String() string {
	var buf bytes.Buffer

	buf.WriteString("Headers:")
	for k, v := range rww.w.Header() {
		buf.WriteString(fmt.Sprintf("%s: %v", k, v))
	}

	buf.WriteString(fmt.Sprintf(" Status Code:%d", rww.statusCode))

	buf.WriteString(" Body:")
	buf.WriteString(rww.body.String())
	return buf.String()
}
