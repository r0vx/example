package admin

import (
	_ "embed"
	"fmt"
	"net/http"

	"example/admin/ui_demo"

	"github.com/go-chi/chi/v5"
	"github.com/r0vx/x/login"
	"github.com/r0vx/x/sitemap"
	"gorm.io/gorm"

	"example/models"

	"github.com/r0vx/x/perm"
)

//go:embed assets/favicon.ico
var favicon []byte

//go:embed assets/logo.svg
var logo []byte

//go:embed assets/logo-icon.svg
var logoIcon []byte

const (
	exportOrdersURL = "/export-orders"
	saveCookieURL   = "/save-cookie"
)

// TestHandlerComplex 创建用于测试的复杂HTTP处理器
func TestHandlerComplex(db *gorm.DB, u *models.User, enableWork bool, opts ...ConfigOption) (http.Handler, Config) {
	mux := http.NewServeMux()
	c := NewConfig(db, enableWork, opts...)
	if u == nil {
		u = &models.User{
			Model: gorm.Model{ID: 888},
			Roles: []perm.Role{
				{
					Name: models.RoleAdmin,
				},
			},
		}
	}
	m := login.MockCurrentUser(u)
	mux.Handle("/page_builder/", m(c.pageBuilder))
	mux.Handle("/", m(c.pb))
	return mux, c
}

// TestHandler 创建用于测试的HTTP处理器
func TestHandler(db *gorm.DB, u *models.User) http.Handler {
	mux, _ := TestHandlerComplex(db, u, false)
	return mux
}

// TestHandlerWorker 创建用于测试的Worker HTTP处理器
func TestHandlerWorker(db *gorm.DB, u *models.User) http.Handler {
	mux, _ := TestHandlerComplex(db, u, true)
	return mux
}

// TestL18nHandler 创建用于测试的国际化HTTP处理器
func TestL18nHandler(db *gorm.DB) (http.Handler, Config) {
	mux := http.NewServeMux()
	c := NewConfig(db, false)
	c.loginSessionBuilder.Secret("test")
	c.loginSessionBuilder.Mount(mux)
	mux.Handle("/", c.pb)
	return mux, c
}

// Router 创建主路由处理器
func Router(db *gorm.DB) http.Handler {
	c := NewConfig(db, true)

	mux := http.NewServeMux()
	c.loginSessionBuilder.Mount(mux)

	// 静态文件服务 - 服务 /system/ 路径下的媒体文件
	// 文件存储在 ./public/ 目录下
	mux.Handle("/system/", http.StripPrefix("/system/", http.FileServer(http.Dir("public/system"))))

	// SEO示例
	mux.Handle("/posts/first", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var post models.Post
		db.First(&post)
		var seodata []byte
		seoBuilder.Render(post, r).MarshalHTML(r.Context(), &seodata)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><head>%s</head><body>%s</body></html>`, seodata, post.Body)
	}))

	mux.Handle("/page_builder/", c.pageBuilder)

	// 帮助中心公开站（免登录）：/help 与 /help/{slug}
	if c.helpCenter != nil {
		hcHandler := c.helpCenter.PublicHandler()
		mux.Handle(c.helpCenter.PublicPrefixPath(), hcHandler)     // /help
		mux.Handle(c.helpCenter.PublicPrefixPath()+"/", hcHandler) // /help/{slug}

		// AI 写作端点：挂在 mux 上（外层 chi 的 loginSessionBuilder.Middleware 已填充请求上下文），
		// AIHandler 内部用 AICheckAuth 校验登录 → 未登录返回 401。
		mux.Handle("/help-ai", c.helpCenter.AIHandler())
	}

	mux.Handle("/", c.pb)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(favicon)
	})

	// 提供本地 Logo SVG
	mux.HandleFunc("/assets/logo.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=31536000") // 缓存1年
		w.Write(logo)
	})

	mux.HandleFunc("/assets/logo-icon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=31536000") // 缓存1年
		w.Write(logoIcon)
	})

	// autocomplete API 端点（如 /complete/products、/complete/categories）
	if c.completeHandler != nil {
		mux.Handle(ui_demo.AutocompletePrefix+"/", c.completeHandler)
	}

	mux.Handle(exportOrdersURL, exportOrders(db))

	// FileInput 文件上传示例端点
	mux.Handle(ui_demo.FileInputUploadPath, ui_demo.FileInputUploadHandler())

	// 头像上传端点
	mux.Handle(ui_demo.AvatarUploadPath, ui_demo.AvatarUploadHandler())

	// 托管上传文件（/uploads/avatars/... → ./public/uploads/avatars/）
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("public/uploads"))))

	// sitemap和robot示例
	sitemap.SiteMap("product").RegisterRawString("https://dev.r0vx.com/admin", "/product").MountTo(mux)
	robot := sitemap.Robots()
	robot.Agent(sitemap.AlexaAgent).Allow("/product1", "/product2").Disallow("/admin")
	robot.Agent(sitemap.GoogleAgent).Disallow("/admin")
	robot.MountTo(mux)

	cr := chi.NewRouter()
	cr.Use(
		c.loginSessionBuilder.Middleware(),
		withRoles(db),
		securityMiddleware(),
	)
	cr.Mount("/", mux)
	return cr
}
