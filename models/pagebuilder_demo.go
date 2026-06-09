package models

import "gorm.io/gorm"

// ============================================================================
// PageBuilder 演示容器模型
// ============================================================================

// PBDemoHero 首屏大图容器
type PBDemoHero struct {
	gorm.Model
	Title    string
	Subtitle string
	BgColor  string
	BgImage  string
}

// PBDemoBanner 横幅容器
type PBDemoBanner struct {
	gorm.Model
	Text      string
	LinkText  string
	LinkURL   string
	Dismissed bool
}

// PBDemoRichText 富文本容器
type PBDemoRichText struct {
	gorm.Model
	Content string `gorm:"type:text"`
}
