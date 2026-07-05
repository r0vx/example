package ui_demo

// ============================================================================
// RowLevelRefresh 行级刷新演示 + 初始渲染开销 benchmark
// ============================================================================
//
// lb.RowLevelRefresh(true) 开启后：收到 NotifModelsUpdated 时只重渲染对应行的
// 单元格并就地 DOM 补丁（cell 包 go-plaid-portal + UpdatePortals 重编译），不整表 reload。
//
// benchmark：同数据同列，/row-refresh-demo(ON) vs /row-refresh-off(OFF) 对比初始渲染开销。
// ============================================================================

import (
	"fmt"
	"time"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// RowRefreshDemo 行级刷新演示模型
type RowRefreshDemo struct {
	gorm.Model
	Name   string
	Status string
}

// rowRefreshBadgeFn 返回渲染 badge 的 cell ComponentFunc（Vue 组件，凑 cell 数测开销）
func rowRefreshBadgeFn(label string) func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	return func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		r := obj.(*RowRefreshDemo)
		v := shadcn.BadgeVariantDefault
		if r.ID%2 == 0 {
			v = shadcn.BadgeVariantSecondary
		}
		return shadcn.Badge(h.Text(fmt.Sprintf("%s-%d", label, r.ID))).Variant(v)
	}
}

// configRowRefreshListing 配置一个 benchmark listing：9 列（6 个 badge Vue 组件列）、perPage 100。
// rowLevel=true 时每 cell 包 portal（被测开销）；false 为对照。
func configRowRefreshListing(mb *presets.ModelBuilder, rowLevel bool) {
	lb := mb.Listing("ID", "Name", "Status", "B1", "B2", "B3", "B4", "B5", "RenderedAt")
	lb.PerPage(100)
	if rowLevel {
		lb.RowLevelRefresh(true)
	}
	lb.Field("Status").ComponentFunc(rowRefreshBadgeFn("Status"))
	for _, c := range []string{"B1", "B2", "B3", "B4", "B5"} {
		lb.Field(c).ComponentFunc(rowRefreshBadgeFn(c))
	}
	lb.Field("RenderedAt").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(h.Text(time.Now().Format("15:04:05.000"))).Class("tabular-nums font-mono text-xs")
	})
	mb.Editing("Name", "Status")
}

// ConfigRowRefreshDemo 注册 ON/OFF 两个 benchmark 列表（同 model 同数据，唯一差 RowLevelRefresh）
func ConfigRowRefreshDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&RowRefreshDemo{}); err != nil {
		panic(err)
	}

	// 种子补到 100 行（benchmark 需足够 cell 看出差异）
	var cnt int64
	db.Model(&RowRefreshDemo{}).Count(&cnt)
	for i := cnt + 1; i <= 100; i++ {
		db.Create(&RowRefreshDemo{Name: fmt.Sprintf("Item %d", i), Status: "active"})
	}

	// ON：RowLevelRefresh(true)，每 cell 包 portal
	mbOn := b.Model(&RowRefreshDemo{}).URIName("row-refresh-demo").Label("RowRefresh ON")
	configRowRefreshListing(mbOn, true)

	// OFF：对照，不包 portal
	mbOff := b.Model(&RowRefreshDemo{}).URIName("row-refresh-off").Label("RowRefresh OFF")
	configRowRefreshListing(mbOff, false)
}
