package ui_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// CTCategory 跨表树演示父模型（分类）
type CTCategory struct {
	ID     uint `gorm:"primarykey"`
	Name   string
	Code   string
	Title  string
	Author string
	Status string
}

// CTArticle 跨表树演示子模型（文档；CategoryID 指向 CTCategory.ID）
type CTArticle struct {
	ID         uint `gorm:"primarykey"`
	CategoryID uint
	Title      string
	Author     string
	Status     string
}

// ConfigCrossTreeListingDemo 注册跨表树演示（分类→文档；展开懒加载+父下新建子项+子行编辑/删除）
func ConfigCrossTreeListingDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&CTCategory{}, &CTArticle{}); err != nil {
		panic(err)
	}
	seedCrossTree(db)

	// 子表 B：独立注册的 ModelBuilder（有自己的 listing/editing）
	artMB := b.Model(&CTArticle{}).URIName("cross-tree-articles")
	artLB := artMB.Listing("Title", "Author", "Status")
	artLB.RowMenu().InlineDefaultsInMenu(true) // 自定义项平铺为按钮 + 默认 Edit/Delete 收进末尾 ⋮（演示第三态）
	artMB.Editing("Title", "Author", "Status", "CategoryID")
	// 富 RowMenuItem——验证跨表树子表格菜单与父行功能对齐：Icon/Tooltip/Variant/Disabled/Confirm 全支持
	artLB.RowMenu().RowMenuItem("preview").
		Icon("eye").
		Tooltip("预览文档").
		Variant(shadcn.ButtonVariantSecondary).
		Disabled(func(obj any, id string, ctx *web.EventContext) bool {
			return obj.(*CTArticle).Status == "草稿" // 草稿禁预览（演示 Disabled 逐行生效）
		}).
		RequiresConfirmation().
		OnClick(func(ctx *web.EventContext, id string) (web.EventResponse, error) {
			var r web.EventResponse
			web.AppendRunScripts(&r, "vars.__window.alert('预览文档 #' + "+id+")")
			return r, nil
		})

	// 父表 A：开启跨表树，引用 artMB
	catMB := b.Model(&CTCategory{}).URIName("cross-tree-listing-demo")
	catLB := catMB.Listing("Name", "Code").SelectableColumns(false).
		ResizableColumns(true).
		ReorderableColumns(true)
	catLB.CrossTreeMode(presets.CrossTree(artMB, "CategoryID"))
	catMB.Editing("Name", "Code")
}

// seedCrossTree 首次启动插入演示数据（2 分类含子 + 1 空分类）
func seedCrossTree(db *gorm.DB) {
	var n int64
	db.Model(&CTCategory{}).Count(&n)
	if n > 0 {
		return
	}
	tech := CTCategory{Name: "技术", Code: "TECH"}
	sales := CTCategory{Name: "销售", Code: "SALES"}
	empty := CTCategory{Name: "空分类", Code: "EMPTY"}
	db.Create(&tech)
	db.Create(&sales)
	db.Create(&empty)
	db.Create(&[]CTArticle{
		{CategoryID: tech.ID, Title: "Go 入门指南", Author: "张三", Status: "已发布"},
		{CategoryID: tech.ID, Title: "Vue 3 实战", Author: "李四", Status: "草稿"},
		{CategoryID: sales.ID, Title: "大客户成单技巧", Author: "王五", Status: "已发布"},
	})
}
