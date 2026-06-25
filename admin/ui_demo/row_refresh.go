package ui_demo

// ============================================================================
// RowLevelRefresh 行级刷新演示
// ============================================================================
//
// lb.RowLevelRefresh(true) 开启后：收到 NotifModelsUpdated 时只重渲染对应行的
// 单元格并就地 DOM 补丁，不整表 reload。新增/删除/换页/被筛掉自动回退整表。
//
// 验证手法：RenderedAt 列显示「服务端渲染时刻」。
//   - 行级刷新生效：编辑某行保存后，只有那一行的 RenderedAt 变，其余行不动。
//   - 若退化成整表 reload：所有行的 RenderedAt 会一起变。
// ============================================================================

import (
	"fmt"
	"time"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"gorm.io/gorm"
)

// RowRefreshDemo 行级刷新演示模型
type RowRefreshDemo struct {
	gorm.Model
	Name   string
	Status string
}

// ConfigRowRefreshDemo 配置行级刷新演示模块（/row-refresh-demo）
func ConfigRowRefreshDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&RowRefreshDemo{}); err != nil {
		panic(err)
	}

	// 种子数据（仅首次）
	var cnt int64
	db.Model(&RowRefreshDemo{}).Count(&cnt)
	if cnt == 0 {
		for i := 1; i <= 8; i++ {
			db.Create(&RowRefreshDemo{Name: fmt.Sprintf("Item %d", i), Status: "active"})
		}
	}

	mb := b.Model(&RowRefreshDemo{}).URIName("row-refresh-demo").Label("Row Refresh Demo")

	lb := mb.Listing("ID", "Name", "Status", "RenderedAt")
	lb.RowLevelRefresh(true) // ← 被测开关

	// RenderedAt：服务端渲染时刻。整表 reload → 全行一起变；行级刷新 → 仅编辑行变。
	lb.Field("RenderedAt").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(h.Text(time.Now().Format("15:04:05.000"))).Class("tabular-nums font-mono text-xs")
	})

	mb.Editing("Name", "Status")
}
