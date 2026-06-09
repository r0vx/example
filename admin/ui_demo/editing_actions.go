package ui_demo

import (
	"fmt"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/presets/actions"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// ============================================================================
// Editing 自定义按钮演示
// ============================================================================
//
// ## 4 种按钮/操作机制
//
// | API                    | 位置               | 行为                              |
// |------------------------|--------------------|---------------------------------|
// | ActionsFunc            | 底部（CardFooter）  | 完全替换默认保存按钮               |
// | TopActionsFunc         | 顶部（Card 上方）   | 设置顶部操作按钮区                 |
// | AppendTopActionsFunc   | 顶部（Card 上方）   | 追加按钮，不覆盖已有               |
// | **Action(name)**       | 顶部按钮 + 嵌套Dialog | 按钮 → 对话框 → 表单 → 提交回调 |
//
// ## 典型场景
//
//   - ActionsFunc: 审批流程（保存草稿 / 提交审核 / 发布）
//   - Action: 编辑页内的额外操作（修改状态、添加备注、导出等），不影响保存按钮
//
// ============================================================================

// configEditingActionsDemo 配置 Editing 自定义按钮演示模块
func ConfigEditingActionsDemo(b *presets.Builder, db *gorm.DB) {
	db.AutoMigrate(&models.EditingActionsDemo{})

	mb := b.Model(&models.EditingActionsDemo{}).URIName("editing-actions-demos")

	// Listing
	lb := mb.Listing("ID", "Title", "Status", "UpdatedAt")
	lb.SearchColumns("title")
	lb.PerPage(10)

	lb.Field("Status").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		demo := obj.(*models.EditingActionsDemo)
		var variant shadcn.BadgeVariant
		switch demo.Status {
		case "published":
			variant = shadcn.BadgeVariantDefault
		case "pending":
			variant = shadcn.BadgeVariantSecondary
		case "draft":
			variant = shadcn.BadgeVariantOutline
		default:
			variant = shadcn.BadgeVariantOutline
		}
		return shadcn.Badge(h.Text(demo.Status)).Variant(variant)
	})

	// Editing（不启用 Detailing，Editing 作为主编辑入口）
	ed := mb.Editing("Title", "Status", "Content")

	ed.Field("Title").
		ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Input().
				Label(field.Label).
				Placeholder("请输入标题").
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	ed.Field("Status").
		ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Select(
				shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("选择状态")),
				shadcn.SelectContent(
					shadcn.SelectItem(h.Text("草稿")).Value("draft"),
					shadcn.SelectItem(h.Text("待审核")).Value("pending"),
					shadcn.SelectItem(h.Text("已发布")).Value("published"),
				),
			).Label(field.Label).
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	ed.Field("Content").
		ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Textarea().
				Label(field.Label).
				Placeholder("请输入内容...").
				Rows(6).
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	ed.ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		demo := obj.(*models.EditingActionsDemo)
		if demo.Title == "" {
			err.FieldError("Title", "标题不能为空")
		}
		return
	})

	// ────────────────────────────────────────────────────────
	// ActionsFunc — 替换底部保存按钮
	// ────────────────────────────────────────────────────────
	//
	// 默认底部只有 [保存] 按钮。这里替换为 [保存草稿] + [提交审核] 两个按钮。
	// 保存按钮复用框架的 actions.Update 事件（触发 ValidateFunc + SaveFunc 流程）。
	// 保存草稿通过自定义 EventFunc 实现，跳过 ValidateFunc 直接写库。
	ed.ActionsFunc(func(obj interface{}, ctx *web.EventContext) h.HTMLComponent {
		demo := obj.(*models.EditingActionsDemo)
		id := fmt.Sprint(demo.ID)

		// 复用框架默认的保存事件（走 Validate + Save 完整流程）
		submitClick := web.Plaid().
			EventFunc(actions.Update).
			URL(mb.Info().ListingHref()).
			Go()

		// 自定义事件：保存草稿（跳过验证，直接更新 status=draft）
		saveDraftClick := web.Plaid().
			EventFunc("editingActionsDemo_saveDraft").
			Query("id", id).
			Go()

		return h.Div(
			shadcn.Button(h.Text("保存草稿")).
				Variant(shadcn.ButtonVariantOutline).
				Attr(":disabled", "isFetching").
				Attr("@click", saveDraftClick),
			h.Div().Class("flex-1"),
			shadcn.Button(h.Text("提交审核")).
				Attr(":disabled", "isFetching").
				Attr(":loading", "isFetching").
				Attr("@click", submitClick),
		).Class("flex items-center gap-2 w-full")
	})

	// 注册"保存草稿"事件
	b.GetWebBuilder().RegisterEventFunc("editingActionsDemo_saveDraft", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.FormValue("id")
		if err = db.Model(&models.EditingActionsDemo{}).Where("id = ?", id).
			Updates(map[string]interface{}{
				"status":  "draft",
				"title":   ctx.R.FormValue("Title"),
				"content": ctx.R.FormValue("Content"),
			}).Error; err != nil {
			return
		}
		presets.ShowMessage(&r, "已保存为草稿", "success")
		r.Emit(
			presets.NotifModelsUpdated(&models.EditingActionsDemo{}),
			presets.PayloadModelsUpdated{Ids: []string{id}},
		)
		return
	})

	// ────────────────────────────────────────────────────────
	// Action — 编辑页对话框操作（类似 Detailing.Action）
	// ────────────────────────────────────────────────────────
	//
	// Action 在编辑页顶部渲染按钮，点击后弹出嵌套对话框。
	// 与 ActionsFunc 的区别：
	//   - ActionsFunc 替换底部保存按钮，用于改变保存行为
	//   - Action 在顶部添加额外操作按钮，不影响保存流程
	//
	// 使用场景：修改状态、添加备注、发起审批等独立操作。
	ed.Action("QuickChangeStatus").
		ComponentFunc(func(id string, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.P(h.Text(fmt.Sprintf("为记录 #%s 快速修改状态（不影响表单其他字段）", id))).
					Class("text-sm text-muted-foreground mb-4"),
				shadcn.Select(
					shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("选择新状态")),
					shadcn.SelectContent(
						shadcn.SelectItem(h.Text("草稿")).Value("draft"),
						shadcn.SelectItem(h.Text("待审核")).Value("pending"),
						shadcn.SelectItem(h.Text("已发布")).Value("published"),
					),
				).Label("新状态").
					Attr(web.VField("NewStatus", "")...),
			)
		}).
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			newStatus := ctx.R.FormValue("NewStatus")
			if newStatus == "" {
				return fmt.Errorf("请选择状态")
			}
			if err := db.Model(&models.EditingActionsDemo{}).
				Where("id = ?", id).
				Update("status", newStatus).Error; err != nil {
				return err
			}
			presets.ShowMessage(r, "状态修改成功", "success")
			r.Emit(
				presets.NotifModelsUpdated(&models.EditingActionsDemo{}),
				presets.PayloadModelsUpdated{Ids: []string{id}},
			)
			return nil
		})

	ed.Action("AddNote").
		ComponentFunc(func(id string, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.P(h.Text(fmt.Sprintf("为记录 #%s 添加备注", id))).
					Class("text-sm text-muted-foreground mb-4"),
				shadcn.Textarea().
					Label("备注内容").
					Placeholder("输入备注...").
					Attr(web.VField("AppendContent", "")...),
			)
		}).
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			content := ctx.R.FormValue("AppendContent")
			if content == "" {
				return fmt.Errorf("请输入备注内容")
			}
			var demo models.EditingActionsDemo
			if err := db.First(&demo, id).Error; err != nil {
				return err
			}
			if demo.Content != "" {
				demo.Content += "\n"
			}
			demo.Content += content
			if err := db.Save(&demo).Error; err != nil {
				return err
			}
			presets.ShowMessage(r, "备注添加成功", "success")
			r.Emit(
				presets.NotifModelsUpdated(&models.EditingActionsDemo{}),
				presets.PayloadModelsUpdated{Ids: []string{id}},
			)
			return nil
		})
}
