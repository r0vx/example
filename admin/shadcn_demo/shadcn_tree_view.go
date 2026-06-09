package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnTreeViewDemo 虚拟模型
type ShadcnTreeViewDemo struct{}

// configTreeView 注册 Tree View demo
func configTreeView(b *presets.Builder) {
	m := b.Model(&ShadcnTreeViewDemo{}).
		Label("Tree View").
		URIName("shadcn-tree-view")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnTreeViewDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Tree View"
		r.Body = shadcnTreeViewBody(ctx)
		return
	})
}

// shadcnTreeViewBody 树形组件演示页面
func shadcnTreeViewBody(ctx *web.EventContext) h.HTMLComponent {
	// 示例树形数据
	treeData := []map[string]interface{}{
		{
			"id":    "1",
			"title": "文档",
			"children": []map[string]interface{}{
				{
					"id":    "1-1",
					"title": "介绍",
				},
				{
					"id":    "1-2",
					"title": "快速开始",
					"children": []map[string]interface{}{
						{"id": "1-2-1", "title": "安装"},
						{"id": "1-2-2", "title": "配置"},
					},
				},
			},
		},
		{
			"id":    "2",
			"title": "组件",
			"children": []map[string]interface{}{
				{"id": "2-1", "title": "Button"},
				{"id": "2-2", "title": "Input"},
				{"id": "2-3", "title": "Dialog"},
			},
		},
		{
			"id":    "3",
			"title": "API 参考",
		},
	}

	return h.Div(
		h.H1("TreeView 树形组件").Style("margin-bottom: 24px;"),

		// 基本用法
		h.Div(
			h.H2("基本用法"),
			h.P(h.Text("展示层级数据的树形结构")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					TreeView().
						Items(treeData).
						Attr("v-model", "form.selected").
						Attr("v-model:expanded", "form.expanded").
						Class("border rounded-md p-4 max-w-md"),
				),
				h.Div(
					h.Span("选中: {{ form.selected || '无' }}").Class("text-sm text-muted-foreground"),
				).Class("mt-2"),
				h.Div(
					h.Span("展开: {{ form.expanded?.join(', ') || '无' }}").Class("text-sm text-muted-foreground"),
				).Class("mt-1"),
			).VSlot("{ form }").FormInit(`{ "selected": null, "expanded": ["1", "1-2"] }`),
		).Class("demo-section"),

		// 多选模式
		h.Div(
			h.H2("多选模式"),
			h.P(h.Text("支持选择多个节点")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					TreeView().
						Items(treeData).
						Multiple(true).
						Attr("v-model", "form.selected").
						Class("border rounded-md p-4 max-w-md"),
				),
				h.Div(
					h.Span("选中: {{ form.selected?.join(', ') || '无' }}").Class("text-sm text-muted-foreground"),
				).Class("mt-2"),
			).VSlot("{ form }").FormInit(`{ "selected": [] }`),
		).Class("demo-section"),

		// 自定义字段
		h.Div(
			h.H2("自定义字段"),
			h.P(h.Text("支持自定义值和子节点字段名")).Class("text-muted-foreground mb-4"),
			TreeView().
				Items([]map[string]interface{}{
					{
						"value": "node1",
						"name":  "节点 1",
						"items": []map[string]interface{}{
							{"value": "node1-1", "name": "子节点 1-1"},
							{"value": "node1-2", "name": "子节点 1-2"},
						},
					},
					{
						"value": "node2",
						"name":  "节点 2",
					},
				}).
				ItemValue("value").
				ItemChildren("items").
				Class("border rounded-md p-4 max-w-md"),
		).Class("demo-section"),
	).Style("max-width: 800px; margin: 0 auto;")
}
