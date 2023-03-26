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
package config

import (
	"github.com/go-playground/validator/v10"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/spf13/viper"

	"google.golang.org/protobuf/proto"
)

// proto.Message is interface
type Config[T proto.Message] struct {
	Proto     T
	Validator *validator.Validate

	opts struct {
		viper *viper.Viper
	}
}

type completedConfig[T proto.Message] struct {
	*Config[T]
	completeError error
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig[T proto.Message] struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig[T]
}

// Validate checks Config.
func (c *completedConfig[T]) Validate() error {
	return c.Validator.Struct(c)
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config[T]) Complete() CompletedConfig[T] {
	err := c.loadViper()
	if err != nil {
		return CompletedConfig[T]{
			&completedConfig[T]{
				Config:        c,
				completeError: err,
			}}
	}

	if c.Validator == nil {
		c.Validator = validator.New()
	}

	return CompletedConfig[T]{&completedConfig[T]{Config: c}}
}

func (c *Config[T]) loadViper() error {
	if c.opts.viper != nil {
		return viper_.UnmarshalProtoMessageWithJsonPb(c.opts.viper, c.Proto)
	}

	return nil
}

func (c completedConfig[T]) New() (T, error) {
	var zero T
	if c.completeError != nil {
		return zero, c.completeError
	}

	return c.Proto, nil
}

func NewConfig[T proto.Message](value T, options ...ConfigOption[T]) *Config[T] {
	var c Config[T]
	c.Proto = value
	c.ApplyOptions(options...)

	return &c
}
