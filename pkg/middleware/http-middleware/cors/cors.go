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
package interceptorcors

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig 跨域资源共享配置。
type CORSConfig struct {
	// AllowOrigins 允许的来源列表。
	// 使用 "*" 表示允许所有来源（不建议在生产环境中使用）。
	// 默认值: ["*"]
	AllowOrigins []string

	// AllowMethods 允许的 HTTP 方法列表。
	// 默认值: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
	AllowMethods []string

	// AllowHeaders 允许的请求头列表。
	// 使用 "*" 表示允许所有请求头。
	// 默认值: ["*"]
	AllowHeaders []string

	// ExposeHeaders 允许浏览器访问的响应头列表。
	ExposeHeaders []string

	// AllowCredentials 是否允许携带凭证（cookies、HTTP认证等）。
	// 注意：当设置为 true 时，AllowOrigins 不能为 "*"。
	AllowCredentials bool

	// MaxAge 预检请求结果的缓存时间（秒）。
	// 默认值: 86400（24小时）
	MaxAge int
}

// DefaultCORSConfig 返回默认的 CORS 配置（允许所有来源）。
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},
		MaxAge:       86400,
	}
}

// CORS 返回一个使用指定配置的 CORS 中间件。
func CORS(config CORSConfig) func(http.Handler) http.Handler {
	// 填充默认值
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = []string{"*"}
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = []string{"*"}
	}
	if config.MaxAge <= 0 {
		config.MaxAge = 86400
	}

	allowMethodsStr := strings.Join(config.AllowMethods, ", ")
	allowHeadersStr := strings.Join(config.AllowHeaders, ", ")
	exposeHeadersStr := strings.Join(config.ExposeHeaders, ", ")
	maxAgeStr := strconv.Itoa(config.MaxAge)

	allowAllOrigins := len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// 设置 CORS 响应头
			if allowAllOrigins {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if origin != "" && isOriginAllowed(origin, config.AllowOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if exposeHeadersStr != "" {
				w.Header().Set("Access-Control-Expose-Headers", exposeHeadersStr)
			}

			// 处理预检请求
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", allowMethodsStr)
				w.Header().Set("Access-Control-Allow-Headers", allowHeadersStr)
				w.Header().Set("Access-Control-Max-Age", maxAgeStr)
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSAllowAll 返回一个允许所有来源的 CORS 中间件（使用默认配置）。
// 适用于开发环境或内部服务。
func CORSAllowAll(next http.Handler) http.Handler {
	return CORS(DefaultCORSConfig())(next)
}

// isOriginAllowed 检查给定的 origin 是否在允许列表中。
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
