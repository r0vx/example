package models

import (
	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	gorm.Model
	Name        string
	Icon        string // lucide 图标名称
	Description string
	Status      string // active, completed, archived
	Featured    bool   // 是否推荐（演示 Checkbox）
	Avatar      string // 项目头像 URL（演示 Avatar）
	Tasks       []Task `gorm:"foreignKey:ProjectID"`
}

// Task 任务模型（子表）
type Task struct {
	gorm.Model
	ProjectID   uint
	Name        string
	Description string
	Priority    string // high, medium, low
	Status      string // pending, in_progress, completed
	Assignee    string
	Icon        string // lucide 图标名称，如 file, folder, star 等
}
