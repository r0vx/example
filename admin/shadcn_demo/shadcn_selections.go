package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// selectFormValue 选择组件表单数据
type selectFormValue struct {
	Country  string
	Language string
	Tags     []string
}

var selectData = &selectFormValue{
	Country:  "cn",
	Language: "go",
	Tags:     []string{"frontend", "backend"},
}

// ShadcnSelectionsDemo 虚拟模型
type ShadcnSelectionsDemo struct{}

// configSelections 注册 Selections demo
func configSelections(b *presets.Builder) {
	m := b.Model(&ShadcnSelectionsDemo{}).
		Label("Selections").
		URIName("shadcn-selections")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnSelectionsDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("submitSelection", submitSelection)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Selections"
		r.Body = shadcnSelectionsBody(ctx)
		return
	})
}

// shadcnSelectionsBody 选择组件演示
func shadcnSelectionsBody(ctx *web.EventContext) h.HTMLComponent {
	var verr web.ValidationErrors
	if ve, ok := ctx.Flash.(web.ValidationErrors); ok {
		verr = ve
	}

	return h.Div(
		h.H1("Shadcn Selections").Style("margin-bottom: 24px;"),

		prettyFormAsJSON(ctx),

		web.Scope(
			// Select
			h.Div(
				h.H2("Select"),
				h.Div(
					h.Div(
						Label(h.Text("Country")).Class("mb-2"),
						Select(
							SelectTrigger(
								SelectValue().Placeholder("Select a country"),
							).Class("w-52"),
							SelectContent(
								SelectItem(h.Text("China")).Value("cn"),
								SelectItem(h.Text("United States")).Value("us"),
								SelectItem(h.Text("Japan")).Value("jp"),
								SelectItem(h.Text("United Kingdom")).Value("uk"),
							),
						).Attr("v-model", "form.Country").
							ErrorMessages(verr.GetFieldErrors("Country")...),
					).Class("mb-4"),
					h.Div(
						Label(h.Text("Programming Language")).Class("mb-2"),
						Select(
							SelectTrigger(
								SelectValue().Placeholder("Select a language"),
							).Class("w-52"),
							SelectContent(
								SelectGroup(
									SelectLabel(h.Text("Backend")),
									SelectItem(h.Text("Go")).Value("go"),
									SelectItem(h.Text("Python")).Value("python"),
									SelectItem(h.Text("Java")).Value("java"),
								),
								SelectSeparator(),
								SelectGroup(
									SelectLabel(h.Text("Frontend")),
									SelectItem(h.Text("JavaScript")).Value("js"),
									SelectItem(h.Text("TypeScript")).Value("ts"),
								),
							),
						).Attr("v-model", "form.Language").
							ErrorMessages(verr.GetFieldErrors("Language")...),
					).Class("mb-4"),
				).Class("demo-row"),
			).Class("demo-section"),

			// Badge (as tags)
			h.Div(
				h.H2("Tags with Badge"),
				h.Div(
					Badge(h.Text("Frontend")).Variant(BadgeVariantSecondary).Class("mr-2"),
					Badge(h.Text("Backend")).Variant(BadgeVariantSecondary).Class("mr-2"),
					Badge(h.Text("DevOps")).Variant(BadgeVariantOutline).Class("mr-2"),
					Badge(h.Text("Deprecated")).Variant(BadgeVariantDestructive),
				).Class("demo-row"),
			).Class("demo-section"),

			// Tabs Selection
			h.Div(
				h.H2("Tabs"),
				Tabs(
					TabsList(
						TabsTrigger(h.Text("Overview")).Value("overview"),
						TabsTrigger(h.Text("Analytics")).Value("analytics"),
						TabsTrigger(h.Text("Reports")).Value("reports"),
						TabsTrigger(h.Text("Notifications")).Value("notifications"),
					),
					TabsContent(
						h.P(h.Text("Overview content here. Make changes to your account.")),
					).Value("overview"),
					TabsContent(
						h.P(h.Text("Analytics content here. View your analytics data.")),
					).Value("analytics"),
					TabsContent(
						h.P(h.Text("Reports content here. Generate and view reports.")),
					).Value("reports"),
					TabsContent(
						h.P(h.Text("Notifications content here. Manage notification settings.")),
					).Value("notifications"),
				).DefaultValue("overview"),
			).Class("demo-section"),

			// Submit Button
			h.Div(
				Button(h.Text("Submit")).On("click", web.POST().EventFunc("submitSelection").Go()),
			).Class("demo-section"),
		).VSlot("{ locals, form }").FormInit(h.JSONString(selectData)),
	).Style("max-width: 600px; margin: 0 auto;")
}

// submitSelection 提交选择
func submitSelection(ctx *web.EventContext) (r web.EventResponse, err error) {
	selectData = &selectFormValue{}
	ctx.MustUnmarshalForm(selectData)

	verr := web.ValidationErrors{}
	if selectData.Country == "" {
		verr.FieldError("Country", "Please select a country")
	}
	if selectData.Language == "" {
		verr.FieldError("Language", "Please select a language")
	}

	ctx.Flash = verr
	r.Reload = true

	return
}
