package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnRangePickerDemo 虚拟模型
type ShadcnRangePickerDemo struct{}

// configRangePicker 注册 Range Picker demo
func configRangePicker(b *presets.Builder) {
	m := b.Model(&ShadcnRangePickerDemo{}).
		Label("Range Picker").
		URIName("shadcn-range-picker")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnRangePickerDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Range Picker"
		r.Body = shadcnRangePickerBody(ctx)
		return
	})
}

// shadcnRangePickerBody 日期范围选择器演示页面
func shadcnRangePickerBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("RangePicker Demo").Style("margin-bottom: 24px;"),
		h.P(h.Text("日期范围选择组件，结合 Popover、Button 和 RangeCalendar 实现")).Class("text-muted-foreground mb-6"),

		// 基础用法
		h.Div(
			h.H2("Basic Usage"),
			h.P(h.Text("基础日期范围选择")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					RangePicker().Placeholder("选择日期范围").Attr("v-model", "form.range"),
				).Class("mb-4"),
				h.Div(
					h.Span("Start: {{ form.range?.start || '-' }}").Class("text-sm text-muted-foreground mr-4"),
					h.Span("End: {{ form.range?.end || '-' }}").Class("text-sm text-muted-foreground"),
				),
			).VSlot("{ form }").FormInit(`{ "range": null }`),
		).Class("demo-section"),

		// 双月视图
		h.Div(
			h.H2("Two Month View"),
			h.P(h.Text("显示两个月份的日历视图")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					RangePicker().Placeholder("选择日期范围").NumberOfMonths(2).Attr("v-model", "form.range2"),
				).Class("mb-4"),
				h.Div(
					h.Span("Start: {{ form.range2?.start || '-' }}").Class("text-sm text-muted-foreground mr-4"),
					h.Span("End: {{ form.range2?.end || '-' }}").Class("text-sm text-muted-foreground"),
				),
			).VSlot("{ form }").FormInit(`{ "range2": null }`),
		).Class("demo-section"),

		// 禁用状态
		h.Div(
			h.H2("Disabled State"),
			h.P(h.Text("禁用的日期范围选择器")).Class("text-muted-foreground mb-4"),
			RangePicker().Placeholder("禁用状态").Disabled(true),
		).Class("demo-section"),

		// 表单中使用
		h.Div(
			h.H2("In Form Context"),
			h.P(h.Text("在表单中使用日期范围选择器")).Class("text-muted-foreground mb-4"),
			web.Scope(
				Card(
					CardHeader(
						CardTitle(h.Text("预订查询")),
						CardDescription(h.Text("选择入住和离店日期")),
					),
					CardContent(
						h.Div(
							h.Div(
								Label(h.Text("客人姓名")).Class("mb-2"),
								Input().Placeholder("请输入姓名").Attr("v-model", "form.name"),
							).Class("mb-4"),
							h.Div(
								Label(h.Text("入住日期范围")).Class("mb-2"),
								RangePicker().Placeholder("选择入住和离店日期").Attr("v-model", "form.dateRange"),
							).Class("mb-4"),
							h.Div(
								Label(h.Text("房间数量")).Class("mb-2"),
								Input().Type("number").Placeholder("1").Attr("v-model", "form.rooms"),
							).Class("mb-4"),
						),
					),
					CardFooter(
						Button(h.Text("查询")).Class("w-full"),
					),
				).Class("max-w-md"),
				h.Div(
					h.H4("Form Data:").Class("text-sm font-medium mt-4 mb-2"),
					h.Pre("{{ JSON.stringify(form, null, 2) }}").Class("text-xs bg-muted p-2 rounded"),
				),
			).VSlot("{ form }").FormInit(`{ "name": "", "dateRange": null, "rooms": 1 }`),
		).Class("demo-section"),

		// 预设范围
		h.Div(
			h.H2("With Presets"),
			h.P(h.Text("带预设快捷选项的日期范围")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					h.Div(
						Button(h.Text("今天")).Variant(ButtonVariantOutline).Size(ButtonSizeSm).
							Attr("@click", "locals.setToday(form)").Class("mr-2"),
						Button(h.Text("最近7天")).Variant(ButtonVariantOutline).Size(ButtonSizeSm).
							Attr("@click", "locals.setLast7Days(form)").Class("mr-2"),
						Button(h.Text("最近30天")).Variant(ButtonVariantOutline).Size(ButtonSizeSm).
							Attr("@click", "locals.setLast30Days(form)").Class("mr-2"),
						Button(h.Text("本月")).Variant(ButtonVariantOutline).Size(ButtonSizeSm).
							Attr("@click", "locals.setThisMonth(form)"),
					).Class("mb-4"),
					RangePicker().Placeholder("选择日期范围").Attr("v-model", "form.presetRange"),
					h.Div(
						h.Span("Start: {{ form.presetRange?.start || '-' }}").Class("text-sm text-muted-foreground mr-4"),
						h.Span("End: {{ form.presetRange?.end || '-' }}").Class("text-sm text-muted-foreground"),
					).Class("mt-4"),
				),
			).VSlot("{ locals, form }").FormInit(`{ "presetRange": null }`).Init(`{
				formatDate(d) {
					const year = d.getFullYear();
					const month = String(d.getMonth() + 1).padStart(2, '0');
					const day = String(d.getDate()).padStart(2, '0');
					return year + '-' + month + '-' + day;
				},
				setToday(form) {
					const today = this.formatDate(new Date());
					form.presetRange = { start: today, end: today };
				},
				setLast7Days(form) {
					const end = new Date();
					const start = new Date();
					start.setDate(start.getDate() - 6);
					form.presetRange = {
						start: this.formatDate(start),
						end: this.formatDate(end)
					};
				},
				setLast30Days(form) {
					const end = new Date();
					const start = new Date();
					start.setDate(start.getDate() - 29);
					form.presetRange = {
						start: this.formatDate(start),
						end: this.formatDate(end)
					};
				},
				setThisMonth(form) {
					const now = new Date();
					const start = new Date(now.getFullYear(), now.getMonth(), 1);
					const end = new Date(now.getFullYear(), now.getMonth() + 1, 0);
					form.presetRange = {
						start: this.formatDate(start),
						end: this.formatDate(end)
					};
				}
			}`),
		).Class("demo-section"),
	).Style("max-width: 800px; margin: 0 auto;")
}
