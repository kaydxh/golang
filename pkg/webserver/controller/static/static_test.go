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
package static

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewController(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		ctrl := NewController(Config{})
		if ctrl.config.IndexFile != "index.html" {
			t.Errorf("expected IndexFile 'index.html', got '%s'", ctrl.config.IndexFile)
		}
		if ctrl.config.EnvKey != "STATIC_ROOT" {
			t.Errorf("expected EnvKey 'STATIC_ROOT', got '%s'", ctrl.config.EnvKey)
		}
		if ctrl.root != "./static" {
			t.Errorf("expected root './static', got '%s'", ctrl.root)
		}
	})

	t.Run("custom config", func(t *testing.T) {
		ctrl := NewController(Config{
			Root:      "/var/www",
			IndexFile: "app.html",
			EnvKey:    "MY_STATIC",
		})
		if ctrl.root != "/var/www" {
			t.Errorf("expected root '/var/www', got '%s'", ctrl.root)
		}
		if ctrl.config.IndexFile != "app.html" {
			t.Errorf("expected IndexFile 'app.html', got '%s'", ctrl.config.IndexFile)
		}
	})

	t.Run("env override", func(t *testing.T) {
		os.Setenv("TEST_STATIC_ROOT", "/tmp/override")
		defer os.Unsetenv("TEST_STATIC_ROOT")

		ctrl := NewController(Config{
			Root:   "/var/www",
			EnvKey: "TEST_STATIC_ROOT",
		})
		if ctrl.root != "/tmp/override" {
			t.Errorf("expected root '/tmp/override', got '%s'", ctrl.root)
		}
	})
}

func TestController_SetRoutes(t *testing.T) {
	// 创建临时目录和文件
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)
	os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte("<html>hello</html>"), 0644)
	os.WriteFile(filepath.Join(assetsDir, "app.js"), []byte("console.log('hi')"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "favicon.ico"), []byte("icon"), 0644)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	ctrl := NewController(Config{
		Root:      tmpDir,
		SPAMode:   true,
		AssetDirs: map[string]string{"assets": "assets"},
		StaticFiles: map[string]string{
			"/favicon.ico": "favicon.ico",
		},
	})
	ctrl.SetRoutes(router, nil)

	t.Run("SPA root returns index.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		if w.Body.String() != "<html>hello</html>" {
			t.Errorf("expected index.html content, got '%s'", w.Body.String())
		}
	})

	t.Run("static asset file served", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		if w.Body.String() != "console.log('hi')" {
			t.Errorf("expected js content, got '%s'", w.Body.String())
		}
	})

	t.Run("static file served", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestController_GetRoot(t *testing.T) {
	ctrl := NewController(Config{Root: "/my/path"})
	if ctrl.GetRoot() != "/my/path" {
		t.Errorf("expected '/my/path', got '%s'", ctrl.GetRoot())
	}
}
