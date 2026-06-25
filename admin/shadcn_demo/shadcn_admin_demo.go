package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
)

// ShadcnAdminDemoDemo 虚拟模型
type ShadcnAdminDemoDemo struct{}

// configAdminDemo 注册 Admin Demo
func configAdminDemo(b *presets.Builder) {
	m := b.Model(&ShadcnAdminDemoDemo{}).
		Label("Admin Demo").
		URIName("shadcn-admin-demo")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnAdminDemoDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Admin Demo"
		r.Body = shadcnAdminDemoBody(ctx)
		return
	})
}

// adminIcon 返回管理后台图标
func adminIcon(name string) h.HTMLComponent {
	icons := map[string]string{
		"dashboard":    `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>`,
		"users":        `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`,
		"products":     `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>`,
		"orders":       `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>`,
		"analytics":    `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>`,
		"settings":     `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>`,
		"menu":         `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="4" y1="12" x2="20" y2="12"/><line x1="4" y1="6" x2="20" y2="6"/><line x1="4" y1="18" x2="20" y2="18"/></svg>`,
		"search":       `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>`,
		"bell":         `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>`,
		"user":         `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>`,
		"chevron-down": `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>`,
		"trending-up":  `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/></svg>`,
		"dollar":       `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>`,
		"activity":     `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>`,
		"credit-card":  `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>`,
		"forms":        `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><line x1="10" y1="9" x2="8" y2="9"/></svg>`,
		"data":         `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>`,
		"feedback":     `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>`,
		"other":        `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg>`,
	}
	if svg, ok := icons[name]; ok {
		return h.RawHTML(svg)
	}
	return h.Span("")
}

