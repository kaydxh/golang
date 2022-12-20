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
package filetransfer

import (
	"context"
	"math/rand"
	"time"

	http_ "github.com/kaydxh/golang/go/net/http"
	time_ "github.com/kaydxh/golang/go/time"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type FileTransferOptions struct {
	// 0 means no timeout
	downloadTimeout time.Duration
	uploadTimeout   time.Duration
	loadBalanceMode Ft_LoadBalanceMode
	retryTimes      int
	// retry interval, 0 means retry immediately
	retryInterval time.Duration

	proxies []*Ft_Proxy
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

func (f *FileTransfer) getProxy() *Ft_Proxy {
	proxies := f.opts.proxies
	if len(proxies) == 0 {
		return &Ft_Proxy{}
	}

	switch f.opts.loadBalanceMode {
	case Ft_load_balance_mode_random:
		return proxies[rand.Intn(len(proxies))]
	default:
		return proxies[0]
	}
}

// short connection
func (f *FileTransfer) Download(ctx context.Context, downloadUrl string) (data []byte, err error) {
	spanName := "Download"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()

	logger := logrus.WithField("trace_id", span.SpanContext().TraceID()).WithField("span_id", span.SpanContext().SpanID()).WithField("download_url", downloadUrl)

	proxy := f.getProxy()

	opts := []http_.ClientOption{http_.WithDisableKeepAlives(true)}
	if proxy.TargetAddr != "" {
		opts = append(opts, http_.WithProxyTarget(proxy.TargetAddr))
	} else if proxy.TargetUrl != "" {
		opts = append(opts, http_.WithProxy(proxy.TargetUrl))
	}
	opts = append(opts, http_.WithTimeout(f.opts.downloadTimeout))

	time_.RetryWithContext(ctx, func(ctx context.Context) error {
		client, err := http_.NewClient(opts...)
		if err != nil {
			logger.WithError(err).Errorf("new http client err: %v", err)
			return err
		}
		data, err = client.Get(ctx, downloadUrl)
		if err != nil {
			logger.WithError(err).Errorf("http client get err: %v", err)
			return err
		}
		return nil

	}, f.opts.retryInterval, f.opts.retryTimes)

	return data, nil
}

// short connection
func (f *FileTransfer) Upload(ctx context.Context, uploadUrl string, body []byte) (data []byte, err error) {
	spanName := "Upload"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()

	logger := logrus.WithField("trace_id", span.SpanContext().TraceID()).WithField("span_id", span.SpanContext().SpanID()).WithField("upload_url", uploadUrl)

	proxy := f.getProxy()

	opts := []http_.ClientOption{http_.WithDisableKeepAlives(true)}
	if proxy.TargetAddr != "" {
		opts = append(opts, http_.WithProxyTarget(proxy.TargetAddr))
	} else if proxy.TargetUrl != "" {
		opts = append(opts, http_.WithProxy(proxy.TargetUrl))
	}
	opts = append(opts, http_.WithTimeout(f.opts.uploadTimeout))

	time_.RetryWithContext(ctx, func(ctx context.Context) error {
		client, err := http_.NewClient(opts...)
		if err != nil {
			logger.WithError(err).Errorf("new http client err: %v", err)
			return err
		}
		data, err = client.Put(ctx, uploadUrl, "", nil, body)
		if err != nil {
			logger.WithError(err).Errorf("http client put err: %v", err)
			return err
		}
		return nil

	}, f.opts.retryInterval, f.opts.retryTimes)

	return data, nil
}
