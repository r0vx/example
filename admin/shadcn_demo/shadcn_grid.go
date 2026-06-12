package shadcn_demo

import (
	"fmt"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnGridDemo 虚拟模型
type ShadcnGridDemo struct{}

// configGrid 注册 Grid demo
func configGrid(b *presets.Builder) {
	m := b.Model(&ShadcnGridDemo{}).
		Label("Grid").
		URIName("shadcn-grid")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnGridDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Grid Layout"
		r.Body = shadcnGridBody(ctx)
		return
	})
}

// shadcnGridBody Grid 布局演示（使用 Tailwind CSS grid）
func shadcnGridBody(ctx *web.EventContext) h.HTMLComponent {
	// 辅助函数：创建网格行
	row := func(cols int, count int, variant BadgeVariant) h.HTMLComponent {
		var children []h.HTMLComponent
		for range count {
			children = append(children,
				Card(
					CardContent(
						h.Div(
							Badge(h.Text(fmt.Sprintf("%d cols", 12/count))).Variant(variant),
						).Class("flex items-center justify-center h-16"),
					).Class("p-4"),
				),
			)
		}
		return h.Div(children...).Class(fmt.Sprintf("grid grid-cols-%d gap-4 mb-4", count))
	}

	return h.Div(
		h.H1("Shadcn Grid Layout").Style("margin-bottom: 24px;"),

		// Grid 演示
		h.Div(
			h.H2("Responsive Grid"),
			h.P(h.Text("Using Tailwind CSS grid system")).Class("text-muted-foreground mb-4"),

			// 1 column
			row(12, 1, BadgeVariantDefault),
			// 2 columns
			row(6, 2, BadgeVariantSecondary),
			// 3 columns
			row(4, 3, BadgeVariantDefault),
			// 4 columns
			row(3, 4, BadgeVariantSecondary),
			// 6 columns
			row(2, 6, BadgeVariantDefault),
		).Class("demo-section"),

		// Card Grid
		h.Div(
			h.H2("Card Grid"),
			h.P(h.Text("Responsive card layout")).Class("text-muted-foreground mb-4"),
			h.Div(
				// Cards in responsive grid
				Card(
					CardHeader(
						CardTitle(h.Text("Card 1")),
						CardDescription(h.Text("First card in grid")),
					),
					CardContent(
						h.P(h.Text("Content goes here")),
					),
				),
				Card(
					CardHeader(
						CardTitle(h.Text("Card 2")),
						CardDescription(h.Text("Second card in grid")),
					),
					CardContent(
						h.P(h.Text("Content goes here")),
					),
				),
				Card(
					CardHeader(
						CardTitle(h.Text("Card 3")),
						CardDescription(h.Text("Third card in grid")),
					),
					CardContent(
						h.P(h.Text("Content goes here")),
					),
				),
				Card(
					CardHeader(
						CardTitle(h.Text("Card 4")),
						CardDescription(h.Text("Fourth card in grid")),
					),
					CardContent(
						h.P(h.Text("Content goes here")),
					),
				),
			).Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4"),
		).Class("demo-section"),

		// Flex Layout
		h.Div(
			h.H2("Flex Layout"),
			h.P(h.Text("Using flex for dynamic layouts")).Class("text-muted-foreground mb-4"),
			h.Div(
				Badge(h.Text("Flex")),
				Badge(h.Text("Items")).Variant(BadgeVariantSecondary),
				Badge(h.Text("With")).Variant(BadgeVariantOutline),
				Badge(h.Text("Gap")).Variant(BadgeVariantDestructive),
			).Class("flex flex-wrap gap-2"),
		).Class("demo-section"),
	).Style("max-width: 1200px; margin: 0 auto;")
}
