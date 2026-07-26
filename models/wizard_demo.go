package models

import (
	"time"

	"gorm.io/gorm"
)

// WizardDemo Wizard 多步向导演示模型
type WizardDemo struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255)"`
	Industry string `gorm:"type:varchar(100)"`
	Phone    string `gorm:"type:varchar(50)"`
	Address  string `gorm:"type:text"`
	Status   string `gorm:"type:varchar(50);default:'draft'"`
}

// NodataDemo 无数据列表显示模型：数据库无对应表，全字段 gorm:"-"，
// 数据由 Listing.SearchFunc 从内存/外部 API 提供（见 ui_demo/nodata_listing.go）。
// ID 主键须保留——presets.ObjectID 靠反射取它做行 key。
type NodataDemo struct {
	ID        uint      `gorm:"-"`
	Name      string    `gorm:"-"`
	Industry  string    `gorm:"-"`
	Phone     string    `gorm:"-"`
	Address   string    `gorm:"-"`
	Status    string    `gorm:"-"`
	UpdatedAt time.Time `gorm:"-"`
}