// buildAdminSidebar 构建管理后台侧边栏（带标签页打开功能）
func buildAdminSidebar() h.HTMLComponent {
	// 菜单点击切换标签页
	openTabJS := func(key string) string {
		return `form.activeTab = '` + key + `'`
	}

	return Sidebar(
		// Logo Header
		SidebarHeader(
			h.Div(
				h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-primary"><path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/></svg>`),
				h.Span("Admin Panel").Class("ml-2 font-bold text-xl group-data-[collapsible=icon]:hidden"),
			).Class("flex items-center px-2"),
		).Class("border-b").Style("height: 60px; display: flex; align-items: center;"),

		// Content
		SidebarContent(
			// 主导航
			SidebarGroup(
				SidebarGroupLabel(h.Text("Main")),
				SidebarGroupContent(
					SidebarMenu(
						// Dashboard
						SidebarMenuItem(
							h.Div(
								SidebarMenuButton(adminIcon("dashboard"), h.Span("Dashboard")).Tooltip("Dashboard"),
							).Attr("@click", openTabJS("dashboard")).Class("cursor-pointer"),
						),
						// Users - 带二级菜单
						SidebarMenuItem(
							Collapsible(
								CollapsibleTrigger(
									SidebarMenuButton(
										adminIcon("users"),
										h.Span("Users"),
										h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="ml-auto transition-transform group-data-[state=open]:rotate-90"><polyline points="9 18 15 12 9 6"/></svg>`),
									).Tooltip("Users"),
								).AsChild(true),
								CollapsibleContent(
									SidebarMenuSub(
										SidebarMenuSubItem(
											h.Div(SidebarMenuSubButton(h.Text("User List"))).
												Attr("@click", openTabJS("users")).Class("cursor-pointer"),
										),
										SidebarMenuSubItem(
											h.Div(SidebarMenuSubButton(h.Text("Add User"))).
												Attr("@click", openTabJS("users")).Class("cursor-pointer"),
										),
										SidebarMenuSubItem(
											h.Div(SidebarMenuSubButton(h.Text("Roles"))).
												Attr("@click", openTabJS("users")).Class("cursor-pointer"),
										),
									),
								),
							).DefaultOpen(true).Class("group/collapsible"),
						),
						// Products - 带二级菜单
						SidebarMenuItem(
							Collapsible(
								CollapsibleTrigger(
									SidebarMenuButton(
										adminIcon("products"),
										h.Span("Products"),
										h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="ml-auto transition-transform group-data-[state=open]:rotate-90"><polyline points="9 18 15 12 9 6"/></svg>`),
									).Tooltip("Products"),
								).AsChild(true),
								CollapsibleContent(
									SidebarMenuSub(
										SidebarMenuSubItem(
											h.Div(SidebarMenuSubButton(h.Text("Product List"))).
												Attr("@click", openTabJS("products")).Class("cursor-pointer"),
										),
										SidebarMenuSubItem(
											h.Div(SidebarMenuSubButton(h.Text("Categories"))).
												Attr("@click", openTabJS("products")).Class("cursor-pointer"),
										),
										SidebarMenuSubItem(
											h.Div(SidebarMenuSubButton(h.Text("Inventory"))).
												Attr("@click", openTabJS("products")).Class("cursor-pointer"),
										),
									),
								),
							).Class("group/collapsible"),
						),
						// Orders
						SidebarMenuItem(
							h.Div(
								SidebarMenuButton(adminIcon("orders"), h.Span("Orders")).Tooltip("Orders"),
							).Attr("@click", openTabJS("orders")).Class("cursor-pointer"),
						),
					),
				),
			),
			// 分析
			SidebarGroup(
				SidebarGroupLabel(h.Text("Analytics")),
				SidebarGroupContent(
					SidebarMenu(
						SidebarMenuItem(
							h.Div(
								SidebarMenuButton(adminIcon("analytics"), h.Span("Reports")).Tooltip("Reports"),
							).Attr("@click", openTabJS("reports")).Class("cursor-pointer"),
						),
					),
				),
			),
			// 系统
			SidebarGroup(
				SidebarGroupLabel(h.Text("System")),
				SidebarGroupContent(
					SidebarMenu(
						SidebarMenuItem(
							h.Div(
								SidebarMenuButton(adminIcon("settings"), h.Span("Settings")).Tooltip("Settings"),
							).Attr("@click", openTabJS("settings")).Class("cursor-pointer"),
						),
					),
				),
			),
			// 组件示例
			SidebarGroup(
				SidebarGroupLabel(h.Text("Components")),
				SidebarGroupContent(
					SidebarMenu(
						// Forms
						SidebarMenuItem(
							Collapsible(
								CollapsibleTrigger(
									SidebarMenuButton(
										adminIcon("forms"),
										h.Span("Forms"),
										h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="ml-auto transition-transform group-data-[state=open]:rotate-90"><polyline points="9 18 15 12 9 6"/></svg>`),
									).Tooltip("Forms"),
								).AsChild(true),
								CollapsibleContent(
									SidebarMenuSub(
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Basic Inputs"))).Attr("@click", openTabJS("basic-inputs")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Selections"))).Attr("@click", openTabJS("selections")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Form Field"))).Attr("@click", openTabJS("form-field")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Autocomplete"))).Attr("@click", openTabJS("autocomplete")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Range Picker"))).Attr("@click", openTabJS("range-picker")).Class("cursor-pointer")),
									),
								),
							).Class("group/collapsible"),
						),
						// Data Display
						SidebarMenuItem(
							Collapsible(
								CollapsibleTrigger(
									SidebarMenuButton(
										adminIcon("data"),
										h.Span("Data"),
										h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="ml-auto transition-transform group-data-[state=open]:rotate-90"><polyline points="9 18 15 12 9 6"/></svg>`),
									).Tooltip("Data"),
								).AsChild(true),
								CollapsibleContent(
									SidebarMenuSub(
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Table"))).Attr("@click", openTabJS("table")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Data Table"))).Attr("@click", openTabJS("data-table")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Grid"))).Attr("@click", openTabJS("grid")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("List"))).Attr("@click", openTabJS("list")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Filter"))).Attr("@click", openTabJS("filter")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Tree View"))).Attr("@click", openTabJS("tree-view")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Cascader"))).Attr("@click", openTabJS("cascader")).Class("cursor-pointer")),
									),
								),
							).Class("group/collapsible"),
						),
						// Feedback
						SidebarMenuItem(
							Collapsible(
								CollapsibleTrigger(
									SidebarMenuButton(
										adminIcon("feedback"),
										h.Span("Feedback"),
										h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="ml-auto transition-transform group-data-[state=open]:rotate-90"><polyline points="9 18 15 12 9 6"/></svg>`),
									).Tooltip("Feedback"),
								).AsChild(true),
								CollapsibleContent(
									SidebarMenuSub(
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Dialog"))).Attr("@click", openTabJS("dialog")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Toast"))).Attr("@click", openTabJS("toast")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Progress"))).Attr("@click", openTabJS("progress")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Popover Menu"))).Attr("@click", openTabJS("popover-menu")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Sheet Drawer"))).Attr("@click", openTabJS("sheet-drawer")).Class("cursor-pointer")),
									),
								),
							).Class("group/collapsible"),
						),
						// Other
						SidebarMenuItem(
							Collapsible(
								CollapsibleTrigger(
									SidebarMenuButton(
										adminIcon("other"),
										h.Span("Other"),
										h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="ml-auto transition-transform group-data-[state=open]:rotate-90"><polyline points="9 18 15 12 9 6"/></svg>`),
									).Tooltip("Other"),
								).AsChild(true),
								CollapsibleContent(
									SidebarMenuSub(
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Sidebar Demo"))).Attr("@click", openTabJS("sidebar-demo")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Invoice List"))).Attr("@click", openTabJS("invoice-list")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Lazy Portals"))).Attr("@click", openTabJS("lazy-portals")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Variant Sub Form"))).Attr("@click", openTabJS("variant-sub-form")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("New Components"))).Attr("@click", openTabJS("new-components")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Display Components"))).Attr("@click", openTabJS("display-components")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("InputOTP"))).Attr("@click", openTabJS("input-otp")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("Item"))).Attr("@click", openTabJS("item")).Class("cursor-pointer")),
										SidebarMenuSubItem(h.Div(SidebarMenuSubButton(h.Text("InputGroup"))).Attr("@click", openTabJS("input-group")).Class("cursor-pointer")),
									),
								),
							).Class("group/collapsible"),
						),
					),
				),
			),
		),

		// Rail
		SidebarRail(),
	).Collapsible(SidebarCollapsibleIcon)
}

