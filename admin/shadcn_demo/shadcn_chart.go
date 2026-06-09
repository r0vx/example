package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/unovis"
	h "github.com/r0vx/htmlgo"
)

// ShadcnChartDemo 虚拟模型
type ShadcnChartDemo struct{}

// configChart 注册 Chart demo
func configChart(b *presets.Builder) {
	m := b.Model(&ShadcnChartDemo{}).
		Label("Chart").
		URIName("shadcn-chart")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnChartDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Chart"
		r.Body = shadcnChartBody(ctx)
		return
	})
}

// shadcnChartBody 图表组件演示页面
func shadcnChartBody(ctx *web.EventContext) h.HTMLComponent {
	// 柱状图配置
	barConfig := unovis.ChartConfig{
		"online":  {Label: "线上", Color: "hsl(var(--chart-1))"},
		"offline": {Label: "线下", Color: "hsl(var(--chart-2))"},
	}
	barData := []map[string]any{
		{"day": "周一", "online": 120, "offline": 60},
		{"day": "周二", "online": 200, "offline": 80},
		{"day": "周三", "online": 150, "offline": 70},
		{"day": "周四", "online": 80, "offline": 50},
		{"day": "周五", "online": 70, "offline": 60},
		{"day": "周六", "online": 110, "offline": 90},
		{"day": "周日", "online": 130, "offline": 100},
	}

	// 折线图配置
	lineConfig := unovis.ChartConfig{
		"newUsers":    {Label: "新用户", Color: "hsl(var(--chart-1))"},
		"activeUsers": {Label: "活跃用户", Color: "hsl(var(--chart-2))"},
	}
	lineData := []map[string]any{
		{"month": "1月", "newUsers": 120, "activeUsers": 80},
		{"month": "2月", "newUsers": 200, "activeUsers": 150},
		{"month": "3月", "newUsers": 150, "activeUsers": 120},
		{"month": "4月", "newUsers": 300, "activeUsers": 250},
		{"month": "5月", "newUsers": 280, "activeUsers": 220},
		{"month": "6月", "newUsers": 400, "activeUsers": 350},
	}

	// 面积图配置
	areaConfig := unovis.ChartConfig{
		"pv": {Label: "PV", Color: "hsl(var(--chart-1))"},
		"uv": {Label: "UV", Color: "hsl(var(--chart-2))"},
	}
	areaData := []map[string]any{
		{"time": "00:00", "pv": 200, "uv": 100},
		{"time": "04:00", "pv": 150, "uv": 80},
		{"time": "08:00", "pv": 400, "uv": 200},
		{"time": "12:00", "pv": 800, "uv": 400},
		{"time": "16:00", "pv": 600, "uv": 300},
		{"time": "20:00", "pv": 900, "uv": 450},
		{"time": "24:00", "pv": 500, "uv": 250},
	}

	// 堆叠柱状图配置
	stackedConfig := unovis.ChartConfig{
		"desktop": {Label: "Desktop", Color: "hsl(var(--chart-1))"},
		"mobile":  {Label: "Mobile", Color: "hsl(var(--chart-2))"},
	}
	stackedData := []map[string]any{
		{"month": "January", "desktop": 186, "mobile": 80},
		{"month": "February", "desktop": 305, "mobile": 200},
		{"month": "March", "desktop": 237, "mobile": 120},
		{"month": "April", "desktop": 73, "mobile": 190},
		{"month": "May", "desktop": 209, "mobile": 130},
		{"month": "June", "desktop": 214, "mobile": 140},
	}

	return h.Div(
		h.H1("Chart 图表组件").Class("text-2xl font-bold mb-6"),
		h.P(h.Text("基于 Unovis 实现，与 shadcn-vue 官方一致")).Class("text-muted-foreground mb-8"),

		// 柱状图
		h.Div(
			h.H2("柱状图 (Bar Chart)").Class("text-lg font-semibold mb-2"),
			h.P(h.Text("用于比较不同类别的数据")).Class("text-muted-foreground mb-4"),
			unovis.BarChart(barConfig, barData, "day", "online", "offline").Class("h-72"),
		).Class("mb-8"),

		// 折线图
		h.Div(
			h.H2("折线图 (Line Chart)").Class("text-lg font-semibold mb-2"),
			h.P(h.Text("展示数据随时间变化的趋势")).Class("text-muted-foreground mb-4"),
			unovis.LineChart(lineConfig, lineData, "month", "newUsers", "activeUsers").Class("h-72"),
		).Class("mb-8"),

		// 面积图
		h.Div(
			h.H2("面积图 (Area Chart)").Class("text-lg font-semibold mb-2"),
			h.P(h.Text("强调数量随时间变化的程度")).Class("text-muted-foreground mb-4"),
			unovis.AreaChart(areaConfig, areaData, "time", "pv", "uv").Class("h-72"),
		).Class("mb-8"),

		// 堆叠柱状图
		h.Div(
			h.H2("堆叠柱状图 (Stacked Bar Chart)").Class("text-lg font-semibold mb-2"),
			h.P(h.Text("多个系列数据堆叠显示")).Class("text-muted-foreground mb-4"),
			unovis.StackedBarChart(stackedConfig, stackedData, "month", "desktop", "mobile").Class("h-72"),
		).Class("mb-8"),

		// 加载状态
		h.Div(
			h.H2("加载状态").Class("text-lg font-semibold mb-2"),
			h.P(h.Text("图表加载时显示 loading 状态")).Class("text-muted-foreground mb-4"),
			unovis.Chart().Config(lineConfig).Data(lineData).XKey("month").YKeys("newUsers").Loading(true).Class("h-48"),
		).Class("mb-8"),
	).Class("max-w-4xl mx-auto p-6")
}
