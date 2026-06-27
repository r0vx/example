package ui_demo

import (
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/vueflow"
)

// 三个虚拟模型（无数据库）—— 每个独立页面：基础画布 / dagre 布局 / 状态卡片
type (
	VueFlowDemo         struct{} // 基础通用画布
	VueFlowDagreDemo    struct{} // dagre 自动布局
	VueFlowStatusDemo   struct{} // status 状态卡片
	VueFlowResizerDemo  struct{} // 节点缩放
	VueFlowToolbarDemo  struct{} // 节点工具栏
	VueFlowDndDemo      struct{} // 拖放建节点
	VueFlowEdgesDemo    struct{} // 边类型/箭头
	VueFlowMathDemo     struct{} // 数学运算流
	VueFlowViewportDemo struct{} // 视口过渡/截图
	VueFlowTeleportDemo struct{} // 传送节点
	VueFlowDragAidsDemo struct{} // 拖拽辅助（辅助线+相交）
	VueFlowConnDemo     struct{} // 连线进阶（自定义线+半径+重连）
)

// ConfigVueFlowDemo 注册 Vue Flow 三个独立演示页（纯 UI，无数据库）
func ConfigVueFlowDemo(b *presets.Builder) {
	wb := b.GetWebBuilder()
	// 事件 toast 出 payload，验证回传链路
	wb.RegisterEventFunc("vueflow_demo_nodeClick", vueflowDemoEcho("点击节点"))
	wb.RegisterEventFunc("vueflow_demo_connect", vueflowDemoEcho("新建连线"))
	wb.RegisterEventFunc("vueflow_demo_dragStop", vueflowDemoEcho("拖动结束"))
	wb.RegisterEventFunc("vueflow_demo_nodeAction", vueflowDemoEcho("卡片按钮"))
	// 轻量补丁：把 svc-api 状态改 error，只该节点变色、画布不重排
	wb.RegisterEventFunc("vueflow_demo_patch", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.RunScript = vueflow.PatchScript("vfstatus", map[string]any{
			"svc-api": map[string]any{"status": "error"},
		})
		return
	})

	wb.RegisterEventFunc("vueflow_demo_drop", vueflowDemoEcho("拖放建节点"))
	wb.RegisterEventFunc("vueflow_demo_edgeUpdate", vueflowDemoEcho("边重连"))
	wb.RegisterEventFunc("vueflow_demo_intersect", vueflowDemoEcho("相交"))

	vueflowUIPage(b, &VueFlowDemo{}, "vueflow-demo", "Vue Flow 基础", vueFlowBasicBody)
	vueflowUIPage(b, &VueFlowDagreDemo{}, "vueflow-dagre-demo", "Vue Flow dagre 布局", vueFlowDagreBody)
	vueflowUIPage(b, &VueFlowStatusDemo{}, "vueflow-status-demo", "Vue Flow 状态卡片", vueFlowStatusBody)
	vueflowUIPage(b, &VueFlowResizerDemo{}, "vueflow-resizer-demo", "Vue Flow 节点缩放", vueFlowResizerBody)
	vueflowUIPage(b, &VueFlowToolbarDemo{}, "vueflow-toolbar-demo", "Vue Flow 节点工具栏", vueFlowToolbarBody)
	vueflowUIPage(b, &VueFlowDndDemo{}, "vueflow-dnd-demo", "Vue Flow 拖放建节点", vueFlowDndBody)
	vueflowUIPage(b, &VueFlowEdgesDemo{}, "vueflow-edges-demo", "Vue Flow 边类型/箭头", vueFlowEdgesBody)
	vueflowUIPage(b, &VueFlowMathDemo{}, "vueflow-math-demo", "Vue Flow 数学运算流", vueFlowMathBody)
	vueflowUIPage(b, &VueFlowViewportDemo{}, "vueflow-viewport-demo", "Vue Flow 视口/截图", vueFlowViewportBody)
	vueflowUIPage(b, &VueFlowTeleportDemo{}, "vueflow-teleport-demo", "Vue Flow 传送节点", vueFlowTeleportBody)
	vueflowUIPage(b, &VueFlowDragAidsDemo{}, "vueflow-dragaids-demo", "Vue Flow 拖拽辅助", vueFlowDragAidsBody)
	vueflowUIPage(b, &VueFlowConnDemo{}, "vueflow-conn-demo", "Vue Flow 连线进阶", vueFlowConnBody)
}

