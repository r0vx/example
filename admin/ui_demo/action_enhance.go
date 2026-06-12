package ui_demo

import (
	"fmt"
	"strconv"

	"example/models"

	"github.com/r0vx/admin/notification"
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// configActionEnhanceDemo 演示 Action / BulkAction 体系增强（2026-05 Filament 对照）：
//   - Icon / Tooltip
//   - Visible / Disabled (动态)
//   - RequiresConfirmation + 4 个文本自定义
//   - 与 notification 集成（UpdateFunc 内部调 notifier 发通知）
//
// URI: /action-enhance-demo
func ConfigActionEnhanceDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&models.WizardDemo{}); err != nil {
		panic(err)
	}

	// 构造一个进程级 notifier（与 notification_demo_config 共享 wizardSessionStore 不同；这里独立）
	notifier := notification.New(
		notification.NewToastChannel(),
		notification.NewDatabaseChannel(db, notificationCurrentUserID),
	)

	// D2 i18n 演示：模型不设中文 Label，用默认 ASCII label "WizardDemos"，
	// 中文显示名注册在 admin/messages.go（key WizardDemos）。
	mb := b.Model(&models.WizardDemo{}).URIName("action-enhance-demo")
	mb.Listing("ID", "Name", "Phone", "Status")
	mb.Editing("Name", "Phone", "Industry", "Address", "Status")

	// --- 单条 Action #1: 简单图标 + tooltip + 直接执行 ---
	mb.Listing().Action("Refresh").
		Label("刷新").
		Icon("refresh-cw").
		Tooltip("重新拉取最新数据").
		ButtonColor("ghost").
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			notifier.Info(ctx, r, "已刷新")
			return nil
		})

	// --- 单条 Action #2: RequiresConfirmation + 4 个自定义文本 ---
	mb.Listing().Action("DangerDelete").
		Label("永久删除").
		Icon("trash-2").
		Tooltip("此操作不可撤销").
		ButtonColor("destructive").
		RequiresConfirmation().
		ConfirmTitle("确认永久删除").
		ConfirmPrompt("此操作会清除该商户所有数据，且无法恢复。继续吗？").
		ConfirmOK("我已知晓，删除").
		ConfirmCancel("放弃").
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			idNum, _ := strconv.ParseUint(id, 10, 64)
			if err := db.Delete(&models.WizardDemo{}, idNum).Error; err != nil {
				notifier.Error(ctx, r, "删除失败："+err.Error())
				return err
			}
			notifier.Success(ctx, r, "已永久删除商户 #"+id)
			r.Reload = true
			return nil
		})

	// --- 单条 Action #3: 动态 Disabled / Visible ---
	// 仅当当前请求包含 ?advanced=1 时才显示这个按钮（演示 Visible）
	mb.Listing().Action("AdvancedOp").
		Label("高级操作").
		Icon("settings-2").
		Visible(func(ctx *web.EventContext) bool {
			return ctx.R.URL.Query().Get("advanced") == "1"
		}).
		Disabled(func(ctx *web.EventContext) bool {
			// 永远不禁用 —— 仅演示 API 形态
			return false
		}).
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			notifier.Warning(ctx, r, "高级操作执行成功 (id="+id+")")
			return nil
		})

	// --- BulkAction: 同款 RequiresConfirmation + Icon ---
	mb.Listing().BulkAction("BulkArchive").
		Label("批量归档").
		Icon("archive").
		Tooltip("将选中项移到归档").
		RequiresConfirmation().
		ConfirmPrompt("将归档选中的所有项，归档后不在列表显示").
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) error {
			notifier.Success(ctx, r, "已归档 "+strconv.Itoa(len(selectedIds))+" 项")
			r.Reload = true
			return nil
		})

	// --- BulkAction: 危险操作 ---
	mb.Listing().BulkAction("BulkDelete").
		Label("批量删除").
		Icon("trash-2").
		ButtonColor("destructive").
		RequiresConfirmation().
		ConfirmTitle("确认批量删除").
		ConfirmPrompt("将永久删除选中的所有项，不可恢复！").
		ConfirmOK("永久删除").
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) error {
			ids := make([]uint64, 0, len(selectedIds))
			for _, s := range selectedIds {
				if n, err := strconv.ParseUint(s, 10, 64); err == nil {
					ids = append(ids, n)
				}
			}
			if err := db.Delete(&models.WizardDemo{}, ids).Error; err != nil {
				notifier.Error(ctx, r, "批量删除失败："+err.Error())
				return err
			}
			notifier.Success(ctx, r, "已删除 "+strconv.Itoa(len(ids))+" 项")
			r.Reload = true
			return nil
		})

	// --- RowMenu 行操作增强演示：per-row Visible/Disabled + Confirm + Tooltip + UpdateFunc ---
	// Inline(true) 让行操作平铺为按钮，便于看到 Tooltip 与按钮 disabled 效果。
	// i18n（D2 实演）：不设 Tooltip/ConfirmTitle/ConfirmPrompt 字面量，由 i18n 解析（前缀 mb.label="WizardDemos"），
	// zh 译文注册在 admin/messages.go 的 Messages_*_ModelsI18nModuleKey：
	//   Tooltip       → WizardDemosUpgradeTooltip / WizardDemosResetTooltip
	//   ConfirmTitle  → WizardDemosUpgradeConfirmTitle
	//   ConfirmPrompt → WizardDemosUpgradeConfirmPrompt
	//   ConfirmOK/Cancel 不设 → 内置 msgr.OK/msgr.Cancel（zh：确定/取消）。
	//   .Label("升级"/"重置") 仍是强制字面量（菜单文字）。
	rm := mb.Listing().RowMenu().InlineDefaultsInMenu(true)
	rm.Inline(true)

	rm.RowMenuItem("Upgrade").
		Icon("arrow-up-circle").
		AlsoInDrawer(true). // 同时显示在编辑(抽屉)/详情顶部操作区，带当前记录 id
		Visible(func(obj any, id string, ctx *web.EventContext) bool {
			d, ok := obj.(*models.WizardDemo)
			return ok && d.Status != "published" // 已发布则隐藏升级
		}).
		RequiresConfirmation(). // 确认文案走 i18n（WizardDemosUpgradeConfirmTitle/Prompt）
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			// 演示：实际项目中此处应更新 d.Status 字段；这里仅发通知
			notifier.Success(ctx, r, "已升级商户 #"+id)
			r.Reload = true
			return nil
		})

	rm.RowMenuItem("Upgrade2").
		Icon("arrow-up-circle").
		AlsoInDrawer(true). // 同时显示在编辑(抽屉)/详情顶部操作区，带当前记录 id
		Visible(func(obj any, id string, ctx *web.EventContext) bool {
			d, ok := obj.(*models.WizardDemo)
			return ok && d.Status != "published" // 已发布则隐藏升级
		}).
		RequiresConfirmation(). // 确认文案走 i18n（WizardDemosUpgradeConfirmTitle/Prompt）
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			// 演示：实际项目中此处应更新 d.Status 字段；这里仅发通知
			notifier.Success(ctx, r, "已升级商户 #"+id)
			r.Reload = true
			return nil
		})

	rm.RowMenuItem("Reset").AlsoInDrawer(true).
		IconOnly(true). // 行内只显示 rotate 图标（hover 出名）；抽屉里仍 icon+文字
		Label("重置").
		Icon("rotate-ccw").
		Disabled(func(obj any, id string, ctx *web.EventContext) bool {
			d, ok := obj.(*models.WizardDemo)
			return ok && d.Status == "draft" // 已是草稿则禁用重置
		}).
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			notifier.Warning(ctx, r, "已重置商户 #"+id)
			r.Reload = true
			return nil
		})

	// OnlyInDrawer 演示：只在编辑/详情抽屉顶部出现，列表行 ⋮ 菜单里不显示
	rm.RowMenuItem("ViewDetail").
		OnlyInDrawer(true).
		Label("查看详情").
		Icon("eye").
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			notifier.Info(ctx, r, "查看详情 #"+id)
			return nil
		})

	// --- OnEvent 演示：行菜单项「声明式」触发自定义事件（非内置 UpdateFunc）---
	// 按钮外观仍走链式（Icon/Label/Tooltip），点击触发下面注册的 eventBalance，自动带当前行 id。
	// 与顶部「预存」按钮复用同一事件 → 逻辑只写一处。需额外参数时链式 .EventQuery("k","v")。
	rm.RowMenuItem("Balance").
		Icon("wallet").
		Label("预存").
		Tooltip("给该商户预存金额").
		OnEvent("eventBalance")

	// eventBalance：顶部按钮与行菜单「预存」共用的自定义事件 handler。
	// 演示用发通知示意；实际可弹 dialog 输入金额 / 改余额字段。
	mb.RegisterEventFunc("eventBalance", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.FormValue(presets.ParamID)
		notifier.Success(ctx, &r, "预存成功 #"+id)
		return
	})

	// 编辑模式
	ed := mb.Editing(
		&presets.FieldsSection{
			Title: "Basic Information",
			Rows: [][]string{
				{"Name"},
				{"Phone"},
			},
		},
	)
	// 顶部操作按钮
	ed.TopActionsFunc(func(obj interface{}, ctx *web.EventContext) h.HTMLComponent {
		var btns []h.HTMLComponent
		// 预存：取当前记录真实 id（与行菜单「预存」复用 eventBalance 同一事件）
		id := ""
		if d, ok := obj.(*models.WizardDemo); ok {
			id = fmt.Sprint(d.ID)
		}
		btns = append(btns,
			Button(h.Text("预存")).
				Attr("@click", web.Plaid().EventFunc("eventBalance").
					Query("id", id).Go()),
		)

		if len(btns) == 0 {
			return nil
		}

		return ButtonGroup(btns...)
	})
}