// buildStatCards 构建统计卡片
func buildStatCards() h.HTMLComponent {
	stats := []struct {
		title    string
		value    string
		change   string
		icon     string
		positive bool
	}{
		{"Total Revenue", "$45,231.89", "+20.1%", "dollar", true},
		{"Active Users", "2,350", "+180.1%", "users", true},
		{"Sales", "+12,234", "+19%", "credit-card", true},
		{"Active Now", "+573", "+201", "activity", true},
	}

	var cards []h.HTMLComponent
	for _, stat := range stats {
		changeColor := "text-green-600"
		if !stat.positive {
			changeColor = "text-red-600"
		}
		cards = append(cards,
			Card(
				CardHeader(
					h.Div(
						h.Div(h.Text(stat.title)).Class("text-sm font-medium text-muted-foreground"),
						h.Div(adminIcon(stat.icon)).Class("text-muted-foreground"),
					).Class("flex justify-between items-center"),
				).Class("p-3 pb-1"),
				CardContent(
					h.Div(h.Text(stat.value)).Class("text-xl font-bold"),
					h.Div(
						h.Span(stat.change).Class(changeColor),
						h.Text(" from last month"),
					).Class("text-xs text-muted-foreground"),
				).Class("p-3 pt-0"),
			).Class("shadow-none border"),
		)
	}

	return h.Div(cards...).Class("mb-4").Style("display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 0.75rem;")
}

// buildUserDataTableWithRowClick 构建带行点击事件的用户数据表格
func buildUserDataTableWithRowClick(users []map[string]any, columns []DataTableColumn) h.HTMLComponent {
	return DataTable().
		Data(users).
		Columns(columns).
		Selectable(true).
		Pagination(true).
		PageSize(20).
		Hover(true).
		SelectedCountLabel("Selected {count} records").
		ClearSelectionLabel("Clear").
		Attr("@row-click", "locals.selectedUser = $event; locals.sheetOpen = true").
		Children(
			// 状态列自定义渲染
			h.Tag("template").Attr("v-slot:cell-status", "{ value }").Children(
				h.Tag("shd-badge").
					Attr(":variant", "value === 'active' ? 'default' : 'secondary'").
					Children(h.RawHTML("{{ value === 'active' ? 'Active' : 'Inactive' }}")),
			),
			// 角色列自定义渲染
			h.Tag("template").Attr("v-slot:cell-role", "{ value }").Children(
				h.Tag("shd-badge").
					Attr("variant", "outline").
					Attr(":class", `{
						'bg-red-100 text-red-800 border-red-200': value === 'Admin',
						'bg-blue-100 text-blue-800 border-blue-200': value === 'Editor',
						'bg-gray-100 text-gray-800 border-gray-200': value === 'Viewer'
					}`).
					Children(h.RawHTML("{{ value }}")),
			),
			// 行菜单
			h.Tag("template").Attr("v-slot:row-menu", "{ row }").Children(
				h.Tag("shd-dropdown-menu-item").
					Attr("@click", "locals.selectedUser = row; locals.sheetOpen = true").
					Children(h.Text("View Details")),
				h.Tag("shd-dropdown-menu-item").
					Attr("@click", "console.log('Edit:', row.name)").
					Children(h.Text("Edit")),
				h.Tag("shd-dropdown-menu-separator"),
				h.Tag("shd-dropdown-menu-item").
					Attr("@click", "console.log('Delete:', row.name)").
					Attr("class", "text-red-600").
					Children(h.Text("Delete")),
			),
		)
}

// buildIframeContent 构建 iframe 加载外部页面内容
func buildIframeContent(title, description, url string) h.HTMLComponent {
	return h.Div(
		h.Div(
			h.H1(title).Class("text-xl font-semibold"),
			h.P(h.Text(description)).Class("text-xs text-muted-foreground"),
		).Class("mb-3"),
		h.Tag("iframe").
			Attr("src", url).
			Attr("frameborder", "0").
			Class("w-full rounded-lg border").
			Style("height: calc(100vh - 140px);"),
	)
}

