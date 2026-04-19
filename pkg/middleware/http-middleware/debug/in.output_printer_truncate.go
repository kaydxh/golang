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
package interceptordebug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	http_ "github.com/kaydxh/golang/go/net/http"
	reflect_ "github.com/kaydxh/golang/go/reflect"
	logs_ "github.com/kaydxh/golang/pkg/logs"
)

// truncateStringForLog 按 reflect_ 包的默认阈值/前缀对单个字符串做截断，
// 输出格式与 reflect_.TruncateBytesAndStrings 对 string 字段的输出保持一致，
// 便于 HTTP 与 gRPC 两侧日志格式统一。
func truncateStringForLog(s string) string {
	threshold := reflect_.DefaultTruncateThreshold
	prefix := reflect_.DefaultTruncatePrefix
	if len(s) <= threshold {
		return s
	}
	if prefix > 0 && prefix < len(s) {
		return fmt.Sprintf("%s...(string len: %d)", s[:prefix], len(s))
	}
	return fmt.Sprintf("string len: %d", len(s))
}

// truncateLogValue 递归遍历 JSON 反序列化后的对象，对所有超长 string 字段做截断。
// map/slice 递归处理；string 调用 truncateStringForLog；其他类型（number/bool/null）原样返回。
func truncateLogValue(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		return truncateStringForLog(val)
	case map[string]interface{}:
		for k, sub := range val {
			val[k] = truncateLogValue(sub)
		}
		return val
	case []interface{}:
		for i, sub := range val {
			val[i] = truncateLogValue(sub)
		}
		return val
	default:
		return v
	}
}

// formatBodyForLog 将原始 body 转成适合打印的字符串：
//   - 尝试按 JSON 解析并递归截断长 string 字段；
//   - 解析失败（非 JSON），如果 body 本身超阈值就整体按 string 截断；短则原样返回。
func formatBodyForLog(buf []byte) string {
	if len(buf) == 0 {
		return ""
	}
	var parsed interface{}
	if err := json.Unmarshal(buf, &parsed); err != nil {
		return truncateStringForLog(string(buf))
	}
	truncated := truncateLogValue(parsed)
	out, err := json.Marshal(truncated)
	if err != nil {
		return truncateStringForLog(string(buf))
	}
	return string(out)
}

// InOutputPrinterWithTruncate HTTP 请求/响应日志中间件，对超长 string 字段做截断。
//
// 与 InOutputPrinter 的差异：
//   - 自动识别 request/response 中的 JSON 结构；
//   - 对 JSON 里超长的 string 字段（例如 base64 图片数据）做截断，
//     格式为 `前 N 字节...(string len: 总长度)`，与 gRPC 侧
//     UnaryServerInterceptorOfInOutputPrinter 的截断输出保持一致；
//   - 非 JSON 的 body 也会按整体超长做截断，避免日志被一次性刷爆；
//   - response 字段仅打印 body（不含 headers），status 作为独立字段输出。
//
// 注意：仅影响日志打印，不会修改实际转发给 handler 的 body。
func InOutputPrinterWithTruncate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs_.GetLogger(r.Context())
		ww := http_.NewResponseWriterWrapper(w)

		calleeMethod := fmt.Sprintf("%v %v", r.Method, r.URL.Path)

		defer func() {
			logger.WithField("method", calleeMethod).
				WithField("status", ww.StatusCode()).
				WithField("response", formatBodyForLog(ww.BodyBytes())).
				Info("send")
		}()

		if r != nil && r.Body != nil {
			buf, err := io.ReadAll(r.Body)
			if err == nil {
				// 把读出来的 body 再塞回去，不影响下游 handler 读取。
				r.Body = io.NopCloser(bytes.NewBuffer(buf))
				logger.WithField("method", calleeMethod).
					WithField("request", formatBodyForLog(buf)).
					Info("recv")
			}
		}

		handler.ServeHTTP(ww, r)
	})
}
