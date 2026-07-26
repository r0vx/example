package ui_demo

// ============================================================================
// 纯展示列表（无数据库表）演示
// ============================================================================
//
// 场景：数据库里【不存在】该表和数据，列表内容全由 SearchFunc 从内存/外部 API
// 提供，只读展示、不支持增删改。
//
// 三个关键点（缺一即坏）：
//  1. 模型全字段 gorm:"-" + 保留 ID 主键字段（行 key / ObjectID 靠反射取 ID）。
//  2. lb.SearchFunc 完全接管数据源——否则默认 Searcher 会 db.Model(&NodataDemo{})
//     拼表名查询 → PostgreSQL relation does not exist 报错。
//  3. 有数据时必须设 PageInfo.StartCursor 非 nil，否则框架 buildDataTableAdditions
//     误判为空 → 表体渲染了行、页脚却显示「无记录」（打架）。空数据则 total=0 自动显示空态。
//
// 不配 Editing / Detailing / RowMenu → 天然只读（点行无编辑、无 ⋮ 菜单）。
// 新建按钮用 NewButtonFunc 返回 nil 隐藏。分页关闭（DisablePagination）。
// ============================================================================

import (
	"fmt"
	"time"

	"example/models"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/theplant/relay"
)

// nodataRows 内存假数据（真实项目里换成外部 API / RPC / 缓存查询）
func nodataRows() []*models.NodataDemo {
	base := time.Date(2026, 7, 11, 10, 0, 0, 0, time.Local)
	return []*models.NodataDemo{
		{ID: 1, Name: "杭州云图科技", Industry: "IT", Phone: "13800000001", Address: "杭州市西湖区", Status: "active", UpdatedAt: base},
		{ID: 2, Name: "北京智造集团", Industry: "制造", Phone: "13800000002", Address: "北京市朝阳区", Status: "active", UpdatedAt: base.Add(-2 * time.Hour)},
		{ID: 3, Name: "深圳海联贸易", Industry: "外贸", Phone: "13800000003", Address: "深圳市南山区", Status: "inactive", UpdatedAt: base.Add(-26 * time.Hour)},
	}
}

// nodataByID 按 id 查内存数据（RowMenuItem.OnClick 的 handler 用；真实项目换成外部 API）
func nodataByID(id string) (*models.NodataDemo, bool) {
	for _, r := range nodataRows() {
		if fmt.Sprint(r.ID) == id {
			return r, true
		}
	}
	return nil, false
}

// ConfigNodataDemo 注册纯展示列表（无 DB 表）
func ConfigNodataDemo(b *presets.Builder) {
	// 注意：不调用 db.AutoMigrate——本模型无表，也不该建表。
	mb := b.Model(&models.NodataDemo{}).URIName("nodata-demo").Label("纯展示列表").MenuIcon("table")

	lb := mb.Listing("ID", "Name", "Phone", "Industry", "Status", "UpdatedAt").
		DisableRowClick(true).
		DisablePagination(true) // 数据量固定、无需分页

	// 时间列短格式（默认 2006-01-02 15:04:05）
	lb.TimeFormat("2006-01-02 15:04")

	// 隐藏「新建」按钮（返回 nil = 不渲染）——纯展示不新建
	lb.NewButtonFunc(func(ctx *web.EventContext) h.HTMLComponent { return nil })

	// 行菜单：⚠️ 默认 Edit/Delete 会炸——它们走 editing.Fetcher/Deleter → gorm 查幻表 → 报错。
	// 故先 Empty() 清掉默认项，再挂【不碰 DB】的自定义项。
	rm := lb.RowMenu()
	rm.Empty()

	// 自定义项 1：OnClick 跑你自己的 handler（handler 里读内存/外部 API，别查幻表）。
	rm.RowMenuItem("查看").Icon("eye").OnClick(func(ctx *web.EventContext, id string) (r web.EventResponse, err error) {
		row, ok := nodataByID(id)
		if !ok {
			presets.ShowError(&r, "记录不存在") // 英文/中文提示随项目；此处 demo 用中文
			return r, nil
		}
		presets.ShowMessage(&r, "查看："+row.Name+"（"+row.Phone+"）", "success")
		return r, nil
	})

	// 自定义项 2：URL 纯前端跳转（{id} 占位，SPA pushState / 外链新标签）——完全不触发后端。
	rm.RowMenuItem("外部资料").Icon("external-link").URL("https://example.com/company/{id}")

	// 完全接管数据源：返回内存数据，绝不查 DB
	lb.SearchFunc(func(ctx *web.EventContext, params *presets.SearchParams) (*presets.SearchResult, error) {
		rows := nodataRows()
		total := len(rows)
		res := &presets.SearchResult{
			Nodes:      rows,
			TotalCount: &total,
		}
		// 有数据时设 StartCursor 非 nil，避开「无记录」误判；空数据留 nil + total=0 → 自动空态
		if total > 0 {
			cursor := "0"
			res.PageInfo = relay.PageInfo{StartCursor: &cursor, EndCursor: &cursor}
		}
		return res, nil
	})
}
