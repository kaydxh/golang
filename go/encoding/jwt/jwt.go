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

// Package jwt 提供轻量级的 JWT（JSON Web Token）解析工具。
//
// 注意：本包仅解析 JWT payload，不验证签名。适用于以下场景：
//   - Token 来自可信的服务端（server-to-server）
//   - 仅需提取 claims 中的字段用于日志、路由等非鉴权用途
//
// 如果 Token 来自不可信客户端或用于鉴权决策，必须使用带签名验证的 JWT 库。
package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// ParsePayload 解析 JWT token 的 payload 部分，返回 claims map。
// 不验证签名，仅做 base64url 解码 + JSON 反序列化。
// token 格式必须为标准 JWT 三段式：header.payload.signature。
func ParsePayload(token string) (map[string]interface{}, error) {
	if token == "" {
		return nil, fmt.Errorf("jwt: empty token")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("jwt: invalid token format, expected 3 parts, got %d", len(parts))
	}

	// base64url 解码 payload（JWT 标准使用无 padding 的 base64url）
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		// 兼容带 padding 的 base64url
		decoded, err = base64.URLEncoding.DecodeString(padBase64(parts[1]))
		if err != nil {
			// 最后尝试标准 base64（某些非标准实现可能使用）
			decoded, err = base64.StdEncoding.DecodeString(padBase64(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("jwt: decode payload: %w", err)
			}
		}
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("jwt: unmarshal payload: %w", err)
	}

	return claims, nil
}

// GetClaimString 从 claims 中安全提取字符串类型的字段值。
// 如果字段不存在或类型不是 string，返回空字符串。
func GetClaimString(claims map[string]interface{}, key string) string {
	if claims == nil {
		return ""
	}
	v, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return v
}

// GetClaimFloat64 从 claims 中安全提取数值类型的字段值。
// JSON 数字默认反序列化为 float64，如果字段不存在或类型不匹配，返回 0。
func GetClaimFloat64(claims map[string]interface{}, key string) float64 {
	if claims == nil {
		return 0
	}
	v, ok := claims[key].(float64)
	if !ok {
		return 0
	}
	return v
}

// padBase64 为 base64 字符串补齐 padding。
func padBase64(s string) string {
	switch len(s) % 4 {
	case 2:
		return s + "=="
	case 3:
		return s + "="
	}
	return s
}
