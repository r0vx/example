package ui_demo

import (
	"fmt"
	"time"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// configNotificationCenter 配置顶部工具栏的通知中心（铃铛图标 + 弹出面板）
//
// 演示 presets.Builder.NotificationFunc 的用法：
//   - countFunc: 返回未读通知数量（显示在铃铛图标的红色 Badge 上）
//   - contentFunc: 返回通知面板的 HTML 内容（点击铃铛弹出的 Popover）
func ConfigNotificationCenter(b *presets.Builder, db *gorm.DB) {
	// contentFunc — 渲染通知面板内容
	contentFunc := func(ctx *web.EventContext) h.HTMLComponent {
		// 查询最近 7 天内新增的订单
		var recentOrders []models.Order
		since := time.Now().AddDate(0, 0, -7)
		db.Where("created_at > ?", since).
			Order("created_at desc").
			Limit(5).
			Find(&recentOrders)

		// 构建通知列表
		var items []h.HTMLComponent
		for _, order := range recentOrders {
			items = append(items, h.Div().Class("flex items-start gap-3 p-3 hover:bg-muted/50 transition-colors").Children(
				// 左侧图标
				h.Div().Class("flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10").Children(
					shadcn.Icon(shadcn.IconShoppingCart).Size(14).Class("text-primary"),
				),
				// 右侧内容
				h.Div().Class("flex-1 space-y-1").Children(
					h.P().Class("text-sm font-medium leading-none").Text(
						fmt.Sprintf("新订单 #%d", order.ID),
					),
					h.P().Class("text-xs text-muted-foreground").Text(
						order.CreatedAt.Format("01-02 15:04"),
					),
				),
			))
		}

		if len(items) == 0 {
			items = append(items, h.Div().Class("p-6 text-center text-sm text-muted-foreground").Children(
				h.Text("暂无新通知"),
			))
		}

		return h.Div().Children(
			// 标题栏
			h.Div().Class("flex items-center justify-between border-b px-4 py-3").Children(
				h.Span("通知").Class("text-sm font-semibold"),
				h.Span(fmt.Sprintf("最近 7 天 · %d 条", len(recentOrders))).Class("text-xs text-muted-foreground"),
			),
			// 通知列表
			h.Div().Class("max-h-80 overflow-y-auto divide-y").Children(items...),
		)
	}

	// countFunc — 返回未读通知数量
	countFunc := func(ctx *web.EventContext) int {
		var count int64
		since := time.Now().AddDate(0, 0, -7)
		db.Model(&models.Order{}).Where("created_at > ?", since).Count(&count)
		return int(count)
	}

	b.NotificationFunc(contentFunc, countFunc)
}
