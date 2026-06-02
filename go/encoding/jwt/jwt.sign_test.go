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
	"errors"
	"testing"
	"time"
)

func TestSignHS256AndVerify(t *testing.T) {
	t.Parallel()

	const secret = "this-is-a-test-secret-at-least-32-bytes!!"
	claims := map[string]interface{}{
		"uid": "alice",
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}

	token, err := SignHS256(claims, secret)
	if err != nil {
		t.Fatalf("SignHS256() error = %v", err)
	}

	got, err := VerifyHS256(token, secret)
	if err != nil {
		t.Fatalf("VerifyHS256() error = %v", err)
	}
	if GetClaimString(got, "uid") != "alice" {
		t.Errorf("VerifyHS256() uid = %q, want alice", GetClaimString(got, "uid"))
	}
}

func TestVerifyHS256Errors(t *testing.T) {
	t.Parallel()

	const secret = "this-is-a-test-secret-at-least-32-bytes!!"
	validToken, _ := SignHS256(map[string]interface{}{"uid": "bob"}, secret)
	expiredToken, _ := SignHS256(map[string]interface{}{
		"uid": "bob",
		"exp": float64(time.Now().Add(-time.Hour).Unix()),
	}, secret)

	tests := []struct {
		name    string
		token   string
		secret  string
		wantErr error
	}{
		{name: "empty secret", token: validToken, secret: "", wantErr: ErrEmptySecret},
		{name: "malformed token", token: "a.b", secret: secret, wantErr: ErrInvalidToken},
		{name: "wrong secret", token: validToken, secret: "another-secret-32-bytes-padding!!", wantErr: ErrSignatureInvalid},
		{name: "tampered payload", token: validToken + "x", secret: secret, wantErr: ErrSignatureInvalid},
		{name: "expired token", token: expiredToken, secret: secret, wantErr: ErrTokenExpired},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := VerifyHS256(tt.token, tt.secret)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("VerifyHS256() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestSignHS256EmptySecret(t *testing.T) {
	t.Parallel()

	if _, err := SignHS256(map[string]interface{}{"uid": "x"}, ""); !errors.Is(err, ErrEmptySecret) {
		t.Errorf("SignHS256() error = %v, want ErrEmptySecret", err)
	}
}