// buildPageTabsContent 使用 shadcn Tabs 组件实现页面标签（支持侧边栏联动）
func buildPageTabsContent(usersContent h.HTMLComponent) h.HTMLComponent {
	return Tabs(
		// Dashboard 内容
		TabsContent(
			h.Div(
				h.Div(
					h.H1("Dashboard").Class("text-xl font-semibold"),
					h.P(h.Text("Welcome to Admin Panel")).Class("text-xs text-muted-foreground"),
				).Class("mb-3"),
				buildStatCards(),
				h.Div(h.Text("Dashboard analytics and charts will be displayed here...")).Class("text-sm text-muted-foreground"),
			).Class("p-4"),
		).Value("dashboard"),

		// Users 内容
		TabsContent(usersContent).Value("users"),

		// Products 内容
		TabsContent(
			h.Div(
				h.Div(
					h.H1("Products").Class("text-xl font-semibold"),
					h.P(h.Text("Manage your products")).Class("text-xs text-muted-foreground"),
				).Class("mb-3"),
				Card(CardContent(h.Text("Products management content...")).Class("p-3")).Class("shadow-none border"),
			).Class("p-4"),
		).Value("products"),

		// Orders 内容
		TabsContent(
			h.Div(
				h.Div(
					h.H1("Orders").Class("text-xl font-semibold"),
					h.P(h.Text("View and manage orders")).Class("text-xs text-muted-foreground"),
				).Class("mb-3"),
				Card(CardContent(h.Text("Orders management content...")).Class("p-3")).Class("shadow-none border"),
			).Class("p-4"),
		).Value("orders"),

		// Reports 内容
		TabsContent(h.Div(buildIframeContent("Reports", "Analytics and reports", "/examples/shadcn-chart")).Class("p-4")).Value("reports"),

		// Settings 内容
		TabsContent(
			h.Div(
				h.Div(
					h.H1("Settings").Class("text-xl font-semibold"),
					h.P(h.Text("System settings")).Class("text-xs text-muted-foreground"),
				).Class("mb-3"),
				Card(CardContent(h.Text("System settings content...")).Class("p-3")).Class("shadow-none border"),
			).Class("p-4"),
		).Value("settings"),

		// === Components - Forms ===
		TabsContent(h.Div(buildIframeContent("Basic Inputs", "Form input components", "/examples/shadcn-basic-inputs")).Class("p-4")).Value("basic-inputs"),
		TabsContent(h.Div(buildIframeContent("Selections", "Selection components", "/examples/shadcn-selections")).Class("p-4")).Value("selections"),
		TabsContent(h.Div(buildIframeContent("Form Field", "Form field component", "/examples/shadcn-form-field")).Class("p-4")).Value("form-field"),
		TabsContent(h.Div(buildIframeContent("Autocomplete", "Autocomplete component", "/examples/shadcn-autocomplete")).Class("p-4")).Value("autocomplete"),
		TabsContent(h.Div(buildIframeContent("Range Picker", "Date range picker", "/examples/shadcn-range-picker")).Class("p-4")).Value("range-picker"),

		// === Components - Data ===
		TabsContent(h.Div(buildIframeContent("Table", "Table component", "/examples/shadcn-table")).Class("p-4")).Value("table"),
		TabsContent(h.Div(buildIframeContent("Data Table", "Advanced data table", "/examples/shadcn-data-table")).Class("p-4")).Value("data-table"),
		TabsContent(h.Div(buildIframeContent("Grid", "Grid layout component", "/examples/shadcn-grid")).Class("p-4")).Value("grid"),
		TabsContent(h.Div(buildIframeContent("List", "List component", "/examples/shadcn-list")).Class("p-4")).Value("list"),
		TabsContent(h.Div(buildIframeContent("Filter", "Filter component", "/examples/shadcn-filter")).Class("p-4")).Value("filter"),
		TabsContent(h.Div(buildIframeContent("Tree View", "Tree view component", "/examples/shadcn-tree-view")).Class("p-4")).Value("tree-view"),
		TabsContent(h.Div(buildIframeContent("Cascader", "Cascader component", "/examples/shadcn-cascader")).Class("p-4")).Value("cascader"),

		// === Components - Feedback ===
		TabsContent(h.Div(buildIframeContent("Dialog", "Dialog component", "/examples/shadcn-dialog")).Class("p-4")).Value("dialog"),
		TabsContent(h.Div(buildIframeContent("Toast", "Toast notifications", "/examples/shadcn-toast")).Class("p-4")).Value("toast"),
		TabsContent(h.Div(buildIframeContent("Progress", "Progress indicators", "/examples/shadcn-progress")).Class("p-4")).Value("progress"),
		TabsContent(h.Div(buildIframeContent("Popover Menu", "Popover menu component", "/examples/shadcn-popover-menu")).Class("p-4")).Value("popover-menu"),
		TabsContent(h.Div(buildIframeContent("Sheet Drawer", "Sheet drawer component", "/examples/shadcn-sheet-drawer")).Class("p-4")).Value("sheet-drawer"),

		// === Components - Other ===
		TabsContent(h.Div(buildIframeContent("Sidebar Demo", "Sidebar component", "/examples/shadcn-sidebar-demo")).Class("p-4")).Value("sidebar-demo"),
		TabsContent(h.Div(buildIframeContent("Invoice List", "Invoice list example", "/examples/shadcn-invoice-list")).Class("p-4")).Value("invoice-list"),
		TabsContent(h.Div(buildIframeContent("Lazy Portals", "Lazy portal loading", "/examples/shadcn-lazy-portals")).Class("p-4")).Value("lazy-portals"),
		TabsContent(h.Div(buildIframeContent("Variant Sub Form", "Variant sub form", "/examples/shadcn-variant-sub-form")).Class("p-4")).Value("variant-sub-form"),
		TabsContent(h.Div(buildIframeContent("New Components", "New components showcase", "/examples/shadcn-new-components")).Class("p-4")).Value("new-components"),
		TabsContent(h.Div(buildIframeContent("Display Components", "Display components", "/examples/shadcn-display-components")).Class("p-4")).Value("display-components"),

		// === Components - InputOTP ===
		TabsContent(buildInputOTPDemo()).Value("input-otp"),

		// === Components - Item ===
		TabsContent(buildItemDemo()).Value("item"),

		// === Components - InputGroup ===
		TabsContent(buildInputGroupDemo()).Value("input-group"),
	).Attr("v-model", "form.activeTab").Class("h-full overflow-y-auto")
}

