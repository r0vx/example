package admin

import (
	"fmt"
	"strconv"
	"time"

	"example/models"

	"github.com/r0vx/admin/dashboard"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// configureDashboard 配置首页仪表盘，使用 dashboard 包声明式构建
func configureDashboard(db *gorm.DB) *dashboard.DashboardBuilder {
	return dashboard.New(
		// 统计卡片：商品总数
		dashboard.Stat("商品总数").
			Value(func(ctx *web.EventContext) (string, error) {
				var count int64
				if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
					return "", fmt.Errorf("查询商品数量失败: %w", err)
				}
				return strconv.FormatInt(count, 10), nil
			}).
			Icon("package").
			Description("全部商品"),

		// 统计卡片：订单总数
		dashboard.Stat("订单总数").
			Value(func(ctx *web.EventContext) (string, error) {
				var count int64
				if err := db.Model(&models.Order{}).Count(&count).Error; err != nil {
					return "", fmt.Errorf("查询订单数量失败: %w", err)
				}
				return strconv.FormatInt(count, 10), nil
			}).
			Icon("shopping-cart").
			Description("全部订单"),

		// 统计卡片：待处理订单
		// 演示 Trend 显式方向 API（2026-05 新增）：
		//   - .Trend("+3") 自动按首字符判断为上升 → 绿色 (text-emerald-600)
		//   - .Trend("...").TrendUp() 显式标记上升（即使字符串没有 + 号）
		//   - .Trend("0%").TrendNeutral() 0%/持平 显式中性 → 默认灰色
		//   - .Trend("-12%").TrendDown() 下降 → 红色 (text-rose-600)
		dashboard.Stat("待处理订单").
			Value(func(ctx *web.EventContext) (string, error) {
				var count int64
				if err := db.Model(&models.Order{}).Where("status = ?", "pending").Count(&count).Error; err != nil {
					return "", fmt.Errorf("查询待处理订单失败: %w", err)
				}
				return strconv.FormatInt(count, 10), nil
			}).
			Icon("clock").
			Trend("+3"). // 自动绿色（首字符 '+'）
			Description("较昨日"),

		// 演示卡片：本周转化率（中性显式标记）
		dashboard.Stat("转化率").
			Value(func(ctx *web.EventContext) (string, error) {
				return "0.0%", nil
			}).
			Icon("trending-up").
			Trend("0.0%").
			TrendNeutral(). // 显式中性，避免被误判
			Description("本周持平"),

		// 图表：订单状态分布
		dashboard.Chart("订单状态分布").
			Description("各状态订单数量统计").
			Type(unovis.ChartTypeBar).
			Config(unovis.ChartConfig{
				"count": {Label: "订单数", Color: "var(--chart-1)"},
			}).
			Data(func(ctx *web.EventContext) (any, error) {
				type StatusData struct {
					Status string `json:"status"`
					Count  int    `json:"count"`
				}
				type statusResult struct {
					Status string
					Count  int
				}
				var results []statusResult
				if err := db.Model(&models.Order{}).Select("status, count(*) as count").Group("status").Scan(&results).Error; err != nil {
					return nil, fmt.Errorf("查询订单状态失败: %w", err)
				}
				statusCount := make(map[models.OrderStatus]int)
				for _, r := range results {
					statusCount[models.OrderStatus(r.Status)] = r.Count
				}
				var data []StatusData
				for _, s := range models.OrderStatuses {
					data = append(data, StatusData{Status: string(s), Count: statusCount[s]})
				}
				return data, nil
			}).
			XKey("status").
			YKeys("count").
			ColSpan(2),

		// 图表：访问量趋势（模拟数据）
		dashboard.Chart("访问量趋势").
			Description("最近 30 天访问量").
			Type(unovis.ChartTypeLine).
			Config(unovis.ChartConfig{
				"desktop": {Label: "Desktop", Color: "var(--chart-2)"},
				"mobile":  {Label: "Mobile", Color: "var(--chart-1)"},
			}).
			Data(func(ctx *web.EventContext) (any, error) {
				type VisitorData struct {
					Date    string `json:"date"`
					Desktop int    `json:"desktop"`
					Mobile  int    `json:"mobile"`
				}
				var data []VisitorData
				now := time.Now()
				for i := 29; i >= 0; i-- {
					day := now.AddDate(0, 0, -i)
					data = append(data, VisitorData{
						Date:    day.Format("01-02"),
						Desktop: 100 + (day.Day()*17+int(day.Month())*23)%400,
						Mobile:  80 + (day.Day()*13+int(day.Month())*31)%350,
					})
				}
				return data, nil
			}).
			XKey("date").
			YKeys("desktop", "mobile").
			ShowLegend(true).
			ColSpan(3),

		// 自定义组件：实时时钟（每 2 秒局部自刷新，演示 RefreshInterval 声明式 API）
		dashboard.Custom("autorefresh-clock").
			Body(func(ctx *web.EventContext) (h.HTMLComponent, error) {
				return shadcn.Card(
					shadcn.CardHeader(shadcn.CardTitle(h.Text("实时时钟"))),
					shadcn.CardContent(
						h.Div().Class("text-2xl font-mono").Text(time.Now().Format("15:04:05")),
					),
				), nil
			}).
			RefreshInterval(2).
			ColSpan(1),

		// 自定义组件：公告
		dashboard.Custom("announcement").
			Body(func(ctx *web.EventContext) (h.HTMLComponent, error) {
				return shadcn.Card(
					shadcn.CardHeader(
						shadcn.CardTitle(h.Text("系统公告")),
					),
					shadcn.CardContent(
						h.Div(
							h.Div().Class("text-sm").Text("Dashboard Widget 系统已上线，支持 Stat / Chart / Table / Custom 四种组件类型。"),
							h.Div().Class("text-xs text-muted-foreground mt-2").Text("2026-04-07"),
						),
					),
				), nil
			}).
			ColSpan(3),
	).Columns(3).Title("r0vx Admin Dashboard").
		Store(dashboardLayoutStore).
		UserID(func(ctx *web.EventContext) string {
			// 简化：使用固定用户 ID（生产环境应从 session 获取）
			return "default"
		})
}

// dashboardLayoutStore 全局布局存储（演示用内存存储）
var dashboardLayoutStore = dashboard.NewMemoryLayoutStore()
