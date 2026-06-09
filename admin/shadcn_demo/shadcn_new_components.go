package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnNewComponentsDemo 虚拟模型
type ShadcnNewComponentsDemo struct{}

// configNewComponents 注册 New Components demo
func configNewComponents(b *presets.Builder) {
	m := b.Model(&ShadcnNewComponentsDemo{}).
		Label("New Components").
		URIName("shadcn-new-components")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnNewComponentsDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "New Components"
		r.Body = shadcnNewComponentsBody(ctx)
		return
	})
}

// shadcnNewComponentsBody 新组件演示页面
func shadcnNewComponentsBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("New Components Demo").Style("margin-bottom: 24px;"),

		// Button Loading
		h.Div(
			h.H2("Button Loading"),
			h.P(h.Text("按钮加载状态")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					Button(h.Text("Default")).Loading(true),
					Button(h.Text("Outline")).Variant(ButtonVariantOutline).Loading(true),
					Button(h.Text("Destructive")).Variant(ButtonVariantDestructive).Loading(true),
					Button(h.Text("Normal")),
				).Class("flex items-center gap-2 mb-4"),
				// 带 icon 的按钮
				h.Div(
					Button(Icon(IconPlus).Size(12).Class("mr-1"), h.Text("添加")).Size(ButtonSizeXs).Loading(true),
					Button(Icon(IconPencil).Size(16).Class("mr-1"), h.Text("编辑")).Size(ButtonSizeDefault).Variant(ButtonVariantOutline).Loading(true),
					Button(Icon(IconPlus).Size(16).Class("mr-1"), h.Text("添加")).Size(ButtonSizeSm),
					Button(Icon(IconPencil).Size(16).Class("mr-1"), h.Text("编辑")).Size(ButtonSizeSm).Variant(ButtonVariantOutline),
				).Class("flex items-center gap-2 mb-4"),
				// 动态切换 loading
				web.Scope(
					h.Div(
						Button(Icon(IconPencil).Size(16).Class("mr-1"), h.Text("编辑")).
							Loading(true).Attr(":loading", "locals.loading").
							Attr("@click", "locals.loading = true; vars.__window.setTimeout(() => locals.loading = false, 2000)"),
						Button(h.Text("Save")).
							Variant(ButtonVariantOutline).
							Attr(":loading", "locals.loading").
							Attr("@click", "locals.loading = true; vars.__window.setTimeout(() => locals.loading = false, 2000)"),
						h.Span("{{ locals.loading ? 'Loading...' : 'Click to test' }}").Class("text-sm text-muted-foreground"),
					).Class("flex items-center gap-2"),
				).VSlot("{ locals }").Init(`{ loading: false }`),
			),
		).Class("demo-section"),

		// Slider
		h.Div(
			h.H2("Slider"),
			h.P(h.Text("数值滑块组件")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					h.Span("Basic Slider").Class("text-sm font-medium"),
					web.Scope(
						h.Div(
							Slider().DefaultValue(50).Max(100).Class("w-72").Attr("v-model", "form.value"),
							h.Span("{{ form.value[0] }}").Class("ml-4 text-sm"),
						).Class("flex items-center"),
					).VSlot("{ form }").FormInit(`{ "value": [50] }`),
				).Class("mb-4"),
				h.Div(
					h.Span("With Label & Tips").Class("text-sm font-medium block mb-2"),
					Slider().Label("音量").Tips("调整系统音量大小").DefaultValue(70).Max(100).Class("w-72"),
				).Class("mb-4"),
				h.Div(
					h.Span("Required with Step").Class("text-sm font-medium block mb-2"),
					Slider().Label("亮度").Required(true).DefaultValue(50).Max(100).Step(10).Class("w-72"),
				).Class("mb-4"),
				h.Div(
					h.Span("With Error").Class("text-sm font-medium block mb-2"),
					Slider().Label("评分").Required(true).DefaultValue(0).Max(100).Step(25).Class("w-72").
						ErrorMessages("请选择评分"),
				).Class("mb-4"),
				h.Div(
					h.Span("Disabled").Class("text-sm font-medium"),
					Slider().DefaultValue(30).Disabled(true).Class("w-72"),
				),
			),
		).Class("demo-section"),

		// Alert
		h.Div(
			h.H2("Alert"),
			h.P(h.Text("警告提示组件")).Class("text-muted-foreground mb-4"),
			h.Div(
				Alert(
					h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>`),
					AlertTitle(h.Text("提示")),
					AlertDescription(h.Text("这是一条默认提示信息。")),
				).Class("mb-4"),
				Alert(
					h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`),
					AlertTitle(h.Text("警告")),
					AlertDescription(h.Text("这是一条警告信息，请注意！")),
				).Variant(AlertVariantDestructive),
			).Class("max-w-md"),
		).Class("demo-section"),

		// Spinner
		h.Div(
			h.H2("Spinner"),
			h.P(h.Text("加载动画组件")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					h.Span("Small").Class("text-sm mr-4"),
					Spinner().Size(SpinnerSizeSm),
				).Class("flex items-center mb-4"),
				h.Div(
					h.Span("Medium").Class("text-sm mr-4"),
					Spinner().Size(SpinnerSizeMd),
				).Class("flex items-center mb-4"),
				h.Div(
					h.Span("Large").Class("text-sm mr-4"),
					Spinner().Size(SpinnerSizeLg),
				).Class("flex items-center mb-4"),
				h.Div(
					h.Span("With Text").Class("text-sm mr-4"),
					Spinner().Size(SpinnerSizeMd),
					h.Span("Loading...").Class("ml-2 text-sm text-muted-foreground"),
				).Class("flex items-center"),
			),
		).Class("demo-section"),

		// TimePicker
		h.Div(
			h.H2("TimePicker"),
			h.P(h.Text("时间选择组件")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					h.Span("Basic").Class("text-sm font-medium block mb-2"),
					web.Scope(
						h.Div(
							TimePicker().Label("选择时间").Attr("v-model", "form.time1"),
							h.Span("{{ form.time1 || '未选择' }}").Class("ml-4 text-sm text-muted-foreground"),
						).Class("flex items-center"),
					).VSlot("{ form }").FormInit(`{ "time1": "" }`),
				).Class("mb-4"),
				h.Div(
					h.Span("With Seconds").Class("text-sm font-medium block mb-2"),
					web.Scope(
						h.Div(
							TimePicker().Label("选择时间").ShowSeconds(true).Attr("v-model", "form.time2"),
							h.Span("{{ form.time2 || '未选择' }}").Class("ml-4 text-sm text-muted-foreground"),
						).Class("flex items-center"),
					).VSlot("{ form }").FormInit(`{ "time2": "" }`),
				).Class("mb-4"),
				h.Div(
					h.Span("With Default").Class("text-sm font-medium block mb-2"),
					web.Scope(
						h.Div(
							TimePicker().Attr("v-model", "form.time3"),
							h.Span("{{ form.time3 }}").Class("ml-4 text-sm text-muted-foreground"),
						).Class("flex items-center"),
					).VSlot("{ form }").FormInit(`{ "time3": "14:30" }`),
				),
			),
		).Class("demo-section"),

		// RangePicker
		h.Div(
			h.H2("RangePicker"),
			h.P(h.Text("日期范围选择组件")).Class("text-muted-foreground mb-4"),
			h.Div(
				web.Scope(
					h.Div(
						RangePicker().Placeholder("选择日期范围").Attr("v-model", "form.range"),
					).Class("mb-2"),
					h.Div(
						h.Span("Start: {{ form.range?.start || '-' }}").Class("text-sm text-muted-foreground mr-4"),
						h.Span("End: {{ form.range?.end || '-' }}").Class("text-sm text-muted-foreground"),
					),
				).VSlot("{ form }").FormInit(`{ "range": null }`),
			),
		).Class("demo-section"),

		// TagInput
		h.Div(
			h.H2("TagInput"),
			h.P(h.Text("标签输入组件")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					h.Span("Basic").Class("text-sm font-medium block mb-2"),
					web.Scope(
						h.Div(
							TagInput().Placeholder("输入后按 Enter 添加标签").Attr("v-model", "form.tags"),
						).Class("max-w-md mb-2"),
						h.Span("Tags: {{ form.tags.join(', ') || '无' }}").Class("text-sm text-muted-foreground"),
					).VSlot("{ form }").FormInit(`{ "tags": [] }`),
				).Class("mb-4"),
				h.Div(
					h.Span("With Default Tags").Class("text-sm font-medium block mb-2"),
					web.Scope(
						TagInput().Attr("v-model", "form.skills").Class("max-w-md"),
					).VSlot("{ form }").FormInit(`{ "skills": ["Vue", "Go", "TypeScript"] }`),
				).Class("mb-4"),
				h.Div(
					h.Span("Max 3 Tags").Class("text-sm font-medium block mb-2"),
					TagInput().MaxTags(3).Placeholder("最多3个标签").Class("max-w-md"),
				),
			),
		).Class("demo-section"),

		// Stepper
		h.Div(
			h.H2("Stepper"),
			h.P(h.Text("步骤条组件")).Class("text-muted-foreground mb-4"),
			web.Scope(
				Stepper(
					StepperItem(
						StepperTrigger(h.Text("1")),
						h.Div(
							StepperTitle(h.Text("账号信息")),
							StepperDescription(h.Text("填写基本信息")),
						),
						StepperSeparator(),
					).Step(1),
					StepperItem(
						StepperTrigger(h.Text("2")),
						h.Div(
							StepperTitle(h.Text("验证邮箱")),
							StepperDescription(h.Text("确认邮箱地址")),
						),
						StepperSeparator(),
					).Step(2),
					StepperItem(
						StepperTrigger(h.Text("3")),
						h.Div(
							StepperTitle(h.Text("完成注册")),
							StepperDescription(h.Text("设置密码")),
						),
					).Step(3),
				).Attr("v-model", "form.step").Class("w-full max-w-lg"),
				h.Div(
					Button(h.Text("上一步")).Variant(ButtonVariantOutline).Attr("@click", "form.step > 1 && form.step--").Class("mr-2"),
					Button(h.Text("下一步")).Attr("@click", "form.step < 3 && form.step++"),
					h.Span("当前步骤: {{ form.step }}").Class("ml-4 text-sm text-muted-foreground"),
				).Class("mt-4"),
			).VSlot("{ form }").FormInit(`{ "step": 1 }`),
		).Class("demo-section"),

		// FileInput
		h.Div(
			h.H2("FileInput"),
			h.P(h.Text("文件上传组件")).Class("text-muted-foreground mb-4"),
			h.Div(
				h.Div(
					h.Span("Basic").Class("text-sm font-medium block mb-2"),
					FileInput().Placeholder("点击或拖拽上传文件").Class("max-w-md"),
				).Class("mb-4"),
				h.Div(
					h.Span("Images Only").Class("text-sm font-medium block mb-2"),
					FileInput().Accept("image/*").Placeholder("仅支持图片文件").Class("max-w-md"),
				).Class("mb-4"),
				h.Div(
					h.Span("Multiple Files").Class("text-sm font-medium block mb-2"),
					FileInput().Multiple(true).Placeholder("可选择多个文件").Class("max-w-md"),
				).Class("mb-4"),
				h.Div(
					h.Span("With Max Size (1MB)").Class("text-sm font-medium block mb-2"),
					FileInput().MaxSize(1024*1024).Placeholder("最大 1MB").Class("max-w-md"),
				),
			),
		).Class("demo-section"),
	).Style("max-width: 800px; margin: 0 auto;")
}
