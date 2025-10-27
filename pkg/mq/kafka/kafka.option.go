/*
 *Copyright (c) 2023, kaydxh
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
package kafka

import "time"

// base options
func WithDialTimeout(dialTimeout time.Duration) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.dialTimeout = dialTimeout
	})
}

// producer options
func WithProducerBatchSize(batchSize int) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.producerOpts.batchSize = batchSize
	})
}

func WithProducerBatchBytes(batchBytes int) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.producerOpts.batchBytes = batchBytes
	})
}

func WithProducerBatchTimeout(batchTimeout time.Duration) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.producerOpts.batchTimeout = batchTimeout
	})
}

// consumer options
func WithConsumerGroupID(groupID string) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.consumerOpts.groupID = groupID
	})
}

func WithConsumerPartition(partition int) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.consumerOpts.partition = partition
	})
}

func WithConsumerMinBytes(minBytes int) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.consumerOpts.minBytes = minBytes
	})
}

func WithConsumerMaxBytes(maxBytes int) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.consumerOpts.maxBytes = maxBytes
	})
}

func WithConsumerMaxWait(maxWait time.Duration) MQOption {
	return MQOptionFunc(func(m *MQ) {
		m.opts.consumerOpts.maxWait = maxWait
	})
}
