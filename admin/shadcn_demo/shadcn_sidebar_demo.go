package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// menuIcon 返回菜单图标
func menuIcon(name string) h.HTMLComponent {
	icons := map[string]string{
		"dashboard": `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>`,
		"inbox":     `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 16 12 14 15 10 15 8 12 2 12"/><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/></svg>`,
		"calendar":  `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>`,
		"settings":  `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>`,
		"chevron":   `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>`,
		"menu":      `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="4" y1="12" x2="20" y2="12"/><line x1="4" y1="6" x2="20" y2="6"/><line x1="4" y1="18" x2="20" y2="18"/></svg>`,
	}
	if svg, ok := icons[name]; ok {
		return h.RawHTML(svg)
	}
	return h.Span("")
}

// buildDemoSidebar 构建演示侧边栏
func buildDemoSidebar() h.HTMLComponent {
	// 菜单项
	menuItems := []struct {
		icon   string
		label  string
		active bool
	}{
		{"dashboard", "Dashboard", true},
		{"inbox", "Inbox", false},
		{"calendar", "Calendar", false},
		{"settings", "Settings", false},
	}

	var menuComponents []h.HTMLComponent
	for _, item := range menuItems {
		menuComponents = append(menuComponents,
			SidebarMenuItem(
				h.A(
					SidebarMenuButton(
						menuIcon(item.icon),
						h.Span(item.label),
					).IsActive(item.active).Tooltip(item.label),
				).Href("#"),
			),
		)
	}

	return Sidebar(
		// Header
		SidebarHeader(
			h.Div(
				h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="text-primary"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>`),
				h.Span("Shadcn Demo").Class("ml-2 font-semibold text-lg group-data-[collapsible=icon]:hidden"),
			).Class("flex items-center"),
		).Class("border-b"),

		// Content
		SidebarContent(
			SidebarGroup(
				SidebarGroupLabel(h.Text("Navigation")),
				SidebarGroupContent(
					SidebarMenu(menuComponents...),
				),
			),
		),

		// Rail
		SidebarRail(),
	).Collapsible(SidebarCollapsibleIcon)
}

// ShadcnSidebarDemoDemo 虚拟模型
type ShadcnSidebarDemoDemo struct{}

// configSidebarDemo 注册 Sidebar Demo
func configSidebarDemo(b *presets.Builder) {
	m := b.Model(&ShadcnSidebarDemoDemo{}).
		Label("Sidebar Demo").
		URIName("shadcn-sidebar-demo")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnSidebarDemoDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Sidebar Demo"
		r.Body = shadcnSidebarDemoBody(ctx)
		return
	})
}

// shadcnSidebarDemoBody Sidebar 演示
func shadcnSidebarDemoBody(ctx *web.EventContext) h.HTMLComponent {
	return SidebarProvider(
		buildDemoSidebar(),
		SidebarInset(
			// Breadcrumb
			h.Div(
				Breadcrumb(
					BreadcrumbList(
						BreadcrumbItem(BreadcrumbLink(h.Text("Home")).Href("#")),
						BreadcrumbSeparator(),
						BreadcrumbItem(BreadcrumbPage(h.Text("Dashboard"))),
					),
				),
			).Class("px-4 py-2 border-b"),

			// Header
			h.Div(
				h.Div(
					SidebarTrigger(menuIcon("menu")).Class("mr-2"),
					h.H1("Sidebar Demo").Class("text-lg font-semibold"),
				).Class("flex items-center"),
			).Class("flex items-center justify-between p-4 border-b"),

			// Content
			h.Div(
				h.H1("Sidebar Layout Demo").Style("margin-bottom: 24px;"),

				h.Div(
					h.H2("Features"),
					h.Ul(
						h.Li(h.Text("Collapsible sidebar with icon mode")),
						h.Li(h.Text("Responsive design")),
						h.Li(h.Text("Breadcrumb navigation")),
						h.Li(h.Text("Menu with tooltips in collapsed mode")),
						h.Li(h.Text("Active state indication")),
					).Class("list-disc list-inside space-y-2 mt-4"),
				).Class("demo-section"),

				h.Div(
					h.H2("Usage"),
					h.Tag("pre").Text(`SidebarProvider(
    Sidebar(
        SidebarHeader(...),
        SidebarContent(...),
        SidebarRail(),
    ).Collapsible(SidebarCollapsibleIcon),
    SidebarInset(
        // Your page content
    ),
)`).Class("bg-muted p-4 rounded text-sm overflow-x-auto"),
				).Class("demo-section"),

				h.Div(
					h.H2("Instructions"),
					h.P(h.Text("Click the sidebar rail (left edge) or the menu button to toggle the sidebar between full and icon mode.")),
				).Class("demo-section"),

				// DatePicker 演示
				h.Div(
					h.H2("DatePicker Demo"),
					h.P(h.Text("Date selection component combining Popover, Button and Calendar")).Class("text-muted-foreground mb-4"),

					// 基础 DatePicker
					h.Div(
						h.H3("Basic DatePicker").Class("text-base font-medium mb-2"),
						web.Scope(
							h.Div(
								DatePicker().Placeholder("选择日期").Attr("v-model", "form.date1"),
								h.Span("{{ form.date1 || '未选择' }}").Class("ml-4 text-sm text-muted-foreground"),
							).Class("flex items-center"),
						).VSlot("{ form }").FormInit(`{ "date1": "" }`),
					).Class("mb-6"),

					// 带默认值的 DatePicker
					h.Div(
						h.H3("With Default Value").Class("text-base font-medium mb-2"),
						web.Scope(
							h.Div(
								DatePicker().Placeholder("选择生日").Attr("v-model", "form.birthday"),
								h.Span("{{ form.birthday }}").Class("ml-4 text-sm text-muted-foreground"),
							).Class("flex items-center"),
						).VSlot("{ form }").FormInit(`{ "birthday": "2024-01-15" }`),
					).Class("mb-6"),

					// 禁用状态
					h.Div(
						h.H3("Disabled State").Class("text-base font-medium mb-2"),
						DatePicker().Placeholder("禁用状态").Disabled(true),
					).Class("mb-6"),

					// 在表单中使用
					h.Div(
						h.H3("In Form Context").Class("text-base font-medium mb-2"),
						Card(
							CardHeader(
								CardTitle(h.Text("Event Registration")),
								CardDescription(h.Text("Please select the event date")),
							),
							CardContent(
								web.Scope(
									h.Div(
										Label(h.Text("Event Name")).Class("mb-2"),
										Input().Placeholder("Enter event name").Attr("v-model", "form.eventName"),
									).Class("mb-4"),
									h.Div(
										Label(h.Text("Event Date")).Class("mb-2"),
										DatePicker().Placeholder("选择活动日期").Attr("v-model", "form.eventDate"),
									).Class("mb-4"),
									h.Div(
										Label(h.Text("Description")).Class("mb-2"),
										Textarea().Placeholder("Event description...").Rows(3).Attr("v-model", "form.desc"),
									).Class("mb-4"),
									h.Div(
										Button(h.Text("Submit")).Attr("@click", `console.log('Event: ' + form.eventName + ', Date: ' + form.eventDate, form)`),
										Button(h.Text("Reset")).Variant(ButtonVariantOutline).Attr("@click", `form.eventName = ''; form.eventDate = ''; form.desc = ''`).Class("ml-2"),
									),
								).VSlot("{ form }").FormInit(`{ "eventName": "", "eventDate": "", "desc": "" }`),
							),
						).Class("max-w-md"),
					).Class("mb-6"),

					// 代码示例
					h.Div(
						h.H3("Usage").Class("text-base font-medium mb-2"),
						h.Tag("pre").Text(`// Basic
DatePicker().Placeholder("选择日期")

// With v-model
DatePicker().Attr("v-model", "form.date")

// With default value (format: yyyy-MM-dd)
DatePicker().ModelValue("2024-01-15")

// Disabled
DatePicker().Disabled(true)`).Class("bg-muted p-4 rounded text-sm overflow-x-auto"),
					),
				).Class("demo-section"),
			).Class("p-6"),
		),
	).Class("min-h-screen")
}
