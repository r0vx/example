package chart_demo

import (
	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/unovis"
	h "github.com/r0vx/htmlgo"
)

// configGraphDemo 配置 Graph 图表组件演示
func ConfigGraphDemo(b *presets.Builder) {

	m := b.Model(&models.GraphDemo{}).
		Label("Graph 图表").
		URIName("graph-demos")

	// 完全自定义数据获取，避免数据库查询
	m.Listing().SearchFunc(func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
		// 返回空数据，不查询数据库
		totalCount := 0
		result = &presets.SearchResult{
			Nodes:      []models.GraphDemo{},
			TotalCount: &totalCount,
		}
		return
	})

	// 禁用新增、编辑、删除等操作
	m.Listing().ActionsAsMenu(false)
	m.Editing().Only()

	// 自定义页面函数
	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.PageTitle = "Graph 图表组件演示"
		r.Body = graphDemoPage()
		return
	})
}

// graphDemoPage 图表演示页面
func graphDemoPage() h.HTMLComponent {
	return web.Scope(
		h.Div(
			h.H1("Graph 图表组件演示").Class("text-2xl font-bold mb-6"),

			// 示例 1: 简单的三节点图（带流动效果）
			h.Div(
				h.H2("示例 1: 基础布局 + 流动动画 (TB)").Class("text-xl font-semibold mb-4"),
				unovis.Graph().
					Direction(unovis.GraphDirectionTB).
					NodeSize(50).
					ShowLabels(true).
					LinkFlow(true).
					LinkFlowParticleSpeed(50).
					Config(unovis.GraphConfig{
						"user": {
							Label: "用户模块",
							Color: "var(--chart-1)",
							Shape: unovis.GraphNodeShapeCircle,
						},
						"order": {
							Label: "订单模块",
							Color: "var(--chart-2)",
							Shape: unovis.GraphNodeShapeSquare,
						},
						"payment": {
							Label: "支付模块",
							Color: "var(--chart-3)",
							Shape: unovis.GraphNodeShapeHexagon,
						},
					}).
					Data(unovis.GraphData{
						Nodes: []map[string]any{
							{"id": "user"},
							{"id": "order"},
							{"id": "payment"},
						},
						Links: []map[string]any{
							{"source": "user", "target": "order"},
							{"source": "order", "target": "payment"},
						},
					}).
					Class("mb-8 h-96"),
			).Class("mb-8"),

			// 示例 2: 左右布局（LR）
			h.Div(
				h.H2("示例 2: 水平布局 (LR)").Class("text-xl font-semibold mb-4"),
				unovis.Graph().
					Direction(unovis.GraphDirectionLR).
					NodeSize(40).
					ShowLabels(true).
					Config(unovis.GraphConfig{
						"frontend": {
							Label: "前端",
							Color: "var(--chart-1)",
							Shape: unovis.GraphNodeShapeCircle,
						},
						"api": {
							Label: "API",
							Color: "var(--chart-2)",
							Shape: unovis.GraphNodeShapeSquare,
						},
						"database": {
							Label: "数据库",
							Color: "var(--chart-3)",
							Shape: unovis.GraphNodeShapeCircle,
						},
						"cache": {
							Label: "缓存",
							Color: "var(--chart-4)",
							Shape: unovis.GraphNodeShapeTriangle,
						},
					}).
					Data(unovis.GraphData{
						Nodes: []map[string]any{
							{"id": "frontend"},
							{"id": "api"},
							{"id": "database"},
							{"id": "cache"},
						},
						Links: []map[string]any{
							{"source": "frontend", "target": "api"},
							{"source": "api", "target": "database"},
							{"source": "api", "target": "cache"},
						},
					}).
					Class("mb-8 h-96"),
			).Class("mb-8"),

			// 示例 3: 复杂依赖图
			h.Div(
				h.H2("示例 3: 复杂依赖关系").Class("text-xl font-semibold mb-4"),
				unovis.Graph().
					Direction(unovis.GraphDirectionTB).
					NodeSize(45).
					ShowLabels(true).
					Config(unovis.GraphConfig{
						"auth": {
							Label: "认证",
							Color: "var(--chart-1)",
						},
						"user": {
							Label: "用户",
							Color: "var(--chart-2)",
						},
						"product": {
							Label: "商品",
							Color: "var(--chart-3)",
						},
						"order": {
							Label: "订单",
							Color: "var(--chart-4)",
						},
						"payment": {
							Label: "支付",
							Color: "var(--chart-5)",
						},
						"notification": {
							Label: "通知",
							Color: "var(--chart-1)",
						},
					}).
					Data(unovis.GraphData{
						Nodes: []map[string]any{
							{"id": "auth"},
							{"id": "user"},
							{"id": "product"},
							{"id": "order"},
							{"id": "payment"},
							{"id": "notification"},
						},
						Links: []map[string]any{
							{"source": "auth", "target": "user"},
							{"source": "user", "target": "order"},
							{"source": "product", "target": "order"},
							{"source": "order", "target": "payment"},
							{"source": "payment", "target": "notification"},
							{"source": "order", "target": "notification"},
						},
					}).
					Class("h-[600px]"),
			).Class("mb-8"),
		).Class("container mx-auto p-6"),
	).VSlot("{ locals }")
}
