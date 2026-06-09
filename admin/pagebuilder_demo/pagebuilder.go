package pagebuilder_demo

import (
	"example/admin/pagebuilder/containers"
	"example/models"

	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/admin/l10n"
	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
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

// configPageBuilderDemo 配置 PageBuilder 演示模块
func ConfigPageBuilderDemo(b *presets.Builder, db *gorm.DB, l10nBuilder ...*l10n.Builder) *pagebuilder.Builder {
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
		PreviewOpenNewTab(true)

	// 设置 l10n 插件（多语言支持）
	if len(l10nBuilder) > 0 && l10nBuilder[0] != nil {
		pb.L10n(l10nBuilder[0])
	}

	// 注册内置演示容器
	registerHeroContainer(pb)
	registerBannerContainer(pb)
	registerRichTextContainer(pb)

	// 注册示例容器（参考 r0vx pagebuilder example）
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

	// 自动填充演示数据（表为空时）
	seedPageBuilderDemo(db)

	// 注册内置 Page 模型到 presets
	b.Use(pb)

	return pb
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
	db.Create(page)

	// 容器引用
	containers := []pagebuilder.Container{
		{PageID: page.ID, PageVersion: page.Version.Version, PageModelName: "Page", ModelName: "Banner", ModelID: banner.ID, DisplayOrder: 1, DisplayName: "Banner"},
		{PageID: page.ID, PageVersion: page.Version.Version, PageModelName: "Page", ModelName: "Hero", ModelID: hero.ID, DisplayOrder: 2, DisplayName: "Hero"},
		{PageID: page.ID, PageVersion: page.Version.Version, PageModelName: "Page", ModelName: "RichText", ModelID: richText.ID, DisplayOrder: 3, DisplayName: "RichText"},
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