// vueflowUIPage 注册一个纯 UI 的 Vue Flow 演示页（空数据 + PageFunc 全量渲染画布）
func vueflowUIPage(b *presets.Builder, model any, uriName, label string, body func() h.HTMLComponent) {
	m := b.Model(model).Label(label).URIName(uriName)
	m.Listing().SearchFunc(func(ctx *web.EventContext, _ *presets.SearchParams) (*presets.SearchResult, error) {
		tc := 0
		return &presets.SearchResult{Nodes: []any{}, TotalCount: &tc}, nil
	})
	m.Editing().Only()
	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.PageTitle = label
		r.Body = body()
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

// vueflowDemoField 构造一个把 $event 整体 JSON 回传给指定 EventFunc 的 Plaid 表达式
func vueflowDemoField(eventFunc string) string {
	return web.Plaid().EventFunc(eventFunc).
		FieldValue("payload", web.Var("JSON.stringify($event)")).Go()
}

// vueFlowBasicBody 基础画布：手动 position + 点击/连线/拖动事件回传
func vueFlowBasicBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "1", "position": map[string]any{"x": 50, "y": 50}, "data": map[string]any{"label": "用户"}},
		{"id": "2", "position": map[string]any{"x": 300, "y": 50}, "data": map[string]any{"label": "订单"}},
		{"id": "3", "position": map[string]any{"x": 300, "y": 200}, "data": map[string]any{"label": "商品"}},
		// 透传测试：group 容器节点 + 两个子节点（parentNode + extent:parent）
		{"id": "g1", "type": "group", "position": map[string]any{"x": 50, "y": 320},
			"style": map[string]any{"width": "320px", "height": "180px", "backgroundColor": "rgba(99,102,241,0.06)"},
			"data":  map[string]any{"label": "分组"}},
		{"id": "g1a", "position": map[string]any{"x": 20, "y": 40}, "parentNode": "g1", "extent": "parent",
			"data": map[string]any{"label": "子 A"}},
		{"id": "g1b", "position": map[string]any{"x": 180, "y": 40}, "parentNode": "g1", "extent": "parent",
			"data": map[string]any{"label": "子 B"}},
	}
	edges := []map[string]any{
		{"id": "e1-2", "source": "1", "target": "2"},
		{"id": "e2-3", "source": "2", "target": "3"},
		// 透传测试：动画 + 标签 + 红色描边
		{"id": "e1-3", "source": "1", "target": "3", "animated": true, "label": "data",
			"style": map[string]any{"stroke": "red"}},
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
			OnNodeClick(vueflowDemoField("vueflow_demo_nodeClick")).
			OnConnect(vueflowDemoField("vueflow_demo_connect")).
			OnNodeDragStop(vueflowDemoField("vueflow_demo_dragStop")),
	).Class("p-4")
}

// vueFlowDagreBody dagre 自动布局：节点不传 position，前端按分层自动排版
func vueFlowDagreBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "a", "data": map[string]any{"label": "根"}},
		{"id": "b", "data": map[string]any{"label": "子 B"}},
		{"id": "c", "data": map[string]any{"label": "子 C"}},
		{"id": "d", "data": map[string]any{"label": "孙 D"}},
		{"id": "e", "data": map[string]any{"label": "孙 E"}},
	}
	edges := []map[string]any{
		{"id": "ea-b", "source": "a", "target": "b"},
		{"id": "ea-c", "source": "a", "target": "c"},
		{"id": "eb-d", "source": "b", "target": "d"},
		{"id": "ec-e", "source": "c", "target": "e"},
	}
	return h.Div(
		h.P(h.Text("节点不传 position，前端用 dagre 自动分层排版（TB 方向）")).Class("mb-2 text-sm text-muted-foreground"),
		vueflow.Flow().
			Nodes(nodes).
			Edges(edges).
			Height("600px").
			Background(true).
			Controls(true).
			MiniMap(true).
			Layout("dagre").
			LayoutDirection("TB").
			Fit(true),
	).Class("p-4")
}

