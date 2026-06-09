package shadcn_demo

import (
	"fmt"
	"time"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// portalState 演示 portal 状态
type portalState struct {
	Company string
	Error   string
}

var portalListItems = []string{"Apple", "Microsoft", "Google"}

// ShadcnLazyPortalsDemo 虚拟模型
type ShadcnLazyPortalsDemo struct{}

// configLazyPortals 注册 Lazy Portals demo
func configLazyPortals(b *presets.Builder) {
	m := b.Model(&ShadcnLazyPortalsDemo{}).
		Label("Lazy Portals").
		URIName("shadcn-lazy-portals")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnLazyPortalsDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("portalMenuItems", portalMenuItems)
	m.RegisterEventFunc("portalAddItemForm", portalAddItemForm)
	m.RegisterEventFunc("portalAddItem", portalAddItem)
	m.RegisterEventFunc("portalContent", portalContent)
	m.RegisterEventFunc("reloadPortalAB", reloadPortalAB)
	m.RegisterEventFunc("updatePortalCD", updatePortalCD)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Lazy Portals and Reload Demo"
		r.Body = shadcnLazyPortalsBody(ctx)
		return
	})
}

// shadcnLazyPortalsBody 演示 Lazy Portal 和 Reload 功能
func shadcnLazyPortalsBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("Lazy Portals and Reload Demo").Style("margin-bottom: 24px;"),

		// Dialog with lazy loaded content
		h.Div(
			h.H2("Dialog with Lazy Portal"),
			Dialog(
				DialogTrigger(
					Button(h.Text("Select Company")),
				),
				DialogContent(
					DialogHeader(
						DialogTitle(h.Text("Select a Company")),
						DialogDescription(h.Text("Choose from the list or create new")),
					),
					web.Portal().Loader(web.POST().EventFunc("portalMenuItems")).Name("menuContent"),
					DialogFooter(
						DialogClose(
							Button(h.Text("Close")).Variant(ButtonVariantOutline),
						),
					),
				),
			),
		).Class("demo-section"),

		// Portal A and B
		h.Div(
			h.H2("Portal A"),
			h.Div(
				web.Portal().Loader(web.POST().EventFunc("portalContent")).Name("portalA"),
			).Style("padding: 16px; border: 2px solid #3b82f6; border-radius: 8px; margin-bottom: 16px;"),

			h.H2("Portal B"),
			h.Div(
				web.Portal().Loader(web.POST().EventFunc("portalContent")).Name("portalB"),
			).Style("padding: 16px; border: 2px solid #ef4444; border-radius: 8px; margin-bottom: 16px;"),

			Button(h.Text("Reload Portal A and B")).
				Variant(ButtonVariantSecondary).
				On("click", web.POST().EventFunc("reloadPortalAB").Go()),
		).Class("demo-section"),

		// Portal C and D (UpdatePortals)
		h.Div(
			h.H2("Portal C"),
			h.Div(
				web.Portal().Name("portalC"),
			).Style("padding: 16px; border: 2px solid #3b82f6; border-radius: 8px; margin-bottom: 16px;"),

			h.H2("Portal D"),
			h.Div(
				web.Portal().Name("portalD"),
			).Style("padding: 16px; border: 2px solid #ef4444; border-radius: 8px; margin-bottom: 16px;"),

			Button(h.Text("Update Portal C and D")).
				On("click", web.POST().EventFunc("updatePortalCD").Go()),
		).Class("demo-section"),
	).Style("max-width: 800px; margin: 0 auto;")
}

// portalMenuItems 加载菜单项列表
func portalMenuItems(ctx *web.EventContext) (r web.EventResponse, err error) {
	s := &portalState{}

	var items []h.HTMLComponent
	for _, item := range portalListItems {
		items = append(items,
			h.Div(
				h.Text(item),
			).Class("p-2 hover:bg-muted rounded cursor-pointer"),
		)
	}

	items = append(items, Separator().Class("my-2"))

	// Add new item dialog
	items = append(items,
		Dialog(
			DialogTrigger(
				Button(h.Text("Create New")).Variant(ButtonVariantOutline).Class("w-full"),
			),
			DialogContent(
				web.Scope(
					web.Portal().Loader(web.POST().EventFunc("portalAddItemForm")).Name("addItemForm").Visible("true"),
				).VSlot("{ locals, form }").FormInit(s),
			).Class("sm:max-w-md"),
		),
	)

	r.Body = h.Div(items...).Class("space-y-1")
	return
}

// portalAddItemForm 添加项目表单
func portalAddItemForm(ctx *web.EventContext) (r web.EventResponse, err error) {
	s := &portalState{}
	ctx.MustUnmarshalForm(s)

	input := Input().
		Placeholder("Enter company name").
		Attr("v-model", "form.Company")

	if len(s.Error) > 0 {
		input = input.ErrorMessages(s.Error)
	}

	r.Body = h.Div(
		DialogHeader(
			DialogTitle(h.Text("Add New Company")),
			DialogDescription(h.Text("Enter the company name below")),
		),
		h.Div(
			Label(h.Text("Company Name")).Class("mb-2"),
			input,
		).Class("py-4"),
		DialogFooter(
			Button(h.Text("Create")).On("click", web.POST().EventFunc("portalAddItem").Go()),
		),
	)
	return
}

// portalAddItem 添加新项目
func portalAddItem(ctx *web.EventContext) (r web.EventResponse, err error) {
	s := &portalState{}
	ctx.MustUnmarshalForm(s)

	if len(s.Company) < 5 {
		r.RunScript = "form.Error = 'Company name must be at least 5 characters'"
		r.ReloadPortals = []string{"addItemForm"}
		return
	}

	portalListItems = append(portalListItems, s.Company)
	s.Company = ""
	s.Error = ""
	r.ReloadPortals = []string{"menuContent"}
	return
}

// portalContent 返回时间戳内容
func portalContent(ctx *web.EventContext) (r web.EventResponse, err error) {
	r.Body = h.Text(fmt.Sprintf("Loaded at: %d", time.Now().UnixNano()))
	return
}

// reloadPortalAB 重新加载 Portal A 和 B
func reloadPortalAB(ctx *web.EventContext) (r web.EventResponse, err error) {
	r.ReloadPortals = []string{"portalA", "portalB"}
	return
}

// updatePortalCD 更新 Portal C 和 D
func updatePortalCD(ctx *web.EventContext) (r web.EventResponse, err error) {
	r.UpdatePortals = append(r.UpdatePortals,
		&web.PortalUpdate{
			Name: "portalC",
			Body: h.Text(fmt.Sprintf("Updated at: %d", time.Now().UnixNano())),
		},
		&web.PortalUpdate{
			Name: "portalD",
			Body: h.Text(fmt.Sprintf("Updated at: %d", time.Now().UnixNano())),
		},
	)
	return
}
