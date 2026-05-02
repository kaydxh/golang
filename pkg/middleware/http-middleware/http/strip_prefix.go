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
package interceptorhttp

import (
	"net/http"
	"strings"
)

// StripPrefix 返回一个 HTTP 中间件，用于去除请求路径中的指定前缀。
//
// 典型使用场景：在 Kubernetes Ingress / 负载均衡器配置了路径前缀路由时，
// 后端服务需要 strip 掉该前缀才能正确匹配内部路由。
//
// 例如：Ingress 配置 /palm-racer 路由到后端，后端路由为 /api/xxx，
// 请求 /palm-racer/api/xxx 经过 StripPrefix("/palm-racer") 后变为 /api/xxx。
//
// 如果请求路径不以 prefix 开头，则不做任何修改直接传递给下一个 handler。
func StripPrefix(prefix string) func(http.Handler) http.Handler {
	// 规范化前缀：确保不以 "/" 结尾
	prefix = strings.TrimRight(prefix, "/")

	return func(next http.Handler) http.Handler {
		// 空前缀或仅为 "/" 时，不做任何处理
		if prefix == "" || prefix == "/" {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := r.URL.Path
			if p == prefix || strings.HasPrefix(p, prefix+"/") {
				newPath := strings.TrimPrefix(p, prefix)
				if newPath == "" {
					newPath = "/"
				}
				r.URL.Path = newPath

				if r.URL.RawPath != "" {
					r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, prefix)
					if r.URL.RawPath == "" {
						r.URL.RawPath = "/"
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
