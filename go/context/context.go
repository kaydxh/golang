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
package context

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// RequestIdKey is metadata key name for request ID
const (
	DefaultHTTPRequestIDKey = "X-Request-ID"
)

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {

	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}

func ExtractStringFromContext(ctx context.Context, key string) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}

	return ""
}

func ExtractIntegerFromContext(ctx context.Context, key string) (int64, error) {
	v, ok := ctx.Value(key).(string)
	if !ok {
		return 0, fmt.Errorf("key[%v] value type is not string", key)
	}

	number, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("key[%v] value type is not number: %v", key, err)
	}

	return number, nil
}

func ExtractFromContext(ctx context.Context, key string) string {
	switch value := ctx.Value(key).(type) {
	case string:
		if value != "" {
			return value
		}
	case []string:
		if len(value) > 0 {
			return value[0]
		}
	default:
		return ""
	}

	return ""
}

func UpdateContext(ctx context.Context, key string, values map[string]interface{}) error {
	currentValues, ok := ctx.Value(key).(map[string]interface{})
	if ok {
		for k, v := range values {
			currentValues[k] = v
		}
		return nil
	}

	return fmt.Errorf("key[%v] is not exist in context", key)
}

func SetPairContext(ctx context.Context, key, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

func AppendContext(ctx context.Context, key string, values ...string) context.Context {
	currentValues, _ := ctx.Value(key).([]string)
	currentValues = append(currentValues, values...)
	return context.WithValue(ctx, key, currentValues)
}

func WithContextRequestId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, DefaultHTTPRequestIDKey, id)
}

func ExtractRequestIDFromContext(ctx context.Context) string {

	if v, ok := ctx.Value(DefaultHTTPRequestIDKey).(string); ok {
		return v
	}

	return ""
}
