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
package opentelemetry

import (
	"context"

	metric_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/metric"
	tracer_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
)

type OpenTelemetryOptions struct {
	meterOptions  []metric_.MeterOption
	tracerOptions []tracer_.TracerOption
}

type OpenTelemetry struct {
	opts OpenTelemetryOptions
}

func NewOpenTelemetry(opts ...OpenTelemetryOption) *OpenTelemetry {
	t := &OpenTelemetry{}
	t.ApplyOptions(opts...)

	return t
}

func (t *OpenTelemetry) Install(ctx context.Context) error {

	if len(t.opts.meterOptions) > 0 {
		meter := metric_.NewMeter(t.opts.meterOptions...)
		err := meter.Install(ctx)
		if err != nil {
			return err
		}
	}

	if len(t.opts.tracerOptions) > 0 {
		tracer := tracer_.NewTracer(t.opts.tracerOptions...)
		err := tracer.Install(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
