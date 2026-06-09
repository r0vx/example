package models

import "time"

// AutocompleteDemo admin/autocomplete 包功能演示模型
//
// 演示通过远程搜索自动补全选择关联的 Product 和 Category。
// ProductName/CategoryName 为冗余字段，保存时自动写入，方便列表直接展示。
type AutocompleteDemo struct {
	ID           uint   `gorm:"primarykey"`
	Title        string // 标题
	ProductID    uint   // 关联产品 ID（单选，通过 autocomplete 选择）
	ProductName  string // 产品名称（冗余，列表展示用）
	ProductIDs   string // 多选产品 IDs（JSON 数组字符串，如 ["1","3"]）
	CategoryID   uint   // 关联分类 ID（通过 autocomplete 选择）
	CategoryName string // 分类名称（冗余，列表展示用）
	AssigneeID   string // 负责人（静态 items + icon 头像演示）
	UpdatedAt    time.Time
	CreatedAt    time.Time
}
