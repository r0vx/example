package pagebuilder_demo

import (
	"example/admin/pagebuilder/containers"
	"example/models"
	"example/themes"
	"io/fs"
	"log"
	"os"

	"github.com/r0vx/admin/activity"
	"github.com/r0vx/admin/l10n"
	"github.com/r0vx/admin/media"
	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/publish"
	"github.com/r0vx/admin/seo"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"gorm.io/gorm"
)

// ============================================================================
// PageBuilder 演示
// ============================================================================
//
// 展示 PageBuilder 的核心功能：
//   - 注册多种容器类型（Hero、Banner、RichText）
//   - 容器渲染函数（预览 + 编辑模式）
//   - 三栏编辑器布局（导航/预览/编辑面板）
//   - 设备切换预览
//
// ============================================================================

// Dependencies reuses the application's configured PageBuilder plugin instances.
type Dependencies struct {
	L10n      *l10n.Builder
	Activity  *activity.Builder
	Media     *media.Builder
	SEO       *seo.Builder
	Publisher *publish.Builder
}

// ConfigPageBuilderDemo 配置 PageBuilder 演示模块。
func ConfigPageBuilderDemo(
	b *presets.Builder,
	db *gorm.DB,
	dependencies Dependencies,
) *pagebuilder.Builder {
	// 迁移容器模型表
	db.AutoMigrate(
		&models.PBDemoHero{},
		&models.PBDemoBanner{},
		&models.PBDemoRichText{},
		// 示例容器模型
		&containers.WebHeader{},
		&containers.WebFooter{},
		&containers.Heading{},
		&containers.VideoBanner{},
		&containers.ImageContainer{},
		&containers.ContactForm{},
		&containers.PageTitle{},
		&containers.BrandGrid{},
		&containers.InNumbers{},
		&containers.ListContent{},
		&containers.ListContentLite{},
		&containers.ListContentWithImage{},
	)

	// 创建 PageBuilder 实例（前缀 + 数据库）
	pb := pagebuilder.New("/page_builder", db).
		AutoMigrate().
		DefaultDevice("computer").
		PreviewOpenNewTab(true).
		PreviewContainer(true).
		ExpendContainers(true).
		PublishBtnColor("primary").
		DuplicateBtnColor("secondary").
		PageStyle(h.Tag("style").
			Attr("data-pagebuilder-style", "example").
			Children(h.RawHTML(`
body { margin: 0; font-family: ui-sans-serif, system-ui, sans-serif; }
`))).
		L10n(dependencies.L10n).
		Activity(dependencies.Activity).
		Media(dependencies.Media).
		SEO(dependencies.SEO).
		Publisher(dependencies.Publisher)

	// 官方示例主题：挂了 .View() 的容器走 themes/demo 的模板，其余照旧走 RenderFunc。
	// 未设 PROD 时为开发模式：模板每次请求重新解析、静态资源不缓存。
	pb.ThemeDev(os.Getenv("PROD") == "").
		RegisterTheme("demo", demoThemeFS()).
		DefaultTheme("demo")

	registerDemoContainers(pb)

	// 主题缺文件或模板语法错误在启动期暴露，不等到用户访问才 500
	if err := pb.ValidateThemes(); err != nil {
		log.Fatalf("主题校验失败: %v", err)
	}

	// 自动填充演示数据（表为空时）
	seedPageBuilderDemo(db)
	repairDemoRelationships(db, pb)

	// 注册内置 Page 模型到 presets
	b.Use(pb)

	return pb
}

// registerDemoContainers 注册全部演示容器，测试复用同一份注册以校验示例主题。
func registerDemoContainers(pb *pagebuilder.Builder) {
	// 内置演示容器
	registerHeroContainer(pb)
	registerBannerContainer(pb)
	registerRichTextContainer(pb)

	// 示例容器（参考 r0vx pagebuilder example）
	containers.RegisterHeader(pb)
	containers.RegisterFooter(pb)
	containers.RegisterHeadingContainer(pb)
	containers.RegisterVideoBannerContainer(pb)
	containers.RegisterImageContainer(pb)
	containers.RegisterContactFormContainer(pb)
	containers.RegisterPageTitleContainer(pb)
	containers.RegisterBrandGridContainer(pb)
	containers.RegisterInNumbersContainer(pb)
	containers.RegisterListContentContainer(pb)
	containers.RegisterListContentLiteContainer(pb)
	containers.RegisterListContentWithImageContainer(pb)
}

// demoThemeDir 示例主题在 example 根下的目录。
const demoThemeDir = "themes/demo"

// demoThemeFS 优先读磁盘上的示例主题（开发时改 HTML/CSS 刷新即见）；
// 工作目录不是 example 根时读不到，回退到编译进二进制的副本。
func demoThemeFS() fs.FS {
	if _, err := os.Stat(demoThemeDir + "/theme.json"); err == nil {
		return os.DirFS(demoThemeDir)
	}
	sub, err := fs.Sub(themes.FS, "demo")
	if err != nil {
		// embed 路径写死在 themes 包里，这里失败只可能是代码写错
		panic(err)
	}
	return sub
}

