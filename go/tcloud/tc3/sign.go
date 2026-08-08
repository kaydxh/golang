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

// Package tc3 实现腾讯云 API 3.0 签名方法 v3（TC3-HMAC-SHA256）。
//
// 参考：https://cloud.tencent.com/document/api/213/30654
// 签名过程：
//  1. 拼接规范请求串 CanonicalRequest
//  2. 拼接待签名字符串 StringToSign
//  3. 派生签名密钥并计算 Signature
//  4. 拼接 Authorization 头
//
// Date 取自 Timestamp 的 UTC 日期（加入本地时区会导致凌晨签名失败）。
package tc3

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	hmac_ "github.com/kaydxh/golang/go/crypto/hmac"
	sha256_ "github.com/kaydxh/golang/go/crypto/sha256"
)

// Algorithm 签名算法，固定为 TC3-HMAC-SHA256。
const Algorithm = "TC3-HMAC-SHA256"

// Options 描述一次 TC3-HMAC-SHA256 签名请求。
type Options struct {
	// SecretID / SecretKey 申请的安全凭证。
	SecretID  string
	SecretKey string

	// Service 产品名，须与调用域名一致，如 "palm"、"cvm"。
	Service string
	// Host 服务地址（裸域名，不带 scheme），同时参与签名。
	Host string

	// Timestamp 秒级 UNIX 时间戳。Date 由其换算为 UTC 日期。
	Timestamp int64

	// Payload 请求体（POST 即为 JSON 串）。GET 传空串。
	Payload string

	// Method HTTP 方法，空默认 POST。
	Method string
	// Path 请求路径，空默认 "/"（API 3.0 固定为 "/"）。
	Path string

	// SignedHeaders 参与签名的 header key（小写），至少包含 content-type 和 host。
	// 多个 key 按 ASCII 升序参与拼接，并以此构造 SignedHeaders 字段。
	SignedHeaders []string

	// HeaderValues 待发送的 header 集合（key 大小写均可），
	// 签名时按 SignedHeaders 的小写 key 取值。
	HeaderValues map[string]string
}

// Authorization 计算 TC3-HMAC-SHA256 签名的 Authorization 头值。
func Authorization(opts Options) string {
	method := opts.Method
	if method == "" {
		method = "POST"
	}
	path := opts.Path
	if path == "" {
		path = "/"
	}

	date := time.Unix(opts.Timestamp, 0).UTC().Format("2006-01-02")

	// 步骤 1：拼接规范请求串
	signed := append([]string{}, opts.SignedHeaders...)
	sort.Strings(signed)
	var canonicalHeaders strings.Builder
	for _, k := range signed {
		canonicalHeaders.WriteString(fmt.Sprintf("%s:%s\n", k, headerValue(opts.HeaderValues, k)))
	}
	signedHeaders := strings.Join(signed, ";")

	hashedPayload := sha256_.SumString(opts.Payload)
	canonicalRequest := strings.Join([]string{
		method, path, "", canonicalHeaders.String(), signedHeaders, hashedPayload,
	}, "\n")

	// 步骤 2：拼接待签名字符串
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, opts.Service)
	hashedCanonical := sha256_.SumString(canonicalRequest)
	stringToSign := strings.Join([]string{
		Algorithm, strconv.FormatInt(opts.Timestamp, 10), credentialScope, hashedCanonical,
	}, "\n")

	// 步骤 3：派生签名密钥并计算签名
	secretDate := hmac_.SumSHA256([]byte("TC3"+opts.SecretKey), date)
	secretService := hmac_.SumSHA256(secretDate, opts.Service)
	secretSigning := hmac_.SumSHA256(secretService, "tc3_request")
	signature := hmac_.SumSHA256Hex(secretSigning, stringToSign)

	// 步骤 4：拼接 Authorization
	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		Algorithm, opts.SecretID, credentialScope, signedHeaders, signature)
}

// headerValue 从 header 集合中按小写 key 取值。
func headerValue(headers map[string]string, lowerKey string) string {
	for k, v := range headers {
		if strings.ToLower(k) == lowerKey {
			return v
		}
	}
	return ""
}
