package models

import "gorm.io/gorm"

// WizardDemo Wizard 多步向导演示模型
type WizardDemo struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255)"`
	Industry string `gorm:"type:varchar(100)"`
	Phone    string `gorm:"type:varchar(50)"`
	Address  string `gorm:"type:text"`
	Status   string `gorm:"type:varchar(50);default:'draft'"`
}