// shadcnAdminDemoBody 后台管理页面演示
func shadcnAdminDemoBody(ctx *web.EventContext) h.HTMLComponent {
	// 用户数据 - 25条
	allUsers := []map[string]any{
		{"id": "1", "name": "Zhang San", "email": "zhangsan@example.com", "role": "Admin", "status": "active", "created_at": "2024-01-15"},
		{"id": "2", "name": "Li Si", "email": "lisi@example.com", "role": "Editor", "status": "active", "created_at": "2024-01-16"},
		{"id": "3", "name": "Wang Wu", "email": "wangwu@example.com", "role": "Viewer", "status": "inactive", "created_at": "2024-01-17"},
		{"id": "4", "name": "Zhao Liu", "email": "zhaoliu@example.com", "role": "Editor", "status": "active", "created_at": "2024-01-18"},
		{"id": "5", "name": "Sun Qi", "email": "sunqi@example.com", "role": "Admin", "status": "active", "created_at": "2024-01-19"},
		{"id": "6", "name": "Zhou Ba", "email": "zhouba@example.com", "role": "Viewer", "status": "inactive", "created_at": "2024-01-20"},
		{"id": "7", "name": "Wu Jiu", "email": "wujiu@example.com", "role": "Editor", "status": "active", "created_at": "2024-01-21"},
		{"id": "8", "name": "Zheng Shi", "email": "zhengshi@example.com", "role": "Viewer", "status": "active", "created_at": "2024-01-22"},
		{"id": "9", "name": "Feng Shiyi", "email": "fengshiyi@example.com", "role": "Admin", "status": "inactive", "created_at": "2024-01-23"},
		{"id": "10", "name": "Chen Shier", "email": "chenshier@example.com", "role": "Editor", "status": "active", "created_at": "2024-01-24"},
		{"id": "11", "name": "Liu Wei", "email": "liuwei@example.com", "role": "Viewer", "status": "active", "created_at": "2024-01-25"},
		{"id": "12", "name": "Yang Ming", "email": "yangming@example.com", "role": "Editor", "status": "active", "created_at": "2024-01-26"},
		{"id": "13", "name": "Huang Lei", "email": "huanglei@example.com", "role": "Admin", "status": "active", "created_at": "2024-01-27"},
		{"id": "14", "name": "Xu Fang", "email": "xufang@example.com", "role": "Viewer", "status": "inactive", "created_at": "2024-01-28"},
		{"id": "15", "name": "Ma Chao", "email": "machao@example.com", "role": "Editor", "status": "active", "created_at": "2024-01-29"},
		{"id": "16", "name": "Guo Jing", "email": "guojing@example.com", "role": "Viewer", "status": "active", "created_at": "2024-01-30"},
		{"id": "17", "name": "Lin Feng", "email": "linfeng@example.com", "role": "Admin", "status": "active", "created_at": "2024-02-01"},
		{"id": "18", "name": "He Tao", "email": "hetao@example.com", "role": "Editor", "status": "inactive", "created_at": "2024-02-02"},
		{"id": "19", "name": "Luo Bin", "email": "luobin@example.com", "role": "Viewer", "status": "active", "created_at": "2024-02-03"},
		{"id": "20", "name": "Xie Yun", "email": "xieyun@example.com", "role": "Editor", "status": "active", "created_at": "2024-02-04"},
		{"id": "21", "name": "Tang Jun", "email": "tangjun@example.com", "role": "Admin", "status": "active", "created_at": "2024-02-05"},
		{"id": "22", "name": "Han Mei", "email": "hanmei@example.com", "role": "Viewer", "status": "inactive", "created_at": "2024-02-06"},
		{"id": "23", "name": "Deng Bo", "email": "dengbo@example.com", "role": "Editor", "status": "active", "created_at": "2024-02-07"},
		{"id": "24", "name": "Cao Peng", "email": "caopeng@example.com", "role": "Viewer", "status": "active", "created_at": "2024-02-08"},
		{"id": "25", "name": "Jiang Hua", "email": "jianghua@example.com", "role": "Admin", "status": "active", "created_at": "2024-02-09"},
	}

	// 按状态筛选用户
	var activeUsers, inactiveUsers []map[string]any
	for _, u := range allUsers {
		if u["status"] == "active" {
			activeUsers = append(activeUsers, u)
		} else {
			inactiveUsers = append(inactiveUsers, u)
		}
	}

	columns := []DataTableColumn{
		{Name: "name", Title: "Name", Sortable: true},
		{Name: "email", Title: "Email", Sortable: true},
		{Name: "role", Title: "Role", Sortable: true},
		{Name: "status", Title: "Status", Sortable: true},
		{Name: "created_at", Title: "Created At", Sortable: true},
	}

	// 构建用户管理内容（带右侧详情 Sheet）
	usersContent := web.Scope(
		h.Div(
			// 数据筛选 Tabs
			Tabs(
				// 第一行：TabsList + 操作按钮（贴顶部导航）
				h.Div(
					TabsList(
						TabsTrigger(h.Text("All"), Badge(h.Text("25")).Variant(BadgeVariantSecondary).Class("ml-2")).Value("all"),
						TabsTrigger(h.Text("Active"), Badge(h.Text("19")).Variant(BadgeVariantSecondary).Class("ml-2")).Value("active"),
						TabsTrigger(h.Text("Inactive"), Badge(h.Text("6")).Variant(BadgeVariantSecondary).Class("ml-2")).Value("inactive"),
					),
					h.Div(
						Button(h.Text("Export")).Variant(ButtonVariantOutline).Size(ButtonSizeSm).Class("mr-2"),
						Button(h.Text("Import")).Variant(ButtonVariantOutline).Size(ButtonSizeSm).Class("mr-2"),
						Button(
							h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-1"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>`),
							h.Text("Add User"),
						).Size(ButtonSizeSm),
					).Class("flex items-center"),
				).Class("flex justify-between items-center py-2 px-4 bg-muted border-b"),
				// 第二行：Filter
				h.Div(
					Filter().Items(
						StringFilterItem("name", "Name"),
						SelectFilterItem("role", "Role", []FilterSelectOption{
							{Text: "Admin", Value: "admin"},
							{Text: "Editor", Value: "editor"},
							{Text: "Viewer", Value: "viewer"},
						}),
						DateRangeFilterItem("created_at", "Created At"),
					).On("change", "console.log('Filter:', $event)"),
				).Class("px-4 py-2"),

				TabsContent(h.Div(buildUserDataTableWithRowClick(allUsers, columns)).Class("px-4")).Value("all"),
				TabsContent(h.Div(buildUserDataTableWithRowClick(activeUsers, columns)).Class("px-4")).Value("active"),
				TabsContent(h.Div(buildUserDataTableWithRowClick(inactiveUsers, columns)).Class("px-4")).Value("inactive"),
			).DefaultValue("all").Class("w-full"),

			// 右侧详情 Sheet
			Sheet(
				SheetContent(
					SheetHeader(
						SheetTitle(h.Text("User Details")),
						SheetDescription(h.Text("View and edit user information")),
					),
					h.Div(
						// 用户详情内容
						h.Div(
							h.Div(h.Text("Name")).Class("text-sm font-medium text-muted-foreground"),
							h.Div(h.RawHTML("{{ locals.selectedUser?.name || '-' }}")).Class("text-base"),
						).Class("mb-4"),
						h.Div(
							h.Div(h.Text("Email")).Class("text-sm font-medium text-muted-foreground"),
							h.Div(h.RawHTML("{{ locals.selectedUser?.email || '-' }}")).Class("text-base"),
						).Class("mb-4"),
						h.Div(
							h.Div(h.Text("Role")).Class("text-sm font-medium text-muted-foreground"),
							h.Div(h.RawHTML("{{ locals.selectedUser?.role || '-' }}")).Class("text-base"),
						).Class("mb-4"),
						h.Div(
							h.Div(h.Text("Status")).Class("text-sm font-medium text-muted-foreground"),
							h.Div(h.RawHTML("{{ locals.selectedUser?.status || '-' }}")).Class("text-base"),
						).Class("mb-4"),
						h.Div(
							h.Div(h.Text("Created At")).Class("text-sm font-medium text-muted-foreground"),
							h.Div(h.RawHTML("{{ locals.selectedUser?.created_at || '-' }}")).Class("text-base"),
						).Class("mb-4"),
					).Class("py-4"),
					SheetFooter(
						Button(h.Text("Edit")).Class("mr-2"),
						SheetClose(Button(h.Text("Close")).Variant(ButtonVariantOutline)),
					),
				).Side("right").Class("w-96"),
			).Attr(":open", "locals.sheetOpen").Attr("@update:open", "locals.sheetOpen = $event"),
		),
	).VSlot("{ locals }").Init(`{ sheetOpen: false, selectedUser: null }`)

	return web.Scope(
		h.Div(
			SidebarProvider(
				buildAdminSidebar(),
				SidebarInset(
					// Top Header
					h.Div(
						h.Div(
							SidebarTrigger(adminIcon("menu")).Class("mr-4"),
							Breadcrumb(
								BreadcrumbList(
									BreadcrumbItem(BreadcrumbLink(h.Text("Home")).Href("#")),
									BreadcrumbSeparator(),
									BreadcrumbItem(BreadcrumbPage(h.Text("Admin"))),
								),
							),
						).Class("flex items-center"),
						h.Div(
							// 搜索框（命令面板样式）- 使用 DialogTrigger 打开搜索对话框
							h.Div(
								Dialog(
									DialogTrigger(
										h.Div(
											h.Div(adminIcon("search")).Class("text-muted-foreground"),
											h.Span("Search...").Class("flex-1 text-sm text-muted-foreground"),
											h.Span("⌘K").Class("text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded border"),
										).Class("flex items-center gap-2 px-3 h-9 w-64 bg-muted/50 hover:bg-muted rounded-md cursor-pointer border border-transparent hover:border-border transition-colors"),
									),
									DialogContent(
										DialogHeader(
											DialogTitle(h.Text("Search")),
											DialogDescription(h.Text("Search for pages, users, and more...")),
										),
										h.Div(
											h.Tag("span").Class("absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground flex items-center justify-center").Children(adminIcon("search")),
											Input().Placeholder("Type to search...").Class("pl-10 h-11 text-base focus:border-primary"),
										).Class("relative mt-2"),
										h.Div(
											h.Div(h.Text("Recent")).Class("text-xs font-medium text-muted-foreground mb-2"),
											h.Div(
												h.Div(
													adminIcon("dashboard"),
													h.Span("Dashboard").Class("ml-2"),
												).Class("flex items-center px-2 py-1.5 rounded hover:bg-muted cursor-pointer text-sm"),
												h.Div(
													adminIcon("users"),
													h.Span("User Management").Class("ml-2"),
												).Class("flex items-center px-2 py-1.5 rounded hover:bg-muted cursor-pointer text-sm"),
												h.Div(
													adminIcon("settings"),
													h.Span("Settings").Class("ml-2"),
												).Class("flex items-center px-2 py-1.5 rounded hover:bg-muted cursor-pointer text-sm"),
											),
										).Class("mt-4"),
									).Class(presets.DialogSizeSm),
								),
							).Class("mr-4"),
							Button(adminIcon("bell")).Variant(ButtonVariantGhost).Size(ButtonSizeIcon).Class("mr-2"),
							DropdownMenu(
								DropdownMenuTrigger(
									h.Div(
										Avatar(AvatarImage().Src("https://api.dicebear.com/7.x/avataaars/svg?seed=admin").Alt("Admin")).Class("w-8 h-8"),
										adminIcon("chevron-down"),
									).Class("flex items-center gap-1 cursor-pointer"),
								),
								DropdownMenuContent(
									DropdownMenuLabel(h.Text("My Account")),
									DropdownMenuSeparator(),
									DropdownMenuItem(h.Text("Profile")),
									DropdownMenuItem(h.Text("Settings")),
									DropdownMenuSeparator(),
									DropdownMenuItem(h.Text("Logout")).Class("text-red-600"),
								).Align("end"),
							),
						).Class("flex items-center"),
					).Class("flex justify-between items-center px-4 border-b bg-background").Style("height: 60px; flex-shrink: 0;"),

					// PageTabs 内容区域
					h.Div(
						buildPageTabsContent(usersContent),
					).Style("height: calc(100vh - 60px);").Class("page-tabs-container"),
				),
			).Class("h-screen overflow-hidden"),
		),
	).VSlot("{ form }").FormInit(`{ "activeTab": "dashboard" }`)
}

