package ec_demo

import (
	"strconv"
	"time"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

type ECDashboard struct{}

func ConfigECDashboard(pb *presets.Builder, db *gorm.DB) {
	b := pb.Model(&ECDashboard{}).Label("EC Dashboard").URIName("ec-dashboard")

	lb := b.Listing()

	// 声明自定义权限资源（fc_ 闸）：dashboard 各区块可在 role 权限树「自定义资源」分区独立授权/拒绝。
	// 执行端在 PageFunc 内用 canSee 校验，无权区块直接不渲染。
	lb.PermResource("stat_cards", "统计卡片")
	lb.PermResource("visitor_chart", "访问量图表")
	lb.PermResource("order_status_chart", "订单状态分布")

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		// canSee 校验当前用户对某自定义资源（fc_<key>）是否有 PermList 权限；越权返回 false → 不渲染该区块。
		canSee := func(key string) bool {
			return b.Info().Verifier().Do(presets.PermList).SnakeOn("fc_"+key).WithReq(ctx.R).IsAllowed() == nil
		}
		// DB query
		var productCount int64
		var orderCount int64

		if err = db.Model(&models.Product{}).Count(&productCount).Error; err != nil {
			r.Body = errorBody(err.Error())
			return
		}

		if err = db.Model(&models.Order{}).Count(&orderCount).Error; err != nil {
			r.Body = errorBody(err.Error())
			return
		}

		// 用 GROUP BY 统计各状态订单数量，避免全表加载
		type statusResult struct {
			Status string
			Count  int
		}
		var statusResults []statusResult
		if err = db.Model(&models.Order{}).Select("status, count(*) as count").Group("status").Scan(&statusResults).Error; err != nil {
			r.Body = errorBody(err.Error())
			return
		}
		statusCount := make(map[models.OrderStatus]int)
		for _, sr := range statusResults {
			statusCount[models.OrderStatus(sr.Status)] = sr.Count
		}

		// 订单状态柱状图数据
		type StatusData struct {
			Status string `json:"status"`
			Count  int    `json:"count"`
		}
		var statusChartData []StatusData
		for _, status := range models.OrderStatuses {
			statusChartData = append(statusChartData, StatusData{
				Status: string(status),
				Count:  statusCount[status],
			})
		}

		// 订单状态图表配置
		statusConfig := unovis.ChartConfig{
			"count": {Label: "订单数", Color: "var(--chart-1)"},
		}

		// 生成访问量数据（模拟 ec.md 的数据结构）
		type VisitorData struct {
			Date    string `json:"date"`
			Desktop int    `json:"desktop"`
			Mobile  int    `json:"mobile"`
		}
		var visitorData []VisitorData
		var totalDesktop, totalMobile int
		now := time.Now()

		// 生成最近90天的数据
		for i := 89; i >= 0; i-- {
			day := now.AddDate(0, 0, -i)
			desktop := 100 + (day.Day()*17+int(day.Month())*23)%400
			mobile := 80 + (day.Day()*13+int(day.Month())*31)%350
			visitorData = append(visitorData, VisitorData{
				Date:    day.Format("2006-01-02"),
				Desktop: desktop,
				Mobile:  mobile,
			})
			totalDesktop += desktop
			totalMobile += mobile
		}

		// 访问量图表配置
		visitorConfig := unovis.ChartConfig{
			"desktop": {Label: "Desktop", Color: "var(--chart-2)"},
			"mobile":  {Label: "Mobile", Color: "var(--chart-1)"},
		}

		var sections []h.HTMLComponent

		// 统计卡片（fc_stat_cards）：无权则整块不渲染
		if canSee("stat_cards") {
			sections = append(sections, h.Div(
				shadcn.Card(
					shadcn.CardHeader(
						shadcn.CardTitle(h.Text(strconv.Itoa(int(productCount)))),
						shadcn.CardDescription(h.Text("商品总数")),
					),
				).Class("flex-1"),
				shadcn.Card(
					shadcn.CardHeader(
						shadcn.CardTitle(h.Text(strconv.Itoa(int(orderCount)))),
						shadcn.CardDescription(h.Text("订单总数")),
					),
				).Class("flex-1"),
			).Class("flex gap-4 mb-6"))
		}

		// 交互式访问量图表（fc_visitor_chart，参考 ec.md）
		if canSee("visitor_chart") {
			sections = append(sections, web.Scope(
				shadcn.Card(
					shadcn.CardHeader(
						h.Div(
							h.Div(
								shadcn.CardTitle(h.Text("Bar Chart - Interactive")),
								shadcn.CardDescription(h.Text("Showing total visitors for the last 3 months")),
							).Class("flex flex-1 flex-col justify-center gap-1 px-6 pt-4 pb-3 sm:py-0"),
							h.Div(
								// Desktop 切换按钮
								h.Tag("button").Children(
									h.Span("Desktop").Class("text-muted-foreground text-xs"),
									h.Span(formatNumber(totalDesktop)).Class("text-lg leading-none font-bold sm:text-3xl"),
								).Class("relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6").
									Attr(":class", "{'bg-muted/50': locals.activeChart === 'desktop'}").
									Attr("@click", "locals.activeChart = 'desktop'"),
								// Mobile 切换按钮
								h.Tag("button").Children(
									h.Span("Mobile").Class("text-muted-foreground text-xs"),
									h.Span(formatNumber(totalMobile)).Class("text-lg leading-none font-bold sm:text-3xl"),
								).Class("relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6").
									Attr(":class", "{'bg-muted/50': locals.activeChart === 'mobile'}").
									Attr("@click", "locals.activeChart = 'mobile'"),
							).Class("flex"),
						).Class("flex flex-col items-stretch border-b p-0 sm:flex-row"),
					).Class("p-0"),
					shadcn.CardContent(
						// 根据 activeChart 显示不同的图表
						h.Div(
							unovis.BarChart(visitorConfig, visitorData, "date", "desktop").
								ShowGrid(true).
								ShowTooltip(true).
								Class("w-full aspect-auto h-64"),
						).Attr("v-show", "locals.activeChart === 'desktop'"),
						h.Div(
							unovis.BarChart(visitorConfig, visitorData, "date", "mobile").
								ShowGrid(true).
								ShowTooltip(true).
								Class("w-full aspect-auto h-64"),
						).Attr("v-show", "locals.activeChart === 'mobile'"),
					).Class("px-2 sm:p-6"),
				).Class("py-0"),
			).Init(`{activeChart: 'desktop'}`).VSlot("{ locals }"))
		}

		// 订单状态分布（fc_order_status_chart）
		if canSee("order_status_chart") {
			sections = append(sections, h.Div(
				// 订单状态 Bar Chart
				shadcn.Card(
					shadcn.CardHeader(
						shadcn.CardTitle(h.Text("订单状态分布")),
						shadcn.CardDescription(h.Text("各状态订单数量统计")),
					),
					shadcn.CardContent(
						unovis.BarChart(statusConfig, statusChartData, "status", "count").
							ShowGrid(true).
							ShowTooltip(true).
							Class("w-full"),
					),
					shadcn.CardFooter(
						h.Div(h.Text("展示各订单状态的数量分布")).Class("leading-none text-muted-foreground"),
					).Class("flex-col items-start gap-2 text-sm"),
				).Class("flex-1"),
			).Class("flex gap-4 mt-6"))
		}

		body := h.Div(sections...).Class("container mx-auto p-4")

		r.Body = body
		r.PageTitle = "EC Dashboard"

		return
	})
}

// formatNumber 格式化数字（添加千位分隔符）
func formatNumber(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	// 计算结果长度：原长度 + 逗号数量
	commas := (len(s) - 1) / 3
	result := make([]byte, len(s)+commas)
	j := len(result) - 1
	for i := len(s) - 1; i >= 0; i-- {
		result[j] = s[i]
		j--
		if (len(s)-i)%3 == 0 && i > 0 {
			result[j] = ','
			j--
		}
	}
	return string(result)
}

func errorBody(msg string) h.HTMLComponent {
	return h.Div(
		h.P().Text(msg),
	).Class("container mx-auto p-4")
}
