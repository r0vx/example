package ui_demo

// ============================================================================
// RelayPagination 游标分页演示
// ============================================================================
//
// cl.RelayPagination(gorm2op.KeysetBasedPagination(true)) 开启 keyset 游标分页：
// 关 COUNT(*)、无总条数，深翻页毫秒级，只有「上一页 / 下一页」。
//
// 默认 DataTable 渲染路径下，框架现已自动外挂游标翻页器（修复前只有配
// dataTableFunc 才显示）。
// ============================================================================

import (
	"fmt"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/presets/gorm2op"
	"gorm.io/gorm"
)

// RelayPaginationDemo 游标分页演示模型
type RelayPaginationDemo struct {
	gorm.Model
	Name string
	Code string
}

// ConfigRelayPaginationDemo 配置游标分页演示（/relay-pagination-demo）
func ConfigRelayPaginationDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&RelayPaginationDemo{}); err != nil {
		panic(err)
	}

	var cnt int64
	db.Model(&RelayPaginationDemo{}).Count(&cnt)
	if cnt == 0 {
		for i := 1; i <= 25; i++ {
			db.Create(&RelayPaginationDemo{Name: fmt.Sprintf("Row %02d", i), Code: fmt.Sprintf("R%03d", i)})
		}
	}

	mb := b.Model(&RelayPaginationDemo{}).URIName("relay-pagination-demo").Label("Relay Pagination Demo")

	cl := mb.Listing("ID", "Name", "Code")
	cl.PerPage(10)
	// 千万级：keyset 游标分页 + 关 COUNT(*)，深翻页毫秒级（只上/下页，无总条数）
	cl.RelayPagination(gorm2op.KeysetBasedPagination(true))

	mb.Editing("Name", "Code")
}
