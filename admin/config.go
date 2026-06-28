package admin

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	// "github.com/aws/aws-sdk-go-v2/service/s3/types" // 移除：使用OSS抽象层
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/r0vx/admin/tiptapeditor"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/web/sse"
	"github.com/r0vx/x/i18n"
	"github.com/r0vx/x/login"
	"github.com/r0vx/x/oss"
	"github.com/r0vx/x/oss/filesystem"
	"github.com/r0vx/x/perm"
	"github.com/r0vx/x/ui/amap"
	"github.com/r0vx/x/ui/codemirror"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/tiptap"
	"github.com/r0vx/x/ui/unovis"
	"github.com/r0vx/x/ui/vueflow"
	"github.com/sunfmin/reflectutils"
	"github.com/theplant/osenv"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"example/admin/chart_demo"
	"example/admin/crud_demo"
	"example/admin/ec_demo"
	"example/admin/pagebuilder_demo"
	"example/admin/shadcn_demo"
	"example/admin/ui_demo"
	"example/admin/wizard_demo"
	"example/models"

	"github.com/r0vx/admin/activity"
	"github.com/r0vx/admin/autosync"
	"github.com/r0vx/admin/erd"
	"github.com/r0vx/admin/helpcenter"
	"github.com/r0vx/admin/l10n"
	plogin "github.com/r0vx/admin/login"
	"github.com/r0vx/admin/media"
	"github.com/r0vx/admin/media/base"
	"github.com/r0vx/admin/media/media_library"
	"github.com/r0vx/admin/microsite"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/presets/gorm2op"
	"github.com/r0vx/admin/publish"
	"github.com/r0vx/admin/redirection"
	adminrole "github.com/r0vx/admin/role"
	"github.com/r0vx/admin/seo"
	"github.com/r0vx/admin/utils"
	"github.com/r0vx/admin/worker"
)

//go:embed assets
var assets embed.FS

// PublishStorage is used to storage static pages published by page builder.
var PublishStorage oss.StorageInterface = filesystem.New("publish")

type Config struct {
	pb                  *presets.Builder
	pageBuilder         *pagebuilder.Builder
	Publisher           *publish.Builder
	loginSessionBuilder *plogin.SessionBuilder
	completeHandler     http.Handler        // autocomplete API 处理器
	helpCenter          *helpcenter.Builder // 帮助中心（公开站 handler + admin CRUD）
	sseHub              *sse.Hub            // SSE 推送中心
}

func (c *Config) GetPresetsBuilder() *presets.Builder {
	return c.pb
}

func (c *Config) GetLoginSessionBuilder() *plogin.SessionBuilder {
	return c.loginSessionBuilder
}

var (
	// s3Bucket                  = osenv.Get("S3_Bucket", "s3-bucket for media library storage", "example")
	// s3Region                  = osenv.Get("S3_Region", "s3-region for media library storage", "ap-northeast-1")
	// s3Endpoint                = osenv.Get("S3_Endpoint", "s3-endpoint for media library storage", "https://s3.ap-northeast-1.amazonaws.com")
	// s3PublishBucket           = osenv.Get("S3_Publish_Bucket", "s3-bucket for publish", "example-publish")
	// s3PublishRegion           = osenv.Get("S3_Publish_Region", "s3-region for publish", "ap-northeast-1")
	publishURL                = osenv.Get("PUBLISH_URL", "publish url", "")
	dbReset                   = osenv.Get("DB_RESET", "db reset for show count down", "")
	resetAndImportInitialData = osenv.GetBool("RESET_AND_IMPORT_INITIAL_DATA",
		"Will reset and import initial data if set to true", false)
)

type ConfigOption func(opts *configOptions)

type configOptions struct {
	StorageWrapper func(oss.StorageInterface) oss.StorageInterface
}

func WithStorageWrapper(fn func(oss.StorageInterface) oss.StorageInterface) ConfigOption {
	return func(opts *configOptions) {
		opts.StorageWrapper = fn
	}
}

