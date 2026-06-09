package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// myFormValue 表单数据结构
type myFormValue struct {
	Username      string
	Email         string
	Password      string
	TextareaValue string
	Gender        string
	Agreed        bool
	Feature1      bool
}

var formData = &myFormValue{
	Username:      "admin",
	Email:         "admin@example.com",
	TextareaValue: "This is textarea value",
	Gender:        "male",
	Agreed:        false,
	Feature1:      true,
}

// ShadcnBasicInputsDemo 虚拟模型
type ShadcnBasicInputsDemo struct{}

// configBasicInputs 注册 Basic Inputs demo
func configBasicInputs(b *presets.Builder) {
	m := b.Model(&ShadcnBasicInputsDemo{}).
		Label("Basic Inputs").
		URIName("shadcn-basic-inputs")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnBasicInputsDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("update", update)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Basic Inputs"
		r.Body = shadcnBasicInputsBody(ctx)
		return
	})
}

// shadcnBasicInputsBody 基础输入组件演示
func shadcnBasicInputsBody(ctx *web.EventContext) h.HTMLComponent {
	var verr web.ValidationErrors
	if ve, ok := ctx.Flash.(web.ValidationErrors); ok {
		verr = ve
	}

	return h.Div(
		h.H1("Shadcn Basic Inputs").Style("margin-bottom: 24px;"),

		prettyFormAsJSON(ctx),

		web.Scope(
			// Input - 使用内置 Label
			h.Div(
				h.H2("Input (内置 Label)"),
				h.Div(
					h.Div(
						Input().
							Label("用户名").
							Required(true).
							Tips("用户名至少3个字符").
							Placeholder("请输入用户名").
							Attr("v-model", "form.Username").
							ErrorMessages(verr.GetFieldErrors("Username")...),
					).Class("mb-4"),
					h.Div(
						Input().
							Label("邮箱").
							Required(true).
							Type("email").
							Placeholder("请输入邮箱").
							Copy(true).
							Attr("v-model", "form.Email").
							ErrorMessages(verr.GetFieldErrors("Email")...),
					).Class("mb-4"),
					h.Div(
						Input().
							Label("密码").
							Type("password").
							Tips("密码是可选的").
							Placeholder("请输入密码").
							Attr("v-model", "form.Password").
							ErrorMessages(verr.GetFieldErrors("Password")...),
					).Class("mb-4"),
				).Class("demo-row"),
			).Class("demo-section"),

			// Textarea - 使用内置 Label
			h.Div(
				h.H2("Textarea (内置 Label)"),
				h.Div(
					Textarea().
						Label("留言").
						Tips("最多100字符").
						Placeholder("请输入留言...").
						Rows(3).
						Attr("v-model", "form.TextareaValue").
						ErrorMessages(verr.GetFieldErrors("TextareaValue")...),
				).Class("demo-row"),
			).Class("demo-section"),

			// RadioGroup
			h.Div(
				h.H2("RadioGroup"),
				h.Div(
					RadioGroup(
						h.Div(
							RadioGroupItem().Value("male").Id("male"),
							Label(h.Text("Male")).For("male").Class("ml-2"),
						).Class("flex items-center"),
						h.Div(
							RadioGroupItem().Value("female").Id("female"),
							Label(h.Text("Female")).For("female").Class("ml-2"),
						).Class("flex items-center"),
					).Attr("v-model", "form.Gender").Class("flex gap-4"),
				).Class("demo-row"),
			).Class("demo-section"),

			// RadioGroup - Choice Card 风格（整卡可点，选中高亮，可带图标/描述）
			h.Div(
				h.H2("RadioGroup (Choice Card)"),
				h.Div(
					RadioGroup(
						RadioCard().Value("starter").Title("Starter").
							Description("适合个人项目，免费").Icon("rocket"),
						RadioCard().Value("pro").Title("Pro").
							Description("$20/月，团队协作与高级功能").Icon("zap"),
						RadioCard().Value("enterprise").Title("Enterprise").
							Description("定制方案与专属支持").Icon("building-2").Disabled(true),
					).Attr("v-model", "form.Plan").Class("grid gap-3 sm:grid-cols-3"),
				).Class("demo-row"),
			).Class("demo-section"),

			// Checkbox - 使用内置 Label + Tips
			h.Div(
				h.H2("Checkbox (内置 Label)"),
				h.Div(
					Checkbox().
						Label("我同意服务条款").
						Tips("请仔细阅读服务条款后勾选").
						Attr("v-model", "form.Agreed"),
				).Class("demo-row"),
			).Class("demo-section"),

			// Switch - 使用内置 Label + Tips
			h.Div(
				h.H2("Switch (内置 Label)"),
				h.Div(
					Switch().
						Label("启用功能").
						Tips("开启后将激活高级功能").
						Attr("v-model", "form.Feature1"),
				).Class("demo-row"),
			).Class("demo-section"),

			// Submit Button
			h.Div(
				Button(h.Text("Update")).On("click", web.POST().EventFunc("update").Go()),
			).Class("demo-section"),
		).VSlot("{ locals, form }").FormInit(h.JSONString(formData)),
	).Style("max-width: 600px; margin: 0 auto;")
}

// update 更新表单数据
func update(ctx *web.EventContext) (r web.EventResponse, err error) {
	formData = &myFormValue{}
	ctx.MustUnmarshalForm(formData)

	verr := web.ValidationErrors{}
	if len(formData.Username) < 3 {
		verr.FieldError("Username", "请输入用户名")
	}

	if len(formData.Email) == 0 {
		verr.FieldError("Email", "Email is required")
	}

	if len(formData.TextareaValue) > 100 {
		verr.FieldError("TextareaValue", "Message is too long")
	}

	if !formData.Agreed {
		verr.FieldError("Agreed", "You must agree to the terms")
	}

	ctx.Flash = verr
	r.Reload = true

	return
}
