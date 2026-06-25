package crud_demo

import (
	"context"
	"fmt"
	"strconv"

	"example/models"

	"github.com/r0vx/admin/activity"
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// 自定义事件名（模型级，挂在 mb 上）
const (
	eventCardEditNumberOpen = "card_editNumber_open" // 打开「改卡号」dialog
	eventCardEditNumberSave = "card_editNumber_save" // 保存卡号（手动 OnEdit 记 diff）
	eventCardRecharge       = "card_recharge"        // 充值（手动 Log 记自定义 action）
)

// ConfigMembershipCard 是 ab.RegisterModel 的全功能参考 demo。
//
// 覆盖：
//  1. ab.RegisterModel(mb) —— 返回 *activity.ModelBuilder，editing 页 CRUD 自动记日志
//  2. ModelBuilder 全部配置链 —— Keys/IgnoredFields/Skip*/LabelFunc/LinkFunc/BeforeCreate/AddTypeHanders
//  3. mb.RegisterEventFunc 内手动记录 —— OnEdit(自动 diff) / Log(自定义 action)
//  4. 时间线两种展示 —— 详情侧栏 NewTimelineCompo + 行内 NewTimelinePopoverBtn
func ConfigMembershipCard(b *presets.Builder, db *gorm.DB, ab *activity.Builder) {
	mb := b.Model(&models.MembershipCard{}).URIName("membership-cards")

	// ========== 1. 注册到 activity ==========
	// 传 presets.ModelBuilder：editing 页的 新建/编辑/删除 会「自动」记日志（WrapperSaveFunc 钩子），
	// 同时装好时间线所需事件，供 NewTimelineCompo / NewTimelinePopoverBtn 加载。
	// 返回的 amb 用于：配置记录行为 + 在自定义事件里手动 OnEdit/Log。
	amb := ab.RegisterModel(mb)

	// ========== 2. ModelBuilder 配置链（全部可选） ==========
	amb.
		// Keys：日志按这些字段值关联到具体记录（默认已用主键；AddKeys 追加业务键）。
		AddKeys("Number").
		// IgnoredFields：这些字段的变更不记进 diff（主键已默认忽略；外键/时间戳常忽略）。
		AddIgnoredFields("CustomerID").
		// SkipCreate/SkipEdit/SkipDelete：关闭某类自动记录（这里只演示，不关）。
		// SkipDelete().
		// LabelFunc：时间线/日志列表里这个模型的显示名。
		LabelFunc(func() string { return "会员卡" }).
		// LinkFunc：点日志条目跳转到的详情链接（obj 是当前记录）。
		LinkFunc(func(v any) string {
			c := v.(*models.MembershipCard)
			return fmt.Sprintf("%s/%d", mb.Info().DetailingHref(""), c.ID)
		}).
		// BeforeCreate：日志落库前的钩子，可改写 log（如补充上下文、改 action 文案）。
		BeforeCreate(func(ctx context.Context, log *activity.ActivityLog) error {
			// 例：给所有会员卡日志打个来源标记（演示用，实际按需）
			return nil
		})
	// AddTypeHanders：自定义某字段类型的 diff 生成。签名：
	//   amb.AddTypeHanders(time.Time{}, func(old, new any, prefixField string) []activity.Diff { ... })
	// ⚠️ 第一参须传该类型的非 nil 值（内部 reflect.Indirect 取类型，传 nil 指针会 panic）。

	// ========== 列表：行内放「改卡号」+「充值」两个自建 dialog/动作 ==========
	lb := mb.Listing("ID", "Number", "CustomerID", "ValidBefore")

	rm := lb.RowMenu()
	// 行项 1：打开自建 dialog 改卡号（→ OnEdit 自动 diff）
	rm.RowMenuItem("EditNumber").ComponentFunc(func(obj any, id string, ctx *web.EventContext) h.HTMLComponent {
		card := obj.(*models.MembershipCard)
		onClick := web.Plaid().EventFunc(eventCardEditNumberOpen).
			Query(presets.ParamID, id).
			Query("number", fmt.Sprint(card.Number)).Go()
		return RowMenuItem("改卡号").SetIcon("pencil").SetOnclick(onClick)
	})
	// 行项 2：充值（→ Log 自定义 action，演示无 dialog 直接动作）
	rm.RowMenuItem("Recharge").ComponentFunc(func(obj any, id string, ctx *web.EventContext) h.HTMLComponent {
		onClick := web.Plaid().EventFunc(eventCardRecharge).Query(presets.ParamID, id).Go()
		return RowMenuItem("充值 ¥100").SetIcon("wallet").SetOnclick(onClick)
	})
	// 行项 3：查看本卡时间线（弹出 Popover）
	rm.RowMenuItem("Timeline").ComponentFunc(func(obj any, id string, ctx *web.EventContext) h.HTMLComponent {
		return amb.NewTimelinePopoverBtn(ctx, obj)
	})

	// ========== 详情页 + 侧栏时间线 ==========
	mb.Detailing("ID", "Number", "CustomerID", "ValidBefore").
		SidePanelFunc(func(obj any, ctx *web.EventContext) h.HTMLComponent {
			if ctx.R.FormValue(presets.ParamID) == "" {
				return nil
			}
			return amb.NewTimelineCompo(ctx, obj, "_side")
		})

	// ========== 3. mb.RegisterEventFunc：模型级事件里手动操作 activity ==========

	// 事件 A：打开改卡号 dialog（受控 Dialog → DialogPortal，预填当前卡号）
	mb.RegisterEventFunc(eventCardEditNumberOpen, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.FormValue(presets.ParamID)
		number := ctx.R.FormValue("number")
		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: presets.DialogPortalName,
			Body: web.Scope(
				Dialog(
					DialogContent(
						DialogHeader(DialogTitle(h.Text("修改卡号"))),
						h.Div(
							Input().Type("number").Label("卡号").
								Attr(web.VField("card_number", number)...),
						).Class("py-4"),
						DialogFooter(
							Button(h.Text("取消")).Variant(ButtonVariantOutline).
								Attr("@click", "locals.show = false"),
							Button(h.Text("保存")).
								Attr("@click", "locals.show = false;"+web.Plaid().
									EventFunc(eventCardEditNumberSave).
									Query(presets.ParamID, id).Go()),
						),
					).Class(presets.DialogSizeXs),
				).Attr(":open", "locals.show").
					OnUpdateOpen("locals.show = $event"),
			).VSlot("{ locals }").Init("{show: true}"),
		})
		return
	})

	// 事件 B：保存卡号 —— 查 old → 改 → 存 → amb.OnEdit 自动算 diff 记日志
	mb.RegisterEventFunc(eventCardEditNumberSave, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.FormValue(presets.ParamID)
		newNumber, _ := strconv.Atoi(ctx.R.FormValue("card_number"))

		var old models.MembershipCard
		if err = db.First(&old, id).Error; err != nil {
			presets.ShowError(&r, "记录不存在")
			return r, nil
		}
		updated := old
		updated.Number = newNumber
		if err = db.Save(&updated).Error; err != nil {
			presets.ShowError(&r, "保存失败")
			return r, nil
		}
		// 自动 diff old→updated（无变化不记）；操作人来自 ctx.R.Context()
		if _, lerr := amb.OnEdit(ctx.R.Context(), &old, &updated); lerr != nil {
			fmt.Printf("activity OnEdit failed: %v\n", lerr)
		}
		presets.ShowSuccess(&r, "卡号已修改")
		r.Emit(presets.NotifModelsUpdated(&models.MembershipCard{}), presets.PayloadModelsUpdated{Ids: []string{id}})
		return
	})

	// 事件 C：充值 —— 直接动作 + amb.Log 记「自定义 action + 任意 detail」（不是 diff 形式）
	mb.RegisterEventFunc(eventCardRecharge, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.FormValue(presets.ParamID)
		var card models.MembershipCard
		if err = db.First(&card, id).Error; err != nil {
			presets.ShowError(&r, "记录不存在")
			return r, nil
		}
		// 业务：这里只演示，假设充值改了某余额字段……（MembershipCard 无余额字段，仅记日志示意）
		// Log(ctx, action, obj, detail)：detail 可是任意可序列化结构，时间线按自定义渲染
		if _, lerr := amb.Log(ctx.R.Context(), "Recharge", &card, map[string]any{
			"amount": 100,
			"note":   "手动充值 ¥100",
		}); lerr != nil {
			fmt.Printf("activity Log failed: %v\n", lerr)
		}
		presets.ShowSuccess(&r, "充值成功（已记录）")
		return
	})
}