func NewConfig(db *gorm.DB, enableWork bool, opts ...ConfigOption) Config {
	options := &configOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if err := db.AutoMigrate(
		&models.Post{},
		&models.InputDemo{},
		&models.User{},
		&models.ListModel{},
		&perm.Role{},
		&perm.DefaultDBPolicy{},
		&models.Customer{},
		&models.Address{},
		&models.Phone{},
		&models.MembershipCard{},
		&models.Product{},
		&models.Order{},
		&models.Category{},
		&models.SubCategory{},
		&models.Project{},
		&models.Task{},
		&ui_demo.NotifDemo{},
		&models.AutocompleteDemo{},
		&models.DialogDemo{},
		&models.FilterDemo{},
		&models.VehicleFilterDemo{},
		&models.TreeSelectDemo{},
		&models.TreeSelectDemoSeries{},
		&models.GraphDemo{},
	); err != nil {
		panic(err)
	}

	// @snippet_begin(ActivityExample)
	ab := activity.New(db, func(ctx context.Context) (*activity.User, error) {
		u := ctx.Value(login.UserKey).(*models.User)
		return &activity.User{
			ID:     fmt.Sprint(u.ID),
			Name:   u.Name,
			Avatar: "",
		}, nil
	}).
		TablePrefix("cms_").
		AutoMigrate()

	// ab.Model(l).SkipDelete().SkipCreate()
	// @snippet_end

	// media_oss.Storage = s3.New(&s3.Config{
	// 	Bucket:   s3Bucket,
	// 	Region:   s3Region,
	// 	ACL:      string(oss.ACLBucketOwnerFullControl),
	// 	Endpoint: s3Endpoint,
	// })
	// s3Client := s3.New(&s3.Config{
	// 	Bucket:   s3PublishBucket,
	// 	Region:   s3PublishRegion,
	// 	ACL:      string(oss.ACLBucketOwnerFullControl),
	// 	Endpoint: publishURL,
	// })
	// PublishStorage = microsite_utils.NewClient(s3Client)
	// if options.StorageWrapper != nil {
	// 	PublishStorage = options.StorageWrapper(PublishStorage)
	// }
	b := presets.New().DataOperator(gorm2op.DataOperator(db))
	// panic 主动告警：配了 PANIC_WEBHOOK_URL 才启用（飞书/钉钉/Slack 机器人均可），空则零行为
	if u := os.Getenv("PANIC_WEBHOOK_URL"); u != "" {
		b.PanicNotifier(presets.WebhookPanicNotifier(u))
	}
	// 移动端底部 Tab 栏：首页 / 订单 / 菜单 / 我的（菜单项 URL 空 → 点击切换侧栏 Sheet）
	b.BottomNav(
		presets.BottomNavItem{Icon: "home", URL: "/"},
		presets.BottomNavItem{Icon: "shopping-cart", URL: "/orders"},
		presets.BottomNavItem{Icon: "menu"},
		presets.BottomNavItem{Icon: "user", URL: "/profile"},
	)
	// 数据隔离全局 resolver：非 admin 只看自己的数据，admin 看全部。
	// ownerValue 用 string（匹配 activity.ActivityLog.UserID 的 string 列）。
	b.DataScopeResolver(func(ctx *web.EventContext) (any, bool) {
		u := getCurrentUser(ctx.R)
		if u == nil {
			return nil, false
		}
		return fmt.Sprint(u.GetID()), slices.Contains(u.GetRoles(), models.RoleAdmin)
	})
	defer b.Build()

	// 添加 shadcn-vue 组件资源
	// b.ExtraAsset("/shadcnx.js", "text/javascript", shadcn.JSComponentsPack())
	// b.ExtraAsset("/shadcnx.css", "text/css", shadcn.CSSComponentsPack())

	// 添加 codemirror 代码编辑器资源（懒加载：744KB 重 chunk 改 LazyAsset，仅渲染 vue-codemirror 的页/抽屉首次用时加载；
	// eager 微 stub 注册异步组件，CSS 保持 eager（小，保证样式）。即使在编辑抽屉 AJAX 也能解析）
	b.LazyAsset("/codemirror.js", "text/javascript", codemirror.JSComponentsPack())
	b.ExtraAsset("/codemirror.css", "text/css", codemirror.CSSComponentsPack())
	b.ExtraAsset("/codemirror-stub.js", "text/javascript", codemirror.StubJSPack())

	// 添加 tiptap 富文本编辑器资源（懒加载：437KB 重 chunk 改 LazyAsset，仅渲染 tiptap-editor 的页/抽屉首次用时加载；
	// eager 微 stub 注册异步组件，CSS + MediaBox 桥接保持 eager。MediaBox bodyObserver 已支持懒渲染 wrapper 晚到）
	b.LazyAsset("/tiptap.js", "text/javascript", tiptap.JSComponentsPack())
	b.ExtraAsset("/tiptap.css", "text/css", tiptap.CSSComponentsPack())
	b.ExtraAsset("/tiptap-stub.js", "text/javascript", tiptap.StubJSPack())
	// tiptap MediaBox 桥接脚本（DOM/事件驱动，独立于 tiptap.js 加载时机）
	b.ExtraAsset("/tiptap-mediabox.js", "text/javascript", tiptapeditor.JSComponentsPack())
	// 帮助中心 Summary AI 摘要 helper 脚本
	b.ExtraAsset("/helpcenter-ai.js", "text/javascript", helpcenter.AIJSPack())

	// 添加 unovis 数据可视化组件资源（懒加载：1.93MB 重 chunk + CSS 改 LazyAsset，eager 微 stub 把 6 个 unovis tag
	// 注册为异步组件，首次渲染时按需加载。异步组件覆盖菜单 pushState 导航/抽屉/初始全部路径）
	b.LazyAsset("/unovis.js", "text/javascript", unovis.JSComponentsPack())
	b.LazyAsset("/unovis.css", "text/css", unovis.CSSComponentsPack())
	b.ExtraAsset("/unovis-stub.js", "text/javascript", unovis.StubJSPack())

	// 添加 Vue Flow 通用画布组件资源（eager：@vue-flow/core 内部 provide/inject + store 对异步组件加载间隙敏感，
	// 懒化后首次菜单 pushState 导航画布渲染失败；故保持 eager（323KB/页换正确性）。unovis 不受此限已懒化）
	b.ExtraAsset("/vueflow.js", "text/javascript", vueflow.JSComponentsPack())
	b.ExtraAsset("/vueflow.css", "text/css", vueflow.CSSComponentsPack())

	// 添加高德地图选点组件资源
	b.ExtraAsset("/amap.js", "text/javascript", web.ComponentsPack(amap.JSComponentsPack()))
	b.ExtraAsset("/amap.css", "text/css", web.ComponentsPack(amap.CSSComponentsPack()))

	initPermission(b, db)

	b.GetI18n().
		SupportLanguages(language.English, language.SimplifiedChinese, language.Japanese).
		RegisterForModule(language.SimplifiedChinese, presets.ModelsI18nModuleKey, Messages_zh_CN_ModelsI18nModuleKey).
		RegisterForModule(language.English, I18nExampleKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nExampleKey, Messages_zh_CN).
		GetSupportLanguagesFromRequestFunc(func(r *http.Request) []language.Tag {
			// // Example:
			// user := getCurrentUser(r)
			// var supportedLanguages []language.Tag
			// for _, role := range user.GetRoles() {
			//	switch role {
			//	case "English Group":
			//		supportedLanguages = append(supportedLanguages, language.English)
			//	case "Chinese Group":
			//		supportedLanguages = append(supportedLanguages, language.SimplifiedChinese)
			//	}
			// }
			// return supportedLanguages
			return b.GetI18n().GetSupportLanguages()
		})
	mediab := media.New(db).AutoMigrate().Activity(ab).CurrentUserID(func(ctx *web.EventContext) (id uint) {
		u := getCurrentUser(ctx.R)
		if u == nil {
			return
		}
		return u.ID
	}).Searcher(func(db *gorm.DB, ctx *web.EventContext) *gorm.DB {
		u := getCurrentUser(ctx.R)
		if u == nil {
			return db
		}
		if rs := u.GetRoles(); !slices.Contains(rs, models.RoleAdmin) && !slices.Contains(rs, models.RoleManager) {
			return db.Where("user_id = ?", u.ID)
		}
		return db
	})
	defer func() {
		mediab.GetPresetsModelBuilder().Use(ab)
		seoBuilder.GetPresetsModelBuilder().Use(ab)
	}()

	l10nBuilder := l10n.New(db)
	l10nBuilder.
		Activity(ab).
		RegisterLocales("International", "international", "International", l10n.InternationalSvg).
		RegisterLocales("Japan", "jp", "Japan", l10n.JapanSvg).
		RegisterLocales("China", "cn", "China", l10n.ChinaSvg).
		SupportLocalesFunc(func(R *http.Request) []string {
			return l10nBuilder.GetSupportLocaleCodes()[:]
		})
	publisher := publish.New(db, PublishStorage).
		ContextValueFuncs(l10nBuilder.ContextValueProvider)
	redirectionBuilder := redirection.New(PublishStorage, db, publisher).AutoMigrate()
	utils.Install(b)

	publisher.Activity(ab)

	// media_view.MediaLibraryPerPage = 3
	// vips.UseVips(vips.Config{EnableGenerateWebp: true})
	configureSeo(b, db, l10nBuilder.GetSupportLocaleCodes()...)
	configMenuOrder(b)

	configPost(b, db, publisher, ab, seoBuilder)

	// 帮助中心：admin CRUD + 公开站 handler + 种子数据
	helpCenterBuilder := ui_demo.ConfigHelpCenter(b, db, publisher)

	roleBuilder := adminrole.New(db).
		// 动作列表：默认 CRUD + 帮助中心 AI 教程生成（helpcenter.PermHelpCenterGen）。
		// 该权限是「动作」而非「资源」：列表页 ✨ AI 教程 按钮的可见性由
		// mb.Verifier().Do(PermHelpCenterGen) 校验，故需在动作下拉里显式可选（授予后配合资源即生效）。
		Actions([]*shadcn.DefaultOptionItem{
			{Text: "All", Value: "*"},
			{Text: "List", Value: presets.PermList},
			{Text: "Get", Value: presets.PermGet},
			{Text: "Create", Value: presets.PermCreate},
			{Text: "Update", Value: presets.PermUpdate},
			{Text: "Delete", Value: presets.PermDelete},
			{Text: "HelpCenterAIGen", Value: helpcenter.PermHelpCenterGen},
		}).
		// Resources 选项改由框架自动枚举（admin v0.4.0：AutoEnumerate 默认开）——
		// 遍历所有注册模型 + 细粒度闸（f_/fl_/fm_/ft_）按模型分组自动生成，免维护这一长串。
		// （旧手维护列表里 *:nested-field-demos:* / *:shadcn-*:* 等 kebab 写法本就与运行时 snake 资源错配，
		//  正是自动枚举要消除的隐患。）如需限定仍可显式 .Resources(...) 覆盖。
		Matrix(true). // 只读权限矩阵总览页 /role-matrix（默认即开，此处显式标注）
		AfterInstall(func(pb *presets.Builder, mb *presets.ModelBuilder) error {
			mb.Listing().SearchFunc(func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
				u := getCurrentUser(ctx.R)
				qdb := db
				// If the current user doesn't has 'admin' role, do not allow them to view admin and manager roles
				// We didn't do this on permission because of we are not supporting the permission on listing page
				if currentRoles := u.GetRoles(); !slices.Contains(currentRoles, models.RoleAdmin) {
					qdb = db.Where("name NOT IN (?)", []string{models.RoleAdmin, models.RoleManager})
				}
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			})
			return nil
		})
	var w *worker.Builder
	if enableWork {
		w = worker.New(db)
		defer w.Listen()
		addJobs(w)
		crud_demo.ConfigProduct(b, db, w, publisher)
		b.Use(w.Activity(ab))
	}
	categoryMB := crud_demo.ConfigCategory(b, db, publisher)
	crud_demo.ConfigSubCategory(b, db) // 演示现有拖拽排序 SortBuilder

	// 给「分类」列表页挂上 ✨ AI 教程 生成按钮（权限 helpcenter:ai-gen 控制）
	if helpCenterBuilder != nil {
		helpCenterBuilder.EnableGenFor(categoryMB)
	}

	// Use m to customize the model, Or config more models here.

	// type Setting struct{}
	// sm := b.Model(&Setting{})
	// sm.RegisterEventFunc(pages.LogInfoEvent, pages.LogInfo)
	// sm.Listing().PageFunc(pages.Settings(db))

	// FIXME: list editor does not support use in page func
	// type ListEditorExample struct{}
	// leem := b.Model(&ListEditorExample{}).Label("List Editor Example")
	// pf, sf := pages.ListEditorExample(db, b)
	// leem.Listing().PageFunc(pf)
	// leem.RegisterEventFunc("save", sf)

	crud_demo.ConfigNestedFieldDemo(b, db, ab)
	crud_demo.ConfigMembershipCard(b, db, ab)

	pageBuilder := pagebuilder_demo.ConfigPageBuilderDemo(b, db, l10nBuilder)

	configListModel(b, ab, publisher)

	microb := microsite.New(db).Publisher(publisher)

	l10nBuilder.Activity(ab)
	l10nM, l10nVM := crud_demo.ConfigL10nModel(db, b)
	l10nM.Use(l10nBuilder)
	l10nVM.Use(l10nBuilder)

	loginSessionBuilder := initLoginSessionBuilder(db, b, ab)

	configBrand(b)

	// 首页仪表盘
	dashboardBuilder := configureDashboard(db)
	b.HomePageFunc(dashboardBuilder.PageFunc()).
		NotFoundPageLayoutConfig(&presets.LayoutConfig{
			NotificationCenterInvisible: true,
		}).ProgressBarColor("success")
	dashboardBuilder.RegisterEvents(b.GetWebBuilder())

	// 主题切换：注册可用主题 + 设置默认主题
	// 5 套 r0vx 内置 + 42 套 tweakcn 移植预设（对标 tweakcn）
	b.Themes(append([]*presets.Theme{
		presets.ThemeNeutral,
		presets.ThemeBlue,
		presets.ThemeGreen,
		presets.ThemeViolet,
		presets.ThemeRose,
		presets.ThemeNavy,
		presets.ThemeAzure,
	}, presets.TweakcnPresets...)...).DefaultTheme("neutral")

	profileBuilder := configProfile(db, ab, loginSessionBuilder)

	ui_demo.ConfigInputDemo(b, db, ab, w)
	shadcn_demo.Configure(b)
	chart_demo.ConfigGraphDemo(b)
	chart_demo.ConfigNetworkGraphDemo(b, db)
	chart_demo.ConfigScatterPlotDemo(b, db)
	chart_demo.ConfigTreemapDemo(b, db)
	chart_demo.ConfigChartRealtimeDemo(b) // 实时图表范式对比（RefreshInterval 定时刷新族：A/B/环形）

	crud_demo.ConfigProject(b, db)
	ec_demo.ConfigECDashboard(b, db)
	ec_demo.ConfigureDemoCase(b, db)
	ui_demo.ConfigNotifDemo(b, db)
	ui_demo.ConfigNotificationCenter(b, db)
	completeHandler := ui_demo.ConfigAutocompleteDemo(b, db)
	ui_demo.ConfigTiptapDemo(b, db)
	ui_demo.ConfigFileInputDemo(b)
	ui_demo.ConfigVueFlowDemo(b)
	ui_demo.ConfigDialogDemo(b, db)
	ui_demo.ConfigNonIDPKDemo(b, db)
	ui_demo.ConfigEditingActionsDemo(b, db)
	ui_demo.ConfigRowRefreshDemo(b, db)
	ui_demo.ConfigRelayPaginationDemo(b, db)
	ui_demo.ConfigPermResourceEventDemo(b, db) // 自定义权限资源 + 裸事件鉴权演示
	wizard_demo.ConfigWizardDemo(b, db)
	wizard_demo.ConfigWizardDeclarativeDemo(b, db)
	wizard_demo.ConfigWizardFullPageDemo(b, db)
	// SSE 推送中心（须在用到 hub 的 demo 之前创建）：DataScope 隔离模型实时刷新 + 通知实时推送。
	b.SSEUserID(func(r *http.Request) (string, bool) {
		u := getCurrentUser(r)
		if u == nil {
			return "", false
		}
		return fmt.Sprint(u.GetID()), true
	})
	sseHub := sse.New(
		sse.Identity(b.SSEIdentity),
		sse.Replay(64, 5*time.Minute),
	)
	b.SSEHub(sseHub)

	crud_demo.ConfigOrder(b, db, sseHub)        // orders 列表实时刷新需 hub 广播，故移到 hub 创建之后
	chart_demo.ConfigStreamChartDemo(b, sseHub) // 滚动流图表 demo（StreamOn 追加流：SSE 推点+客户端追加，需 hub）
	ui_demo.ConfigNotificationDemo(b, db, sseHub)
	ui_demo.ConfigActionEnhanceDemo(b, db, sseHub)
	ui_demo.ConfigDisableRowClickDemo(b, db)
	ui_demo.ConfigFilterDemo(b, db)
	ui_demo.ConfigVehicleFilterDemo(b, db)
	ui_demo.ConfigTreeSelectDemo(b, db)
	ui_demo.ConfigListingWrapDemo(b, db)
	ui_demo.ConfigTreeListingDemo(b, db)
	ui_demo.ConfigCrossTreeListingDemo(b, db)
	ui_demo.ConfigAvatarUploadDemo(b, db)
	ui_demo.ConfigSiteSettingDemo(b, db) // Singleton(true) 单例配置页范例
	ConfigDataScopeDemo(b, db)           // 数据隔离：同角色跨表不同隔离字段（agent_id/user_id/parent_id）
	ui_demo.ConfigRecordGraphDemo(b, db) // 记录关系图（ego 图）演示
	// ERD 数据模型图
	erd.New(b, db).MountTo(b, "/erd", "数据模型图")

	crud_demo.ConfigUser(b, ab, db, publisher, loginSessionBuilder)
	b.Use(
		mediab,
		microb,
		ab,
		publisher,
		l10nBuilder,
		roleBuilder,
		loginSessionBuilder,
		profileBuilder,
		redirectionBuilder,
	)

	if resetAndImportInitialData {
		tbs := GetNonIgnoredTableNames(db)
		EmptyDB(db, tbs)
		InitDB(db, tbs)
	}

	return Config{
		pb:                  b,
		pageBuilder:         pageBuilder,
		Publisher:           publisher,
		loginSessionBuilder: loginSessionBuilder,
		completeHandler:     completeHandler,
		helpCenter:          helpCenterBuilder,
		sseHub:              sseHub,
	}
}

