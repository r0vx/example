package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
)

// dialogFormValue 对话框表单数据
type dialogFormValue struct {
	Name  string
	Email string
}

var dialogData = &dialogFormValue{
	Name:  "John Doe",
	Email: "john@example.com",
}

// ShadcnDialogDemo 虚拟模型
type ShadcnDialogDemo struct{}

// configDialog 注册 Dialog demo
func configDialog(b *presets.Builder) {
	m := b.Model(&ShadcnDialogDemo{}).
		Label("Dialog").
		URIName("shadcn-dialog")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnDialogDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("saveDialog", saveDialog)
	m.RegisterEventFunc("deleteAccount", deleteAccount)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Dialog"
		r.Body = shadcnDialogBody(ctx)
		return
	})
}

// shadcnDialogBody 对话框组件演示
func shadcnDialogBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("Shadcn Dialog").Style("margin-bottom: 24px;"),

		prettyFormAsJSON(ctx),

		web.Scope(
			// Basic Dialog
			h.Div(
				h.H2("Basic Dialog"),
				h.Div(
					Dialog(
						DialogTrigger(
							Button(h.Text("Open Dialog")),
						),
						DialogContent(
							DialogHeader(
								DialogTitle(h.Text("Edit Profile")),
								DialogDescription(h.Text("Make changes to your profile here. Click save when you're done.")),
							),
							h.Div(
								h.Div(
									Label(h.Text("Name")).For("name"),
									Input().Id("name").Attr("v-model", "form.Name"),
								).Style("display: grid; gap: 8px;"),
								h.Div(
									Label(h.Text("Email")).For("email"),
									Input().Id("email").Type("email").Attr("v-model", "form.Email"),
								).Style("display: grid; gap: 8px;"),
							).Style("display: grid; gap: 16px; padding: 16px 0;"),
							DialogFooter(
								DialogClose(
									Button(h.Text("Cancel")).Variant(ButtonVariantOutline),
								),
								Button(h.Text("Save changes")).On("click", web.POST().EventFunc("saveDialog").Go()),
							),
						),
					),
				).Class("demo-row"),
			).Class("demo-section"),

			// Alert Dialog
			h.Div(
				h.H2("Alert Dialog"),
				h.Div(
					AlertDialog(
						AlertDialogTrigger(
							Button(h.Text("Delete Account")).Variant(ButtonVariantDestructive),
						),
						AlertDialogContent(
							AlertDialogHeader(
								AlertDialogTitle(h.Text("Are you absolutely sure?")),
								AlertDialogDescription(h.Text("This action cannot be undone. This will permanently delete your account and remove your data from our servers.")),
							),
							AlertDialogFooter(
								AlertDialogCancel(h.Text("Cancel")),
								AlertDialogAction(h.Text("Delete")).On("click", web.POST().EventFunc("deleteAccount").Go()),
							),
						),
					),
				).Class("demo-row"),
			).Class("demo-section"),

			// Sheet
			h.Div(
				h.H2("Sheet"),
				h.Div(
					Sheet(
						SheetTrigger(
							Button(h.Text("Open Left Sheet")).Variant(ButtonVariantOutline),
						),
						SheetContent(
							SheetHeader(
								SheetTitle(h.Text("Edit Profile")),
								SheetDescription(h.Text("Make changes to your profile here.")),
							),
							h.Div(
								h.Div(
									Label(h.Text("Name")).For("sheet-name"),
									Input().Id("sheet-name").Attr("v-model", "form.Name"),
								).Style("display: grid; gap: 8px;"),
								h.Div(
									Label(h.Text("Email")).For("sheet-email"),
									Input().Id("sheet-email").Type("email").Attr("v-model", "form.Email"),
								).Style("display: grid; gap: 8px;"),
							).Style("display: grid; gap: 16px; padding: 16px 0;"),
							SheetFooter(
								SheetClose(
									Button(h.Text("Save changes")),
								),
							),
						).Side(SheetSideLeft),
					),
					Sheet(
						SheetTrigger(
							Button(h.Text("Open Right Sheet")).Variant(ButtonVariantOutline),
						),
						SheetContent(
							SheetHeader(
								SheetTitle(h.Text("Settings")),
								SheetDescription(h.Text("Configure your application settings.")),
							),
							h.Div(
								h.P(h.Text("Settings content goes here.")),
							).Style("padding: 16px 0;"),
						).Side(SheetSideRight),
					),
				).Class("demo-row"),
			).Class("demo-section"),

			// Popover
			h.Div(
				h.H2("Popover"),
				h.Div(
					Popover(
						PopoverTrigger(
							Button(h.Text("Open Popover")).Variant(ButtonVariantOutline),
						),
						PopoverContent(
							h.Div(
								h.H4("Dimensions").Style("font-weight: 500; margin-bottom: 8px;"),
								h.P(h.Text("Set the dimensions for the layer.")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),
								h.Div(
									h.Div(
										Label(h.Text("Width")),
										Input().Placeholder("100%"),
									).Style("display: grid; gap: 4px;"),
									h.Div(
										Label(h.Text("Height")),
										Input().Placeholder("auto"),
									).Style("display: grid; gap: 4px;"),
								).Style("display: grid; grid-template-columns: 1fr 1fr; gap: 16px;"),
							),
						).Class("w-80"),
					),
				).Class("demo-row"),
			).Class("demo-section"),

			// Dropdown Menu
			h.Div(
				h.H2("Dropdown Menu"),
				h.Div(
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
				).Class("demo-row"),
			).Class("demo-section"),

			// Tooltip
			h.Div(
				h.H2("Tooltip"),
				h.Div(
					TooltipProvider(
						Tooltip(
							TooltipTrigger(
								Button(h.Text("Hover me")).Variant(ButtonVariantOutline),
							),
							TooltipContent(
								h.P(h.Text("This is a tooltip with useful information")),
							),
						),
					).DelayDuration(200), // 设置延迟为 200ms，默认是 700ms
				).Class("demo-row"),
			).Class("demo-section"),
		).VSlot("{ locals, form }").FormInit(h.JSONString(dialogData)),
	).Style("max-width: 600px; margin: 0 auto;")
}

// saveDialog 保存对话框数据
func saveDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	dialogData = &dialogFormValue{}
	ctx.MustUnmarshalForm(dialogData)
	r.Reload = true
	return
}

// deleteAccount 删除账户
func deleteAccount(ctx *web.EventContext) (r web.EventResponse, err error) {
	// 模拟删除操作
	dialogData = &dialogFormValue{}
	r.Reload = true
	return
}
