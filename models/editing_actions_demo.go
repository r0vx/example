package models

import "time"

// EditingActionsDemo 编辑页自定义按钮演示模型
//
// 演示 Editing.ActionsFunc / TopActionsFunc / AppendTopActionsFunc 的用法：
//   - ActionsFunc: 替换编辑页底部的保存/取消按钮
//   - TopActionsFunc: 表单 Card 上方的操作按钮区
//   - AppendTopActionsFunc: 追加顶部按钮
type EditingActionsDemo struct {
	ID        uint   `gorm:"primarykey"`
	Title     string // 标题
	Status    string // 状态（draft / pending / published）
	Content   string // 内容
	UpdatedAt time.Time
	CreatedAt time.Time
}