func configListModel(b *presets.Builder, ab *activity.Builder, publisher *publish.Builder) *presets.ModelBuilder {
	mb := b.Model(&models.ListModel{})
	defer mb.Use(ab, publisher)
	{
		mb.Listing("ID", "Title", "Status")
		mb.Editing("Title")

		detailing := mb.Detailing(publish.VersionsPublishBar, "Title", "DetailPath", "ListPath").Drawer(true)

		titleSection := presets.NewSectionBuilder(mb, "Title").Editing("Title")

		detailPathSection := presets.NewSectionBuilder(mb, "DetailPath").
			ComponentFunc(
				func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (r h.HTMLComponent) {
					this := obj.(*models.ListModel)

					if this.Status.Status != publish.StatusOnline {
						return nil
					}

					domain := PublishStorage.GetEndpoint(ctx.R.Context())
					if this.OnlineUrl == "" {
						return nil
					}

					return h.Div(
						h.Label(i18n.PT(ctx.R, presets.ModelsI18nModuleKey, mb.Info().Label(), field.Label)).Class("text-sm font-medium mb-1 block"),
						h.A(h.Text(this.OnlineUrl)).Href(domain+this.OnlineUrl).Class("text-primary hover:underline"),
					).Class("mb-4")
				},
			)

		listPathSection := presets.NewSectionBuilder(mb, "ListPath").ComponentFunc(
			func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (r h.HTMLComponent) {
				this := obj.(*models.ListModel)

				if this.Status.Status != publish.StatusOnline || this.PageNumber == 0 {
					return nil
				}

				domain := PublishStorage.GetEndpoint(ctx.R.Context())
				if this.OnlineUrl == "" {
					return nil
				}

				p := this.GetListUrl(strconv.Itoa(this.PageNumber))
				return h.Div(
					h.Label(i18n.PT(ctx.R, presets.ModelsI18nModuleKey, mb.Info().Label(), field.Label)).Class("text-sm font-medium mb-1 block"),
					h.A(h.Text(p)).Href(domain+p).Class("text-primary hover:underline"),
				).Class("mb-4")
			},
		)
		detailing.Section(titleSection, detailPathSection, listPathSection)
	}
	return mb
}

