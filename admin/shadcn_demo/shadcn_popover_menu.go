package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// menuFormData 菜单表单数据
type menuFormData struct {
	EnableMessages bool
	EnableHints    bool
}

var globalMenuFavored bool

const menuFavoredIconPortalName = "menuFavoredIcon"

// ShadcnPopoverMenuDemo 虚拟模型
type ShadcnPopoverMenuDemo struct{}

// configPopoverMenu 注册 Popover Menu demo
func configPopoverMenu(b *presets.Builder) {
	m := b.Model(&ShadcnPopoverMenuDemo{}).
		Label("Popover Menu").
		URIName("shadcn-popover-menu")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnPopoverMenuDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("menuSubmit", menuSubmit)
	m.RegisterEventFunc("menuToggleFavored", menuToggleFavored)
	m.RegisterEventFunc("menuItemClick", menuItemClick)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Popover Menu"
		r.Body = shadcnPopoverMenuBody(ctx)
		return
	})
}

// shadcnPopoverMenuBody Popover 菜单演示
func shadcnPopoverMenuBody(ctx *web.EventContext) h.HTMLComponent {
	var fv menuFormData
	_ = ctx.UnmarshalForm(&fv)

	return h.Div(
		h.H1("Shadcn Popover Menu").Style("margin-bottom: 24px;"),

		prettyFormAsJSON(ctx),

		// Popover as Menu
		h.Div(
			h.H2("Popover as Menu"),
			web.Scope(
				Popover(
					PopoverTrigger(
						Button(h.Text("Menu as Popover")),
					),
					PopoverContent(
						Card(
							// User Info
							h.Div(
								Avatar(
									AvatarImage().Src("https://cdn.vuetifyjs.com/images/john.jpg").Alt("John"),
									AvatarFallback(h.Text("JL")),
								),
								h.Div(
									h.Div(h.Text("John Leider")).Class("font-medium"),
									h.Div(h.Text("Founder of Vuetify")).Class("text-sm text-muted-foreground"),
								),
								web.Portal(
									menuFavoredIcon(),
								).Name(menuFavoredIconPortalName),
							).Class("flex items-center gap-3 p-4"),

							Separator(),

							// Switches
							h.Div(
								h.Div(
									h.Div(
										Label(h.Text("Enable messages")).For("messages"),
									).Class("flex-1"),
									Switch().Id("messages").Attr("v-model", "form.EnableMessages"),
								).Class("flex items-center justify-between py-2"),
								h.Div(
									h.Div(
										Label(h.Text("Enable hints")).For("hints"),
									).Class("flex-1"),
									Switch().Id("hints").Attr("v-model", "form.EnableHints"),
								).Class("flex items-center justify-between py-2"),
							).Class("p-4"),

							Separator(),

							// Actions
							h.Div(
								Button(h.Text("Cancel")).Variant(ButtonVariantGhost).
									On("click", "locals.menuOpen = false"),
								Button(h.Text("Save")).On("click", web.POST().EventFunc("menuSubmit").Go()),
							).Class("flex justify-end gap-2 p-4"),
						).Class("w-72"),
					).Class("p-0"),
				).Attr("v-model:open", "locals.menuOpen"),
			).VSlot("{ locals, form }").Init("{ menuOpen: false }").FormInit(h.JSONString(fv)),
		).Class("demo-section"),

		// Dropdown Menu
		h.Div(
			h.H2("Dropdown Menu"),
			DropdownMenu(
				DropdownMenuTrigger(
					Button(h.Text("Open Menu")).Variant(ButtonVariantOutline),
				),
				DropdownMenuContent(
					DropdownMenuLabel(h.Text("My Account")),
					DropdownMenuSeparator(),
					DropdownMenuItem(h.Text("Profile")).On("click", web.POST().EventFunc("menuItemClick").Query("item", "profile").Go()),
					DropdownMenuItem(h.Text("Settings")).On("click", web.POST().EventFunc("menuItemClick").Query("item", "settings").Go()),
					DropdownMenuItem(h.Text("Billing")).On("click", web.POST().EventFunc("menuItemClick").Query("item", "billing").Go()),
					DropdownMenuSeparator(),
					DropdownMenuItem(h.Text("Logout")).Class("text-destructive"),
				).SideOffset(4),
			),
			web.Portal().Name("menuClickResult"),
		).Class("demo-section"),

		// Context Menu (Right Click)
		h.Div(
			h.H2("Context Menu (Right Click)"),
			ContextMenu(
				ContextMenuTrigger(
					h.Div(
						h.Text("Right click here"),
					).Class("flex h-36 w-72 items-center justify-center rounded-md border border-dashed text-sm"),
				),
				ContextMenuContent(
					ContextMenuItem(h.Text("Back")),
					ContextMenuItem(h.Text("Forward")),
					ContextMenuItem(h.Text("Reload")),
					ContextMenuSeparator(),
					ContextMenuItem(h.Text("Save As...")),
					ContextMenuItem(h.Text("Print...")),
					ContextMenuSeparator(),
					ContextMenuLabel(h.Text("More Tools")),
					ContextMenuItem(h.Text("Developer Tools")),
					ContextMenuItem(h.Text("View Source")),
				).Class("w-64"),
			),
		).Class("demo-section"),
	).Style("max-width: 600px; margin: 0 auto;")
}

// menuFavoredIcon 收藏图标
func menuFavoredIcon() h.HTMLComponent {
	variant := ButtonVariantGhost
	text := "♡"
	if globalMenuFavored {
		text = "♥"
	}
	return Button(h.Text(text)).Variant(variant).Size(ButtonSizeSm).
		On("click", web.POST().EventFunc("menuToggleFavored").Go())
}

// menuToggleFavored 切换收藏状态
func menuToggleFavored(ctx *web.EventContext) (er web.EventResponse, err error) {
	globalMenuFavored = !globalMenuFavored
	er.UpdatePortals = append(er.UpdatePortals, &web.PortalUpdate{
		Name: menuFavoredIconPortalName,
		Body: menuFavoredIcon(),
	})
	return
}

// menuSubmit 提交菜单表单
func menuSubmit(ctx *web.EventContext) (er web.EventResponse, err error) {
	er.Reload = true
	er.RunScript = "locals.menuOpen = false"
	return
}

// menuItemClick 菜单项点击
func menuItemClick(ctx *web.EventContext) (er web.EventResponse, err error) {
	item := ctx.R.URL.Query().Get("item")
	er.UpdatePortals = append(er.UpdatePortals, &web.PortalUpdate{
		Name: "menuClickResult",
		Body: h.Div(
			Badge(h.Text("Clicked: "+item)).Variant(BadgeVariantSecondary),
		).Class("mt-4"),
	})
	return
}