// repairDemoRelationships safely relinks only unambiguous invalid demo rows.
func repairDemoRelationships(db *gorm.DB, builder *pagebuilder.Builder) {
	report, err := pagebuilder.RepairInvalidRelationships(
		db,
		builder,
		pagebuilder.RepairOptions{DemoOnly: true},
	)
	if err != nil {
		log.Printf("pagebuilder demo relationship repair failed: %v", err)
		return
	}
	if report.Repaired > 0 {
		log.Printf("pagebuilder repaired %d demo relationship(s)", report.Repaired)
	}
	for _, issue := range report.Issues {
		log.Printf(
			"pagebuilder left demo relationship %s/%d (%s) unchanged: %s",
			issue.Table,
			issue.ID,
			issue.ModelName,
			issue.Reason,
		)
	}
}

// seedPageBuilderDemo 如果表为空则插入演示数据
func seedPageBuilderDemo(db *gorm.DB) {
	var count int64
	db.Model(&pagebuilder.Page{}).Count(&count)
	if count > 0 {
		return
	}

	// 容器数据
	hero := &models.PBDemoHero{Title: "Welcome to r0vx", Subtitle: "Build enterprise admin systems with Go + Vue", BgColor: "#1a1a2e"}
	db.Create(hero)

	banner := &models.PBDemoBanner{Text: "New: PageBuilder is now available!", LinkText: "Learn more", LinkURL: "/page_builder/"}
	db.Create(banner)

	richText := &models.PBDemoRichText{Content: "<h2>About r0vx</h2><p>r0vx is a Go-based Admin framework for building enterprise management systems. It uses server-side rendered HTML with Vue 3 hydration.</p><ul><li>PageBuilder for visual page composition</li><li>CRUD presets with listing, editing, and detailing</li><li>Media library with image processing</li><li>Multi-language support (l10n)</li><li>Publishing workflow with versioning</li></ul>"}
	db.Create(richText)

	// 页面
	page := &pagebuilder.Page{Title: "Demo Homepage", Slug: "/"}
	page.Version.Version = "2024-01-01-v01"
	page.LocaleCode = "International"
	db.Create(page)

	// 容器引用
	containers := []pagebuilder.Container{
		{PageID: page.ID, PageVersion: page.Version.Version, PageModelName: "Page", ModelName: "Banner", ModelID: banner.ID, DisplayOrder: 1, DisplayName: "Banner", Locale: l10n.Locale{LocaleCode: "International"}},
		{PageID: page.ID, PageVersion: page.Version.Version, PageModelName: "Page", ModelName: "Hero", ModelID: hero.ID, DisplayOrder: 2, DisplayName: "Hero", Locale: l10n.Locale{LocaleCode: "International"}},
		{PageID: page.ID, PageVersion: page.Version.Version, PageModelName: "Page", ModelName: "RichText", ModelID: richText.ID, DisplayOrder: 3, DisplayName: "RichText", Locale: l10n.Locale{LocaleCode: "International"}},
	}
	db.Create(&containers)
}

// registerHeroContainer 注册 Hero 首屏大图容器
func registerHeroContainer(pb *pagebuilder.Builder) {
	pb.RegisterContainer("Hero").
		Model(&models.PBDemoHero{}).
		Group("Headers").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) h.HTMLComponent {
			hero := &models.PBDemoHero{
				Title:    "Welcome to r0vx",
				Subtitle: "Build enterprise admin systems with Go + Vue",
				BgColor:  "#1a1a2e",
			}
			if obj != nil {
				hero = obj.(*models.PBDemoHero)
			}

			bgStyle := "background-color: " + hero.BgColor + ";"
			if hero.BgImage != "" {
				bgStyle = "background-image: url(" + hero.BgImage + "); background-size: cover; background-position: center;"
			}

			return h.Div(
				h.Div(
					h.H1(hero.Title).Class("text-4xl font-bold text-white mb-4"),
					h.P(h.Text(hero.Subtitle)).Class("text-xl text-white/80"),
				).Class("max-w-3xl mx-auto text-center"),
			).Class("py-24 px-6").
				Style(bgStyle).
				Attr("data-container-id", input.ContainerId)
		})
}

// registerBannerContainer 注册 Banner 横幅容器
func registerBannerContainer(pb *pagebuilder.Builder) {
	pb.RegisterContainer("Banner").
		Model(&models.PBDemoBanner{}).
		Group("Headers").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) h.HTMLComponent {
			banner := &models.PBDemoBanner{
				Text:     "New feature available!",
				LinkText: "Learn more",
				LinkURL:  "#",
			}
			if obj != nil {
				banner = obj.(*models.PBDemoBanner)
			}

			return h.Div(
				h.Div(
					h.Span(banner.Text).Class("text-sm font-medium"),
					h.If(banner.LinkText != "",
						h.A(h.Text(banner.LinkText)).
							Href(banner.LinkURL).
							Class("ml-2 text-sm font-semibold underline"),
					),
				).Class("flex items-center justify-center gap-2"),
			).Class("bg-primary text-primary-foreground py-3 px-4").
				Attr("data-container-id", input.ContainerId)
		})
}

// registerRichTextContainer 注册富文本容器
func registerRichTextContainer(pb *pagebuilder.Builder) {
	pb.RegisterContainer("RichText").
		Model(&models.PBDemoRichText{}).
		Group("Content").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) h.HTMLComponent {
			rt := &models.PBDemoRichText{
				Content: "<p>This is a rich text container. Edit to add your content.</p>",
			}
			if obj != nil {
				rt = obj.(*models.PBDemoRichText)
			}

			return h.Div(
				h.Div(
					h.RawHTML(rt.Content),
				).Class("prose prose-sm max-w-none"),
			).Class("py-8 px-6 max-w-4xl mx-auto").
				Attr("data-container-id", input.ContainerId)
		})
}
