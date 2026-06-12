package ui_demo

import (
	"fmt"
	"time"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// Department 树形列表演示模型（部门管理：子.ParentID = 父.ID，0 = 根）
type Department struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	Code      string
	ParentID  uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ConfigTreeListingDemo 注册树形列表演示（懒加载展开/折叠 + 搜索退化平铺）
func ConfigTreeListingDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&Department{}); err != nil {
		panic(err)
	}
	seedDepartments(db)

	mb := b.Model(&Department{}).URIName("tree-listing-demo")
	lb := mb.Listing("ID", "Name", "Code")
	// ExpandInActions(true)：展开按钮放行尾 action 列（默认在第一列跟随缩进，删掉该调用即恢复）
	lb.TreeMode(presets.Tree("ParentID").ExpandInActions(true))
	lb.SearchColumns("name", "code")

	rmb := lb.RowMenu().InlineDefaultsInMenu(true)

	// test
	rmb.RowMenuItem("readmeDoc").ComponentFunc(func(obj any, id string, ctx *web.EventContext) h.HTMLComponent {
		cu := obj.(*Department)

		return shadcn.RowMenuItem("test").SetOnclick(
			fmt.Sprintf("vars.__window.open('/readme?id=%d', '_blank')", cu.ID),
			//web.GET().URL("/readme").PushState(true).Query("id", fmt.Sprint(cu.ID)).Go(),
		)

	})

	mb.Editing("Name", "Code", "ParentID")
}

// seedDepartments 首次启动插入演示树（3 根 × 三级）
func seedDepartments(db *gorm.DB) {
	var count int64
	db.Model(&Department{}).Count(&count)
	if count > 0 {
		return
	}
	depts := []Department{
		{ID: 1, Name: "技术中心", Code: "TECH", ParentID: 0},
		{ID: 2, Name: "销售中心", Code: "SALES", ParentID: 0},
		{ID: 3, Name: "人力资源", Code: "HR", ParentID: 0},
		{ID: 4, Name: "后端组", Code: "TECH-BE", ParentID: 1},
		{ID: 5, Name: "前端组", Code: "TECH-FE", ParentID: 1},
		{ID: 6, Name: "平台组", Code: "TECH-BE-PLAT", ParentID: 4},
		{ID: 7, Name: "基础设施组", Code: "TECH-BE-INFRA", ParentID: 4},
		{ID: 8, Name: "国内销售", Code: "SALES-CN", ParentID: 2},
		{ID: 9, Name: "海外销售", Code: "SALES-INTL", ParentID: 2},
		{ID: 10, Name: "华东区", Code: "SALES-CN-EAST", ParentID: 8},
	}
	db.Create(&depts)
}