// buildInputOTPDemo 构建 InputOTP 组件演示
func buildInputOTPDemo() h.HTMLComponent {
	return h.Div(
		h.Div(
			h.H1("InputOTP Component").Class("text-xl font-semibold"),
			h.P(h.Text("One-Time Password input component")).Class("text-xs text-muted-foreground"),
		).Class("mb-6"),

		h.Div(
			Card(
				CardHeader(CardTitle(h.Text("6-Digit OTP"))),
				CardContent(
					h.Div(
						InputOTP().Maxlength(6).Children(
							InputOTPGroup().Children(
								InputOTPSlot(0),
								InputOTPSlot(1),
								InputOTPSlot(2),
							),
							InputOTPSeparator(),
							InputOTPGroup().Children(
								InputOTPSlot(3),
								InputOTPSlot(4),
								InputOTPSlot(5),
							),
						).Attr("v-model", "form.otp1").Class("justify-center"),
						h.P(h.Text("Value: ")).Class("text-sm text-muted-foreground mt-4"),
						h.P(h.Text("{{ form.otp1 }}")).Class("text-sm font-mono"),
					).Class("flex flex-col items-center gap-2"),
				),
			).Class("shadow-none border"),
		).Class("mb-6"),

		h.Div(
			Card(
				CardHeader(CardTitle(h.Text("4-Digit OTP with Custom Styling"))),
				CardContent(
					h.Div(
						InputOTP().Maxlength(4).Children(
							InputOTPGroup().Children(
								InputOTPSlot(0),
								InputOTPSlot(1),
								InputOTPSlot(2),
								InputOTPSlot(3),
							),
						).Attr("v-model", "form.otp2").Class("justify-center gap-4"),
						h.P(h.Text("Value: ")).Class("text-sm text-muted-foreground mt-4"),
						h.P(h.Text("{{ form.otp2 }}")).Class("text-sm font-mono"),
					).Class("flex flex-col items-center gap-2"),
				),
			).Class("shadow-none border"),
		).Class("mb-6"),
	).Class("p-4")
}

