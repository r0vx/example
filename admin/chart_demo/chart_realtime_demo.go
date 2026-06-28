package chart_demo

import (
	"math"
	"time"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
)

// 实时图表范式对比 Demo（/chart-realtime-demo）—— 展示 RefreshInterval（定时重取替换族）：
//
//	A 刷新更新：固定「近 7 天按天」窗口，新数据只让当天的点原位涨落，X 轴标签不动、线不横移。
//	B 滑动：    滚动细粒度时间窗口（每 2 秒一个桶），窗口随时间前移 → 线左滑。
//	饼图：      环形图各切片随数据涨缩、中心 Total 实时更新。
//
// 三图都用 DataURL + RefreshInterval(2)（定时重取替换、document.hidden 时暂停）。数据为时间驱动的合成序列。
// 对照 /stream-chart-demo（StreamOn 客户端追加流）：本页=重取替换族。
const (
	evRefreshUpdateData = "chartRealtime_refreshUpdate"
	evSlidingStreamData = "chartRealtime_slidingStream"
	evDonutLiveData     = "chartRealtime_donutLive"
)

// donutCategories 环形图固定分类（名称 + 相位偏移让切片错峰涨缩 + 颜色）。
var donutCategories = []struct {
	name  string
	phase float64
	color string
}{
	{"Pending", 0, "var(--chart-1)"},
	{"Paid", 1.3, "var(--chart-2)"},
	{"Sending", 2.6, "var(--chart-3)"},
	{"Refunded", 3.9, "var(--chart-4)"},
	{"Cancelled", 5.2, "var(--chart-5)"},
}

// refreshUpdateData 范式 A：固定 7 天窗口；历史天值稳定，今天随 now 振荡 → 今天的点原位更新、线不滑。
func refreshUpdateData() []map[string]any {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	out := make([]map[string]any, 0, 7)
	for d := 6; d >= 0; d-- {
		day := todayStart.AddDate(0, 0, -d)
		v := 40 + int(day.Unix()/86400)%25 // 历史天：按天稳定
		if d == 0 {
			v = 50 + int(25*math.Sin(float64(now.Unix())/4.0)) // 今天：随时间振荡
		}
		out = append(out, map[string]any{"day": day.Format("01-02"), "count": v})
	}
	return out
}

// slidingStreamData 范式 B：最近 30 个 2 秒桶（升序，最新在右）。窗口随 now.Truncate(step) 前移 → 每刷新整体左移。
func slidingStreamData() []map[string]any {
	const n = 30
	step := 2 * time.Second
	end := time.Now().Truncate(step)
	out := make([]map[string]any, 0, n)
	for i := n - 1; i >= 0; i-- {
		t := end.Add(-time.Duration(i) * step)
		secs := float64(t.Unix())
		v := 50 + 30*math.Sin(secs/8.0) + 12*math.Sin(secs/3.0)
		out = append(out, map[string]any{"t": t.Format("04:05"), "v": math.Round(v)})
	}
	return out
}

// donutLiveData 环形图实时数据：各分类值随时间错峰振荡（恒正）→ 切片平滑涨缩、中心 Total 变。
func donutLiveData() []map[string]any {
	secs := float64(time.Now().Unix())
	out := make([]map[string]any, 0, len(donutCategories))
	for _, c := range donutCategories {
		v := 20 + math.Round(15*math.Sin(secs/5.0+c.phase)+15) // [20,50] 恒正
		out = append(out, map[string]any{"category": c.name, "count": v})
	}
	return out
}

