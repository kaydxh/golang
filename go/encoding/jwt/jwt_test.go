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
	"testing"
)

func TestParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
		check   func(claims map[string]interface{}) bool
	}{
		{
			name:    "空 token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "非法格式（只有一段）",
			token:   "abc",
			wantErr: true,
		},
		{
			name:    "非法格式（只有两段）",
			token:   "abc.def",
			wantErr: true,
		},
		{
			name:  "合法 JWT（base64url 无 padding）",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWxpY2UiLCJleHAiOjE3MTk5OTk5OTl9.signature",
			check: func(claims map[string]interface{}) bool {
				return GetClaimString(claims, "user_id") == "alice"
			},
		},
		{
			name:  "payload 为空 JSON 对象",
			token: "eyJhbGciOiJIUzI1NiJ9.e30.signature",
			check: func(claims map[string]interface{}) bool {
				return len(claims) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParsePayload(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.check != nil && !tt.check(claims) {
				t.Errorf("ParsePayload() claims check failed, claims = %v", claims)
			}
		})
	}
}

func TestGetClaimString(t *testing.T) {
	claims := map[string]interface{}{
		"user_id": "bob",
		"count":   float64(42),
	}

	if got := GetClaimString(claims, "user_id"); got != "bob" {
		t.Errorf("GetClaimString(user_id) = %q, want %q", got, "bob")
	}
	if got := GetClaimString(claims, "missing"); got != "" {
		t.Errorf("GetClaimString(missing) = %q, want empty", got)
	}
	if got := GetClaimString(claims, "count"); got != "" {
		t.Errorf("GetClaimString(count) = %q, want empty (type mismatch)", got)
	}
	if got := GetClaimString(nil, "user_id"); got != "" {
		t.Errorf("GetClaimString(nil) = %q, want empty", got)
	}
}

func TestGetClaimFloat64(t *testing.T) {
	claims := map[string]interface{}{
		"exp":     float64(1719999999),
		"user_id": "alice",
	}

	if got := GetClaimFloat64(claims, "exp"); got != 1719999999 {
		t.Errorf("GetClaimFloat64(exp) = %v, want %v", got, 1719999999)
	}
	if got := GetClaimFloat64(claims, "user_id"); got != 0 {
		t.Errorf("GetClaimFloat64(user_id) = %v, want 0 (type mismatch)", got)
	}
	if got := GetClaimFloat64(nil, "exp"); got != 0 {
		t.Errorf("GetClaimFloat64(nil) = %v, want 0", got)
	}
}
