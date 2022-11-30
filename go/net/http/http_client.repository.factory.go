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
package http

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/proto"
)

type FactoryConfigFunc func(c *FactoryConfig) error

type FactoryConfig struct {
	Validator *validator.Validate
	Url       string
	Timeout   time.Duration //接口处理超时时间
	*Client
	RetryTimes    int
	RetryInterval time.Duration
}

func (fc *FactoryConfig) ApplyOptions(configFuncs ...FactoryConfigFunc) error {

	for _, f := range configFuncs {
		err := f(fc)
		if err != nil {
			return fmt.Errorf("failed to apply factory config, err: %v", err)
		}
	}

	return nil
}

func (fc FactoryConfig) Validate() error {
	valid := fc.Validator
	if valid == nil {
		valid = validator.New()
	}
	return valid.Struct(fc)
}

type ProtoMessage interface {
	proto.Message
}

type Repository[REQ any, RESP any] struct {
	Url     string
	Timeout time.Duration

	*Client

	// not include the first call
	RetryTimes    int
	RetryInterval time.Duration
}

func newRepository[REQ any, RESP any](ctx context.Context, fc FactoryConfig) (*Repository[REQ, RESP], error) {
	repo := &Repository[REQ, RESP]{
		Url:           fc.Url,
		Timeout:       fc.Timeout,
		Client:        fc.Client,
		RetryTimes:    fc.RetryTimes,
		RetryInterval: fc.RetryInterval,
	}

	return repo, nil

}

type Factory[REQ any, RESP any] struct {
	fc FactoryConfig
}

func NewFactory[REQ any, RESP any](fc FactoryConfig, configFuncs ...FactoryConfigFunc) (Factory[REQ, RESP], error) {
	err := fc.ApplyOptions(configFuncs...)
	if err != nil {
		return Factory[REQ, RESP]{}, err
	}

	err = fc.Validate()
	if err != nil {
		return Factory[REQ, RESP]{}, err
	}

	return Factory[REQ, RESP]{fc: fc}, nil
}

func (f Factory[REQ, RESP]) NewClient(ctx context.Context) (*Repository[REQ, RESP], error) {
	return newRepository[REQ, RESP](ctx, f.fc)
}
