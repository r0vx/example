package crud_demo

import (
	"example/models"

	"github.com/r0vx/admin/presets"
	"gorm.io/gorm"
)

// ConfigL10nModel 配置多语言模型模块
func ConfigL10nModel(db *gorm.DB, b *presets.Builder) (*presets.ModelBuilder, *presets.ModelBuilder) {
	l10nM := b.Model(&models.L10nModel{}).URIName("l10n-models")
	l10nM.Listing("ID", "Title", "LocaleCode").SearchColumns("title")
	l10nM.Editing("Title")

	l10nVM := b.Model(&models.L10nModelWithVersion{}).URIName("l10n-model-with-versions")
	l10nVM.Listing("ID", "Title", "LocaleCode", "Version").SearchColumns("title")
	l10nVM.Editing("Title")

	return l10nM, l10nVM
}