// vueFlowStatusBody status 状态卡片：色条/metrics/badges + 卡片内按钮发 Plaid 事件
func vueFlowStatusBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "svc-api", "type": "status", "data": map[string]any{
			"title": "API 服务", "subtitle": "api.r0vx.io", "status": "ok", "icon": "●",
			"metrics": []map[string]any{{"label": "QPS", "value": "1.2k"}, {"label": "延迟", "value": "23ms"}},
			"badges":  []map[string]any{{"text": "prod"}, {"text": "v2.1"}},
			"actions": []map[string]any{{"id": "restart", "label": "重启", "icon": "↻"}, {"id": "logs", "label": "日志", "icon": "▤"}},
		}},
		{"id": "svc-db", "type": "status", "data": map[string]any{
			"title": "数据库", "subtitle": "postgres-primary", "status": "warn", "icon": "●",
			"metrics": []map[string]any{{"label": "连接", "value": "87/100"}, {"label": "磁盘", "value": "78%"}},
			"badges":  []map[string]any{{"text": "primary"}},
			"actions": []map[string]any{{"id": "restart", "label": "重启", "icon": "↻"}, {"id": "logs", "label": "日志", "icon": "▤"}},
		}},
	}
	edges := []map[string]any{
		{"id": "es-api-db", "source": "svc-api", "target": "svc-db"},
	}
	return h.Div(
		h.P(h.Text("type:\"status\" 卡片，点卡片内按钮经 @node-action 回 Go（带 {id, action}）")).Class("mb-2 text-sm text-muted-foreground"),
		// 轻量补丁演示：点按钮经 PatchScript 只改 svc-api 的 status，不重建整图
		h.Button("模拟 API 故障（patch status=error）").
			Attr("@click", web.Plaid().EventFunc("vueflow_demo_patch").Go()).
			Class("mb-2 inline-flex items-center px-3 py-1.5 text-sm border rounded-md hover:bg-accent cursor-pointer"),
		vueflow.Flow().
			Nodes(nodes).
			Edges(edges).
			Height("600px").
			Background(true).
			Controls(true).
			MiniMap(true).
			Layout("dagre").
			LayoutDirection("LR").
			Fit(true).
			CanvasKey("vfstatus").
			OnNodeAction(vueflowDemoField("vueflow_demo_nodeAction")),
	).Class("p-4")
}

// vueFlowResizerBody 节点缩放：type:"resizable" 节点拖四角/四边手柄改尺寸
func vueFlowResizerBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "r1", "type": "resizable", "position": map[string]any{"x": 80, "y": 80},
			"style": map[string]any{"width": "200px", "height": "120px"},
			"data":  map[string]any{"label": "拖我四角缩放", "minWidth": 120, "minHeight": 60}},
		{"id": "r2", "type": "resizable", "position": map[string]any{"x": 380, "y": 120},
			"style": map[string]any{"width": "160px", "height": "90px"},
			"data":  map[string]any{"label": "可缩放节点 2"}},
	}
	edges := []map[string]any{{"id": "er1-2", "source": "r1", "target": "r2"}}
	return h.Div(
		h.P(h.Text("type:\"resizable\" 节点，选中后拖四角/四边手柄改尺寸（@vue-flow/node-resizer）")).Class("mb-2 text-sm text-muted-foreground"),
		vueflow.Flow().
			Nodes(nodes).Edges(edges).
			Height("600px").Background(true).Controls(true).MiniMap(true).Fit(true),
	).Class("p-4")
}