// ConfigChartRealtimeDemo 注册「实时图表范式对比」演示页（/chart-realtime-demo）。
func ConfigChartRealtimeDemo(b *presets.Builder) {
	wb := b.GetWebBuilder()
	wb.RegisterEventFunc(evRefreshUpdateData, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.Data = refreshUpdateData()
		return
	})
	wb.RegisterEventFunc(evSlidingStreamData, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.Data = slidingStreamData()
		return
	})
	wb.RegisterEventFunc(evDonutLiveData, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.Data = donutLiveData()
		return
	})

	cb := presets.NewCustomPage(b).
		PageTitleFunc(func(*web.EventContext) string { return "实时图表范式对比" }).
		Body(func(ctx *web.EventContext) h.HTMLComponent {
			cfgA := unovis.ChartConfig{"count": {Label: "订单数", Color: "var(--chart-1)"}}
			cfgB := unovis.ChartConfig{"v": {Label: "实时值", Color: "var(--chart-2)"}}
			cfgDonut := unovis.ChartConfig{}
			for _, c := range donutCategories {
				cfgDonut[c.name] = unovis.ChartConfigItem{Label: c.name, Color: c.color}
			}

			donutCard := shadcn.Card(
				shadcn.CardHeader(
					h.Div(h.Text("环形图 · 定时刷新")).Class("text-sm font-medium"),
					h.Div(h.Text("各切片随数据平滑涨缩、中心 Total 实时更新（AutoCentralTotal + RefreshInterval）。")).
						Class("text-xs text-muted-foreground mt-1"),
				).Class("p-4 pb-0"),
				shadcn.CardContent(
					unovis.RadialChart(cfgDonut, donutLiveData(), "count", "category").
						AutoCentralTotal(true).
						CentralSubLabel("Total").
						DataURL("?__execute_event__="+evDonutLiveData).
						RefreshInterval(2).
						Class("w-full h-72"),
				).Class("p-3"),
			)

			cardA := shadcn.Card(
				shadcn.CardHeader(
					h.Div(h.Text("A · 刷新更新（固定近 7 天窗口）")).Class("text-sm font-medium"),
					h.Div(h.Text("新数据让今天的点原位涨落，X 轴不动、线不横移。按天聚合业务量的主流范式。")).
						Class("text-xs text-muted-foreground mt-1"),
				).Class("p-4 pb-0"),
				shadcn.CardContent(
					unovis.AreaChart(cfgA, refreshUpdateData(), "day", "count").
						CurveType(unovis.ChartCurveBasis).
						ShowYAxis(false).ShowGrid(false).ShowLegend(false).
						DataURL("?__execute_event__="+evRefreshUpdateData).
						RefreshInterval(2).
						Class("w-full h-56"),
				).Class("p-3"),
			)

			cardB := shadcn.Card(
				shadcn.CardHeader(
					h.Div(h.Text("B · 滑动（滚动 2 秒桶窗口，定时重取替换）")).Class("text-sm font-medium"),
					h.Div(h.Text("窗口随时间前移，线左滑。注意：这是定时重取整窗口（轮询），非客户端追加——追加流见 /stream-chart-demo。")).
						Class("text-xs text-muted-foreground mt-1"),
				).Class("p-4 pb-0"),
				shadcn.CardContent(
					unovis.AreaChart(cfgB, slidingStreamData(), "t", "v").
						CurveType(unovis.ChartCurveBasis).
						ShowYAxis(false).ShowGrid(false).ShowLegend(false).
						DataURL("?__execute_event__="+evSlidingStreamData).
						RefreshInterval(2).
						Class("w-full h-56"),
				).Class("p-3"),
			)

			return h.Div(
				h.H1("实时图表范式对比").Class("text-2xl font-bold mb-1"),
				h.Div(h.Text("均每 2 秒自动重取（document.hidden 时暂停）。左：环形图定时刷新；右：折线「刷新更新」vs「滑动」。")).
					Class("text-sm text-muted-foreground mb-6"),
				h.Div(
					donutCard,
					h.Div(cardA, cardB).Class("flex flex-col gap-6"),
				).Class("grid grid-cols-1 lg:grid-cols-2 gap-6 items-start"),
			).Class("container mx-auto p-6")
		})
	b.HandleCustomPage("chart-realtime-demo", cb)
}
