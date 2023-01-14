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
package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type FactoryConfigFunc[T any] func(c *FactoryConfig[T]) error

type FactoryConfig[T any] struct {
	Validator        *validator.Validate
	Addr             string
	Timeout          time.Duration //接口处理超时时间
	NewServiceClient func(*grpc.ClientConn) T

	// not include the first call
	RetryTimes    int
	RetryInterval time.Duration

	DisablePrintInoutMethods []string
}

func (fc *FactoryConfig[T]) ApplyOptions(configFuncs ...FactoryConfigFunc[T]) error {

	for _, f := range configFuncs {
		err := f(fc)
		if err != nil {
			return fmt.Errorf("failed to apply factory config, err: %v", err)
		}
	}

	return nil
}

func (fc FactoryConfig[T]) Validate() error {
	valid := fc.Validator
	if valid == nil {
		valid = validator.New()
	}
	return valid.Struct(fc)
}

type Repository[T any] struct {
	Timeout time.Duration

	Addr             string
	NewServiceClient func(*grpc.ClientConn) T
	Conn             *grpc.ClientConn
	Client           T

	// not include the first call
	RetryTimes    int
	RetryInterval time.Duration

	DisablePrintInoutMethods []string
}

func (r *Repository[T]) NewConnect() (client T, conn *grpc.ClientConn, err error) {

	var zeroClient T
	conn, err = GetGrpcClientConn(r.Addr, r.Timeout)
	if err != nil {
		return zeroClient, nil, err
	}

	return r.NewServiceClient(conn), conn, nil
}

func (r *Repository[T]) Close(conn *grpc.ClientConn) (err error) {
	return conn.Close()
}

func newRepository[T any](ctx context.Context, fc FactoryConfig[T]) (Repository[T], error) {
	conn, err := GetGrpcClientConn(fc.Addr, fc.Timeout, fc.DisablePrintInoutMethods...)
	if err != nil {
		return Repository[T]{}, err
	}

	repo := Repository[T]{
		Timeout:                  fc.Timeout,
		Addr:                     fc.Addr,
		NewServiceClient:         fc.NewServiceClient,
		Conn:                     conn,
		Client:                   fc.NewServiceClient(conn),
		RetryTimes:               fc.RetryTimes,
		RetryInterval:            fc.RetryInterval,
		DisablePrintInoutMethods: fc.DisablePrintInoutMethods,
	}

	return repo, nil
}

type Factory[T any] struct {
	fc FactoryConfig[T]
}

func NewFactory[T any](fc FactoryConfig[T], configFuncs ...FactoryConfigFunc[T]) (Factory[T], error) {
	err := fc.ApplyOptions(configFuncs...)
	if err != nil {
		return Factory[T]{}, err
	}

	err = fc.Validate()
	if err != nil {
		return Factory[T]{}, err
	}

	return Factory[T]{fc: fc}, nil
}

func (f Factory[T]) NewClient(ctx context.Context) (Repository[T], error) {
	return newRepository(ctx, f.fc)
}
