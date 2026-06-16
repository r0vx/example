package ui_demo

import (
	"context"
	"crypto/sha1"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"example/models"

	"github.com/r0vx/admin/helpcenter"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/publish"
	"github.com/r0vx/x/login"
	"gorm.io/gorm"
)

// helpGenImageDir AI 配图落盘目录（相对工作目录）。
// router.go 已把 public/system 挂在 /system/ 下做静态服务，故此目录天然可公开访问。
const helpGenImageDir = "public/system/helpcenter-gen"

// exampleStoreGenImage 把 AI 文生图的字节落盘到 public/system 静态目录，返回可公开访问的相对 URL。
// 选用本地文件系统而非 S3：example 默认无 S3 凭据，且 router.go 已有 /system/ 静态路由，
// 写盘即可直出，零外部依赖、运行时即可用。生产环境可换成 OSS/媒体库回调。
func exampleStoreGenImage(ctx context.Context, data []byte, contentType string) (string, error) {
	if err := os.MkdirAll(helpGenImageDir, 0o755); err != nil {
		return "", fmt.Errorf("helpcenter: 创建配图目录失败: %w", err)
	}
	// 文件名用「时间戳 + 内容哈希」保证唯一且幂等（同图不重复占盘）。
	ext := extByContentType(contentType)
	sum := sha1.Sum(data)
	name := fmt.Sprintf("%d-%x%s", time.Now().UnixNano(), sum[:8], ext)
	if err := os.WriteFile(filepath.Join(helpGenImageDir, name), data, 0o644); err != nil {
		return "", fmt.Errorf("helpcenter: 写入配图失败: %w", err)
	}
	// 返回相对 URL：admin 与公开站同源，<img src> 用相对路径即可正常加载。
	return "/system/helpcenter-gen/" + name, nil
}

// extByContentType 由 Content-Type 推断文件扩展名，缺省回退 .png。
func extByContentType(contentType string) string {
	if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
		return exts[0]
	}
	return ".png"
}

// configHelpCenter 装配帮助中心：admin CRUD + 公开站 handler + 种子数据
func ConfigHelpCenter(b *presets.Builder, db *gorm.DB, publisher *publish.Builder) *helpcenter.Builder {
	db.AutoMigrate(&helpcenter.Article{})

	hc := helpcenter.New(db, publisher).
		PublicPrefix("/help").
		AdminMenu("帮助文档").
		AIEndpoint("/help-ai").
		// AI 端点鉴权：复用 example 的 getCurrentUser（登录中间件已填充请求上下文）
		AICheckAuth(func(r *http.Request) bool { return login.GetCurrentUser(r).(*models.User) != nil }).
		// 文生图落盘到本地静态目录，存档教程时把临时图转存为持久 URL
		GenImageStore(exampleStoreGenImage)
	hc.Install(b)

	seedHelpArticles(db)
	return hc
}

// seedHelpArticles 幂等插入示例文档（便于测试公开站）
func seedHelpArticles(db *gorm.DB) {
	var count int64
	db.Model(&helpcenter.Article{}).Count(&count)
	if count > 0 {
		return
	}
	samples := []struct {
		title, slug, icon, summary, content string
		tags                                []string
		pos                                 int
	}{
		{"安装指南", "install", "📦", "如何安装与初始化 r0vx。", "<h2>环境要求</h2><ul><li>Go 1.22+</li><li>PostgreSQL</li></ul><h2>安装步骤</h2><pre><code>go get github.com/r0vx/admin</code></pre><h3>验证</h3><p>运行 <code>go build ./...</code> 应无报错。</p>", []string{"入门"}, 1},
		{"第一个页面", "first-page", "📄", "创建你的第一个管理页面。", "<h2>创建模型</h2><pre><code>b.Model(&Product{})</code></pre><h2>列表配置</h2><p>用 <code>Listing</code> 指定展示列。</p>", []string{"入门"}, 2},
		{"权限配置", "permissions", "🔐", "角色与权限的配置方式。", "<h2>角色</h2><p>定义角色并分配权限资源。</p><h2>权限点</h2><p>通过 <code>PermissionRN</code> 暴露资源名。</p>", []string{"进阶"}, 3},
	}
	for _, s := range samples {
		a := &helpcenter.Article{
			Title: s.title, Slug: s.slug, Icon: s.icon,
			Summary: s.summary, Content: s.content, Tags: s.tags, Position: s.pos,
		}
		a.Status.Status = publish.StatusOnline // 种子直接上线，便于公开站可见
		a.Version.Version = "v1"
		a.Version.VersionName = "v1" // 正常创建流程会自动设 VersionName；种子直插需手动设，否则版本列表弹窗版本列空白
		db.Create(a)
	}
}
