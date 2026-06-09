package chart_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// ScatterPlotDemo 散点图演示
type ScatterPlotDemo struct{}

// configScatterPlotDemo 配置散点图演示
func ConfigScatterPlotDemo(pb *presets.Builder, db *gorm.DB) {
	b := pb.Model(&ScatterPlotDemo{}).Label("Scatter Plot Demo").URIName("scatter-plot-demo")

	lb := b.Listing()

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		// 构建散点图数据（模拟不同产品的销售数据）
		scatterData := []unovis.ScatterDataItem{
			// 电子产品系列
			{X: 10, Y: 30, Size: 150, Label: "笔记本电脑", Color: "hsl(217 91% 60%)", Meta: map[string]any{"category": "电子产品", "units": "150台"}},
			{X: 15, Y: 45, Size: 230, Label: "智能手机", Color: "hsl(217 91% 60%)", Meta: map[string]any{"category": "电子产品", "units": "230台"}},
			{X: 20, Y: 60, Size: 310, Label: "平板电脑", Color: "hsl(217 91% 60%)", Meta: map[string]any{"category": "电子产品", "units": "310台"}},
			{X: 25, Y: 55, Size: 180, Label: "耳机", Color: "hsl(217 91% 60%)", Meta: map[string]any{"category": "电子产品", "units": "180台"}},

			// 家居用品系列
			{X: 30, Y: 40, Size: 120, Label: "沙发", Color: "hsl(142 76% 73%)", Meta: map[string]any{"category": "家居用品", "units": "120件"}},
			{X: 35, Y: 50, Size: 200, Label: "床垫", Color: "hsl(142 76% 73%)", Meta: map[string]any{"category": "家居用品", "units": "200件"}},
			{X: 40, Y: 45, Size: 160, Label: "书架", Color: "hsl(142 76% 73%)", Meta: map[string]any{"category": "家居用品", "units": "160件"}},

			// 服装系列
			{X: 45, Y: 70, Size: 450, Label: "T恤", Color: "hsl(48 96% 89%)", Meta: map[string]any{"category": "服装", "units": "450件"}},
			{X: 50, Y: 65, Size: 380, Label: "牛仔裤", Color: "hsl(48 96% 89%)", Meta: map[string]any{"category": "服装", "units": "380件"}},
			{X: 55, Y: 75, Size: 520, Label: "运动鞋", Color: "hsl(48 96% 89%)", Meta: map[string]any{"category": "服装", "units": "520件"}},
			{X: 60, Y: 68, Size: 420, Label: "外套", Color: "hsl(48 96% 89%)", Meta: map[string]any{"category": "服装", "units": "420件"}},

			// 食品系列
			{X: 65, Y: 35, Size: 280, Label: "零食礼包", Color: "hsl(351 83% 82%)", Meta: map[string]any{"category": "食品", "units": "280份"}},
			{X: 70, Y: 42, Size: 340, Label: "咖啡豆", Color: "hsl(351 83% 82%)", Meta: map[string]any{"category": "食品", "units": "340包"}},
			{X: 75, Y: 38, Size: 290, Label: "茶叶", Color: "hsl(351 83% 82%)", Meta: map[string]any{"category": "食品", "units": "290盒"}},

			// 运动器材系列
			{X: 80, Y: 55, Size: 90, Label: "跑步机", Color: "hsl(280 60% 60%)", Meta: map[string]any{"category": "运动器材", "units": "90台"}},
			{X: 85, Y: 62, Size: 150, Label: "健身车", Color: "hsl(280 60% 60%)", Meta: map[string]any{"category": "运动器材", "units": "150台"}},
			{X: 90, Y: 58, Size: 110, Label: "哑铃组", Color: "hsl(280 60% 60%)", Meta: map[string]any{"category": "运动器材", "units": "110套"}},
		}

		// 配置散点图
		scatterConfig := unovis.ScatterPlotConfig{
			"电子产品": {
				Label: "电子产品",
				Color: "hsl(217 91% 60%)",
			},
			"家居用品": {
				Label: "家居用品",
				Color: "hsl(142 76% 73%)",
			},
			"服装": {
				Label: "服装",
				Color: "hsl(48 96% 89%)",
			},
			"食品": {
				Label: "食品",
				Color: "hsl(351 83% 82%)",
			},
			"运动器材": {
				Label: "运动器材",
				Color: "hsl(280 60% 60%)",
			},
		}

		// 创建散点图（带气泡大小）
		scatterPlot := unovis.ScatterPlot().
			Data(scatterData).
			Config(scatterConfig).
			XLabel("广告投入（千元）").
			YLabel("销售利润率（%）").
			SizeRange(5, 40).
			Class("w-full h-[600px]")

		// 创建固定大小散点图（不使用 Size 字段）
		fixedSizeData := []unovis.ScatterDataItem{
			{X: 10, Y: 20, Label: "数据点 1", Color: "hsl(217 91% 60%)"},
			{X: 20, Y: 30, Label: "数据点 2", Color: "hsl(142 76% 73%)"},
			{X: 30, Y: 25, Label: "数据点 3", Color: "hsl(48 96% 89%)"},
			{X: 40, Y: 35, Label: "数据点 4", Color: "hsl(351 83% 82%)"},
			{X: 50, Y: 40, Label: "数据点 5", Color: "hsl(280 60% 60%)"},
			{X: 60, Y: 38, Label: "数据点 6", Color: "hsl(217 91% 60%)"},
			{X: 70, Y: 45, Label: "数据点 7", Color: "hsl(142 76% 73%)"},
		}

		fixedScatterPlot := unovis.ScatterPlot().
			Data(fixedSizeData).
			XLabel("X 值").
			YLabel("Y 值").
			ShowLabels(true).
			Class("w-full h-[400px]")

		body := h.Div(
			h.H1("Scatter Plot Demo").Class("text-2xl font-bold mb-4"),
			h.P(h.Text("展示散点图（气泡图）的使用，支持三维数据可视化（X 轴、Y 轴、点大小）。")).Class("text-muted-foreground mb-6"),

			// 气泡图（带大小变化）
			shadcn.Card(
				shadcn.CardHeader(
					shadcn.CardTitle(h.Text("产品销售分析（气泡图）")),
					shadcn.CardDescription(h.Text("展示不同产品类别的广告投入、销售利润率和销量关系，点的大小代表销量")),
				),
				shadcn.CardContent(
					scatterPlot,
				),
			).Class("mb-6"),

			// 固定大小散点图
			shadcn.Card(
				shadcn.CardHeader(
					shadcn.CardTitle(h.Text("标准散点图（固定大小）")),
					shadcn.CardDescription(h.Text("不使用 Size 字段，所有点大小一致，适合简单的二维数据展示")),
				),
				shadcn.CardContent(
					fixedScatterPlot,
				),
			).Class("mb-6"),

			// 说明文档
			shadcn.Card(
				shadcn.CardHeader(
					shadcn.CardTitle(h.Text("功能说明")),
				),
				shadcn.CardContent(
					h.Ul(
						h.Li(h.Text("支持三维数据可视化：X 轴、Y 轴和点大小（气泡图）")),
						h.Li(h.Text("自动将 Size 数值映射到像素范围（默认 5-30px，可自定义）")),
						h.Li(h.Text("支持固定大小模式：不提供 Size 字段时，所有点大小一致")),
						h.Li(h.Text("支持自定义点颜色，自动循环使用主题色")),
						h.Li(h.Text("支持显示点标签（ShowLabels）")),
						h.Li(h.Text("鼠标悬停显示详细信息（tooltip）")),
						h.Li(h.Text("支持自定义 X/Y 轴标签")),
						h.Li(h.Text("可通过 Meta 字段添加额外数据，tooltip 会自动显示")),
					).Class("list-disc list-inside space-y-2 text-sm text-muted-foreground"),
				),
			),
		).Class("container mx-auto p-4")

		r.Body = body
		r.PageTitle = "Scatter Plot Demo"

		return
	})
}
