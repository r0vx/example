package chart_demo

import (
	"math"
	"sync"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
)

// evStreamDemoPoint 滚动流 demo 的推点事件名（与图表 StreamOn 经 strcase.ToCamel 对齐）。
const evStreamDemoPoint = "streamDemoPoint"

// streamDemoOnce 保证后台推点定时器只启一次（config 多次调用也不重复起 goroutine）。
var streamDemoOnce sync.Once

// streamDemoPointAt 按时刻生成一个合成数据点（正弦叠加，恒正），t=HH:MM:SS、v=值。
func streamDemoPointAt(t time.Time) map[string]any {
	ms := float64(t.UnixMilli())
	v := 50 + 30*math.Sin(ms/4000.0) + 12*math.Sin(ms/1300.0)
	return map[string]any{"t": t.Format("15:04:05"), "v": math.Round(v)}
}

// streamDemoInitial SSR 初始基线：近 40 个点（每 1.5s 一个），首屏即有一条完整曲线。
func streamDemoInitial() []map[string]any {
	now := time.Now()
	const n = 40
	step := 1500 * time.Millisecond
	out := make([]map[string]any, 0, n)
	for i := n - 1; i >= 0; i-- {
		out = append(out, streamDemoPointAt(now.Add(-time.Duration(i)*step)))
	}
	return out
}

// ConfigStreamChartDemo 滚动流图表演示（/stream-chart-demo）—— 展示 StreamOn（客户端追加流）：
// 服务端后台定时器每 1.5s 经 SSE 推一个合成点 → 客户端 StreamOn 追加到缓冲、裁到 StreamMax →
// 面积图持续左滑滚动（真·追加：不重取、不重渲、不轮询）。
// 注意：demo 用常驻定时器自驱动只为单独展示效果；真实业务的滚动流应由业务事件驱动推点。
func ConfigStreamChartDemo(b *presets.Builder, sseHub presets.SSEHub) {
	// ponytail: 演示用后台定时器，持续向所有 SSE 客户端推点（即使没人看也推）。仅 example 演示可接受。
	if sseHub != nil {
		streamDemoOnce.Do(func() {
			go func() {
				tk := time.NewTicker(1500 * time.Millisecond)
				defer tk.Stop()
				for range tk.C {
					sseHub.Broadcast(strcase.ToCamel(evStreamDemoPoint), streamDemoPointAt(time.Now()))
				}
			}()
		})
	}

	cb := presets.NewCustomPage(b).
		PageTitleFunc(func(*web.EventContext) string { return "滚动流图表 Demo" }).
		Body(func(ctx *web.EventContext) h.HTMLComponent {
			cfg := unovis.ChartConfig{"v": {Label: "实时值", Color: "var(--chart-2)"}}
			card := shadcn.Card(
				shadcn.CardHeader(
					h.Div(h.Text("滚动流（SSE 推点 + 客户端追加）")).Class("text-sm font-medium"),
					h.Div(h.Text("服务端每 1.5 秒经 SSE 推一个点；客户端追加到缓冲、裁到 40 个 → 线持续左滑。不重取、不重渲、不轮询。")).
						Class("text-xs text-muted-foreground mt-1"),
				).Class("p-4 pb-0"),
				shadcn.CardContent(
					unovis.AreaChart(cfg, streamDemoInitial(), "t", "v").
						CurveType(unovis.ChartCurveBasis).
						ShowYAxis(false).ShowGrid(false).ShowLegend(false).ShowTooltip(true).
						StreamOn(evStreamDemoPoint).
						StreamMax(40).
						Class("w-full h-72"),
				).Class("p-3"),
			)
			return h.Div(
				h.H1("滚动流图表").Class("text-2xl font-bold mb-1"),
				h.Div(h.Text("SSE 推送数据点 → 客户端自己追加 → 左滑滚动。真高频实时流范式（StreamOn）。")).
					Class("text-sm text-muted-foreground mb-6"),
				card,
			).Class("container mx-auto p-6")
		})
	b.HandleCustomPage("stream-chart-demo", cb)
}
