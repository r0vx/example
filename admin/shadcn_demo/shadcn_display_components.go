package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
)

// ShadcnDisplayComponentsDemo 虚拟模型
type ShadcnDisplayComponentsDemo struct{}

// configDisplayComponents 注册 Display Components demo
func configDisplayComponents(b *presets.Builder) {
	m := b.Model(&ShadcnDisplayComponentsDemo{}).
		Label("Display Components").
		URIName("shadcn-display-components")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnDisplayComponentsDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Display Components"
		r.Body = shadcnDisplayComponentsBody(ctx)
		return
	})
}

// shadcnDisplayComponentsBody 展示类组件演示
func shadcnDisplayComponentsBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		// 页面标题
		h.Div(
			h.H1("Display Components").Style("font-size: 28px; font-weight: bold; margin-bottom: 8px;"),
			h.P(h.Text("展示类组件:Avatar、Card、Separator、Accordion、AlertDialog、DropdownMenu、Popover、Tooltip")).Style("color: #666; margin-bottom: 24px;"),
		),

		// Avatar
		h.Div(
			h.H2("Avatar").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("用户头像组件,支持图片和文字后备")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			h.Div(
				Avatar(
					AvatarImage().Src("https://github.com/shadcn.png").Alt("@shadcn"),
					AvatarFallback(h.Text("CN")),
				),
				Avatar(
					AvatarFallback(h.Text("JD")),
				),
				Avatar(
					AvatarFallback(h.Text("AB")),
				).Class("h-12 w-12"),
			).Class("demo-row"),
		).Class("demo-section"),

		// Card
		h.Div(
			h.H2("Card").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("卡片容器组件,包含 Header、Content、Footer 区域")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			Card(
				CardHeader(
					CardTitle(h.Text("Card Title")),
					CardDescription(h.Text("Card description goes here.")),
				),
				CardContent(
					h.P(h.Text("This is the main content of the card. You can put any content here.")),
				),
				CardFooter(
					Button(h.Text("Cancel")).Variant(ButtonVariantOutline),
					Button(h.Text("Save")),
				).Class("flex justify-end gap-2"),
			).Class("w-80"),
		).Class("demo-section"),

		// Separator
		h.Div(
			h.H2("Separator").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("分隔线组件,用于分隔内容区域")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			h.Div(
				h.Div(h.Text("Section 1")).Style("padding: 8px 0;"),
				Separator(),
				h.Div(h.Text("Section 2")).Style("padding: 8px 0;"),
				Separator(),
				h.Div(h.Text("Section 3")).Style("padding: 8px 0;"),
			),
		).Class("demo-section"),

		// Accordion
		h.Div(
			h.H2("Accordion").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("手风琴组件,可折叠展开的内容区域")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			// 单选模式
			h.Div(
				h.H3("Single 模式").Style("font-size: 14px; font-weight: 500; margin-bottom: 8px;"),
				Accordion(
					AccordionItem(
						AccordionTrigger(h.Text("Is it accessible?")),
						AccordionContent(h.Text("Yes. It adheres to the WAI-ARIA design pattern.")),
					).Value("item-1"),
					AccordionItem(
						AccordionTrigger(h.Text("Is it styled?")),
						AccordionContent(h.Text("Yes. It comes with default styles that match the other components.")),
					).Value("item-2"),
					AccordionItem(
						AccordionTrigger(h.Text("Is it animated?")),
						AccordionContent(h.Text("Yes. It's animated by default, but you can disable it if you prefer.")),
					).Value("item-3"),
				).Type("single").Collapsible(true).DefaultValue("item-1").Class("w-full"),
			).Style("margin-bottom: 24px;"),
			// 多选模式
			h.Div(
				h.H3("Multiple 模式").Style("font-size: 14px; font-weight: 500; margin-bottom: 8px;"),
				Accordion(
					AccordionItem(
						AccordionTrigger(h.Text("功能特性")),
						AccordionContent(h.Text("支持单选和多选模式，可设置默认展开项，支持禁用状态。")),
					).Value("feature-1"),
					AccordionItem(
						AccordionTrigger(h.Text("使用场景")),
						AccordionContent(h.Text("适用于 FAQ、设置面板、分类内容展示等场景。")),
					).Value("feature-2"),
				).Type("multiple").Class("w-full"),
			),
		).Class("demo-section"),

		// AlertDialog
		h.Div(
			h.H2("AlertDialog").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("警告对话框,用于确认危险操作")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			AlertDialog(
				AlertDialogTrigger(
					Button(h.Text("Delete Account")).Variant(ButtonVariantDestructive),
				),
				AlertDialogContent(
					AlertDialogHeader(
						AlertDialogTitle(h.Text("Are you absolutely sure?")),
						AlertDialogDescription(h.Text("This action cannot be undone. This will permanently delete your account.")),
					),
					AlertDialogFooter(
						AlertDialogCancel(h.Text("Cancel")),
						AlertDialogAction(h.Text("Delete")),
					),
				),
			),
		).Class("demo-section"),

		// DropdownMenu
		h.Div(
			h.H2("DropdownMenu").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("下拉菜单组件,支持分组和分隔线")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			DropdownMenu(
				DropdownMenuTrigger(
					Button(h.Text("Open Menu")).Variant(ButtonVariantOutline),
				),
				DropdownMenuContent(
					DropdownMenuLabel(h.Text("My Account")),
					DropdownMenuSeparator(),
					DropdownMenuItem(h.Text("Profile")),
					DropdownMenuItem(h.Text("Settings")),
					DropdownMenuItem(h.Text("Billing")),
					DropdownMenuSeparator(),
					DropdownMenuItem(h.Text("Logout")),
				).SideOffset(4),
			),
		).Class("demo-section"),

		// Popover
		h.Div(
			h.H2("Popover").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("弹出框组件,用于显示额外信息")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			Popover(
				PopoverTrigger(
					Button(h.Text("Open Popover")).Variant(ButtonVariantOutline),
				),
				PopoverContent(
					h.Div(
						h.H4("Dimensions").Style("font-weight: 500; margin-bottom: 8px;"),
						h.P(h.Text("Set the dimensions for the layer.")).Style("color: #666; font-size: 14px;"),
					),
				).Class("w-80"),
			),
		).Class("demo-section"),

		// Tooltip
		h.Div(
			h.H2("Tooltip").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("工具提示组件,鼠标悬停显示")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			TooltipProvider(
				Tooltip(
					TooltipTrigger(
						Button(h.Text("Hover me")).Variant(ButtonVariantOutline),
					),
					TooltipContent(
						h.P(h.Text("This is a tooltip")),
					),
				),
			),
		).Class("demo-section"),

		// DisplayField
		h.Div(
			h.H2("DisplayField").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("只读展示字段,类似 disabled Input 但可交互（点击选中、复制）")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			h.Div(
				// 基础用法
				DisplayField().Value("只读内容").Label("基础用法"),
			).Class("mb-4"),
			h.Div(
				// 带复制按钮
				DisplayField().Value("abc-123-def-456").Label("带复制").Copy(true),
			).Class("mb-4"),
			h.Div(
				// 无 label
				DisplayField().Value("无标签的展示值"),
			).Class("mb-4"),
			h.Div(
				// 代码模式
				DisplayField().Value("{\n  \"name\": \"r0vx\",\n  \"version\": \"1.0.0\",\n  \"description\": \"Admin framework\"\n}").Label("JSON 内容").Code(true).Copy(true),
			).Class("mb-4"),
		).Class("demo-section"),

		// CopyButton
		h.Div(
			h.H2("CopyButton").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("独立复制按钮:点击复制指定文本,2s 内图标变 ✓。内部复用 Button,继承 Variant/Size;非 https 自动回落 execCommand")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
			h.Div(
				CopyButton("https://r0vx.dev/copy-button-test"),               // 纯图标 ghost
				CopyButton("https://r0vx.dev/copy-button-test").Label("复制链接"), // 带文字
				CopyButton("abc-123-def-456").Label("Outline").Variant(ButtonVariantOutline),
				CopyButton("abc-123-def-456").Label("Secondary").Variant(ButtonVariantSecondary),
				CopyButton("small").Variant(ButtonVariantOutline).Size(ButtonSizeIconSm),
				CopyButton("https://r0vx.dev/with-tooltip").Tooltip("复制链接"), // 带悬停提示
				CopyButton("disabled").Disabled(true),
			).Class("demo-row"),
			// ButtonGroup 内分段：CopyButton 渲染裸 <button>，自动吃 [&>button] 分段样式
			h.Div(
				ButtonGroup(
					Button(h.Text("https://r0vx.dev/share")).Variant(ButtonVariantOutline),
					CopyButton("https://r0vx.dev/share").Variant(ButtonVariantOutline),
				),
				ButtonGroup(
					Button(h.Text("编辑")).Variant(ButtonVariantOutline),
					Button(h.Text("分享")).Variant(ButtonVariantOutline),
					CopyButton("token-xyz-789").Variant(ButtonVariantOutline).Tooltip("复制 Token"),
				),
			).Class("demo-row"),
		).Class("demo-section"),

		// 功能特性说明
		h.Div(
			h.H2("组件特性").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.Ul(
				h.Li(h.Text("Avatar - 支持图片和文字后备")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("Card - 完整的 Header/Content/Footer 结构")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("Separator - 水平分隔线")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("Accordion - 可折叠展开的手风琴组件")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("AlertDialog - 危险操作确认对话框")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("DropdownMenu - 支持分组和分隔的下拉菜单")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("Popover - 弹出式信息展示")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("Tooltip - 悬停提示信息")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("DisplayField - 只读展示字段,点击选中,支持复制")).Style("padding: 8px 0;"),
			).Style("list-style: none; padding: 0; margin: 0;"),
		).Class("demo-section").Style("background: #eff6ff;"),
	)
}
