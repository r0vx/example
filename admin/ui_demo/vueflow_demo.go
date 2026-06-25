package ui_demo

import (
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/vueflow"
)

// VueFlowDemo 虚拟模型（无数据库）
type VueFlowDemo struct{}

// emptyVueFlowSearchFunc 返回空数据，避免数据库查询
func emptyVueFlowSearchFunc() presets.SearchFunc {
	return func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
		totalCount := 0
		result = &presets.SearchResult{Nodes: []VueFlowDemo{}, TotalCount: &totalCount}
		return
	}
}

// ConfigVueFlowDemo 配置 Vue Flow 通用画布演示页（纯 UI，无数据库）
func ConfigVueFlowDemo(b *presets.Builder) {
	m := b.Model(&VueFlowDemo{}).
		Label("Vue Flow").
		URIName("vueflow-demo")
	m.Listing().SearchFunc(emptyVueFlowSearchFunc())
	m.Editing().Only()

	// 三个事件 toast 出 payload，验证回传链路
	b.GetWebBuilder().RegisterEventFunc("vueflow_demo_nodeClick", vueflowDemoEcho("点击节点"))
	b.GetWebBuilder().RegisterEventFunc("vueflow_demo_connect", vueflowDemoEcho("新建连线"))
	b.GetWebBuilder().RegisterEventFunc("vueflow_demo_dragStop", vueflowDemoEcho("拖动结束"))

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.PageTitle = "Vue Flow Demo"
		r.Body = vueFlowDemoBody()
		return
	})
}

// vueflowDemoEcho 返回一个把收到的 payload 弹 toast 的 EventFunc
func vueflowDemoEcho(label string) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		payload := ctx.R.FormValue("payload")
		// demo 仅透传原始 payload，不做结构校验
		presets.ShowMessage(&r, label+"："+payload, "success")
		return
	}
}

// vueFlowDemoBody 渲染演示画布
func vueFlowDemoBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "1", "position": map[string]any{"x": 50, "y": 50}, "data": map[string]any{"label": "用户"}},
		{"id": "2", "position": map[string]any{"x": 300, "y": 50}, "data": map[string]any{"label": "订单"}},
		{"id": "3", "position": map[string]any{"x": 300, "y": 200}, "data": map[string]any{"label": "商品"}},
	}
	edges := []map[string]any{
		{"id": "e1-2", "source": "1", "target": "2"},
		{"id": "e2-3", "source": "2", "target": "3"},
	}

	field := func(eventFunc string) string {
		return web.Plaid().EventFunc(eventFunc).
			FieldValue("payload", web.Var("JSON.stringify($event)")).Go()
	}

	return h.Div(
		vueflow.Flow().
			Nodes(nodes).
			Edges(edges).
			Height("600px").
			Background(true).
			Controls(true).
			MiniMap(true).
			Fit(true).
			OnNodeClick(field("vueflow_demo_nodeClick")).
			OnConnect(field("vueflow_demo_connect")).
			OnNodeDragStop(field("vueflow_demo_dragStop")),
	).Class("p-4")
}