// vueFlowToolbarBody 节点工具栏：type:"toolbar" 节点选中时上方浮出按钮，点击经 node-action 回 Go
func vueFlowToolbarBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "t1", "type": "toolbar", "position": map[string]any{"x": 120, "y": 100},
			"data": map[string]any{"label": "选中我看工具栏", "toolbarPosition": "top",
				"actions": []map[string]any{{"id": "edit", "label": "编辑", "icon": "✎"}, {"id": "delete", "label": "删除", "icon": "🗑"}}}},
		{"id": "t2", "type": "toolbar", "position": map[string]any{"x": 420, "y": 100},
			"data": map[string]any{"label": "常驻工具栏", "toolbarPosition": "right", "alwaysVisible": true,
				"actions": []map[string]any{{"id": "run", "label": "运行", "icon": "▶"}}}},
	}
	edges := []map[string]any{{"id": "et1-2", "source": "t1", "target": "t2"}}
	return h.Div(
		h.P(h.Text("type:\"toolbar\" 节点：选中浮出工具栏（t2 常驻），按钮经 @node-action 回 Go（@vue-flow/node-toolbar）")).Class("mb-2 text-sm text-muted-foreground"),
		vueflow.Flow().
			Nodes(nodes).Edges(edges).
			Height("600px").Background(true).Controls(true).MiniMap(true).Fit(true).
			OnNodeAction(vueflowDemoField("vueflow_demo_nodeAction")),
	).Class("p-4")
}

// vueFlowDndBody 拖放建节点：左侧调色板拖到画布，按投影坐标建节点并 emit node-drop
func vueFlowDndBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "seed", "position": map[string]any{"x": 250, "y": 60}, "data": map[string]any{"label": "已有节点"}},
	}
	// 调色板项：draggable + dragstart 写 dataTransfer（type 名）
	paletteItem := func(typ, label string) h.HTMLComponent {
		return h.Div(h.Text(label)).
			Attr("draggable", "true").
			Attr("@dragstart", "$event.dataTransfer.setData('application/vueflow', '"+typ+"'); $event.dataTransfer.effectAllowed='move'").
			Class("px-3 py-2 mb-2 text-sm border rounded-md bg-background cursor-grab select-none hover:bg-accent")
	}
	return h.Div(
		h.P(h.Text("把左侧节点拖到画布——按投影坐标建节点并经 node-drop 回 Go")).Class("mb-2 text-sm text-muted-foreground"),
		h.Div(
			// 左侧调色板
			h.Div(
				h.Div(h.Text("调色板")).Class("mb-2 text-xs font-medium text-muted-foreground"),
				paletteItem("default", "默认节点"),
				paletteItem("input", "输入节点"),
				paletteItem("output", "输出节点"),
			).Class("w-32 shrink-0 pr-3 border-r"),
			// 右侧画布
			h.Div(
				vueflow.Flow().
					Nodes(nodes).Edges([]map[string]any{}).
					Height("600px").Background(true).Controls(true).MiniMap(true).Fit(true).
					Droppable(true).
					OnNodeDrop(vueflowDemoField("vueflow_demo_drop")),
			).Class("flex-1 min-w-0"),
		).Class("flex"),
	).Class("p-4")
}

// vueFlowEdgesBody 边类型/箭头：default/step/smoothstep/straight + markerEnd 箭头（原生透传，无源码改动）
func vueFlowEdgesBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "a", "position": map[string]any{"x": 60, "y": 40}, "data": map[string]any{"label": "A"}},
		{"id": "b", "position": map[string]any{"x": 360, "y": 40}, "data": map[string]any{"label": "B"}},
		{"id": "c", "position": map[string]any{"x": 60, "y": 200}, "data": map[string]any{"label": "C"}},
		{"id": "d", "position": map[string]any{"x": 360, "y": 200}, "data": map[string]any{"label": "D"}},
		{"id": "e", "position": map[string]any{"x": 60, "y": 360}, "data": map[string]any{"label": "E"}},
	}
	arrow := map[string]any{"type": "arrowclosed", "color": "#2563eb"}
	edges := []map[string]any{
		{"id": "ed-default", "source": "a", "target": "b", "type": "default", "label": "default", "markerEnd": arrow},
		{"id": "ed-step", "source": "a", "target": "c", "type": "step", "label": "step", "markerEnd": "arrowclosed"},
		{"id": "ed-smooth", "source": "b", "target": "d", "type": "smoothstep", "label": "smoothstep", "markerEnd": "arrow"},
		{"id": "ed-straight", "source": "c", "target": "d", "type": "straight", "label": "straight"},
		{"id": "ed-loop", "source": "e", "target": "e", "type": "default", "label": "loopback", "markerEnd": arrow},
	}
	return h.Div(
		h.P(h.Text("边类型 default/step/smoothstep/straight + markerEnd 箭头 + 自环（loopback），全靠原生透传")).Class("mb-2 text-sm text-muted-foreground"),
		vueflow.Flow().
			Nodes(nodes).Edges(edges).
			Height("600px").Background(true).Controls(true).MiniMap(true).Fit(true),
	).Class("p-4")
}