func configMenuOrder(b *presets.Builder) {
	// 使用 MenuSection 创建带标签的菜单分区
	// 参考 shadcn-admin-demo 的分组样式

	// 设计原则：每个菜单可见模型都归入某个 MenuGroup 作为子项（子项无图标=干净，
	// 与 shadcn-admin-demo 一致），不再留裸顶层项，避免掉进底部未分类的扁平列表。
	// uriName 以实际注册名为准（系统模块多为结构体名推导的 kebab 复数）。

	// Main 分区 - 核心业务功能
	mainSection := b.MenuSection("Main").Items(
		b.MenuGroup("EC Management").SubItems(
			"ec-dashboard",
			"orders",
			"products",
			"categories",
			"sub-categories",
			"customers",
			"membership-cards",
		).Icon("shopping-cart"),
		b.MenuGroup("Content").SubItems(
			"posts",        // Post → 文章
			"articles",     // helpcenter Article → 帮助文档
			"seo-settings", // seo SEOSetting → SEO 管理
			"micro-sites",  // microsite MicroSite → 微站点
		).Icon("file-text"),
		b.MenuGroup("Project Management").SubItems(
			"projects",
			"tasks",
		).Icon("folder"),
		b.MenuGroup("Localization").SubItems(
			"l10n-models",
			"l10n-model-with-versions",
		).Icon("languages"),
		b.MenuGroup("Page Builder").SubItems(
			"page-builder-pages",
			"shared-containers",
			"demo-containers",
			"page-builder-templates",
			"page-builder-categories",
		).Icon("layout-grid"),
	)

	// Analytics 分区 - 图表与数据可视化
	analyticsSection := b.MenuSection("Analytics").Items(
		b.MenuGroup("Charts").SubItems(
			"graph-demos",
			"network-graph-demo",
			"scatter-plot-demo",
			"treemap-demo",
			"shadcn-chart",
			"chart-realtime-demo", // 定时刷新族（RefreshInterval：A/B/环形）
			"stream-chart-demo",   // 追加流（StreamOn：SSE 推点+客户端追加左滑）
		).Icon("bar-chart"),
	)

	// System 分区 - 系统管理功能
	systemSection := b.MenuSection("System").Items(
		b.MenuGroup("User & Access").SubItems(
			"users",
			"roles",          // adminrole perm.Role
			"login-sessions", // plogin 会话（若 InMenu(false) 自动跳过）
		).Icon("users"),
		b.MenuGroup("Operations").SubItems(
			"workers",         // worker
			"activity-logs",   // activity ActivityLog
			"media-libraries", // media MediaLibrary
			"redirects",       // redirection Redirect
			"site-setting",    // 站点设置单例
		).Icon("settings"),
	)

	// Components 分区 - 示例组件（对齐 demo：单分区下多个可折叠分组）
	componentsSection := b.MenuSection("Components").Items(
		b.MenuGroup("Form & Input").SubItems(
			"input-demos",
			"file-input-demo",
			"autocomplete-demos",
			"tiptap-demos",
			"filter-demos",
			"vehicle-filter-demos",
			"tree-select-demos",
		).Icon("list"),
		b.MenuGroup("Data & Listing").SubItems(
			"list-models",
			"tree-listing-demo",
			"cross-tree-listing-demo",
			"cross-tree-articles",
			"listing-wrap-demo",
			"readonly-list",
			"editing-actions-demos",
			"action-enhance-demo",
			"dialog-demos",
			"notif-demos",
			"notification-demo", // SSE Toast 通知页
			"avatar-upload-demo",
			"sc-member-demo",        // 非 ID 主键
			"row-refresh-demo",      // 行级局部刷新
			"relay-pagination-demo", // relay 游标分页
		).Icon("table"),
		b.MenuGroup("Workflow & Scope").SubItems(
			"wizard-demos",
			"wizard-declarative",
			"scope-agent-deal",
			"scope-user-note",
			"scope-org-doc",
			"perm-resource-event-demo", // 自定义权限资源 + 裸事件鉴权
			"demo-cases",
		).Icon("workflow"),
		b.MenuGroup("Vue Flow").SubItems(
			"vueflow-demo",          // 基础通用画布
			"vueflow-dagre-demo",    // dagre 自动布局
			"vueflow-status-demo",   // status 状态卡片
			"vueflow-resizer-demo",  // 节点缩放
			"vueflow-toolbar-demo",  // 节点工具栏
			"vueflow-dnd-demo",      // 拖放建节点
			"vueflow-edges-demo",    // 边类型/箭头
			"vueflow-math-demo",     // 数学运算流
			"vueflow-viewport-demo", // 视口/截图
			"vueflow-teleport-demo", // 传送节点
			"vueflow-dragaids-demo", // 拖拽辅助
			"vueflow-conn-demo",     // 连线进阶
			"record-graph",          // 记录关系图（ego 图）
			"erd",                   // 数据模型图（ERD）
			"rg-demo-users",         // RG 用户（关系图数据源）
			"rg-demo-orders",        // RG 订单（关系图数据源）
		).Icon("git-fork"),
		b.MenuGroup("Shadcn UI").SubItems(
			"shadcn-basic-inputs",
			"shadcn-selections",
			"shadcn-dialog",
			"shadcn-table",
			"shadcn-data-table",
			"shadcn-lazy-portals",
			"shadcn-grid",
			"shadcn-list",
			"shadcn-popover-menu",
			"shadcn-sheet-drawer",
			"shadcn-variant-subform",
			"shadcn-progress",
			"shadcn-sidebar-demo",
			"shadcn-invoice-list",
			"shadcn-new-components",
			"shadcn-range-picker",
			"shadcn-form-field",
			"shadcn-autocomplete",
			"shadcn-display-components",
			"shadcn-filter",
			"shadcn-tree-view",
			"shadcn-cascader",
			"shadcn-timeline",
			// shadcn-chart 已归入 Analytics / Charts 组（图表 demo 集中管理），此处不再重复
			"shadcn-admin-demo",
		).Icon("component"),
	)

	b.MenuOrder(
		mainSection,
		analyticsSection,
		systemSection,
		componentsSection,
	)

	b.SidebarBrandFunc(func(ctx *web.EventContext) h.HTMLComponent {
		// 系统名（回退到默认品牌名）
		name := b.GetBrandTitle()
		if name == "" {
			name = "r0vx Admin"
		}
		// 副标题：当前登录用户名（折叠态随文字块一起隐藏）
		subtitle := ""
		if u := getCurrentUser(ctx.R); u != nil {
			subtitle = u.Name
		}
		// 参考 payManage 排版：圆角方形 logo 块（bg-primary 实色 + 白色 logo-icon）+ 双行文字
		return h.A().Href("/").Class("flex w-full items-center gap-2 px-2 py-1 group-data-[collapsible=icon]:px-0 group-data-[collapsible=icon]:justify-center").Children(
			// 圆角方形 logo 块（折叠态常驻）
			h.Div(
				h.Img("/assets/logo-icon.svg").Class("size-5").Style("filter: brightness(0) invert(1);"),
			).Class("flex aspect-square size-8 items-center justify-center rounded-lg bg-primary shrink-0"),
			// 文字块：系统名(粗) + 当前用户名(小灰)；折叠态隐藏只剩图标
			h.Div(
				h.Span(name).Class("block truncate text-sm font-semibold leading-tight"),
				h.Span(subtitle).Class("block truncate text-xs text-muted-foreground leading-tight"),
			).Class("grid flex-1 text-left group-data-[collapsible=icon]:hidden"),
		)
	})
}

