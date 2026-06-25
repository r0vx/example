package ui_demo

import (
	"strconv"
	"time"

	"example/models"

	"github.com/r0vx/admin/notification"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/login"
	"github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// notificationCurrentUserID 从 request 取登录用户 ID（数字转字符串），用于按用户隔离通知
func notificationCurrentUserID(ctx *web.EventContext) string {
	u, _ := login.GetCurrentUser(ctx.R).(*models.User)
	if u == nil {
		return ""
	}
	return strconv.FormatUint(uint64(u.ID), 10)
}

// configNotificationDemo 演示 notification 模块全链路集成：
//  1. AutoMigrate 通知表
//  2. 注册 Toast + Database 双通道 Notifier
//  3. 接通 Bell UI（presets NotificationFunc）
//  4. 注册 mark_read / mark_all_read 事件
//  5. 提供发送通知的 demo 按钮
func ConfigNotificationDemo(b *presets.Builder, db *gorm.DB, sseHub notification.Pusher) {
	if err := db.AutoMigrate(&notification.Model{}); err != nil {
		panic(err)
	}

	notifier := notification.New(
		notification.NewToastChannel(),
		notification.NewDatabaseChannel(db, notificationCurrentUserID),
		// SSE：把通知实时推给接收者，触发其铃铛刷新（事件名与 presets 铃铛监听一致）
		notification.NewSSEChannel(sseHub, notificationCurrentUserID, presets.NotifNotificationUpdated),
	)

	// 注册通知面板 i18n（en/zh），否则面板回退英文文案
	notification.RegisterI18n(b)

	// 接通顶栏 Bell：未读数 + 内容渲染
	b.NotificationFunc(
		notification.MakeContentFunc(db, notificationCurrentUserID, 10),
		notification.MakeCountFunc(db, notificationCurrentUserID),
	)
	notification.RegisterEvents(b.GetWebBuilder(), db, notificationCurrentUserID)

	// demo 按钮事件：发一条带 Action 的通知
	b.GetWebBuilder().RegisterEventFunc("demo_notification_send", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		typ := ctx.R.FormValue("type")
		if typ == "" {
			typ = string(notification.TypeSuccess)
		}
		notifier.Send(ctx, &r, &notification.Notification{
			Type:  notification.Type(typ),
			Title: "演示通知 " + time.Now().Format("15:04:05"),
			Body:  "这是从 demo 按钮触发的通知；同时存数据库 + 弹 toast",
			Action: &notification.Action{
				Label: "查看商品",
				URL:   "/products",
			},
		})
		return
	})

	// 演示页面：4 个按钮分别发 4 种类型通知
	cb := presets.NewCustomPage(b).
		PageTitleFunc(func(*web.EventContext) string { return "通知中心 Demo" }).
		Body(func(ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.H1("Notification Demo").Class("text-2xl font-bold mb-4"),
				h.Div(h.Text("点击按钮发送对应类型通知，同时落库（铃铛红点 +1）+ 弹 toast")).
					Class("text-sm text-muted-foreground mb-6"),
				h.Div(
					notifDemoButton("Success", string(notification.TypeSuccess), shadcn.ButtonVariantDefault),
					notifDemoButton("Error", string(notification.TypeError), shadcn.ButtonVariantDestructive),
					notifDemoButton("Warning", string(notification.TypeWarning), shadcn.ButtonVariantOutline),
					notifDemoButton("Info", string(notification.TypeInfo), shadcn.ButtonVariantSecondary),
				).Class("flex gap-2 mb-6"),
				h.Div(h.Text("点击右上角 🔔 查看通知列表（默认显示未读）")).
					Class("text-sm text-muted-foreground"),
			).Class("p-6")
		})
	b.HandleCustomPage("notification-demo", cb)
}

// notifDemoButton 渲染一个发送通知按钮
func notifDemoButton(label, typ string, variant shadcn.ButtonVariant) h.HTMLComponent {
	return shadcn.Button(h.Text(label)).
		Variant(variant).
		Attr("@click", web.Plaid().EventFunc("demo_notification_send").Query("type", typ).Go())
}