// vueFlowMathBody 数学运算流：input 节点编辑数值 → op 节点对入边上游值实时运算
func vueFlowMathBody() h.HTMLComponent {
	mathNode := func(id string, x, y int, label, kind, op string, val int) map[string]any {
		data := map[string]any{"label": label, "kind": kind}
		if kind == "input" {
			data["value"] = val
		} else {
			data["op"] = op
		}
		return map[string]any{"id": id, "type": "math", "position": map[string]any{"x": x, "y": y}, "data": data}
	}
	nodes := []map[string]any{
		mathNode("in1", 40, 40, "输入 A", "input", "", 3),
		mathNode("in2", 40, 160, "输入 B", "input", "", 4),
		mathNode("in3", 40, 300, "输入 C", "input", "", 5),
		mathNode("sum", 280, 100, "求和", "op", "add", 0),
		mathNode("prod", 520, 180, "乘积", "op", "mul", 0),
	}
	edges := []map[string]any{
		{"id": "m1", "source": "in1", "target": "sum"},
		{"id": "m2", "source": "in2", "target": "sum"},
		{"id": "m3", "source": "sum", "target": "prod"},
		{"id": "m4", "source": "in3", "target": "prod"},
	}
	return h.Div(
		h.P(h.Text("改输入框数值，下游 op 节点（求和/乘积）实时重算——节点式计算流")).Class("mb-2 text-sm text-muted-foreground"),
		vueflow.Flow().
			Nodes(nodes).Edges(edges).
			Height("600px").Background(true).Controls(true).MiniMap(true).Fit(true),
	).Class("p-4")
}

// vueFlowViewportBody 视口过渡 + 截图：按钮带动画 fitView / 聚焦远节点 / 导出 PNG
func vueFlowViewportBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "near", "position": map[string]any{"x": 60, "y": 60}, "data": map[string]any{"label": "近节点"}},
		{"id": "mid", "position": map[string]any{"x": 500, "y": 300}, "data": map[string]any{"label": "中间"}},
		{"id": "far", "position": map[string]any{"x": 1400, "y": 900}, "data": map[string]any{"label": "远节点"}},
	}
	edges := []map[string]any{
		{"id": "v1", "source": "near", "target": "mid"},
		{"id": "v2", "source": "mid", "target": "far"},
	}
	btn := func(label, script string) h.HTMLComponent {
		return h.Button(label).Attr("@click", script).
			Class("mr-2 mb-2 inline-flex items-center px-3 py-1.5 text-sm border rounded-md hover:bg-accent cursor-pointer")
	}
	return h.Div(
		h.P(h.Text("视口带动画过渡 + 画布导出 PNG")).Class("mb-2 text-sm text-muted-foreground"),
		h.Div(
			btn("全览(fit)", vueflow.ViewportFitScript("vpshot", 800)),
			btn("聚焦远节点", vueflow.ViewportToNodeScript("vpshot", "far", 2.0, 800)),
			btn("截图下载 PNG", vueflow.ScreenshotScript("vpshot", "vueflow.png")),
		),
		vueflow.Flow().
			Nodes(nodes).Edges(edges).
			Height("560px").Background(true).Controls(true).MiniMap(true).Fit(true).
			CanvasKey("vpshot"),
	).Class("p-4")
}

