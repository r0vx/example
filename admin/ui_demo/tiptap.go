package ui_demo

import (
	"fmt"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/tiptapeditor"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// ============================================================================
// Tiptap 富文本编辑器演示
// ============================================================================
//
// ## 概述
//
// x/ui/tiptap 包提供基于 tiptap (ProseMirror) 的富文本编辑器，核心功能：
//   - 所见即所得的 HTML 编辑（标题、列表、粗体、链接、表格等）
//   - 工具栏按钮分组（文本格式/标题/列表/对齐/插入/操作）
//   - 图片上传（通过 imageUploadUrl 配置）
//   - 任务列表、代码块、引用等丰富内容类型
//   - 表单集成（hidden input 自动提交 HTML 内容）
//
// ## API 用法
//
//	tiptap.TiptapEditor().
//		Value(htmlContent).           // 设置 HTML 内容
//		FieldName(field.FormKey).     // 表单字段名
//		Label(field.Label).           // 字段标签
//		Placeholder("请输入内容..."). // 占位符
//		MinHeight("200px").           // 最小高度
//		MaxHeight("500px").           // 最大高度
//		ImageUploadURL("/api/upload").// 图片上传接口
//		Disabled(false).              // 是否禁用
//		ErrorMessages(field.Errors...)// 错误信息
//
// ============================================================================

// configTiptapDemo 配置 tiptap 富文本编辑器演示模块
func ConfigTiptapDemo(b *presets.Builder, db *gorm.DB) {
	// 数据库迁移
	db.AutoMigrate(&models.TiptapDemo{})

	mb := b.Model(&models.TiptapDemo{})

	// 列表配置
	mb.Listing("ID", "Title", "UpdatedAt").
		SearchColumns("title", "body").
		PerPage(20)

	// 编辑配置
	ed := mb.Editing("Title", "Body", "Summary")

	// Title — 普通输入框
	ed.Field("Title").Label("标题").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Input().
				Label(field.Label).
				Placeholder("请输入标题").
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	// Body — tiptap 富文本编辑器（含 MediaBox 图片集成）
	ed.Field("Body").Label("正文（Tiptap 富文本）").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return tiptapeditor.TiptapEditor(db, field.FormKey).
				Value(fmt.Sprint(field.Value(obj))).
				Label(field.Label).
				Placeholder("请输入正文内容...").
				MinHeight("300px").
				ErrorMessages(field.Errors...)
		})

	// Summary — 普通 Textarea（对比）
	ed.Field("Summary").Label("摘要（纯文本）").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Textarea().
				Label(field.Label).
				Placeholder("请输入摘要...").
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	// 验证
	ed.ValidateFunc(func(obj any, ctx *web.EventContext) (err web.ValidationErrors) {
		demo := obj.(*models.TiptapDemo)
		if demo.Title == "" {
			err.FieldError("Title", "标题不能为空")
		}
		return
	})
}
