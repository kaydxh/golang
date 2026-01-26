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
package otlp

import (
	"time"

	meter_ "github.com/kaydxh/golang/pkg/opentelemetry/metric"
)

// WithOTLPPushExporter creates an OTLP push exporter builder
// for Tencent Cloud Prometheus Remote Write
func WithOTLPPushExporter(
	endpoint string,
	opts ...OTLPExporterBuilderOption,
) meter_.MeterOption {
	allOpts := append([]OTLPExporterBuilderOption{WithEndpoint(endpoint)}, opts...)
	return meter_.WithPushExporter(NewOTLPExporterBuilder(allOpts...))
}

// WithTencentCloudPrometheus creates an OTLP exporter configured for Tencent Cloud Prometheus
// endpoint: Remote Write URL from Tencent Cloud Prometheus console
// token: Authentication token from Tencent Cloud Prometheus console
// insecure: Whether to use insecure connection (typically false for Tencent Cloud)
func WithTencentCloudPrometheus(
	endpoint string,
	token string,
	insecure bool,
) meter_.MeterOption {
	headers := make(map[string]string)
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	}

	return meter_.WithPushExporter(NewOTLPExporterBuilder(
		WithEndpoint(endpoint),
		WithHeaders(headers),
		WithProtocol(ProtocolHTTP),
		WithInsecure(insecure),
		WithTimeout(30*time.Second),
	))
}
