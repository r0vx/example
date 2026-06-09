package crud_demo

import (
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"

	"example/models"

	"gorm.io/gorm"
)

// ConfigOrder 配置订单管理模块
func ConfigOrder(pb *presets.Builder, db *gorm.DB) {
	b := pb.Model(&models.Order{}).URIName("orders")

	lb := b.Listing("ID", "CreatedAt", "ConfirmedAt", "PaymentMethod", "Status", "Source").
		SearchColumns("source")

	lb.Field("CreatedAt").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		order := obj.(*models.Order)
		v := order.CreatedAt.Local().Format("2006-01-02 15:04:05")
		return h.Text(v)
	}).Label("Date Created")

	lb.Field("ConfirmedAt").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
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

	// detailing
	b.Detailing(
		&presets.FieldsSection{
			Title: "Basic Information",
			Rows:  [][]string{{"ID", "CreatedAt"}, {"Status", "ConfirmedAt"}, {"PaymentMethod", "DeliveryMethod"}, {"Source"}},
		},
	).Drawer(true)

	b.Editing("Status", "Source", "DeliveryMethod", "PaymentMethod")
}

// GetColoredStatus 返回带颜色的状态组件
func GetColoredStatus(status models.OrderStatus) h.HTMLComponent {
	color := models.OrderStatusColorMap[status]
	if color == "" {
		return shadcn.Badge(h.Text(string(status)))
	}
	return shadcn.Badge(h.Text(string(status))).Class("text-white").Attr("style", "background-color: "+color+";")
}