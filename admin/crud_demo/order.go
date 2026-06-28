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

	lb := b.Listing("ID", "CreatedAt", "ConfirmedAt", "PaymentMethod", "Status", "Source").
		SearchColumns("source")

	// 行级刷新：SSE 推送的「更新」事件只就地补丁对应行的单元格，不整表重渲（消除闪屏）。
	// 新增/删除（行数变化）仍自动回退整表 reload。代价：状态分布图表头在行级更新时不实时刷新，
	// 接受——闪屏体验优先（图表随下次整表 reload/手动刷新更新）。
	lb.RowLevelRefresh(true)

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
		}
	})

	// 快捷筛选 Tab（按状态）：Query 对接 status MultipleSelect 筛选项，格式 `status.in=值`（逗号分隔多值）。
	// 用于测试：切 Tab / 增删订单触发表格 portal 刷新时，Tab 栏与筛选栏保持挂载、不重渲、不闪。
	lb.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		return []*presets.FilterTab{
			{ID: "all", Label: "全部", Query: url.Values{"all": []string{"1"}}},
			{ID: "pending", Label: "待处理", Query: url.Values{"status.in": []string{string(models.OrderStatus_Pending)}}},
			{ID: "paid", Label: "已支付", Query: url.Values{"status.in": []string{string(models.OrderStatus_Paid)}}},
			{ID: "sending", Label: "配送中", Query: url.Values{"status.in": []string{string(models.OrderStatus_Sending)}}},
			{ID: "cancelled", Label: "已取消", Query: url.Values{"status.in": []string{string(models.OrderStatus_Cancelled)}}},
		}
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

// GetColoredStatus 返回带颜色的状态组件
func GetColoredStatus(status models.OrderStatus) h.HTMLComponent {
	color := models.OrderStatusColorMap[status]
	if color == "" {
		return shadcn.Badge(h.Text(string(status)))
	}
	return shadcn.Badge(h.Text(string(status))).Class("text-white").Attr("style", "background-color: "+color+";")
}