// vueFlowTeleportBody 传送节点：节点详情经 Vue Teleport 渲到画布外的侧边面板
func vueFlowTeleportBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "tp1", "type": "teleport", "position": map[string]any{"x": 80, "y": 60},
			"data": map[string]any{"label": "服务 A", "to": "#vf-teleport-target", "detail": "CPU 32% / 内存 1.2G"}},
		{"id": "tp2", "type": "teleport", "position": map[string]any{"x": 80, "y": 200},
			"data": map[string]any{"label": "服务 B", "to": "#vf-teleport-target", "detail": "CPU 8% / 内存 512M"}},
	}
	return h.Div(
		h.P(h.Text("节点主体在画布，详情经 Teleport 渲到右侧外部面板")).Class("mb-2 text-sm text-muted-foreground"),
		h.Div(
			h.Div(
				vueflow.Flow().
					Nodes(nodes).Edges([]map[string]any{}).
					Height("560px").Background(true).Controls(true).MiniMap(false).Fit(true),
			).Class("flex-1 min-w-0"),
			// 画布外的 teleport 目标面板
			h.Div(
				h.Div(h.Text("外部详情面板（teleport 目标）")).Class("mb-2 text-xs font-medium text-muted-foreground"),
				h.Div().Id("vf-teleport-target"),
			).Class("w-64 shrink-0 pl-3 border-l"),
		).Class("flex"),
	).Class("p-4")
}

// vueFlowDragAidsBody 拖拽辅助：辅助线对齐 + 相交高亮
func vueFlowDragAidsBody() h.HTMLComponent {
	nodes := func() []map[string]any {
		return []map[string]any{
			{"id": "n1", "position": map[string]any{"x": 120, "y": 80}, "data": map[string]any{"label": "拖我对齐"}},
			{"id": "n2", "position": map[string]any{"x": 360, "y": 80}, "data": map[string]any{"label": "参照 1"}},
			{"id": "n3", "position": map[string]any{"x": 120, "y": 300}, "data": map[string]any{"label": "参照 2"}},
		}
	}
	return h.Div(
		h.P(h.Text("左：拖节点与其他节点左/顶缘对齐时出蓝色辅助线；右：拖节点压到别的节点上高亮橙框")).Class("mb-2 text-sm text-muted-foreground"),
		h.Div(
			h.Div(
				h.Div(h.Text("辅助线（helperLines）")).Class("mb-1 text-xs font-medium"),
				vueflow.Flow().
					Nodes(nodes()).Edges([]map[string]any{}).
					Height("460px").Background(true).Controls(true).Fit(true).
					HelperLines(true),
			).Class("flex-1 min-w-0"),
			h.Div(
				h.Div(h.Text("相交高亮（intersections）")).Class("mb-1 text-xs font-medium"),
				vueflow.Flow().
					Nodes(nodes()).Edges([]map[string]any{}).
					Height("460px").Background(true).Controls(true).Fit(true).
					Intersections(true).
					OnNodeIntersect(vueflowDemoField("vueflow_demo_intersect")),
			).Class("flex-1 min-w-0 ml-3"),
		).Class("flex"),
	).Class("p-4")
}

// vueFlowConnBody 连线进阶：拖边端点重连 + 自定义连接线 + 加大吸附半径
func vueFlowConnBody() h.HTMLComponent {
	nodes := []map[string]any{
		{"id": "s1", "position": map[string]any{"x": 80, "y": 80}, "data": map[string]any{"label": "源 1"}},
		{"id": "s2", "position": map[string]any{"x": 80, "y": 240}, "data": map[string]any{"label": "源 2"}},
		{"id": "t1", "position": map[string]any{"x": 420, "y": 160}, "data": map[string]any{"label": "目标"}},
	}
	edges := []map[string]any{
		{"id": "c1", "source": "s1", "target": "t1", "label": "拖端点重连我", "updatable": true},
	}
	return h.Div(
		h.P(h.Text("拖已有边端点改连源 1/源 2（重连回 Go）；从 handle 拖出新连线显示蓝色虚线预览；吸附半径加大到 60px")).Class("mb-2 text-sm text-muted-foreground"),
		vueflow.Flow().
			Nodes(nodes).Edges(edges).
			Height("560px").Background(true).Controls(true).MiniMap(true).Fit(true).
			EdgesUpdatable(true).
			CustomConnectionLine(true).
			ConnectionRadius(60).
			OnEdgeUpdate(vueflowDemoField("vueflow_demo_edgeUpdate")),
	).Class("p-4")
}
