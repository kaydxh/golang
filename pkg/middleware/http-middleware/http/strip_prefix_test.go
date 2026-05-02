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
	"net/http/httptest"
	"testing"
)

func TestStripPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		requestPath  string
		expectedPath string
	}{
		{
			name:         "strip prefix from path",
			prefix:       "/palm-racer",
			requestPath:  "/palm-racer/api/v1/users",
			expectedPath: "/api/v1/users",
		},
		{
			name:         "strip prefix exact match becomes root",
			prefix:       "/palm-racer",
			requestPath:  "/palm-racer",
			expectedPath: "/",
		},
		{
			name:         "strip prefix with trailing slash",
			prefix:       "/palm-racer",
			requestPath:  "/palm-racer/",
			expectedPath: "/",
		},
		{
			name:         "no match leaves path unchanged",
			prefix:       "/palm-racer",
			requestPath:  "/other/api/v1",
			expectedPath: "/other/api/v1",
		},
		{
			name:         "partial match not stripped",
			prefix:       "/palm",
			requestPath:  "/palm-racer/api",
			expectedPath: "/palm-racer/api",
		},
		{
			name:         "empty prefix does nothing",
			prefix:       "",
			requestPath:  "/api/v1/users",
			expectedPath: "/api/v1/users",
		},
		{
			name:         "root prefix does nothing",
			prefix:       "/",
			requestPath:  "/api/v1/users",
			expectedPath: "/api/v1/users",
		},
		{
			name:         "prefix with trailing slash normalized",
			prefix:       "/palm-racer/",
			requestPath:  "/palm-racer/api/v1",
			expectedPath: "/api/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			handler := StripPrefix(tt.prefix)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if gotPath != tt.expectedPath {
				t.Errorf("StripPrefix(%q): path %q -> got %q, want %q",
					tt.prefix, tt.requestPath, gotPath, tt.expectedPath)
			}
		})
	}
}

func TestStripPrefixRawPath(t *testing.T) {
	handler := StripPrefix("/app")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1" {
			t.Errorf("expected Path '/api/v1', got '%s'", r.URL.Path)
		}
		if r.URL.RawPath != "/api/v1" {
			t.Errorf("expected RawPath '/api/v1', got '%s'", r.URL.RawPath)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/app/api/v1", nil)
	req.URL.RawPath = "/app/api/v1"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}
