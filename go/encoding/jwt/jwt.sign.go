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

package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// 签名与校验相关的 sentinel error，便于上层用 errors.Is 做错误码映射。
var (
	// ErrEmptySecret 表示签名密钥为空。
	ErrEmptySecret = errors.New("jwt: empty secret")
	// ErrInvalidToken 表示 token 格式非法（非标准三段式或解码失败）。
	ErrInvalidToken = errors.New("jwt: invalid token")
	// ErrSignatureInvalid 表示签名校验不通过。
	ErrSignatureInvalid = errors.New("jwt: signature invalid")
	// ErrTokenExpired 表示 token 已过期（exp 声明早于当前时间）。
	ErrTokenExpired = errors.New("jwt: token expired")
	// ErrUnsupportedAlg 表示 token 的 alg 不是 HS256。
	ErrUnsupportedAlg = errors.New("jwt: unsupported alg")
)

// hs256Header 是固定的 HS256 JWT 头部 base64url 编码结果。
// header = {"alg":"HS256","typ":"JWT"}
const hs256Header = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

// SignHS256 使用 HMAC-SHA256 将 claims 签发为标准三段式 JWT（header.payload.signature）。
//
// claims 为业务自定义载荷，可包含 exp（unix 秒）等标准字段；secret 为签名密钥，不能为空。
// 返回的 token 适用于鉴权场景，需配合 VerifyHS256 校验签名。
func SignHS256(claims map[string]interface{}, secret string) (string, error) {
	if secret == "" {
		return "", ErrEmptySecret
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadSeg := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signingInput := hs256Header + "." + payloadSeg
	sig := signHS256Segment(signingInput, secret)

	return signingInput + "." + sig, nil
}

// VerifyHS256 校验 HS256 JWT 的签名，并在存在 exp 声明时校验是否过期。
//
// 校验通过后返回 payload 中的 claims。任何失败都返回对应的 sentinel error，
// 上层可用 errors.Is 区分 ErrSignatureInvalid / ErrTokenExpired 等。
func VerifyHS256(token, secret string) (map[string]interface{}, error) {
	if secret == "" {
		return nil, ErrEmptySecret
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	// 重算签名并做恒定时间比较，防止时序攻击。
	expectedSig := signHS256Segment(parts[0]+"."+parts[1], secret)
	if !hmac.Equal([]byte(expectedSig), []byte(parts[2])) {
		return nil, ErrSignatureInvalid
	}

	claims, err := ParsePayload(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 存在 exp 时校验过期（exp 为 unix 秒，0 视为不设过期）。
	if exp := GetClaimFloat64(claims, "exp"); exp > 0 {
		if time.Now().Unix() >= int64(exp) {
			return nil, ErrTokenExpired
		}
	}

	return claims, nil
}

// signHS256Segment 计算 signingInput 的 HMAC-SHA256 签名，返回 base64url（无 padding）编码结果。
func signHS256Segment(signingInput, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
