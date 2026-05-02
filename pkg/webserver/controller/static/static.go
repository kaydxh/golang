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
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
)

// Config 静态文件服务配置。
type Config struct {
	// Root 静态文件根目录路径。
	// 可通过 EnvKey 环境变量覆盖。
	Root string

	// IndexFile SPA 入口文件名（相对于 Root）。
	// 默认值: "index.html"
	IndexFile string

	// AssetDirs 额外的静态资源目录映射。
	// key 为 URL 路径前缀，value 为相对于 Root 的子目录名。
	// 例如: {"assets": "assets", "models": "models"}
	AssetDirs map[string]string

	// StaticFiles 单个静态文件映射。
	// key 为 URL 路径，value 为相对于 Root 的文件路径。
	// 例如: {"/favicon.ico": "favicon.ico"}
	StaticFiles map[string]string

	// EnvKey 用于覆盖 Root 的环境变量名。
	// 如果设置了该环境变量，其值将覆盖 Root 配置。
	// 默认值: "STATIC_ROOT"
	EnvKey string

	// SPAMode 是否启用 SPA 模式。
	// 启用后，所有未匹配的 GET 请求将返回 IndexFile（Vue/React SPA fallback）。
	SPAMode bool
}

// Controller 静态文件服务控制器。
// 实现 WebHandler 接口，用于在 GenericWebServer 中注册静态文件路由。
type Controller struct {
	config Config
	root   string // 解析后的实际根目录
}

// NewController 创建静态文件服务控制器。
func NewController(config Config) *Controller {
	// 填充默认值
	if config.IndexFile == "" {
		config.IndexFile = "index.html"
	}
	if config.EnvKey == "" {
		config.EnvKey = "STATIC_ROOT"
	}
	if config.Root == "" {
		config.Root = "./static"
	}

	// 环境变量覆盖
	root := config.Root
	if v := os.Getenv(config.EnvKey); v != "" {
		root = v
	}

	return &Controller{
		config: config,
		root:   root,
	}
}

// SetRoutes 注册静态文件服务路由到 Gin 路由器。
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *gw_.GRPCGateway) {
	// 注册额外的静态资源目录
	for urlPath, dirName := range c.config.AssetDirs {
		absDir := filepath.Join(c.root, dirName)
		ginRouter.Static("/"+urlPath, absDir)
	}

	// 注册单个静态文件
	for urlPath, filePath := range c.config.StaticFiles {
		absFile := filepath.Join(c.root, filePath)
		ginRouter.StaticFile(urlPath, absFile)
	}

	// SPA 模式：根路径返回 index.html。
	// 使用 recover 捕获重复注册 panic，因为其他控制器（如 healthz）
	// 可能已经注册了 GET "/"，Gin 不允许重复注册同一路径。
	if c.config.SPAMode {
		indexFile := filepath.Join(c.root, c.config.IndexFile)
		func() {
			defer func() {
				if r := recover(); r != nil {
					// GET "/" 已被其他控制器注册，跳过 SPA 根路径注册
				}
			}()
			ginRouter.GET("/", func(ctx *gin.Context) {
				ctx.File(indexFile)
			})
		}()
	}
}

// GetRoot 返回解析后的静态文件根目录路径。
func (c *Controller) GetRoot() string {
	return c.root
}
