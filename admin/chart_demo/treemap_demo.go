package chart_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// TreemapDemo 树状图演示
type TreemapDemo struct{}

// configTreemapDemo 配置树状图演示
func ConfigTreemapDemo(pb *presets.Builder, db *gorm.DB) {
	b := pb.Model(&TreemapDemo{}).Label("Treemap Demo").URIName("treemap-demo")

	lb := b.Listing()

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		// 双层级数据：区域 -> 省份
		doubleLayerData := []unovis.TreemapDataItem{
			// 华北区域
			{Group: "华北区域", Label: "北京市", Value: 2850, Meta: map[string]any{"orders": "2,850单", "growth": "+15%"}},
			{Group: "华北区域", Label: "天津市", Value: 1420, Meta: map[string]any{"orders": "1,420单", "growth": "+8%"}},
			{Group: "华北区域", Label: "河北省", Value: 1680, Meta: map[string]any{"orders": "1,680单", "growth": "+12%"}},
			{Group: "华北区域", Label: "山西省", Value: 980, Meta: map[string]any{"orders": "980单", "growth": "+5%"}},

			// 华东区域
			{Group: "华东区域", Label: "上海市", Value: 3200, Meta: map[string]any{"orders": "3,200单", "growth": "+18%"}},
			{Group: "华东区域", Label: "江苏省", Value: 2600, Meta: map[string]any{"orders": "2,600单", "growth": "+16%"}},
			{Group: "华东区域", Label: "浙江省", Value: 2400, Meta: map[string]any{"orders": "2,400单", "growth": "+14%"}},
			{Group: "华东区域", Label: "安徽省", Value: 1350, Meta: map[string]any{"orders": "1,350单", "growth": "+10%"}},
			{Group: "华东区域", Label: "福建省", Value: 1580, Meta: map[string]any{"orders": "1,580单", "growth": "+11%"}},
			{Group: "华东区域", Label: "江西省", Value: 920, Meta: map[string]any{"orders": "920单", "growth": "+6%"}},
			{Group: "华东区域", Label: "山东省", Value: 2100, Meta: map[string]any{"orders": "2,100单", "growth": "+13%"}},

			// 华南区域
			{Group: "华南区域", Label: "广东省", Value: 3500, Meta: map[string]any{"orders": "3,500单", "growth": "+20%"}},
			{Group: "华南区域", Label: "广西省", Value: 1150, Meta: map[string]any{"orders": "1,150单", "growth": "+9%"}},
			{Group: "华南区域", Label: "海南省", Value: 650, Meta: map[string]any{"orders": "650单", "growth": "+7%"}},

			// 华中区域
			{Group: "华中区域", Label: "湖北省", Value: 1800, Meta: map[string]any{"orders": "1,800单", "growth": "+12%"}},
			{Group: "华中区域", Label: "湖南省", Value: 1620, Meta: map[string]any{"orders": "1,620单", "growth": "+11%"}},
			{Group: "华中区域", Label: "河南省", Value: 1950, Meta: map[string]any{"orders": "1,950单", "growth": "+10%"}},

			// 西南区域
			{Group: "西南区域", Label: "重庆市", Value: 1750, Meta: map[string]any{"orders": "1,750单", "growth": "+13%"}},
			{Group: "西南区域", Label: "四川省", Value: 2200, Meta: map[string]any{"orders": "2,200单", "growth": "+14%"}},
			{Group: "西南区域", Label: "贵州省", Value: 850, Meta: map[string]any{"orders": "850单", "growth": "+8%"}},
			{Group: "西南区域", Label: "云南省", Value: 1100, Meta: map[string]any{"orders": "1,100单", "growth": "+9%"}},

			// 西北区域
			{Group: "西北区域", Label: "陕西省", Value: 1320, Meta: map[string]any{"orders": "1,320单", "growth": "+10%"}},
			{Group: "西北区域", Label: "甘肃省", Value: 680, Meta: map[string]any{"orders": "680单", "growth": "+5%"}},
			{Group: "西北区域", Label: "宁夏省", Value: 420, Meta: map[string]any{"orders": "420单", "growth": "+4%"}},
			{Group: "西北区域", Label: "新疆省", Value: 780, Meta: map[string]any{"orders": "780单", "growth": "+6%"}},

			// 东北区域
			{Group: "东北区域", Label: "辽宁省", Value: 1450, Meta: map[string]any{"orders": "1,450单", "growth": "+9%"}},
			{Group: "东北区域", Label: "吉林省", Value: 820, Meta: map[string]any{"orders": "820单", "growth": "+6%"}},
			{Group: "东北区域", Label: "黑龙江省", Value: 950, Meta: map[string]any{"orders": "950单", "growth": "+7%"}},
		}

		// 双层级树状图配置
		doubleLayerConfig := unovis.TreemapConfig{
			"华北区域": {Color: "hsl(217 91% 60%)", Label: "华北区域"},
			"华东区域": {Color: "hsl(142 76% 73%)", Label: "华东区域"},
			"华南区域": {Color: "hsl(48 96% 89%)", Label: "华南区域"},
			"华中区域": {Color: "hsl(351 83% 82%)", Label: "华中区域"},
			"西南区域": {Color: "hsl(280 60% 60%)", Label: "西南区域"},
			"西北区域": {Color: "hsl(20 90% 65%)", Label: "西北区域"},
			"东北区域": {Color: "hsl(200 80% 70%)", Label: "东北区域"},
		}

		// 创建双层级树状图
		doubleLayerTreemap := unovis.Treemap().
			Data(doubleLayerData).
			Config(doubleLayerConfig).
			Layers("double").
			TilePadding(3).
			TileBorderRadius(6).
			Class("w-full h-[600px]")

		// 单层级数据：产品类别订单统计
		singleLayerData := []unovis.TreemapDataItem{
			{Label: "电子产品", Value: 8500, Color: "hsl(217 91% 60%)", Meta: map[string]any{"orders": "8,500单", "revenue": "¥425万"}},
			{Label: "服装鞋帽", Value: 6200, Color: "hsl(142 76% 73%)", Meta: map[string]any{"orders": "6,200单", "revenue": "¥186万"}},
			{Label: "家居用品", Value: 4800, Color: "hsl(48 96% 89%)", Meta: map[string]any{"orders": "4,800单", "revenue": "¥240万"}},
			{Label: "食品饮料", Value: 5500, Color: "hsl(351 83% 82%)", Meta: map[string]any{"orders": "5,500单", "revenue": "¥165万"}},
			{Label: "美妆护肤", Value: 3900, Color: "hsl(280 60% 60%)", Meta: map[string]any{"orders": "3,900单", "revenue": "¥195万"}},
			{Label: "图书音像", Value: 2800, Color: "hsl(20 90% 65%)", Meta: map[string]any{"orders": "2,800单", "revenue": "¥84万"}},
			{Label: "运动户外", Value: 3200, Color: "hsl(200 80% 70%)", Meta: map[string]any{"orders": "3,200单", "revenue": "¥160万"}},
			{Label: "母婴用品", Value: 2600, Color: "hsl(320 70% 75%)", Meta: map[string]any{"orders": "2,600单", "revenue": "¥130万"}},
		}

		// 创建单层级树状图
		singleLayerTreemap := unovis.Treemap().
			Data(singleLayerData).
			Layers("single").
			TilePadding(2).
			TileBorderRadius(4).
			Class("w-full h-[400px]")

		body := h.Div(
			h.H1("Treemap Demo").Class("text-2xl font-bold mb-4"),
			h.P(h.Text("树状图（矩形树图）：使用方块大小表示数值占比，适合展示层级数据和占比关系。")).Class("text-muted-foreground mb-6"),

			// 双层级树状图 - 省份订单统计
			shadcn.Card(
				shadcn.CardHeader(
					shadcn.CardTitle(h.Text("全国省份订单统计（双层级）")),
					shadcn.CardDescription(h.Text("按区域-省份展示订单分布，方块大小代表订单数量")),
				),
				shadcn.CardContent(
					doubleLayerTreemap,
				),
			).Class("mb-6"),

			// 单层级树状图 - 产品类别
			shadcn.Card(
				shadcn.CardHeader(
					shadcn.CardTitle(h.Text("产品类别订单统计（单层级）")),
					shadcn.CardDescription(h.Text("展示各产品类别的订单占比，方块大小代表订单数量")),
				),
				shadcn.CardContent(
					singleLayerTreemap,
				),
			).Class("mb-6"),

			// 说明文档
			shadcn.Card(
				shadcn.CardHeader(
					shadcn.CardTitle(h.Text("功能说明")),
				),
				shadcn.CardContent(
					h.Ul(
						h.Li(h.Text("支持单层级和双层级数据展示")),
						h.Li(h.Text("方块大小自动根据数值比例计算")),
						h.Li(h.Text("支持自定义颜色，同组数据自动亮度变化")),
						h.Li(h.Text("鼠标悬停显示详细信息（tooltip）")),
						h.Li(h.Text("支持自定义方块间距和圆角")),
						h.Li(h.Text("适合展示：省份统计、产品分类、销售占比等数据")),
						h.Li(h.Text("双层级模式：第一层按区域分组，第二层显示省份")),
					).Class("list-disc list-inside space-y-2 text-sm text-muted-foreground"),
				),
			),
		).Class("container mx-auto p-4")

		r.Body = body
		r.PageTitle = "Treemap Demo"

		return
	})
}
