package models

import "time"

// DialogDemo 对话框类型演示模型
//
// 用于展示 admin 框架 4 种 Dialog 类型的区别：
//   - Global Dialog（全局对话框）：编辑表单默认弹出方式
//   - ListingCompo Dialog（列表选择对话框）：在对话框中展示列表供选择
//   - BulkAction Dialog（批量操作对话框）：勾选多行后触发
//   - Action Dialog（操作对话框）：列表/详情页的自定义操作按钮
type DialogDemo struct {
	ID        uint   `gorm:"primarykey"`
	Title     string // 标题
	Status    string // 状态（active / inactive / pending）
	Priority  int    // 优先级（1-5）
	Notes     string // 备注
	RelatedID uint   // 关联记录 ID（通过 ListingCompo Dialog 选择）
	UpdatedAt time.Time
	CreatedAt time.Time
}
