package models

import "time"

// VehicleFilterDemo 车型级联筛选演示模型
//
// 展示 LinkageSelectItem 在汽车品牌/厂商/车系三级联动场景下的筛选效果
type VehicleFilterDemo struct {
	ID        uint      `gorm:"primarykey"`
	Title     string    // 配件名称
	Brand     string    // 品牌（级联第一级）
	Maker     string    // 厂商（级联第二级）
	Series    string    // 车系（级联第三级）
	Price     float64   // 价格
	UpdatedAt time.Time // 更新时间
	CreatedAt time.Time // 创建时间
}