// buildItemDemo 构建 Item 组件演示
func buildItemDemo() h.HTMLComponent {
	return h.Div(
		h.Div(
			h.H1("Item Component").Class("text-xl font-semibold"),
			h.P(h.Text("Flexible list item component with variants and slots")).Class("text-xs text-muted-foreground"),
		).Class("mb-6"),

		h.Div(
			Card(
				CardHeader(CardTitle(h.Text("Item with Icon"))),
				CardContent(
					ItemGroup().Children(
						Item().Variant(ItemVariantDefault).Children(
							ItemMedia().Variant(ItemMediaVariantIcon).Children(
								h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 11l3 3L22 4"/><path d="M20.84 4.61a2.5 2.5 0 0 0-3.54 0l-2.12 2.12a2.5 2.5 0 0 0 0 3.54l2.12 2.12a2.5 2.5 0 0 0 3.54 0l2.12-2.12a2.5 2.5 0 0 0 0-3.54z"/></svg>`),
							),
							ItemContent().Children(
								ItemTitle(h.Text("Task 1")),
								ItemDescription(h.Text("Complete the project setup")),
							),
							ItemActions().Children(
								Badge(h.Text("In Progress")).Variant(BadgeVariantSecondary),
							),
						),
						Item().Variant(ItemVariantDefault).Children(
							ItemMedia().Variant(ItemMediaVariantIcon).Children(
								h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>`),
							),
							ItemContent().Children(
								ItemTitle(h.Text("Task 2")),
								ItemDescription(h.Text("Review pull requests")),
							),
							ItemActions().Children(
								Badge(h.Text("Pending")).Variant(BadgeVariantOutline),
							),
						),
						Item().Variant(ItemVariantDefault).Children(
							ItemMedia().Variant(ItemMediaVariantIcon).Children(
								h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>`),
							),
							ItemContent().Children(
								ItemTitle(h.Text("Task 3")),
								ItemDescription(h.Text("Deploy to production")),
							),
							ItemActions().Children(
								Badge(h.Text("Completed")).Variant(BadgeVariantDefault),
							),
						),
					),
				),
			).Class("shadow-none border"),
		).Class("mb-6"),

		h.Div(
			Card(
				CardHeader(CardTitle(h.Text("Item Variants"))),
				CardContent(
					h.Div(
						h.H3("Default Variant").Class("text-sm font-medium mb-2"),
						Item().Variant(ItemVariantDefault).Children(
							ItemContent().Children(
								ItemTitle(h.Text("Default Item")),
								ItemDescription(h.Text("Hover to see effect")),
							),
						).Class("mb-3"),

						h.H3("Outline Variant").Class("text-sm font-medium mb-2 mt-4"),
						Item().Variant(ItemVariantOutline).Children(
							ItemContent().Children(
								ItemTitle(h.Text("Outline Item")),
								ItemDescription(h.Text("With border styling")),
							),
						).Class("mb-3"),

						h.H3("Muted Variant").Class("text-sm font-medium mb-2 mt-4"),
						Item().Variant(ItemVariantMuted).Children(
							ItemContent().Children(
								ItemTitle(h.Text("Muted Item")),
								ItemDescription(h.Text("With muted background")),
							),
						),
					),
				),
			).Class("shadow-none border"),
		),
	).Class("p-4")
}

