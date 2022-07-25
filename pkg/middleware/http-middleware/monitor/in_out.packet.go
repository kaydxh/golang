package interceptormonitor

import (
	"bytes"
	"io/ioutil"
	"net/http"

	http_ "github.com/kaydxh/golang/go/net/http"
	logs_ "github.com/kaydxh/golang/pkg/logs"
)

func InOutPacket(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs_.GetLogger(r.Context())
		ww := http_.NewResponseWriterWrapper(w)
		defer func() {
			logger.WithField("response", ww.String()).Info("send")
		}()
		if r != nil {
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr
			logger.WithField("request", string(buf)).Info("recv")

		}

		handler.ServeHTTP(ww, r)

	})
}

func HeaderPacket(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs_.GetLogger(r.Context())
		logger.WithField("request headers", r.Header).Info("send")

		handler.ServeHTTP(w, r)
	})
}
