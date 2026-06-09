package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnTimelineDemo 虚拟模型
type ShadcnTimelineDemo struct{}

// configTimeline 注册 Timeline demo
func configTimeline(b *presets.Builder) {
	m := b.Model(&ShadcnTimelineDemo{}).
		Label("Timeline").
		URIName("shadcn-timeline")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnTimelineDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Timeline"
		r.Body = shadcnTimelineBody(ctx)
		return
	})
}

// shadcnTimelineBody Timeline 组件演示
func shadcnTimelineBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("Shadcn Timeline").Style("margin-bottom: 24px;"),

		// 基础 Timeline
		h.Div(
			h.H2("Basic Timeline"),
			Timeline(
				TimelineItem(
					TimelineIndicator(
						Avatar(
							AvatarImage().Src("https://cdn.vuetifyjs.com/images/lists/1.jpg").Alt("Ali"),
							AvatarFallback(h.Text("AC")),
						).Class("h-7 w-7"),
					),
					TimelineContent(
						h.Div().Class("flex items-center justify-between gap-2 mb-1").Children(
							h.Span("Ali Connors").Class("text-sm font-semibold"),
							h.Span("Just now").Class("text-xs text-muted-foreground"),
						),
						h.Div(h.Text("Created a new project")).Class("text-xs text-muted-foreground"),
					),
				),
				TimelineItem(
					TimelineIndicator(
						Avatar(
							AvatarImage().Src("https://cdn.vuetifyjs.com/images/lists/2.jpg").Alt("Scott"),
							AvatarFallback(h.Text("SC")),
						).Class("h-7 w-7"),
					),
					TimelineContent(
						h.Div().Class("flex items-center justify-between gap-2 mb-1").Children(
							h.Span("Scott Adams").Class("text-sm font-semibold"),
							h.Span("2 hours ago").Class("text-xs text-muted-foreground"),
						),
						h.Div(h.Text("Edited 3 fields")).Class("text-xs text-muted-foreground"),
					),
				),
				TimelineItem(
					TimelineIndicator(
						Avatar(
							AvatarFallback(h.Text("JD")).Class("text-xs font-medium text-primary bg-primary/10"),
						).Class("h-7 w-7"),
					),
					TimelineContent(
						h.Div().Class("flex items-center justify-between gap-2 mb-1").Children(
							h.Span("Jennifer Davis").Class("text-sm font-semibold"),
							h.Span("Yesterday").Class("text-xs text-muted-foreground"),
						),
						h.Div(h.Text("Added a note")).Class("text-xs text-muted-foreground"),
						h.Div(h.Text("This is a sample note content that can span multiple lines.")).
							Class("mt-1 text-xs p-2 bg-muted/50 rounded border"),
					),
				),
				TimelineItem(
					TimelineIndicator(
						Avatar(
							AvatarImage().Src("https://cdn.vuetifyjs.com/images/lists/3.jpg").Alt("Sandra"),
							AvatarFallback(h.Text("SA")),
						).Class("h-7 w-7"),
					),
					TimelineContent(
						h.Div().Class("flex items-center justify-between gap-2 mb-1").Children(
							h.Span("Sandra Adams").Class("text-sm font-semibold"),
							h.Span("3 days ago").Class("text-xs text-muted-foreground"),
						),
						h.Div(h.Text("Deleted a record")).Class("text-xs text-muted-foreground"),
					),
				),
			),
		).Class("demo-section"),

		// 自定义样式 Timeline
		h.Div(
			h.H2("Custom Styled Timeline"),
			Timeline(
				TimelineItem(
					TimelineIndicator(
						h.Div().Class("w-3 h-3 rounded-full bg-green-500"),
					).Class("w-5 h-5"),
					TimelineContent(
						h.Div().Class("flex items-center gap-2 mb-1").Children(
							h.Span("Deployment Successful").Class("text-sm font-semibold text-green-700"),
							Badge(h.Text("Success")).Variant(BadgeVariantSecondary).Class("text-green-700 bg-green-100"),
						),
						h.Div(h.Text("v2.1.0 deployed to production")).Class("text-xs text-muted-foreground"),
					).Class("bg-green-50/50 border-green-200/50"),
				),
				TimelineItem(
					TimelineIndicator(
						h.Div().Class("w-3 h-3 rounded-full bg-blue-500"),
					).Class("w-5 h-5"),
					TimelineContent(
						h.Div().Class("flex items-center gap-2 mb-1").Children(
							h.Span("Build Started").Class("text-sm font-semibold text-blue-700"),
							Badge(h.Text("In Progress")).Variant(BadgeVariantSecondary).Class("text-blue-700 bg-blue-100"),
						),
						h.Div(h.Text("Running CI/CD pipeline...")).Class("text-xs text-muted-foreground"),
					).Class("bg-blue-50/50 border-blue-200/50"),
				),
				TimelineItem(
					TimelineIndicator(
						h.Div().Class("w-3 h-3 rounded-full bg-red-500"),
					).Class("w-5 h-5"),
					TimelineContent(
						h.Div().Class("flex items-center gap-2 mb-1").Children(
							h.Span("Test Failed").Class("text-sm font-semibold text-red-700"),
							Badge(h.Text("Failed")).Variant(BadgeVariantDestructive),
						),
						h.Div(h.Text("3 test cases failed in integration suite")).Class("text-xs text-muted-foreground"),
					).Class("bg-red-50/50 border-red-200/50"),
				),
			),
		).Class("demo-section"),
	).Style("max-width: 600px; margin: 0 auto;")
}