// buildInputGroupDemo 构建 InputGroup 组件演示
func buildInputGroupDemo() h.HTMLComponent {
	return h.Div(
		h.Div(
			h.H1("InputGroup Component").Class("text-xl font-semibold"),
			h.P(h.Text("Flexible input field grouping component")).Class("text-xs text-muted-foreground"),
		).Class("mb-6"),

		h.Div(
			Card(
				CardHeader(CardTitle(h.Text("Search Input Group"))),
				CardContent(
					h.Div(
						InputGroup().Children(
							InputGroupItem().Class("pl-3").Children(
								h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-muted-foreground"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>`),
							),
							InputGroupItem().Children(
								Input().Placeholder("Search...").Class("border-0 focus-visible:ring-0"),
							),
							InputGroupItem().Class("pr-3").Children(
								Button(h.Text("Search")).Size(ButtonSizeIcon).Variant(ButtonVariantGhost),
							),
						).Class("rounded-lg border border-input shadow-sm"),
					).Class("gap-0"),
				),
			).Class("shadow-none border"),
		).Class("mb-6"),

		h.Div(
			Card(
				CardHeader(CardTitle(h.Text("URL Input Group"))),
				CardContent(
					h.Div(
						InputGroup().Children(
							InputGroupItem().Class("pl-3").Children(
								h.Text("https://"),
							).Class("text-sm text-muted-foreground font-medium"),
							InputGroupItem().Children(
								Input().Placeholder("example.com").Class("border-0 focus-visible:ring-0"),
							),
						).Class("rounded-lg border border-input shadow-sm"),
					),
				),
			).Class("shadow-none border"),
		).Class("mb-6"),

		h.Div(
			Card(
				CardHeader(CardTitle(h.Text("Price Input Group"))),
				CardContent(
					h.Div(
						InputGroup().Children(
							InputGroupItem().Class("pl-3").Children(
								h.Text("$"),
							).Class("text-sm font-medium"),
							InputGroupItem().Children(
								Input().Placeholder("0.00").Class("border-0 focus-visible:ring-0 text-right").Type("number"),
							),
							InputGroupItem().Class("pr-3").Children(
								h.Text("USD"),
							).Class("text-sm text-muted-foreground font-medium"),
						).Class("rounded-lg border border-input shadow-sm"),
					),
				),
			).Class("shadow-none border"),
		).Class("mb-6"),

		h.Div(
			Card(
				CardHeader(CardTitle(h.Text("Form Input Group"))),
				CardContent(
					h.Div(
						h.Div(
							h.P(h.Text("Email Input:")).Class("text-sm font-medium mb-2"),
							InputGroup().Children(
								InputGroupItem().Class("pl-3").Children(
									h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-muted-foreground"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>`),
								),
								InputGroupItem().Children(
									Input().Placeholder("you@example.com").Class("border-0 focus-visible:ring-0"),
								),
							).Class("rounded-lg border border-input shadow-sm"),
						).Class("mb-4"),
						h.P(h.Text("Password Input:")).Class("text-sm font-medium mb-2"),
						InputGroup().Children(
							InputGroupItem().Class("pl-3").Children(
								h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-muted-foreground"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>`),
							),
							InputGroupItem().Children(
								Input().Placeholder("Password").Class("border-0 focus-visible:ring-0").Type("password"),
							),
							InputGroupItem().Class("pr-3").Children(
								Button(h.Text("Show")).Size(ButtonSizeIcon).Variant(ButtonVariantGhost),
							),
						).Class("rounded-lg border border-input shadow-sm"),
					),
				),
			).Class("shadow-none border"),
		),
	).Class("p-4")
}
