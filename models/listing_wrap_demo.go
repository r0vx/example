package models

import "time"

// ListingWrapDemo 列表 Wrap* 方法演示模型
//
// 演示 ListingBuilder 的 WrapColumns、WrapCell、WrapRow、
// WrapFilterDataFunc、WrapNewButtonFunc、WrapSearchFunc 用法。
type ListingWrapDemo struct {
	ID        uint      `gorm:"primarykey"`
	Title     string    // 标题
	Status    string    // 状态：draft / active / archived
	Priority  int       // 优先级：1=低 2=中 3=高
	Assignee  string    // 负责人
	Category  string    // 分类
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}
