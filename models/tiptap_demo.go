package models

import "time"

// TiptapDemo tiptap 富文本编辑器演示模型
//
// 演示 x/ui/tiptap 包的功能：
//   - 基于 tiptap (ProseMirror) 的富文本编辑
//   - 工具栏操作（标题、列表、表格、图片等）
//   - 表单集成（hidden input 提交）
type TiptapDemo struct {
	ID        uint   `gorm:"primarykey"`
	Title     string // 标题
	Body      string // 富文本正文（HTML）
	Summary   string // 摘要（纯文本 Textarea 对比）
	UpdatedAt time.Time
	CreatedAt time.Time
}
