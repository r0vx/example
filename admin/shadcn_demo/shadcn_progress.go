package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnProgressDemo 虚拟模型
type ShadcnProgressDemo struct{}

// configProgress 注册 Progress demo
func configProgress(b *presets.Builder) {
	m := b.Model(&ShadcnProgressDemo{}).
		Label("Progress").
		URIName("shadcn-progress")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnProgressDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("simulateFetch", simulateFetch)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Progress"
		r.Body = shadcnProgressBody(ctx)
		return
	})
}

// shadcnProgressBody Progress 进度条演示
func shadcnProgressBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("Shadcn Progress").Style("margin-bottom: 24px;"),

		// Basic Progress
		h.Div(
			h.H2("Basic Progress"),
			h.P(h.Text("Standard progress bar with different values")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					h.Span("0%").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(0),
				).Class("space-y-1"),
				h.Div(
					h.Span("25%").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(25),
				).Class("space-y-1"),
				h.Div(
					h.Span("50%").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(50),
				).Class("space-y-1"),
				h.Div(
					h.Span("75%").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(75),
				).Class("space-y-1"),
				h.Div(
					h.Span("100%").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(100),
				).Class("space-y-1"),
			).Class("space-y-4 max-w-md"),
		).Class("demo-section"),

		// Indeterminate Progress
		h.Div(
			h.H2("Indeterminate Progress"),
			h.P(h.Text("Loading state with continuous animation")).Class("text-muted-foreground mb-4"),
			h.Div(
				Progress().Indeterminate(true),
			).Class("max-w-md"),
			h.P(h.Text("常用于页面加载、数据获取等场景")).Class("mt-4 text-xs text-muted-foreground"),
		).Class("demo-section"),

		// Progress with Active Control
		h.Div(
			h.H2("Active Control"),
			h.P(h.Text("Control visibility with active prop")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					Progress().Indeterminate(true).Attr(":active", "locals.loading").Height(2),
				).Class("max-w-md mb-4"),
				h.Div(
					Button(h.Text("Toggle Loading")).
						Variant(ButtonVariantOutline).
						Attr("@click", "locals.loading = !locals.loading"),
					h.Span("{{ locals.loading ? 'Loading...' : 'Idle' }}").Class("ml-4 text-sm text-muted-foreground"),
				).Class("flex items-center"),
			).VSlot("{ locals }").Init(`{ loading: true }`),
		).Class("demo-section"),

		// Progress Heights
		h.Div(
			h.H2("Different Heights"),
			h.P(h.Text("Progress bars with different heights")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					h.Span("Height: 2px").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(60).Height(2),
				).Class("space-y-1"),
				h.Div(
					h.Span("Height: 4px (default)").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(60),
				).Class("space-y-1"),
				h.Div(
					h.Span("Height: 8px").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(60).Height(8),
				).Class("space-y-1"),
				h.Div(
					h.Span("Height: 16px").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(60).Height(16),
				).Class("space-y-1"),
			).Class("space-y-4 max-w-md"),
		).Class("demo-section"),

		// Progress Colors
		h.Div(
			h.H2("Colors"),
			h.P(h.Text("Progress bars with different colors")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					h.Span("Primary").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(70).Color("primary"),
				).Class("space-y-1"),
				h.Div(
					h.Span("Success").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(70).Color("success"),
				).Class("space-y-1"),
				h.Div(
					h.Span("Warning").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(70).Color("warning"),
				).Class("space-y-1"),
				h.Div(
					h.Span("Error").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(70).Color("error"),
				).Class("space-y-1"),
				h.Div(
					h.Span("Info").Class("text-sm text-muted-foreground"),
					Progress().ModelValue(70).Color("info"),
				).Class("space-y-1"),
			).Class("space-y-4 max-w-md"),
		).Class("demo-section"),

		// ShowValue - 在进度条中间显示百分比
		h.Div(
			h.H2("Show Value"),
			h.P(h.Text("Display percentage text inside the progress bar")).Class("text-muted-foreground mb-4"),
			h.Div(
				Progress().ModelValue(35).ShowValue(true).Height(20),
				Progress().ModelValue(65).ShowValue(true).Height(20).Color("success"),
				Progress().ModelValue(90).ShowValue(true).Height(20).Color("warning"),
			).Class("space-y-4 max-w-md"),
		).Class("demo-section"),

		// Tooltip - 鼠标悬停显示百分比
		h.Div(
			h.H2("Tooltip"),
			h.P(h.Text("Hover to see percentage in tooltip")).Class("text-muted-foreground mb-4"),
			h.Div(
				Progress().ModelValue(42).Tooltip(true),
				Progress().ModelValue(78).Tooltip(true).Color("info"),
			).Class("space-y-4 max-w-md"),
		).Class("demo-section"),

		// Animated Progress
		h.Div(
			h.H2("Animated Progress"),
			h.P(h.Text("Progress bar with animation")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					Progress().Attr(":model-value", "locals.progress"),
				).Class("max-w-md mb-4"),
				h.Div(
					Button(h.Text("Start")).
						Variant(ButtonVariantOutline).
						Size(ButtonSizeSm).
						Attr("@click", `
							locals.progress = 0;
							const interval = setInterval(() => {
								locals.progress += 10;
								if (locals.progress >= 100) {
									clearInterval(interval);
								}
							}, 300);
						`),
					Button(h.Text("Reset")).
						Variant(ButtonVariantOutline).
						Size(ButtonSizeSm).
						Attr("@click", "locals.progress = 0"),
					h.Span("{{ locals.progress }}%").Class("ml-4 text-sm font-medium"),
				).Class("flex items-center gap-2"),
			).VSlot("{ locals }").Init(`{ progress: 0 }`),
		).Class("demo-section"),

		// Use Case: Page Loading
		h.Div(
			h.H2("Use Case: Page Loading"),
			h.P(h.Text("Fixed position progress bar for page loading")).Class("text-muted-foreground mb-4"),
			web.Scope(
				// 固定在顶部的进度条
				h.Div(
					Progress().
						Indeterminate(true).
						Height(2).
						Attr(":active", "locals.fetching").
						Attr("style", "position: fixed; top: 0; left: 0; right: 0; z-index: 99; border-radius: 0;"),
				),
				h.Div(
					Button(h.Text("Simulate Fetch")).
						On("click", web.POST().EventFunc("simulateFetch").Go()),
					h.Span("{{ locals.fetching ? 'Fetching...' : 'Ready' }}").Class("ml-4 text-sm text-muted-foreground"),
				).Class("flex items-center"),
			).VSlot("{ locals }").Init(`{ fetching: false }`),
			h.P(h.Text("这类似于 Vuetify VProgressLinear 在 admin presets 中的用法")).Class("mt-4 text-xs text-muted-foreground"),
		).Class("demo-section"),
	).Style("max-width: 600px; margin: 0 auto;")
}

// simulateFetch 模拟数据获取
func simulateFetch(ctx *web.EventContext) (er web.EventResponse, err error) {
	er.RunScript = `
		locals.fetching = true;
		setTimeout(() => {
			locals.fetching = false;
		}, 2000);
	`
	return
}
