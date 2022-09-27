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

// A OpenTelemetryOption sets options.
type OpenTelemetryOption interface {
	apply(*OpenTelemetry)
}

// EmptyOpenTelemetryOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyOpenTelemetryOption struct{}

func (EmptyOpenTelemetryOption) apply(*OpenTelemetry) {}

// OpenTelemetryOptionFunc wraps a function that modifies Client into an
// implementation of the OpenTelemetryOption interface.
type OpenTelemetryOptionFunc func(*OpenTelemetry)

func (f OpenTelemetryOptionFunc) apply(do *OpenTelemetry) {
	f(do)
}

// sample code for option, default for nothing to change
func _OpenTelemetryOptionWithDefault() OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(*OpenTelemetry) {
		// nothing to change
	})
}
func (o *OpenTelemetry) ApplyOptions(options ...OpenTelemetryOption) *OpenTelemetry {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
