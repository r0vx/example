package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// subFormValue 子表单数据
type subFormValue struct {
	Type  string
	Form1 struct {
		Gender string
	}
	Form2 struct {
		Feature1 bool
		Progress int
	}
}

// ShadcnVariantSubFormDemo 虚拟模型
type ShadcnVariantSubFormDemo struct{}

// configVariantSubForm 注册 Variant SubForm demo
func configVariantSubForm(b *presets.Builder) {
	m := b.Model(&ShadcnVariantSubFormDemo{}).
		Label("Variant SubForm").
		URIName("shadcn-variant-subform")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnVariantSubFormDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("subformSwitchForm", subformSwitchForm)
	m.RegisterEventFunc("subformSubmit", subformSubmit)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Variant Sub Form"
		r.Body = shadcnVariantSubFormBody(ctx)
		return
	})
}

// shadcnVariantSubFormBody 变体子表单演示
func shadcnVariantSubFormBody(ctx *web.EventContext) h.HTMLComponent {
	var fv subFormValue
	_ = ctx.UnmarshalForm(&fv)
	if fv.Type == "" {
		fv.Type = "Type1"
	}
	var verr web.ValidationErrors

	return h.Div(
		h.H1("Shadcn Variant Sub Form").Style("margin-bottom: 24px;"),

		prettyFormAsJSON(ctx),

		web.Scope(
			// Type Select
			h.Div(
				h.H2("Select Form Type"),
				Select(
					SelectTrigger(
						SelectValue().Placeholder("Select type"),
					).Class("w-52"),
					SelectContent(
						SelectItem(h.Text("Type 1")).Value("Type1"),
						SelectItem(h.Text("Type 2")).Value("Type2"),
					),
				).Attr("v-model", "form.Type").
					On("update:model-value", web.POST().EventFunc("subformSwitchForm").Go()),
			).Class("demo-section"),

			// Sub Form Portal
			web.Portal(
				h.If(fv.Type == "Type1",
					subForm1(ctx, &fv),
				).Else(
					subForm2(ctx, &fv, &verr),
				),
			).Name("variantSubform"),

			// Submit
			h.Div(
				Button(h.Text("Submit")).On("click", web.POST().EventFunc("subformSubmit").Go()),
			).Class("demo-section"),
		).VSlot("{ locals, form }").FormInit(h.JSONString(fv)),
	).Style("max-width: 600px; margin: 0 auto;")
}

// subForm1 类型1表单
func subForm1(ctx *web.EventContext, fv *subFormValue) h.HTMLComponent {
	return h.Div(
		h.H2("Form Type 1: Gender Selection"),
		h.Div(
			Label(h.Text("Gender")),
			RadioGroup(
				h.Div(
					RadioGroupItem().Value("F").Id("female"),
					Label(h.Text("Female")).For("female").Class("ml-2"),
				).Class("flex items-center"),
				h.Div(
					RadioGroupItem().Value("M").Id("male"),
					Label(h.Text("Male")).For("male").Class("ml-2"),
				).Class("flex items-center"),
			).Attr("v-model", "form.Form1.Gender").Class("flex gap-4 mt-2"),
		),
	).Class("demo-section")
}

// subForm2 类型2表单
func subForm2(ctx *web.EventContext, fv *subFormValue, verr *web.ValidationErrors) h.HTMLComponent {
	return h.Div(
		h.H2("Form Type 2: Feature Settings"),
		h.Div(
			// Switch
			h.Div(
				h.Div(
					Label(h.Text("Feature 1")),
					h.P(h.Text("Enable this feature")).Class("text-sm text-muted-foreground"),
				).Class("flex-1"),
				Switch().Attr("v-model", "form.Form2.Feature1"),
			).Class("flex items-center justify-between"),

			// Slider (using input range as fallback)
			h.Div(
				Label(h.Text("Progress")),
				h.Div(
					h.Input("range").
						Attr("v-model", "form.Form2.Progress").
						Attr("min", "0").
						Attr("max", "100").
						Attr("step", "1").
						Class("w-full"),
					h.Span("").Attr("v-text", "form.Form2.Progress + '%'").Class("ml-2 text-sm"),
				).Class("flex items-center mt-2"),
				func() h.HTMLComponent {
					if errors := verr.GetFieldErrors("Progress"); len(errors) > 0 {
						return h.Div(h.Text(errors[0])).Class("text-destructive text-sm mt-1")
					}
					return nil
				}(),
			).Class("mt-4"),
		),
	).Class("demo-section")
}

// subformSwitchForm 切换表单类型
func subformSwitchForm(ctx *web.EventContext) (r web.EventResponse, err error) {
	var verr web.ValidationErrors

	var fv subFormValue
	ctx.MustUnmarshalForm(&fv)

	form := subForm1(ctx, &fv)
	if fv.Type == "Type2" {
		form = subForm2(ctx, &fv, &verr)
	}

	r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
		Name: "variantSubform",
		Body: form,
	})

	return
}

// subformSubmit 提交表单
func subformSubmit(ctx *web.EventContext) (r web.EventResponse, err error) {
	r.Reload = true
	return
}
