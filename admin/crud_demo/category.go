package crud_demo

import (
	"github.com/r0vx/admin/publish"

	"example/models"

	"github.com/r0vx/admin/presets"
	"gorm.io/gorm"
)

// ConfigCategory 配置分类管理模块
func ConfigCategory(b *presets.Builder, db *gorm.DB, publisher *publish.Builder) *presets.ModelBuilder {
	mb := b.Model(&models.Category{}).Use(publisher)
	mb.Listing("ID", "Name", "Position", "Status").SearchColumns("name")
	mb.Editing("Name", "Products", "Position")
	return mb
}

// ConfigNestedFieldDemo 配置嵌套字段演示（Customer → Address → Phone / MembershipCard）
func ConfigNestedFieldDemo(b *presets.Builder, db *gorm.DB) {
	mb := b.Model(&models.Customer{}).URIName("customers")
	mb.Listing("ID", "Name").SearchColumns("name")
	mb.Editing(
		&presets.FieldsSection{
			Title: "Basic Info",
			Rows:  [][]string{{"Name", "Titile"}, {"Addresses"}},
		},
		&presets.FieldsSection{
			Title: "Membership",
			Rows:  [][]string{{"MembershipCard"}},
		},
	)
	mb.Detailing("ID", "Name", "Addresses", "MembershipCard")
}

// ConfigProject 配置项目管理模块
func ConfigProject(pb *presets.Builder, db *gorm.DB) {
	mb := pb.Model(&models.Project{}).URIName("projects")
	mb.Listing("ID", "Name", "Status", "Featured").SearchColumns("name", "description")

	ConfigTask(pb, db)

	mb.Editing(
		&presets.FieldsSection{
			Title: "Basic Info",
			Rows:  [][]string{{"Name", "Status"}, {"Description"}, {"Featured", "Avatar"}},
		},
		&presets.FieldsSection{
			Title: "Tasks",
			Rows:  [][]string{{"Tasks"}},
		},
	)
}

// ConfigTask 配置任务子表模块
func ConfigTask(pb *presets.Builder, _ *gorm.DB) {
	mb := pb.Model(&models.Task{}).URIName("tasks")
	mb.Listing("ID", "Name", "Priority", "Status", "Assignee")
	mb.Editing("Name", "Description", "Priority", "Status", "Assignee")
}