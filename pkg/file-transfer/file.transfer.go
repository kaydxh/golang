package filetransfer

import (
	"context"
	"time"

	context_ "github.com/kaydxh/golang/go/context"
	http_ "github.com/kaydxh/golang/go/net/http"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type FileTransferOptions struct {
	// 0 means no timeout
	downloadTimeout time.Duration
	uploadTimeout   time.Duration
}

func defaultFileTransferOptions() FileTransferOptions {
	return FileTransferOptions{}
}

type FileTransfer struct {
	opts FileTransferOptions
}

func NewFileTransfer(opts ...FileTransferOption) *FileTransfer {
	ft := &FileTransfer{}
	ft.ApplyOptions(opts...)

	return ft
}

func (f *FileTransfer) Download(ctx context.Context, downloadUrl string) (data []byte, err error) {
	spanName := "Download"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()

	logger := logrus.WithField("trace_id", span.SpanContext().TraceID()).WithField("span_id", span.SpanContext().SpanID()).WithField("download_url", downloadUrl)

	ctx, cancel := context_.WithTimeout(ctx, f.opts.downloadTimeout)
	defer cancel()

	client, err := http_.NewClient(http_.WithTimeout(f.opts.downloadTimeout))
	if err != nil {
		logger.WithError(err).Errorf("new http client err: %v", err)
		return nil, err
	}
	data, err = client.Get(downloadUrl)
	if err != nil {
		logger.WithError(err).Errorf("http client get err: %v", err)
		return nil, err
	}

	return data, nil
}
