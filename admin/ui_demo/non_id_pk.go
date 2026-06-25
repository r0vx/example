package ui_demo

import (
	"fmt"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"gorm.io/gorm"
)

// ConfigNonIDPKDemo 非 ID 主键 CRUD 演示（URI: /sc-member-demo）。
//
// SCMember 主键是 ServiceUserID（无 ID 字段）。本 demo 验证：
//   - 列表能渲染行 id（选择框 v-model / 编辑链接不为空，不再触发 Vue #38）
//   - 行点击进编辑抽屉、详情、删除、批量选择全走真实主键
//   - 列表分页排序用 ServiceUserID（schema 自动解析，无需 PrimaryField）
func ConfigNonIDPKDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&models.SCMember{}); err != nil {
		panic(err)
	}
	seedSCMember(db)

	// 注意：不写 PrimaryField("ServiceUserID")，验证 schema 自动解析。需手动指定也可：
	//   mb := b.Model(&models.SCMember{}).URIName("sc-member-demo").PrimaryField("ServiceUserID")
	mb := b.Model(&models.SCMember{}).URIName("sc-member-demo").Label("SCMember (非 ID 主键)")

	lb := mb.Listing("ServiceUserID", "Name", "Phone", "Balance", "CreatedAt")
	lb.SearchColumns("name", "phone")
	lb.PerPage(10)

	// 批量删除：渲染行选择框（v-model 绑行 id），专门验证 ObjectID 非空
	// RequiresConfirmation 提供确认 Dialog（compFunc），否则点击报 "bulk.compFunc not set"
	lb.BulkAction("BulkDelete").Label("批量删除").
		RequiresConfirmation().
		ConfirmPrompt("确认删除选中的会员？").
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) error {
			if err := db.Delete(&models.SCMember{}, "service_user_id IN ?", selectedIds).Error; err != nil {
				return err
			}
			r.Reload = true
			return nil
		})

	mb.Editing("Name", "Phone", "Balance")
	mb.Detailing("ServiceUserID", "Name", "Phone", "Balance", "CreatedAt").Drawer(true)
}

// seedSCMember 幂等种子：空表才插入
func seedSCMember(db *gorm.DB) {
	var n int64
	db.Model(&models.SCMember{}).Count(&n)
	if n > 0 {
		return
	}
	for i := 1; i <= 23; i++ {
		db.Create(&models.SCMember{
			Name:    fmt.Sprintf("商户 %02d", i),
			Phone:   fmt.Sprintf("138%08d", i),
			Balance: float64(i) * 100.5,
		})
	}
}
