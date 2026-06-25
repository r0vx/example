package crud_demo

import (
	"fmt"
	"strconv"

	"github.com/r0vx/admin/publish"

	"example/models"

	"github.com/r0vx/admin/activity"
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// Customer 列表「改会员卡」事件名
const (
	eventCustomerEditCardOpen = "customer_editCard_open"
	eventCustomerEditCardSave = "customer_editCard_save"
)

// ConfigCategory 配置分类管理模块
func ConfigCategory(b *presets.Builder, db *gorm.DB, publisher *publish.Builder) *presets.ModelBuilder {
	mb := b.Model(&models.Category{}).Use(publisher)
	mb.Listing("ID", "Name", "Position", "Status").SearchColumns("name")
	mb.Editing("Name", "Products", "Position")
	return mb
}

// ConfigSubCategory 演示现有 SortBuilder：列表工具栏「Sort」按钮 → 独立排序页 → 拖动改 Position
func ConfigSubCategory(b *presets.Builder, db *gorm.DB) *presets.ModelBuilder {
	mb := b.Model(&models.SubCategory{}).URIName("sub-categories").Label("二级分类（排序演示）")
	mb.Listing("ID", "Name", "Position").SearchColumns("name").SelectableColumns(true)
	mb.Editing("Name", "Position")
	// 启用拖拽排序：工具栏图标按钮，点击弹排序 Dialog，保存后按新顺序写回 Position。
	// 不设 Label → 标题/tooltip 走 i18n 默认（中文「排序」/英文 Sort）。
	mb.Sorting(db).PositionField("Position").Display("Name").InstallToolbarButton()
	return mb
}

// ConfigNestedFieldDemo 配置嵌套字段演示（Customer → Address → Phone / MembershipCard）
func ConfigNestedFieldDemo(b *presets.Builder, db *gorm.DB, ab *activity.Builder) {
	mb := b.Model(&models.Customer{}).URIName("customers")

	// 确保 MembershipCard 已注册到 activity（idempotent：已注册则返回现有），供下面 OnEdit 记 diff。
	cardAMB := ab.RegisterModel(&models.MembershipCard{})

	lb := mb.Listing("ID", "Name").SearchColumns("name")

	// 行内「改会员卡」按钮：从 Customer 行打开 dialog，改这个客户关联的会员卡（按 customer_id 查）
	lb.RowMenu().RowMenuItem("EditCard").ComponentFunc(func(obj any, id string, ctx *web.EventContext) h.HTMLComponent {
		onClick := web.Plaid().EventFunc(eventCustomerEditCardOpen).Query(presets.ParamID, id).Go()
		return RowMenuItem("改会员卡").SetIcon("credit-card").SetOnclick(onClick)
	})

	mb.Editing(
		&presets.FieldsSection{
			Title: "Basic Info",
			Rows:  [][]string{{"Name", "Titile"}, {"Addresses"}},
		},
		&presets.FieldsSection{
			Title: "Membership",
			Rows:  [][]string{{"MembershipCard"}},
		},
	)
	mb.Detailing("ID", "Name", "Addresses", "MembershipCard")

	// 事件 A：打开 dialog —— 按 customer_id 查该客户的会员卡，预填卡号
	mb.RegisterEventFunc(eventCustomerEditCardOpen, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		customerID := ctx.R.FormValue(presets.ParamID)
		var card models.MembershipCard
		// 无卡则给一张空卡（CustomerID 预填），保存时 db.Save 自动建
		cid, _ := strconv.Atoi(customerID)
		if e := db.Where("customer_id = ?", customerID).First(&card).Error; e != nil {
			card = models.MembershipCard{CustomerID: uint(cid)}
		}
		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: presets.DialogPortalName,
			Body: web.Scope(
				Dialog(
					DialogContent(
						DialogHeader(DialogTitle(h.Text("修改会员卡"))),
						h.Div(
							Input().Type("number").Label("卡号").
								Attr(web.VField("card_number", fmt.Sprint(card.Number))...),
						).Class("py-4"),
						DialogFooter(
							Button(h.Text("取消")).Variant(ButtonVariantOutline).
								Attr("@click", "locals.show = false"),
							Button(h.Text("保存")).
								Attr("@click", "locals.show = false;"+web.Plaid().
									EventFunc(eventCustomerEditCardSave).
									Query(presets.ParamID, customerID).Go()),
						),
					).Class(presets.DialogSizeXs),
				).Attr(":open", "locals.show").
					OnUpdateOpen("locals.show = $event"),
			).VSlot("{ locals }").Init("{show: true}"),
		})
		return
	})

	// 事件 B：保存 —— 查 old → 改 → 存 → cardAMB.OnEdit 自动 diff 记日志（记到 MembershipCard）
	mb.RegisterEventFunc(eventCustomerEditCardSave, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		customerID := ctx.R.FormValue(presets.ParamID)
		newNumber, _ := strconv.Atoi(ctx.R.FormValue("card_number"))
		cid, _ := strconv.Atoi(customerID)

		var old models.MembershipCard
		isNew := db.Where("customer_id = ?", customerID).First(&old).Error != nil
		if isNew {
			old = models.MembershipCard{CustomerID: uint(cid)}
		}
		updated := old
		updated.Number = newNumber
		if err = db.Save(&updated).Error; err != nil {
			presets.ShowError(&r, "保存失败")
			return r, nil
		}
		// 新建走 OnCreate、改走 OnEdit（操作人来自 ctx.R.Context()）
		if isNew {
			_, _ = cardAMB.OnCreate(ctx.R.Context(), &updated)
		} else {
			_, _ = cardAMB.OnEdit(ctx.R.Context(), &old, &updated)
		}
		presets.ShowSuccess(&r, "会员卡已更新")
		r.Emit(presets.NotifModelsUpdated(&models.Customer{}), presets.PayloadModelsUpdated{Ids: []string{customerID}})
		return
	})
}

// ConfigProject 配置项目管理模块
func ConfigProject(pb *presets.Builder, db *gorm.DB) {
	mb := pb.Model(&models.Project{}).URIName("projects")
	mb.Listing("ID", "Name", "Status", "Featured").SearchColumns("name", "description")

	ConfigTask(pb, db)

	mb.Editing(
		&presets.FieldsSection{
			Title: "Basic Info",
			Rows:  [][]string{{"Name", "Status"}, {"Description"}, {"Featured", "Avatar"}},
		},
		&presets.FieldsSection{
			Title: "Tasks",
			Rows:  [][]string{{"Tasks"}},
		},
	)
}

// ConfigTask 配置任务子表模块
func ConfigTask(pb *presets.Builder, _ *gorm.DB) {
	mb := pb.Model(&models.Task{}).URIName("tasks")
	mb.Listing("ID", "Name", "Priority", "Status", "Assignee")
	mb.Editing("Name", "Description", "Priority", "Status", "Assignee")
}
