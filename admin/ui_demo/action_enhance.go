package ui_demo

import (
	"strconv"

	"example/models"

	"github.com/r0vx/admin/notification"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
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
	rm := mb.Listing().RowMenu()
	rm.Inline(true)

	rm.RowMenuItem("Upgrade").
		Icon("arrow-up-circle").
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

	rm.RowMenuItem("Reset").
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
}
