package crud_demo

import (
	"time"

	"example/models"

	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
	"gorm.io/gorm"
)

// orderStatusDatum 环形图数据项：category=状态值，count=该状态订单数。AutoCentralTotal 按 count 求和显示中心总数。
type orderStatusDatum struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

// orderDailyDatum A（原位更新）数据项：day=MM-DD，count=当日订单数。
type orderDailyDatum struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

// orderMinuteDatum B（滑动）数据项：t=HH:MM（分钟桶），count=该分钟订单数。
type orderMinuteDatum struct {
	T     string `json:"t"`
	Count int64  `json:"count"`
}

// orderStatusData 各状态订单数（环形图）。db.Model 应用软删除作用域。
func orderStatusData(db *gorm.DB) []orderStatusDatum {
	var rows []orderStatusDatum
	db.Model(&models.Order{}).
		Select("status as category, COUNT(*) as count").
		Group("status").
		Scan(&rows)
	return rows
}

// orderDailyData A：近 7 天按天订单数（固定窗口，当天值随新单变 → 原位更新、线不横移）。
func orderDailyData(db *gorm.DB) []orderDailyDatum {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	since := todayStart.AddDate(0, 0, -6)

	type dayCount struct {
		Day   string
		Count int64
	}
	var rows []dayCount
	db.Model(&models.Order{}).
		Select("to_char(created_at, 'MM-DD') as day, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("day").
		Scan(&rows)
	m := make(map[string]int64, len(rows))
	for _, r := range rows {
		m[r.Day] = r.Count
	}
	out := make([]orderDailyDatum, 0, 7)
	for i := 6; i >= 0; i-- {
		key := todayStart.AddDate(0, 0, -i).Format("01-02")
		out = append(out, orderDailyDatum{Day: key, Count: m[key]})
	}
	return out
}

// orderSlidingData B：近 30 分钟按分钟订单数（滚动窗口，窗口随 now 前移 → 重取时整体左滑）。
func orderSlidingData(db *gorm.DB) []orderMinuteDatum {
	const n = 30
	step := time.Minute
	end := time.Now().Truncate(step)
	since := end.Add(-time.Duration(n-1) * step)

	var rows []struct {
		Bucket time.Time `gorm:"column:bucket"`
		Count  int64     `gorm:"column:count"`
	}
	db.Model(&models.Order{}).
		Select("date_trunc('minute', created_at) as bucket, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("bucket").
		Scan(&rows)
	m := make(map[string]int64, len(rows))
	for _, r := range rows {
		m[r.Bucket.Local().Format("15:04")] = r.Count
	}
	out := make([]orderMinuteDatum, 0, n)
	for i := n - 1; i >= 0; i-- {
		key := end.Add(-time.Duration(i)*step).Format("15:04")
		out = append(out, orderMinuteDatum{T: key, Count: m[key]})
	}
	return out
}

// orderStatusChartHeader 渲染订单列表顶部图表区（chart-realtime 三件套，全部 DataURL + RefreshOn 订单事件驱动）：
//   - 环形图：状态分布，切片随新单平滑涨缩、中心 Total 实时更新（AutoCentralTotal）
//   - A 原位更新：近 7 天按天，固定窗口、当天的点原位涨落
//   - B 滑动：    近 30 分钟按分钟，滚动窗口、重取时左滑
//
// SSE 只推「通知」（订单增删改 NotifModels*，载荷 {ids}）；三图 RefreshOn 收到通知 → 重取 DataURL → 平滑更新。
// 无定时器、无轮询、无追加（与表格行级刷新解耦，各自更新、互不闪屏）。
func orderStatusChartHeader(db *gorm.DB, statusURL, dailyURL, slidingURL string, refreshEvents []string) func(ctx *web.EventContext) h.HTMLComponent {
	return func(ctx *web.EventContext) h.HTMLComponent {
		statusConfig := unovis.ChartConfig{}
		for _, st := range models.OrderStatuses {
			statusConfig[string(st)] = unovis.ChartConfigItem{Label: string(st), Color: models.OrderStatusColorMap[st]}
		}

		statusCard := shadcn.Card(
			shadcn.CardHeader(
				h.Div(h.Text("Order Status")).Class("text-sm font-medium text-muted-foreground"),
			).Class("p-3 pb-0"),
			shadcn.CardContent(
				unovis.RadialChart(statusConfig, orderStatusData(db), "count", "category").
					AutoCentralTotal(true).
					CentralSubLabel("Total").
					DataURL(statusURL).
					RefreshOn(refreshEvents...).
					Class("w-full h-40"),
			).Class("p-2"),
		).Class("overflow-hidden")

		dailyConfig := unovis.ChartConfig{"count": {Label: "Orders", Color: models.OrderStatusColor_Blue}}
		dailyCard := shadcn.Card(
			shadcn.CardHeader(
				h.Div(h.Text("Last 7 Days")).Class("text-sm font-medium text-muted-foreground"),
			).Class("p-3 pb-0"),
			shadcn.CardContent(
				unovis.AreaChart(dailyConfig, orderDailyData(db), "day", "count").
					CurveType(unovis.ChartCurveBasis).
					ShowXAxis(true).ShowYAxis(false).ShowGrid(false).ShowTooltip(true).ShowLegend(false).
					DataURL(dailyURL).
					RefreshOn(refreshEvents...).
					Class("w-full h-40"),
			).Class("p-2"),
		).Class("overflow-hidden")

		slidingConfig := unovis.ChartConfig{"count": {Label: "Orders/min", Color: models.OrderStatusColor_Green}}
		slidingCard := shadcn.Card(
			shadcn.CardHeader(
				h.Div(h.Text("Last 30 Min")).Class("text-sm font-medium text-muted-foreground"),
			).Class("p-3 pb-0"),
			shadcn.CardContent(
				unovis.AreaChart(slidingConfig, orderSlidingData(db), "t", "count").
					CurveType(unovis.ChartCurveBasis).
					ShowXAxis(true).ShowYAxis(false).ShowGrid(false).ShowTooltip(true).ShowLegend(false).
					DataURL(slidingURL).
					RefreshOn(refreshEvents...).
					Class("w-full h-40"),
			).Class("p-2"),
		).Class("overflow-hidden")

		return h.Div(
			statusCard,
			dailyCard,
			slidingCard,
		).Class("grid grid-cols-1 md:grid-cols-3 gap-4 px-4 pt-3 pb-1")
	}
}
