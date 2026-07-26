package crud_demo

import (
	"fmt"
	"net/url"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"

	"example/models"

	"gorm.io/gorm"
)

// seedOrders 首次启动插入演示订单：近 7 天每天 4 单、状态循环铺满枚举，
// 让趋势图（按天）与环形图（按状态）都有数据，也提供可 Edit 的行用于演示实时刷新。
func seedOrders(db *gorm.DB) {
	var count int64
	db.Model(&models.Order{}).Count(&count)
	if count > 0 {
		return
	}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	sources := []string{"web", "app", "pos", "phone"}
	var orders []models.Order
	for day := 0; day < 7; day++ {
		for k := 0; k < 4; k++ {
			i := day*4 + k
			orders = append(orders, models.Order{
				// 显式 CreatedAt 回填到对应天（gorm 仅在零值时自动填，故此处会被尊重）
				Model:          gorm.Model{CreatedAt: todayStart.AddDate(0, 0, -day).Add(time.Duration(9+k*3) * time.Hour)},
				Source:         sources[i%len(sources)],
				Status:         models.OrderStatuses[i%len(models.OrderStatuses)],
				PaymentMethod:  "card",
				DeliveryMethod: "pickup",
			})
		}
	}
	db.Create(&orders)
}

// ConfigOrder 配置订单管理模块。sseHub 用于变更后向所有客户端广播，实现订单列表实时刷新。
func ConfigOrder(pb *presets.Builder, db *gorm.DB, sseHub presets.SSEHub) {
	seedOrders(db) // 演示数据：近 7 天订单，供图表与实时刷新演示
	b := pb.Model(&models.Order{}).URIName("orders")

	lb := b.Listing("ID", "CreatedAt", "ConfirmedAt", "PaymentMethod", "DeliveryMethod", "Status", "Source").
		SearchColumns("source").
		SelectableColumns(true).                // 列可勾选显隐（列设置按钮）
		ResizableColumns(true).                 // 列宽可拖拽（localStorage 持久化）
		ReorderableColumns(true).               // 列可拖拽排序（localStorage 持久化）
		DefaultHiddenColumns("DeliveryMethod"). // 默认隐藏该列（可在列设置里打开）
		OrderableFields("ID", "CreatedAt").     // 可排序列（表头出现排序图标）
		PerPage(20)

	// 行级刷新：SSE 推送的「更新」事件只就地补丁对应行的单元格，不整表重渲（消除闪屏）。
	// 新增/删除（行数变化）仍自动回退整表 reload。代价：状态分布图表头在行级更新时不实时刷新，
	// 接受——闪屏体验优先（图表随下次整表 reload/手动刷新更新）。
	lb.RowLevelRefresh(true)

	// 有筛选/搜索时暂停列表 SSE/通知自动刷新：只保留 Updated（自身改动就地补丁），去掉 Created/Deleted
	// —— 外部新单/删单不打乱正在看的筛选视图；清空筛选/搜索后自动恢复全监听。
	lb.PauseRefreshWhenFiltered(true)

	// 三个图表数据 GET 事件函数：组件经 DataURL 拉取，返回 r.Data（{data:[...]} 信封，前端读 .data）。
	const (
		evStatusData  = "orderStatusChartData"  // 环形图：状态分布
		evDailyData   = "orderDailyChartData"   // A 原位更新：近 7 天按天
		evSlidingData = "orderSlidingChartData" // B 滑动：近 30 分钟按分钟
	)
	b.RegisterEventFunc(evStatusData, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.Data = orderStatusData(db)
		return
	})
	b.RegisterEventFunc(evDailyData, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.Data = orderDailyData(db)
		return
	})
	b.RegisterEventFunc(evSlidingData, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.Data = orderSlidingData(db)
		return
	})

	// 列表顶部图表（chart-realtime 三件套，全部 DataURL + RefreshOn 订单增删改事件）：
	// SSE 只推「通知」（订单事件，载荷 {ids}）；三图收到通知 → 重取 DataURL → 平滑更新。
	// 没有订单事件就不动——无定时器、无轮询（与表格行级刷新解耦、各自更新、互不闪屏）。
	base := b.Info().ListingHref()
	refreshEvents := []string{b.NotifModelsUpdated(), b.NotifModelsCreated(), b.NotifModelsDeleted()}
	// ContentHeaderFunc：图表渲染在 Tab 栏下方、筛选器上方（正确位置）。配合 RowLevelRefresh：
	// 编辑=行级补丁、新增/删除=只刷表格 portal —— 整个 compo 都不重渲，故图表（在 ContentHeader）保持挂载、不闪，
	// 只靠 RefreshOn 的 :data 平滑更新。Tab/筛选栏同样不受刷新影响。
	lb.ContentHeaderFunc(orderStatusChartHeader(db,
		base+"?__execute_event__="+evStatusData,
		base+"?__execute_event__="+evDailyData,
		base+"?__execute_event__="+evSlidingData,
		refreshEvents,
	))

	lb.Field("CreatedAt").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		order := obj.(*models.Order)
		v := order.CreatedAt.Local().Format("2006-01-02 15:04:05")
		return h.Text(v)
	}).Label("Date Created")

	lb.Field("ConfirmedAt").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		order := obj.(*models.Order)
		if order.ConfirmedAt != nil {
			return h.Text(order.ConfirmedAt.Local().Format("2006-01-02 15:04:05"))
		}
		return h.Text("")
	}).Label("Check In Date")

	lb.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		statusOptions := make([]shadcn.FilterSelectOption, 0)
		for _, status := range models.OrderStatuses {
			statusOptions = append(statusOptions, shadcn.FilterSelectOption{Value: string(status), Text: string(status)})
		}
		return []*shadcn.FilterItem{
			{
				Key:          "created_at",
				Label:        "Created At",
				ItemType:     shadcn.FilterItemTypeDatetimeRangePicker,
				SQLCondition: `created_at %s ?`,
			},
			{
				Key:          "status",
				Label:        "Status",
				ItemType:     shadcn.FilterItemTypeMultipleSelect,
				SQLCondition: `status %s ?`,
				Options:      statusOptions,
			},
			// 时间 Tab 用的隐藏筛选项（参考 payManage）：不进面板，仅由下方 FilterTab 的 Query 激活。
			// SQLCondition 是整日 created_at 范围的字面量（服务端算好、无占位符），命中即整日过滤。
			{Key: "today", Invisible: true, SQLCondition: orderDayRangeSQL(0)},
			{Key: "yesterday", Invisible: true, SQLCondition: orderDayRangeSQL(1)},
			{Key: "beforeyesterday", Invisible: true, SQLCondition: orderDayRangeSQL(2)},
		}
	})

	// 快捷筛选 Tab：时间维度（今日/昨日/前日，激活隐藏日范围筛选）+ 状态维度（对接 status MultipleSelect，
	// 格式 `status.in=值`）。切 Tab / 增删订单触发表格 portal 刷新时，Tab 栏与筛选栏保持挂载、不重渲、不闪。
	lb.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		return []*presets.FilterTab{
			{ID: "all", Label: "全部", Query: url.Values{"all": []string{"1"}}},
			{ID: "today", Label: "今日", Query: url.Values{"today": []string{"1"}}},
			{ID: "yesterday", Label: "昨日", Query: url.Values{"yesterday": []string{"1"}}},
			{ID: "beforeyesterday", Label: "前日", Query: url.Values{"beforeyesterday": []string{"1"}}},
			{ID: "pending", Label: "待处理", Query: url.Values{"status.in": []string{string(models.OrderStatus_Pending)}}},
			{ID: "paid", Label: "已支付", Query: url.Values{"status.in": []string{string(models.OrderStatus_Paid)}}},
			{ID: "sending", Label: "配送中", Query: url.Values{"status.in": []string{string(models.OrderStatus_Sending)}}},
			{ID: "cancelled", Label: "已取消", Query: url.Values{"status.in": []string{string(models.OrderStatus_Cancelled)}}},
		}
	})

	// 手动刷新按钮（参考 payManage）：静音（SSE 实时关）或想立即拉新时点一下。
	// 走 handReload 发 NotifModelsCreated → 复用 SSE 同款就地重载链路：RowLevelRefresh 下只重渲
	// 表格区 portal（按当前筛选/搜索/翻页 locals 重查），图表监听同键平滑重取 —— 不整页 reload、
	// 不丢筛选、Tab/筛选栏 DOM 保持挂载不闪。
	lb.Action("Reload").ButtonCompFunc(func(ctx *web.EventContext) h.HTMLComponent {
		// 关键：本列表开了 PauseRefreshWhenFiltered —— 有筛选/搜索/Tab 时 Created 监听被掐，
		// 若仍走 handReload（发 Created）刷新按钮会「点了没反应」。故按渲染时是否处于筛选态分流：
		//   · 筛选态 → 直接 PushState(false).Reload() 重载当前列表（保留筛选/翻页、绕过 pause，
		//     不改地址栏防 keep-alive 多标签漂移）；
		//   · 无筛选 → handReload 发 Created（行级刷新、图表不闪，最省）。
		click := web.Plaid().EventFunc("handReload").Go()
		if orderListingFiltered(ctx) {
			click = web.Plaid().PushState(false).Reload().Go()
		}
		return shadcn.Button(shadcn.Icon("refresh-cw").Size(16).Class("mr-1"), h.Text("刷新")).
			Variant(shadcn.ButtonVariantOutline).Size(shadcn.ButtonSizeSm).
			Toast("刷新中...", shadcn.ToasterPositionTopCenter).
			Click(click).Busy()
	})
	b.RegisterEventFunc("handReload", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.Emit(b.NotifModelsCreated())
		return
	})

	// 顶栏「实时」开关：控全局 SSE 静音（localStorage 持久化）。静音时后台推送被丢弃、列表不自动刷新，
	// 想看最新就点上面的「刷新」按钮拉一次；恢复实时后又自动跟随推送更新。
	lb.ToolbarTrailing(func(ctx *web.EventContext) h.HTMLComponent {
		return presets.SSERealtimeToggle("实时")
	})

	// detailing
	b.Detailing(
		&presets.FieldsSection{
			Title: "Basic Information",
			Rows:  [][]string{{"ID", "CreatedAt"}, {"Status", "ConfirmedAt"}, {"PaymentMethod", "DeliveryMethod"}, {"Source"}},
		},
	).Drawer(true)

	eb := b.Editing("Status", "Source", "DeliveryMethod", "PaymentMethod")

	// ===== SSE 实时刷新 =====
	// orders 无 DataScope（无属主隔离），框架的 DataScope 自动 SSE 推送不覆盖它，
	// 故在保存/删除成功后向**所有连接的客户端**广播 Notif*。其它标签/用户的订单列表
	// 经内置 web.Listen 监听器自动整表 reload（compo 级 ReloadAction，已 PushState(false)，
	// keep-alive 多标签安全）。本地保存自身已 emit 刷新，此处补的是跨客户端推送。
	// 事件名经 strcase.ToCamel 规范化，与列表 web.Listen 的监听键对齐（同 publishScopeUpdate）。
	broadcast := func(notifKey, id string) {
		if sseHub == nil {
			return
		}
		sseHub.Broadcast(strcase.ToCamel(notifKey), presets.PayloadModelsUpdated{Ids: []string{id}})
	}
	eb.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj any, id string, ctx *web.EventContext) error {
			created := id == "" // 须在 in() 前判定：保存后 obj 才有 ID
			if err := in(obj, id, ctx); err != nil {
				return err
			}
			newID := id
			if order, ok := obj.(*models.Order); ok {
				newID = fmt.Sprint(order.ID)
			}
			if created {
				broadcast(b.NotifModelsCreated(), newID)
			} else {
				broadcast(b.NotifModelsUpdated(), newID)
			}
			return nil
		}
	})
	eb.WrapDeleteFunc(func(in presets.DeleteFunc) presets.DeleteFunc {
		return func(obj any, id string, ctx *web.EventContext) error {
			if err := in(obj, id, ctx); err != nil {
				return err
			}
			broadcast(b.NotifModelsDeleted(), id)
			return nil
		}
	})
}

