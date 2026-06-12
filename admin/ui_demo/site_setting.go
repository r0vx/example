package ui_demo

import (
	"time"

	"github.com/r0vx/admin/avatarupload"
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// SiteSetting 站点设置（单例模型：全站只保留一条记录）
type SiteSetting struct {
	ID           uint   `gorm:"primarykey"`
	SiteName     string // 站点名称
	LogoURL      string // Logo 图片 URL
	ContactEmail string // 联系邮箱
	ICPNumber    string // ICP 备案号
	Copyright    string // 版权信息

	Description     string    // 站点描述（多行文本）
	MaintenanceMode bool      // 维护模式开关
	UpdatedAt       time.Time // 最后更新时间
}

// ConfigSiteSettingDemo 注册「站点设置」单例配置页范例
//
// Singleton(true) 的核心作用：该模型的 URL 路由（/site-setting）不再走列表页，
// 而是改走 editing.singletonPageFunc —— 直接渲染「唯一一条记录」的编辑表单。
// 首次访问时若表中无记录，框架会自动创建一条空记录再渲染（见 admin/presets/editing.go:270）。
// 因此单例模型只需配置 Editing 表单，无需配置 Listing 列表。
func ConfigSiteSettingDemo(b *presets.Builder, db *gorm.DB) {
	// 建表
	if err := db.AutoMigrate(&SiteSetting{}); err != nil {
		panic(err)
	}
	// 预置一条带默认值的记录（单例只需一条；不预置时框架也会自动建空记录，这里给出更友好的默认值）
	seedSiteSetting(db)

	// 关键：Singleton(true) 把此 Model 变成「单例配置页」
	mb := b.Model(&SiteSetting{}).URIName("site-setting").Singleton(true)

	// 单例模型不展示列表，只配置 Editing 表单字段；用两个分区分组展示
	ed := mb.Editing(
		&presets.FieldsSection{
			Title: "基本信息",
			Rows: [][]string{
				{"SiteName"},
				{"LogoURL"},
				{"ContactEmail"},
				{"ICPNumber"},
				{"Copyright"},
			},
		},
		&presets.FieldsSection{
			Title: "高级设置",
			Rows: [][]string{
				{"Description"},
				{"MaintenanceMode"},
			},
		},
	)

	// 站点描述 —— 多行文本（Textarea）；表单绑定一律用 web.VField
	ed.Field("Description").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return shadcn.Textarea().
			Label(field.Label).
			Rows(4).
			Placeholder("一句话介绍你的站点").
			Attr(web.VField(field.Name, field.Value(obj))...).
			ErrorMessages(field.Errors...)
	})

	// 维护模式 —— 开关（Switch），bool 字段
	ed.Field("MaintenanceMode").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		checked, _ := field.Value(obj).(bool)
		return shadcn.Switch().
			Label(field.Label).
			Tips("开启后前台显示维护页").
			Checked(checked).
			Disabled(field.Disabled).
			Attr(web.VField(field.Name, field.Value(obj))...)
	})

	// Logo —— 复用头像上传组件（点击→cropperjs 裁剪→上传→回显 URL），方形展示
	// 复用 avatar-upload-demo 已挂载的上传端点（router.go: AvatarUploadPath，同 ui_demo 包）
	avatarupload.Configure(mb, "LogoURL", avatarupload.Config{
		UploadURL: AvatarUploadPath,
		Shape:     "square",
		Size:      96,
	})
}

// seedSiteSetting 首次插入一条默认站点设置（单例只保留一条记录）
func seedSiteSetting(db *gorm.DB) {
	var count int64
	db.Model(&SiteSetting{}).Count(&count)
	if count > 0 {
		return
	}
	db.Create(&SiteSetting{
		SiteName:     "r0vx 示例站点",
		ContactEmail: "admin@example.com",
		Copyright:    "© 2026 r0vx",
		Description:  "基于 r0vx 框架构建的企业级管理后台示例。",
	})
}