func configProfile(db *gorm.DB, ab *activity.Builder, lsb *plogin.SessionBuilder) *plogin.ProfileBuilder {
	return plogin.NewProfileBuilder(
		func(ctx context.Context) (*plogin.Profile, error) {
			evCtx := web.MustGetEventContext(ctx)
			u := getCurrentUser(evCtx.R)
			if u == nil {
				return nil, perm.PermissionDenied
			}
			notifiCounts, err := ab.GetNotesCounts(ctx, "", nil)
			if err != nil {
				return nil, err
			}
			user := &plogin.Profile{
				ID:   fmt.Sprint(u.ID),
				Name: u.Name,
				// Avatar: "",
				Roles:  u.GetRoles(),
				Status: strcase.ToCamel(u.Status),
				Fields: []*plogin.ProfileField{
					{Name: "Email", Key: "email", Value: u.Account, Icon: "mail", Editable: true},
					{Name: "Company", Key: "company", Value: u.Company, Icon: "building", Editable: true},
				},
				NotifCounts: notifiCounts,
			}
			user.Avatar = u.Avatar
			if user.Avatar == "" && u.OAuthAvatar != "" {
				user.Avatar = u.OAuthAvatar // 回退到 OAuth 头像
			}
			return user, nil
		},
		func(ctx context.Context, newName string) error {
			evCtx := web.MustGetEventContext(ctx)
			u := getCurrentUser(evCtx.R)
			if u == nil {
				return perm.PermissionDenied
			}
			u.Name = newName
			if err := db.Save(u).Error; err != nil {
				return errors.Wrap(err, "failed to update user name")
			}
			return nil
		},
	).SessionBuilder(lsb).
		AvatarUpload(ui_demo.AvatarUploadPath).
		UpdateProfileFunc(func(ctx context.Context, changes map[string]string) error {
			evCtx := web.MustGetEventContext(ctx)
			u := getCurrentUser(evCtx.R)
			if u == nil {
				return perm.PermissionDenied
			}
			if v, ok := changes["name"]; ok {
				u.Name = v
			}
			if v, ok := changes["email"]; ok {
				u.Account = v // 注意：Account 即登录邮箱，改动需保证唯一
			}
			if v, ok := changes["company"]; ok {
				u.Company = v
			}
			if v, ok := changes["avatar"]; ok {
				u.Avatar = v
			}
			if err := db.Save(u).Error; err != nil {
				return errors.Wrap(err, "failed to update profile")
			}
			return nil
		}) // .DisableNotification(true).LogoutURL(lsb.GetLoginBuilder().LogoutURL)
	// 		CustomizeButtons(func(ctx context.Context, buttons ...h.HTMLComponent) ([]h.HTMLComponent, error) {
	// 	evCtx := web.MustGetEventContext(ctx)
	// 	u := getCurrentUser(evCtx.R)
	// 	if u == nil {
	// 		return nil, perm.PermissionDenied
	// 	}
	// 	if u.GetAccountName() == loginInitialUserEmail {
	// 		return buttons, nil
	// 	}
	// 	msgr := i18n.MustGetModuleMessages(evCtx.R, I18nExampleKey, Messages_en_US).(*Messages)
	// 	return slices.Concat([]h.HTMLComponent{
	// 		v.VBtn(msgr.ChangePassword).Variant(v.VariantTonal).Color(v.ColorSecondary).OnClick(plogin.OpenChangePasswordDialogEvent),
	// 	}, buttons), nil
	// }).
	// PrependCompos(func(ctx context.Context, profileCompo *plogin.ProfileCompo) ([]h.HTMLComponent, error) {
	// 	return []h.HTMLComponent{
	// 		web.Listen(presets.NotifModelsUpdated(&models.User{}), stateful.ReloadAction(ctx, profileCompo, nil).Go()),
	// 	}, nil
	// })
}

