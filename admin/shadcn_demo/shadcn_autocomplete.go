package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnAutocompleteDemo 虚拟模型
type ShadcnAutocompleteDemo struct{}

// configAutocomplete 注册 Autocomplete demo
func configAutocomplete(b *presets.Builder) {
	m := b.Model(&ShadcnAutocompleteDemo{}).
		Label("Autocomplete").
		URIName("shadcn-autocomplete")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnAutocompleteDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Autocomplete"
		r.Body = shadcnAutocompleteBody(ctx)
		return
	})
}

// shadcnAutocompleteBody Autocomplete 组件完整演示
func shadcnAutocompleteBody(ctx *web.EventContext) h.HTMLComponent {
	// 基本用法的选项（Go 结构体，Items() 会正确 JSON 序列化）
	basicItems := []shadcn.DefaultOptionItem{
		{Value: "1", Text: "Apple"},
		{Value: "2", Text: "Banana"},
		{Value: "3", Text: "Cherry"},
		{Value: "4", Text: "Date"},
		{Value: "5", Text: "Elderberry"},
	}

	// 多选模式的选项
	frameworkItems := []shadcn.DefaultOptionItem{
		{Value: "vue", Text: "Vue.js"},
		{Value: "react", Text: "React"},
		{Value: "angular", Text: "Angular"},
		{Value: "svelte", Text: "Svelte"},
		{Value: "solid", Text: "SolidJS"},
	}

	// 分组显示的选项（含 group 字段，用 map）
	groupedItems := []map[string]string{
		{"value": "1", "text": "Apple", "group": "水果"},
		{"value": "2", "text": "Banana", "group": "水果"},
		{"value": "3", "text": "Carrot", "group": "蔬菜"},
		{"value": "4", "text": "Celery", "group": "蔬菜"},
		{"value": "5", "text": "Chicken", "group": "肉类"},
		{"value": "6", "text": "Beef", "group": "肉类"},
	}

	return h.Div(
		// 页面标题
		h.Div(
			h.H1("Autocomplete 组件").Style("font-size: 28px; font-weight: bold; margin-bottom: 8px;"),
			h.P(h.Text("高级自动完成组件,支持远程搜索、分组显示和创建新项")).Style("color: #666; margin-bottom: 24px;"),
		),

		// 1. 基本用法
		h.Div(
			h.H2("1. 基本用法").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("从预定义列表中选择项目")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			web.Scope(
				shadcn.Autocomplete().
					Attr("v-model", "locals.basicValue").
					Items(basicItems).
					Placeholder("选择水果..."),
				h.Div(
					h.Text("选中的值: "),
					h.Tag("code").Children(h.Text("{{locals.basicValue}}")).Style("background: #f3f4f6; padding: 2px 8px; border-radius: 4px; font-family: monospace;"),
				).Style("margin-top: 12px; font-size: 14px;"),
			).Init(`{basicValue: ''}`).VSlot("{ locals }"),
		).Class("demo-section"),

		// 2. 多选模式
		h.Div(
			h.H2("2. 多选模式").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("支持选择多个选项，以标签形式展示")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			web.Scope(
				shadcn.Autocomplete().
					Attr("v-model", "locals.multipleValue").
					Items(frameworkItems).
					Multiple(true).
					Placeholder("选择前端框架..."),
				h.Div(
					h.Text("选中的值: "),
					h.Tag("code").Children(h.Text("{{locals.multipleValue}}")).Style("background: #f3f4f6; padding: 2px 8px; border-radius: 4px; font-family: monospace;"),
				).Style("margin-top: 12px; font-size: 14px;"),
			).Init(`{multipleValue: []}`).VSlot("{ locals }"),
		).Class("demo-section"),

		// 3. 分组显示
		h.Div(
			h.H2("3. 分组显示").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("按类别分组的选项列表")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			web.Scope(
				shadcn.Autocomplete().
					Attr("v-model", "locals.groupedValue").
					Items(groupedItems).
					Grouped(true).
					Placeholder("选择食物..."),
				h.Div(
					h.Text("选中的值: "),
					h.Tag("code").Children(h.Text("{{locals.groupedValue}}")).Style("background: #f3f4f6; padding: 2px 8px; border-radius: 4px; font-family: monospace;"),
				).Style("margin-top: 12px; font-size: 14px;"),
			).Init(`{groupedValue: ''}`).VSlot("{ locals }"),
		).Class("demo-section"),

		// 4. 创建新项
		h.Div(
			h.H2("4. 创建新项").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("允许用户创建不在列表中的新项目")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			web.Scope(
				shadcn.Autocomplete().
					Attr("v-model", "locals.createValue").
					ItemsExpr("locals.tags").
					AllowCreate(true).
					CreateText("创建标签").
					OnCreate(`(newValue) => {
						locals.tags.push({value: newValue, text: newValue});
						locals.createValue = newValue;
						console.log('创建新标签:', newValue);
					}`).
					Placeholder("输入或创建标签..."),
				h.Div(
					h.Text("选中的值: "),
					h.Tag("code").Children(h.Text("{{locals.createValue}}")).Style("background: #f3f4f6; padding: 2px 8px; border-radius: 4px; font-family: monospace;"),
				).Style("margin-top: 12px; font-size: 14px;"),
				h.Div(
					h.Text("所有标签: "),
					h.Tag("code").Children(h.Text("{{locals.tags.map(t => t.text).join(', ')}}")).Style("background: #f3f4f6; padding: 2px 8px; border-radius: 4px; font-family: monospace;"),
				).Style("margin-top: 8px; font-size: 14px; color: #666;"),
			).Init(`{
				createValue: '',
				tags: [
					{value: 'vue', text: 'Vue'},
					{value: 'react', text: 'React'},
					{value: 'angular', text: 'Angular'}
				]
			}`).VSlot("{ locals }"),
		).Class("demo-section"),

		// 5. 远程搜索 + 防抖
		h.Div(
			h.H2("5. 远程搜索 + 防抖 (300ms)").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("从模拟 API 异步获取搜索结果,带防抖优化")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			web.Scope(
				shadcn.Autocomplete().
					Attr("v-model", "locals.remoteValue").
					Remote(true).
					RemoteMethod(`async (query) => {
						if (!query) return [];

						// 模拟 API 请求
						console.log('搜索:', query);
						await new Promise(resolve => setTimeout(resolve, 500));

						// 模拟搜索结果
						const allUsers = [
							{value: '1', text: 'Alice Johnson', email: 'alice@example.com'},
							{value: '2', text: 'Bob Smith', email: 'bob@example.com'},
							{value: '3', text: 'Charlie Brown', email: 'charlie@example.com'},
							{value: '4', text: 'David Wilson', email: 'david@example.com'},
							{value: '5', text: 'Eve Davis', email: 'eve@example.com'}
						];

						return allUsers.filter(u =>
							u.text.toLowerCase().includes(query.toLowerCase()) ||
							u.email.toLowerCase().includes(query.toLowerCase())
						);
					}`).
					Debounce(300).
					Placeholder("搜索用户...").
					LoadingText("搜索中...").
					EmptyText("未找到用户"),
				h.Div(
					h.Text("选中的值: "),
					h.Tag("code").Children(h.Text("{{locals.remoteValue}}")).Style("background: #f3f4f6; padding: 2px 8px; border-radius: 4px; font-family: monospace;"),
				).Style("margin-top: 12px; font-size: 14px;"),
				h.Div(
					h.Text("提示: 输入用户名进行搜索,防抖延迟 300ms,模拟 API 延迟 500ms"),
				).Style("margin-top: 8px; font-size: 12px; color: #2563eb;"),
			).Init(`{remoteValue: ''}`).VSlot("{ locals }"),
		).Class("demo-section"),

		// 6. 综合示例：远程搜索 + 分组 + 创建
		h.Div(
			h.H2("6. 综合示例 (远程搜索 + 分组 + 创建新项)").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("所有功能的综合应用")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			web.Scope(
				shadcn.Autocomplete().
					Attr("v-model", "locals.advancedValue").
					Remote(true).
					RemoteMethod(`async (query) => {
						if (!query) return [];

						console.log('综合搜索:', query);
						await new Promise(resolve => setTimeout(resolve, 400));

						const allItems = [
							{value: 'frontend-vue', text: 'Vue.js', group: '前端框架'},
							{value: 'frontend-react', text: 'React', group: '前端框架'},
							{value: 'backend-go', text: 'Go', group: '后端语言'},
							{value: 'backend-python', text: 'Python', group: '后端语言'},
							{value: 'db-postgres', text: 'PostgreSQL', group: '数据库'},
							{value: 'db-redis', text: 'Redis', group: '数据库'}
						];

						return allItems.filter(item =>
							item.text.toLowerCase().includes(query.toLowerCase())
						);
					}`).
					Grouped(true).
					AllowCreate(true).
					CreateText("创建新技术").
					Debounce(300).
					OnCreate(`(newValue) => {
						locals.advancedValue = newValue;
						console.log('创建新技术:', newValue);
					}`).
					Placeholder("搜索或创建技术栈...").
					LoadingText("搜索中..."),
				h.Div(
					h.Text("选中的值: "),
					h.Tag("code").Children(h.Text("{{locals.advancedValue}}")).Style("background: #f3f4f6; padding: 2px 8px; border-radius: 4px; font-family: monospace;"),
				).Style("margin-top: 12px; font-size: 14px;"),
				h.Div(
					h.Text("此示例结合了远程搜索、分组显示和创建新项三大功能"),
				).Style("margin-top: 8px; font-size: 12px; color: #16a34a;"),
			).Init(`{advancedValue: ''}`).VSlot("{ locals }"),
		).Class("demo-section"),

		// 功能特性说明
		h.Div(
			h.H2("功能特性").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.Ul(
				h.Li(h.Text("多选模式 (multiple selection with tags)")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("远程搜索 + 防抖 (configurable debounce time)")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("分组显示 (grouped items with separators)")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("创建新项功能 (allow creating new items)")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("加载状态显示 (loading spinner during async operations)")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("错误消息支持 (error messages display)")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("完全类型安全 (TypeScript interfaces)")).Style("padding: 8px 0;"),
			).Style("list-style: none; padding: 0; margin: 0;"),
		).Class("demo-section").Style("background: #eff6ff;"),
	)
}
