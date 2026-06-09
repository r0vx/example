package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ListItemData 列表项数据
type ListItemData struct {
	Title       string
	Description string
	Avatar      string
	Subtitle    string
	IsHeader    bool
	IsDivider   bool
}

// ShadcnListDemo 虚拟模型
type ShadcnListDemo struct{}

// configList 注册 List demo
func configList(b *presets.Builder) {
	m := b.Model(&ShadcnListDemo{}).
		Label("List").
		URIName("shadcn-list")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnListDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn List"
		r.Body = shadcnListBody(ctx)
		return
	})
}

// shadcnListBody 列表组件演示
func shadcnListBody(ctx *web.EventContext) h.HTMLComponent {
	items := []ListItemData{
		{IsHeader: true, Title: "Today"},
		{
			Title:       "Brunch this weekend?",
			Avatar:      "https://cdn.vuetifyjs.com/images/lists/1.jpg",
			Subtitle:    "Ali Connors",
			Description: "I'll be in your neighborhood doing errands this weekend. Do you want to hang out?",
		},
		{IsDivider: true},
		{
			Title:       "Summer BBQ",
			Avatar:      "https://cdn.vuetifyjs.com/images/lists/2.jpg",
			Subtitle:    "to Alex, Scott, Jennifer",
			Description: "Wish I could come, but I'm out of town this weekend.",
		},
		{IsDivider: true},
		{
			Title:       "Oui oui",
			Avatar:      "https://cdn.vuetifyjs.com/images/lists/3.jpg",
			Subtitle:    "Sandra Adams",
			Description: "Do you have Paris recommendations? Have you ever been?",
		},
		{IsDivider: true},
		{
			Title:       "Birthday gift",
			Avatar:      "https://cdn.vuetifyjs.com/images/lists/4.jpg",
			Subtitle:    "Trevor Hansen",
			Description: "Have any ideas about what we should get Heidi for her birthday?",
		},
		{IsDivider: true},
		{
			Title:       "Recipe to try",
			Avatar:      "https://cdn.vuetifyjs.com/images/lists/5.jpg",
			Subtitle:    "Britta Holt",
			Description: "We should eat this: Grate, Squash, Corn, and tomatillo Tacos.",
		},
	}

	var listItems []h.HTMLComponent
	for _, item := range items {
		if item.IsHeader {
			listItems = append(listItems,
				h.Div(
					h.Text(item.Title),
				).Class("px-4 py-2 text-sm font-semibold text-muted-foreground"),
			)
		} else if item.IsDivider {
			listItems = append(listItems, Separator().Class("my-1"))
		} else {
			listItems = append(listItems,
				h.Div(
					Avatar(
						AvatarImage().Src(item.Avatar).Alt(item.Subtitle),
						AvatarFallback(h.Text(item.Subtitle[:2])),
					).Class("h-10 w-10"),
					h.Div(
						h.Div(
							h.Span(item.Title).Class("font-medium"),
						),
						h.Div(
							h.Span(item.Subtitle).Class("text-primary text-sm"),
							h.Span(" — ").Class("text-muted-foreground"),
							h.Span(item.Description).Class("text-muted-foreground text-sm"),
						).Class("line-clamp-2"),
					).Class("flex-1 min-w-0"),
				).Class("flex items-start gap-3 p-4 hover:bg-muted/50 cursor-pointer"),
			)
		}
	}

	return h.Div(
		h.H1("Shadcn List").Style("margin-bottom: 24px;"),

		// Inbox List
		h.Div(
			h.H2("Inbox List"),
			Card(
				// Header
				h.Div(
					h.Div(
						h.H3("Inbox").Class("text-lg font-semibold text-white"),
					),
					h.Div(
						Button(h.Text("Search")).Variant(ButtonVariantGhost).Size(ButtonSizeSm).Class("text-white"),
					),
				).Class("flex items-center justify-between p-4 bg-cyan-600 rounded-t-lg"),
				// List Items
				h.Div(listItems...).Class("divide-y"),
			).Class("overflow-hidden"),
		).Class("demo-section"),

		// Simple List
		h.Div(
			h.H2("Simple List"),
			Card(
				h.Div(
					h.Div(
						h.Span("Dashboard").Class("flex-1"),
						Badge(h.Text("New")).Variant(BadgeVariantSecondary),
					).Class("flex items-center p-3 hover:bg-muted/50 cursor-pointer"),
					Separator(),
					h.Div(
						h.Span("Settings").Class("flex-1"),
					).Class("flex items-center p-3 hover:bg-muted/50 cursor-pointer"),
					Separator(),
					h.Div(
						h.Span("Profile").Class("flex-1"),
					).Class("flex items-center p-3 hover:bg-muted/50 cursor-pointer"),
					Separator(),
					h.Div(
						h.Span("Logout").Class("flex-1 text-destructive"),
					).Class("flex items-center p-3 hover:bg-muted/50 cursor-pointer"),
				),
			),
		).Class("demo-section"),
	).Style("max-width: 600px; margin: 0 auto;")
}