func configBrand(b *presets.Builder) {
	b.BrandFunc(func(ctx *web.EventContext) h.HTMLComponent {
		// msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)

		// 完整Logo SVG (桌面端)
		fullLogoSVG := `<svg width="164" height="24" viewBox="0 0 164 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M0.000210938 2.84508C0.000210938 5.29849 0.800211 6.49909 4.00021 8.6915C7.55021 11.1449 8.00021 11.9279 8.00021 15.5819C8.00021 20.6975 9.50021 23.7252 12.0502 23.7252C13.8002 23.7252 14.0002 22.9944 14.0002 15.6341L13.9502 7.54309L7.55021 3.62808C-0.149789 -1.06994 0.000210938 -1.06994 0.000210938 2.84508Z" fill="white"/>
<path d="M22.2502 3.68028L16.0502 7.5431L16.0002 15.6341C16.0002 22.9944 16.2002 23.7252 17.9502 23.7252C20.5002 23.7252 22.0002 20.6975 22.0002 15.5819C22.0002 11.9279 22.4502 11.1449 26.0002 8.6915C29.2002 6.49909 30.0002 5.29849 30.0002 2.84508C30.0002 -1.06994 29.8502 -1.01774 22.2502 3.68028Z" fill="white"/>
<path d="M44.3469 5.0042H46.5793L49.8664 14.3997L53.1535 5.0042H55.386L50.7629 17.8011H48.9699L44.3469 5.0042ZM43.1428 5.0042H45.3664L45.7707 14.1624V17.8011H43.1428V5.0042ZM54.3664 5.0042H56.5988V17.8011H53.9621V14.1624L54.3664 5.0042ZM64.2717 7.19267L60.7912 17.8011H57.9875L62.7424 5.0042H64.5266L64.2717 7.19267ZM67.1633 17.8011L63.674 7.19267L63.3928 5.0042H65.1945L69.9758 17.8011H67.1633ZM67.0051 13.0374V15.1028H60.2463V13.0374H67.0051ZM75.7942 5.0042V17.8011H73.1662V5.0042H75.7942ZM79.7317 5.0042V7.06963H69.2903V5.0042H79.7317ZM89.8215 15.7444V17.8011H83.01V15.7444H89.8215ZM83.8713 5.0042V17.8011H81.2346V5.0042H83.8713ZM88.9338 10.2161V12.22H83.01V10.2161H88.9338ZM89.8127 5.0042V7.06963H83.01V5.0042H89.8127ZM91.3596 5.0042H96.132C97.1106 5.0042 97.9514 5.15068 98.6545 5.44365C99.3635 5.73662 99.9084 6.17021 100.289 6.74443C100.67 7.31865 100.861 8.0247 100.861 8.86259C100.861 9.54814 100.743 10.137 100.509 10.6292C100.28 11.1155 99.9553 11.5228 99.5334 11.8509C99.1174 12.1731 98.6281 12.431 98.0656 12.6243L97.2307 13.0638H93.0822L93.0647 11.0071H96.1496C96.6125 11.0071 96.9963 10.9251 97.301 10.761C97.6057 10.597 97.8342 10.3685 97.9865 10.0755C98.1447 9.78252 98.2238 9.44267 98.2238 9.05595C98.2238 8.6458 98.1477 8.29131 97.9953 7.99248C97.843 7.69365 97.6115 7.46513 97.301 7.30693C96.9904 7.14873 96.6008 7.06963 96.132 7.06963H93.9963V17.8011H91.3596V5.0042ZM98.5139 17.8011L95.5959 12.097L98.382 12.0794L101.335 17.678V17.8011H98.5139ZM105.695 5.0042V17.8011H103.067V5.0042H105.695ZM113.526 7.19267L110.045 17.8011H107.241L111.996 5.0042H113.78L113.526 7.19267ZM116.417 17.8011L112.928 7.19267L112.647 5.0042H114.448L119.23 17.8011H116.417ZM116.259 13.0374V15.1028H109.5V13.0374H116.259ZM128.643 15.7444V17.8011H122.2V15.7444H128.643ZM123.053 5.0042V17.8011H120.416V5.0042H123.053ZM134.988 12.5892H131.341V11.6663H134.988C135.75 11.6663 136.368 11.5433 136.843 11.2972C137.323 11.0452 137.672 10.7054 137.889 10.2776C138.112 9.8499 138.223 9.36943 138.223 8.83623C138.223 8.31474 138.112 7.83427 137.889 7.39482C137.672 6.95537 137.323 6.60381 136.843 6.34013C136.368 6.0706 135.75 5.93584 134.988 5.93584H131.719V17.8011H130.638V5.0042H134.988C135.926 5.0042 136.714 5.16533 137.353 5.48759C137.997 5.80986 138.484 6.2581 138.812 6.83232C139.14 7.40654 139.304 8.06865 139.304 8.81865C139.304 9.60381 139.14 10.2806 138.812 10.8489C138.484 11.4114 138 11.8421 137.362 12.1409C136.723 12.4397 135.932 12.5892 134.988 12.5892ZM141.703 5.0042H145.861C146.745 5.0042 147.513 5.14775 148.163 5.43486C148.814 5.72197 149.315 6.14677 149.666 6.70927C150.024 7.26591 150.202 7.95146 150.202 8.76591C150.202 9.36943 150.073 9.91728 149.816 10.4095C149.564 10.9017 149.215 11.3147 148.77 11.6487C148.324 11.9769 147.809 12.2024 147.223 12.3255L146.845 12.4661H142.406L142.389 11.5433H146.107C146.775 11.5433 147.331 11.4144 147.777 11.1565C148.222 10.8987 148.556 10.5589 148.779 10.137C149.007 9.70927 149.121 9.25224 149.121 8.76591C149.121 8.18584 148.995 7.68486 148.743 7.26299C148.497 6.83525 148.131 6.50713 147.645 6.27861C147.158 6.05009 146.564 5.93584 145.861 5.93584H142.784V17.8011H141.703V5.0042ZM149.605 17.8011L146.291 12.0794L147.451 12.0706L150.756 17.6868V17.8011H149.605ZM162.621 10.7171V12.0882C162.621 12.9847 162.504 13.7962 162.27 14.5228C162.041 15.2435 161.707 15.8616 161.268 16.3772C160.834 16.8929 160.313 17.2884 159.703 17.5638C159.094 17.8392 158.408 17.9769 157.647 17.9769C156.903 17.9769 156.223 17.8392 155.608 17.5638C154.998 17.2884 154.474 16.8929 154.034 16.3772C153.595 15.8616 153.255 15.2435 153.015 14.5228C152.775 13.7962 152.655 12.9847 152.655 12.0882V10.7171C152.655 9.8206 152.772 9.01201 153.006 8.29131C153.246 7.56474 153.586 6.94365 154.026 6.42802C154.465 5.9124 154.989 5.51689 155.599 5.2415C156.208 4.96611 156.885 4.82841 157.629 4.82841C158.391 4.82841 159.076 4.96611 159.686 5.2415C160.295 5.51689 160.82 5.9124 161.259 6.42802C161.698 6.94365 162.035 7.56474 162.27 8.29131C162.504 9.01201 162.621 9.8206 162.621 10.7171ZM161.549 12.0882V10.6995C161.549 9.94365 161.461 9.26396 161.285 8.66045C161.115 8.05693 160.861 7.54131 160.521 7.11357C160.187 6.68584 159.777 6.35771 159.29 6.1292C158.804 5.90068 158.25 5.78642 157.629 5.78642C157.026 5.78642 156.484 5.90068 156.003 6.1292C155.523 6.35771 155.113 6.68584 154.773 7.11357C154.439 7.54131 154.181 8.05693 153.999 8.66045C153.823 9.26396 153.736 9.94365 153.736 10.6995V12.0882C153.736 12.8499 153.823 13.5354 153.999 14.1448C154.181 14.7483 154.442 15.2669 154.781 15.7005C155.121 16.1282 155.531 16.4563 156.012 16.6849C156.498 16.9134 157.043 17.0276 157.647 17.0276C158.274 17.0276 158.827 16.9134 159.308 16.6849C159.788 16.4563 160.196 16.1282 160.53 15.7005C160.863 15.2669 161.115 14.7483 161.285 14.1448C161.461 13.5354 161.549 12.8499 161.549 12.0882Z" fill="white"/>
</svg>`

		// 小Logo SVG (移动端)
		smallLogoSVG := `<svg width="30" height="23" viewBox="0 0 30 23" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M0.000210938 2.72516C0.000210938 5.07516 0.800211 6.22516 4.00021 8.32516C7.55021 10.6752 8.00021 11.4252 8.00021 14.9252C8.00021 19.8252 9.50021 22.7252 12.0502 22.7252C13.8002 22.7252 14.0002 22.0252 14.0002 14.9752L13.9502 7.22516L7.55021 3.47516C-0.149789 -1.02484 0.000210938 -1.02484 0.000210938 2.72516Z" fill="white"/>
<path d="M22.2502 3.52516L16.0502 7.22516L16.0002 14.9752C16.0002 22.0252 16.2002 22.7252 17.9502 22.7252C20.5002 22.7252 22.0002 19.8252 22.0002 14.9252C22.0002 11.4252 22.4502 10.6752 26.0002 8.32516C29.2002 6.22516 30.0002 5.07516 30.0002 2.72516C30.0002 -1.02484 29.8502 -0.974842 22.2502 3.52516Z" fill="white"/>
</svg>`

		// now := time.Now()
		// nextEvenHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1+(now.Hour()%2), 0, 0, 0, now.Location())
		// diff := int(nextEvenHour.Sub(now).Seconds())
		// hours := diff / 3600
		// minutes := (diff % 3600) / 60
		// seconds := diff % 60
		// countdown := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

		return h.Div(
			// 响应式Logo设计
			h.Div(
				// 桌面端Logo (完整版)
				h.Div(
					h.A(
						h.RawHTML(fullLogoSVG),
					).Href("/").Class("logo"),
				).Class("pt-2 hidden sm:flex"),

				// 移动端Logo (简化版)
				h.Div(
					h.A(
						h.RawHTML(smallLogoSVG),
					).Href("/"),
				).Class("pt-2 pr-2 flex sm:hidden"),
			).Class(""),

			// 环境指示器
			// h.If(dbReset != "",
			// 	h.Div(
			// 		v.VAlert().
			// 			Type("warning").
			// 			Variant("tonal").
			// 			Density("compact").
			// 			Class("ma-2 text-center").
			// 			Children(
			// 				h.Div(
			// 					v.VIcon("mdi-clock-outline").Size("small").Class("mr-1"),
			// 					h.Span(msgr.DBResetTipLabel).Class("text-caption mr-1"),
			// 					h.Span(countdown).Id("countdown").Class("font-weight-bold"),
			// 				).Class("d-flex align-center justify-center"),
			// 			),
			// 	),
			// 	h.Script("function updateCountdown(){const now=new Date();const nextEvenHour=new Date(now);nextEvenHour.setHours(nextEvenHour.getHours()+(nextEvenHour.getHours()%2===0?2:1),0,0,0);const timeLeft=nextEvenHour-now;const hours=Math.floor(timeLeft/(60*60*1000));const minutes=Math.floor((timeLeft%(60*60*1000))/(60*1000));const seconds=Math.floor((timeLeft%(60*1000))/1000);const countdownElem=document.getElementById(\"countdown\");if(countdownElem){countdownElem.innerText=`${hours.toString().padStart(2,\"0\")}:${minutes.toString().padStart(2,\"0\")}:${seconds.toString().padStart(2,\"0\")}`}}updateCountdown();setInterval(updateCountdown,1000);"),
			// ),
		)
	})
}

