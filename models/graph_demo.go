package models

import "time"

// GraphDemo Graph 图表演示模型
// 用于菜单展示，不实际使用数据库
type GraphDemo struct {
	ID        uint `gorm:"primarykey"`
	UpdatedAt time.Time
	CreatedAt time.Time
}
