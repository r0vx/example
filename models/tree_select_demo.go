package models

import "time"

// TreeSelectDemo 树形选择演示模型（汽配商经营车型）
type TreeSelectDemo struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    // 汽配商名称
	UpdatedAt time.Time // 更新时间
	CreatedAt time.Time // 创建时间
}

// TreeSelectDemoSeries 汽配商经营车系关联表（多对多）
type TreeSelectDemoSeries struct {
	ID               uint   `gorm:"primarykey"`
	TreeSelectDemoID uint   `gorm:"index"` // 汽配商 ID
	SeriesID         string // 车系 ID（叶子节点）
}
