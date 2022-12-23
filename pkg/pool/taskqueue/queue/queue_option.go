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
package queue

/*
// A QueueOption sets options.
type QueueOption interface {
	apply(*Queue)
}

// EmptyQueueUrlOption does not alter the Queueuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyQueueOption struct{}

func (EmptyQueueOption) apply(*Queue) {}

// QueueOptionFunc wraps a function that modifies Queue into an
// implementation of the QueueOption interface.
type QueueOptionFunc func(*Queue)

func (f QueueOptionFunc) apply(do *Queue) {
	f(do)
}

// sample code for option, default for nothing to change
func _QueueOptionWithDefault() QueueOption {
	return QueueOptionFunc(func(*Queue) {
		// nothing to change
	})
}
func (o *Queue) ApplyOptions(options ...QueueOption) *Queue {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
*/