// orderListingFiltered 判断当前列表请求是否处于筛选态（有搜索词 / 激活了 FilterTab / 任一 f_ 筛选）。
// 用于「刷新」按钮分流：筛选态下 PauseRefreshWhenFiltered 会掐掉 Created 监听，须改走直接重载。
func orderListingFiltered(ctx *web.EventContext) bool {
	q := ctx.R.URL.Query()
	if q.Get("keyword") != "" || q.Get("active_filter_tab") != "" {
		return true
	}
	for k := range q {
		if len(k) >= 2 && k[:2] == "f_" {
			return true
		}
	}
	return false
}

// orderDayRangeSQL 生成「n 天前那一整天」的 created_at 范围 SQL（n=0 今日 /1 昨日 /2 前日）。
// 时间点由服务端 time.Now() 本地算出、写进字面量（非用户输入，无注入面），供隐藏筛选项直接使用。
func orderDayRangeSQL(n int) string {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -n)
	end := start.AddDate(0, 0, 1)
	const layout = "2006-01-02 15:04:05"
	return fmt.Sprintf("created_at >= '%s' AND created_at < '%s'", start.Format(layout), end.Format(layout))
}

// GetColoredStatus 返回带颜色的状态组件
func GetColoredStatus(status models.OrderStatus) h.HTMLComponent {
	color := models.OrderStatusColorMap[status]
	if color == "" {
		return shadcn.Badge(h.Text(string(status)))
	}
	return shadcn.Badge(h.Text(string(status))).Class("text-white").Attr("style", "background-color: "+color+";")
}
