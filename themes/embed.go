// Package themes 把示例主题编译进二进制。
//
// 开发时 pagebuilder_demo 优先读磁盘上的 themes/demo（改 HTML/CSS 刷新即见）；
// 工作目录不是 example 根（集成测试、部署后的二进制）时回退到这里的副本。
package themes

import "embed"

// FS 内含 demo/ 整个目录（all: 前缀连带以 . 或 _ 开头的文件）。
//
//go:embed all:demo
var FS embed.FS
