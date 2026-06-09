package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnFormFieldDemo 虚拟模型
type ShadcnFormFieldDemo struct{}

// configFormField 注册 Form Field demo
func configFormField(b *presets.Builder) {
	m := b.Model(&ShadcnFormFieldDemo{}).
		Label("Form Field").
		URIName("shadcn-form-field")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnFormFieldDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Form Field"
		r.Body = shadcnFormFieldBody(ctx)
		return
	})
}

// shadcnFormFieldBody FormField 演示页面
func shadcnFormFieldBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("FormField Demo").Style("margin-bottom: 24px;"),
		h.P(h.Text("统一表单字段包装器，支持多种类型、Label、错误消息、必填标记等")).Class("text-muted-foreground mb-6"),

		// 基础文本输入
		h.Div(
			h.H2("Basic Text Input"),
			h.P(h.Text("基础文本输入字段")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					FormField().
						Label("用户名").
						Placeholder("请输入用户名").
						Required(true).
						Attr("v-model", "form.username"),
				).Class("mb-4"),
				h.Div(
					h.Span("Value: {{ form.username || '-' }}").Class("text-sm text-muted-foreground"),
				),
			).VSlot("{ form }").FormInit(`{ "username": "" }`),
		).Class("demo-section"),

		// 不同类型的输入
		h.Div(
			h.H2("Different Input Types"),
			h.P(h.Text("多种输入类型演示")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					FormField().
						Label("邮箱").
						Type("email").
						Placeholder("example@domain.com").
						Required(true).
						Attr("v-model", "form.email"),
				).Class("mb-4"),
				h.Div(
					FormField().
						Label("电话").
						Type("tel").
						Placeholder("+86 188 8888 8888").
						Attr("v-model", "form.phone"),
				).Class("mb-4"),
				h.Div(
					FormField().
						Label("网址").
						Type("url").
						Placeholder("https://example.com").
						Attr("v-model", "form.website"),
				).Class("mb-4"),
				h.Div(
					FormField().
						Label("数字").
						Type("number").
						Placeholder("请输入数字").
						Attr("v-model", "form.age"),
				).Class("mb-4"),
				h.Div(
					h.Pre("{{ JSON.stringify(form, null, 2) }}").Class("text-xs bg-muted p-2 rounded"),
				),
			).VSlot("{ form }").FormInit(`{ "email": "", "phone": "", "website": "", "age": null }`),
		).Class("demo-section"),

		// Textarea
		h.Div(
			h.H2("Textarea"),
			h.P(h.Text("多行文本输入")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					FormField().
						Label("描述").
						Type("textarea").
						Placeholder("请输入描述信息...").
						Description("最多500字").
						Attr("v-model", "form.description"),
				).Class("mb-4"),
				h.Div(
					h.Span("{{ form.description?.length || 0 }} / 500").Class("text-sm text-muted-foreground"),
				),
			).VSlot("{ form }").FormInit(`{ "description": "" }`),
		).Class("demo-section"),

		// Password
		h.Div(
			h.H2("Password Input"),
			h.P(h.Text("密码输入，支持可见性切换")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					FormField().
						Label("密码").
						Type("password").
						Placeholder("请输入密码").
						PasswordVisibleToggle(true).
						Required(true).
						Attr("v-model", "form.password"),
				).Class("mb-4"),
				h.Div(
					FormField().
						Label("确认密码").
						Type("password").
						Placeholder("再次输入密码").
						PasswordVisibleToggle(true).
						Required(true).
						Attr("v-model", "form.confirmPassword").
						Attr(":error-messages", "form.password !== form.confirmPassword ? ['两次密码不一致'] : []"),
				).Class("mb-4"),
			).VSlot("{ form }").FormInit(`{ "password": "", "confirmPassword": "" }`),
		).Class("demo-section"),

		// 带提示信息
		h.Div(
			h.H2("With Tips"),
			h.P(h.Text("带帮助提示图标的字段")).Class("text-muted-foreground mb-4"),
			web.Scope(
				FormField().
					Label("API Key").
					Placeholder("请输入 API Key").
					Tips("从控制台获取您的 API Key\n格式：sk-xxxxxxxx").
					Required(true).
					Attr("v-model", "form.apiKey"),
			).VSlot("{ form }").FormInit(`{ "apiKey": "" }`),
		).Class("demo-section"),

		// 错误状态
		h.Div(
			h.H2("With Error Messages"),
			h.P(h.Text("带错误消息显示")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					FormField().
						Label("手机号").
						Placeholder("请输入11位手机号").
						Required(true).
						Attr("v-model", "form.mobile").
						Attr(":error-messages", "form.mobile && form.mobile.length !== 11 ? ['手机号必须为11位'] : []"),
				).Class("mb-4"),
				h.Div(
					FormField().
						Label("邮箱验证码").
						Placeholder("请输入6位验证码").
						Attr("v-model", "form.code").
						Attr(":error-messages", "form.code && form.code.length !== 6 ? ['验证码必须为6位'] : []"),
				),
			).VSlot("{ form }").FormInit(`{ "mobile": "", "code": "" }`),
		).Class("demo-section"),

		// 禁用和只读状态
		h.Div(
			h.H2("Disabled & Readonly"),
			h.P(h.Text("禁用和只读状态")).Class("text-muted-foreground mb-4"),
			h.Div(
				FormField().
					Label("禁用字段").
					Placeholder("此字段已禁用").
					Disabled(true).
					ModelValue("这是一个禁用的字段"),
			).Class("mb-4"),
			FormField().
				Label("只读字段").
				Placeholder("此字段只读").
				Readonly(true).
				ModelValue("这是一个只读字段"),
		).Class("demo-section"),

		// 完整表单示例
		h.Div(
			h.H2("Complete Form Example"),
			h.P(h.Text("完整的表单示例")).Class("text-muted-foreground mb-4"),
			web.Scope(
				Card(
					CardHeader(
						CardTitle(h.Text("用户注册")),
						CardDescription(h.Text("请填写以下信息完成注册")),
					),
					CardContent(
						h.Div(
							FormField().
								Label("用户名").
								Placeholder("请输入用户名").
								Required(true).
								Tips("用户名长度需要在3-20个字符之间").
								Attr("v-model", "form.regUsername").
								Attr(":error-messages", "form.regUsername && (form.regUsername.length < 3 || form.regUsername.length > 20) ? ['用户名长度需要在3-20个字符之间'] : []"),
						).Class("mb-4"),
						h.Div(
							FormField().
								Label("邮箱").
								Type("email").
								Placeholder("example@domain.com").
								Required(true).
								Attr("v-model", "form.regEmail"),
						).Class("mb-4"),
						h.Div(
							FormField().
								Label("密码").
								Type("password").
								Placeholder("请输入密码").
								PasswordVisibleToggle(true).
								Required(true).
								Tips("密码长度至少8位，需包含字母和数字").
								Attr("v-model", "form.regPassword"),
						).Class("mb-4"),
						h.Div(
							FormField().
								Label("确认密码").
								Type("password").
								Placeholder("再次输入密码").
								PasswordVisibleToggle(true).
								Required(true).
								Attr("v-model", "form.regConfirmPassword").
								Attr(":error-messages", "form.regPassword !== form.regConfirmPassword ? ['两次密码不一致'] : []"),
						).Class("mb-4"),
						h.Div(
							FormField().
								Label("个人简介").
								Type("textarea").
								Placeholder("介绍一下自己吧...").
								Description("选填，最多200字").
								Attr("v-model", "form.regBio"),
						).Class("mb-4"),
					),
					CardFooter(
						Button(h.Text("提交注册")).Class("w-full").
							Attr("@click", "console.log('表单已提交', form)"),
					),
				).Class("max-w-md"),
			).VSlot("{ form }").FormInit(`{
				"regUsername": "",
				"regEmail": "",
				"regPassword": "",
				"regConfirmPassword": "",
				"regBio": ""
			}`),
		).Class("demo-section"),
	).Style("max-width: 800px; margin: 0 auto;")
}