func configPost(
	b *presets.Builder,
	db *gorm.DB,
	publisher *publish.Builder,
	ab *activity.Builder,
	seoBuilder *seo.Builder,
) *presets.ModelBuilder {
	m := b.Model(&models.Post{})
	defer m.Use(publisher, ab, seoBuilder)

	mListing := m.Listing("ID", "Title", "TitleWithSlug", "HeroImage", "Body", activity.ListFieldNotes).
		SearchColumns("title", "body").
		PerPage(10)

	mListing.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		item, err := ab.MustGetModelBuilder(m).NewHasUnreadNotesFilterItem(ctx.R.Context(), "")
		if err != nil {
			panic(err)
		}
		return []*shadcn.FilterItem{
			item,
			{
				Key:          "created",
				Label:        "Create Time",
				ItemType:     shadcn.FilterItemTypeDatetimeRangePicker,
				SQLCondition: `created_at %s ?`,
			},
			{
				Key:          "title",
				Label:        "Title",
				ItemType:     shadcn.FilterItemTypeString,
				SQLCondition: `title %s ?`,
			},
			{
				Key:      "status",
				Label:    "Status",
				ItemType: shadcn.FilterItemTypeSelect,
				Options: []shadcn.FilterSelectOption{
					{Text: publish.StatusDraft, Value: publish.StatusDraft},
					{Text: publish.StatusOnline, Value: publish.StatusOnline},
					{Text: publish.StatusOffline, Value: publish.StatusOffline},
				},
				SQLCondition: `status %s ?`,
			},
			{
				Key:      "multi_statuses",
				Label:    "Multiple Statuses",
				ItemType: shadcn.FilterItemTypeMultipleSelect,
				Options: []shadcn.FilterSelectOption{
					{Text: publish.StatusDraft, Value: publish.StatusDraft},
					{Text: publish.StatusOnline, Value: publish.StatusOnline},
					{Text: publish.StatusOffline, Value: publish.StatusOffline},
				},
				SQLCondition: `status %s ?`,
				Folded:       true,
			},
			{
				Key:          "id",
				Label:        "ID",
				ItemType:     shadcn.FilterItemTypeNumber,
				SQLCondition: `id %s ?`,
				Folded:       true,
			},
		}
	})

	mListing.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)

		tab, err := ab.MustGetModelBuilder(m).NewHasUnreadNotesFilterTab(ctx.R.Context())
		if err != nil {
			panic(err)
		}
		return []*presets.FilterTab{
			{
				Label: msgr.FilterTabsAll,
				ID:    "all",
				Query: url.Values{"all": []string{"1"}},
			},
			tab,
		}
	})

	lazyWrapperEditCompoSync := autosync.NewLazyWrapComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) *autosync.Config {
		return &autosync.Config{
			SyncFromFromKey: strings.TrimSuffix(field.FormKey, "WithSlug"),
			InitialChecked:  autosync.InitialCheckedAuto,
			CheckboxLabel:   "Auto Sync",
			SyncEndpoint:    autosync.SlugEndpointPath, // 服务端 slug（gosimple/slug），替代客户端 unidecode
		}
	})
	m.Editing().Field("TitleWithSlug").LazyWrapComponentFunc(lazyWrapperEditCompoSync)
	m.Editing().ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		p := obj.(*models.Post)
		if p.Title == "" {
			err.FieldError("Title", "Title Is Required")
		}
		if p.TitleWithSlug == "" {
			err.FieldError("TitleWithSlug", "TitleWithSlug Is Required")
		}
		return
	})
	dp := m.Detailing(publish.VersionsPublishBar, "Detail", seo.SeoDetailFieldName).Drawer(true)
	detailSection := presets.NewSectionBuilder(m, "Detail").
		Editing("Title", "TitleWithSlug", "HeroImage", "Body", "BodyImage")
	detailSection.EditingField("TitleWithSlug").LazyWrapComponentFunc(lazyWrapperEditCompoSync)
	// TODO: need viewing field setting
	detailSection.EditingField("HeroImage").
		WithContextValue(
			media.MediaBoxConfig,
			&media_library.MediaBoxConfig{
				AllowType: "image",
				Sizes: map[string]*base.Size{
					"thumb": {
						Width:  400,
						Height: 300,
					},
					"main": {
						Width:  800,
						Height: 500,
					},
				},
			})
	detailSection.EditingField("BodyImage").
		WithContextValue(
			media.MediaBoxConfig,
			&media_library.MediaBoxConfig{})
	detailSection.EditingField("Body").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return shadcn.Textarea().
			Label(field.Label).
			ErrorMessages(field.Errors...).
			Attr(web.VField(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)))...).
			Rows(10).
			Disabled(field.Disabled)
	})
	dp.Section(detailSection)
	return m
}
