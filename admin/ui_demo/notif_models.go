package ui_demo

// ============================================================================
// NotifModels 事件系统演示
// ============================================================================
//
// NotifModels 是 r0vx admin 框架中的前端事件通知机制。
// 它允许组件之间通过 Emit（发送）和 Listen（监听）进行通信。
//
// 核心用途：当编辑页面创建/更新/删除数据后，自动通知列表页面刷新。
//
// 事件流程图：
//
//   [编辑页面 Save]
//         ↓
//   r.Emit(NotifModelsCreated)    ← editing.go 自动发送
//         ↓
//   [前端 Vue 事件总线传递]
//         ↓
//   web.Listen(NotifModelsCreated) ← listing_compo.go 自动监听
//         ↓
//   [列表页面自动刷新]
//
// 三种内置事件：
//   1. NotifModelsCreated - 新建记录后触发
//   2. NotifModelsUpdated - 更新记录后触发
//   3. NotifModelsDeleted - 删除记录后触发
//
// ============================================================================

import (
	"fmt"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// NotifDemo 演示用模型
type NotifDemo struct {
	gorm.Model
	Title  string
	Status string
}

// configNotifDemo 配置 NotifModels 事件演示模块
func ConfigNotifDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&NotifDemo{}); err != nil {
		panic(err)
	}
	mb := b.Model(&NotifDemo{})

	// ================================================================
	// 一、自动事件机制（框架内置，无需额外代码）
	// ================================================================
	//
	// 1) Listing 自动监听：listing_compo.go 默认注册了三个事件监听器
	//
	//    web.Listen(
	//        mb.NotifModelsCreated(), ReloadAction,  // 创建后刷新列表
	//        mb.NotifModelsUpdated(), ReloadAction,  // 更新后刷新列表
	//        mb.NotifModelsDeleted(), ReloadAction,  // 删除后刷新列表
	//    )
	//
	// 2) Editing 自动发送：editing.go 在 Save/Delete 操作后自动 Emit 事件
	//
	//    新建时: r.Emit(mb.NotifModelsCreated(), PayloadModelsCreated{Models: []any{obj}})
	//    更新时: r.Emit(mb.NotifModelsUpdated(), PayloadModelsUpdated{Ids: []string{id}})
	//    删除时: r.Emit(mb.NotifModelsDeleted(), PayloadModelsDeleted{Ids: deletedIds})
	//
	// 总结：只要正常配置 Listing + Editing，列表就会自动刷新，无需任何额外代码。

	mb.Listing("ID", "Title", "Status", "UpdatedAt").SearchColumns("title")
	mb.Editing("Title", "Status")

	// ================================================================
	// 二、手动发送事件（在自定义事件处理函数中）
	// ================================================================
	//
	// 某些场景下，你需要在自定义逻辑中手动触发刷新。
	// 例如：批量操作、自定义 Action、后台任务完成后。
	//
	// 关键 API：
	//   r.Emit(事件名称, 事件载荷)
	//
	// 事件名称生成方式：
	//   presets.NotifModelsCreated(&NotifDemo{})  → "presets_NotifModelsCreated_*admin.NotifDemo"
	//   presets.NotifModelsUpdated(&NotifDemo{})  → "presets_NotifModelsUpdated_*admin.NotifDemo"
	//   presets.NotifModelsDeleted(&NotifDemo{})  → "presets_NotifModelsDeleted_*admin.NotifDemo"
	//
	// 也可以通过 ModelBuilder 生成：
	//   mb.NotifModelsCreated()  → 同上
	//   mb.NotifModelsUpdated()  → 同上
	//   mb.NotifModelsDeleted()  → 同上

	// 示例：注册自定义批量操作，手动 Emit 事件通知列表刷新
	mb.Listing().BulkAction("BulkActivate").
		ComponentFunc(func(selectedIds []string, ctx *web.EventContext) h.HTMLComponent {
			// 批量操作的 UI 组件（确认对话框内容）
			return h.Div(
				h.Text(fmt.Sprintf("确认激活 %d 条记录？", len(selectedIds))),
			)
		}).
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) error {
			// 执行批量更新
			if err := db.Model(&NotifDemo{}).
				Where("id IN (?)", selectedIds).
				Update("status", "active").Error; err != nil {
				return err
			}

			// ★ 关键：手动 Emit 事件，通知列表刷新
			// 不 Emit 的话，列表不会自动更新，需要用户手动刷新页面
			r.Emit(
				presets.NotifModelsUpdated(&NotifDemo{}),
				presets.PayloadModelsUpdated{Ids: selectedIds},
			)

			// 显示成功提示
			presets.ShowMessage(r, "批量激活成功", "success")
			return nil
		})

	// 示例：注册自定义事件处理函数，手动 Emit 事件
	b.GetWebBuilder().RegisterEventFunc("notifDemoToggleStatus",
		func(ctx *web.EventContext) (r web.EventResponse, err error) {
			id := ctx.R.FormValue("id")
			var demo NotifDemo
			if err = db.First(&demo, id).Error; err != nil {
				return
			}

			// 切换状态
			if demo.Status == "active" {
				demo.Status = "inactive"
			} else {
				demo.Status = "active"
			}
			if err = db.Save(&demo).Error; err != nil {
				return
			}

			// ★ 手动发送更新事件 → 列表自动刷新
			r.Emit(
				presets.NotifModelsUpdated(&NotifDemo{}),
				presets.PayloadModelsUpdated{
					Ids:    []string{id},
					Models: map[string]any{id: demo},
				},
			)

			presets.ShowMessage(&r, "激活成功", "success")
			return
		})

	// ================================================================
	// 三、自定义列表字段中使用事件
	// ================================================================
	//
	// 在列表的自定义字段中，可以绑定事件触发按钮。
	// 事件处理函数中 Emit 后，列表会自动刷新。

	mb.Listing().Field("Status").StopClick().ComponentFunc(
		func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			demo := obj.(*NotifDemo)
			onclick := web.Plaid().
				EventFunc("notifDemoToggleStatus").
				Query("id", fmt.Sprint(demo.ID)).
				Go()
			return shadcn.Switch().
				Checked(demo.Status == "active").
				OnChange(onclick)
		})

	// ================================================================
	// 四、禁用自动监听（DisableModelListeners）
	// ================================================================
	//
	// 在某些场景下，你可能不希望列表自动刷新：
	//   - 高频更新的数据（避免频繁刷新影响用户体验）
	//   - 需要用户手动控制刷新时机
	//   - 特殊的列表页面（如仪表盘中嵌入的列表）
	//
	// 使用方法：
	//   mb.Listing().DisableModelListeners(true)
	//
	// 禁用后，listing_compo.go 中的 web.Listen 代码不会被渲染，
	// 列表将不再响应 NotifModels 事件。
	//
	// 注意：这只影响自动监听，手动调用 r.Emit 的逻辑不受影响。

	// ================================================================
	// 五、跨模型事件通知
	// ================================================================
	//
	// NotifModels 事件是按模型类型区分的。
	// 你可以在一个模型的操作中，Emit 另一个模型的事件。
	//
	// 例如：删除 Category 后，通知 Product 列表刷新
	//
	//   r.Emit(
	//       presets.NotifModelsUpdated(&Product{}),  // ← 注意：Product 而非 Category
	//       presets.PayloadModelsUpdated{},
	//   )
	//
	// 这样，Product 列表会收到通知并刷新（如果它正在显示的话）。
	// 这在处理关联数据时非常有用。
}
