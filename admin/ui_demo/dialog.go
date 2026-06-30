package ui_demo

import (
	"fmt"
	"strings"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/presets/actions"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
	"gorm.io/gorm"
)

// ============================================================================
// Dialog 类型演示
// ============================================================================
//
// ## 4 种 Dialog 对比
//
// | 类型               | 触发方式               | 回调签名                          | 用途                        |
// |--------------------|----------------------|----------------------------------|----------------------------|
// | Global Dialog      | 点击新建/编辑按钮       | —（框架内置的 Edit 流程）           | 表单编辑（默认 Dialog 模式）   |
// | ListingCompo Dialog| EventFunc(OpenListingDialog) | —（展示另一个 Model 的列表）  | 关联记录选择                  |
// | BulkAction Dialog  | 勾选多行 → 点击按钮     | func(selectedIds, ctx) Component | 批量修改状态/属性             |
// | Action Dialog      | 点击列表顶部操作按钮     | func(id, ctx) Component         | 全局操作（导出、导入等）       |
//
// ## 关键区别
//
//   - Global Dialog: 由框架自动管理，使用 DialogPortalName Portal，支持离开确认
//   - ListingCompo Dialog: 使用 ListingDialogPortalName Portal，内部含完整列表（搜索、筛选、分页）
//   - BulkAction Dialog: 需要先勾选行，接收 selectedIds []string，适合批量操作
//   - Action Dialog: 接收 id string（列表级别为空，详情级别为记录 ID），适合全局或单条操作
//
// ============================================================================

const (
	// ListingCompo Dialog 用的 URI 名称（独立 Model，不显示在菜单中）
	uriNameDialogDemoSelector = "dialog-demo-selector"

	// 自定义事件名
	eventSelectDialogDemo = "dialog_demo_selectRelated"
)

// configDialogDemo 配置 Dialog 类型演示模块
func ConfigDialogDemo(b *presets.Builder, db *gorm.DB) {
	db.AutoMigrate(&models.DialogDemo{})

	mb := b.Model(&models.DialogDemo{}).URIName("dialog-demos")

	// ========================================================================
	// 1. Global Dialog（全局对话框）— 编辑表单
	// ========================================================================
	//
	// 这是 admin 框架默认的编辑模式。
	// 不设置 .Drawer(true) 时，编辑表单在全局 Dialog 中打开。
	// Dialog 由 presets.Builder.dialog() 方法创建，使用 DialogPortalName Portal。
	// 特性：
	//   - 支持表单未保存离开确认（VarsPresetsDataChanged）
	//   - 宽度由 presets.Builder.DialogContentClass() 控制（默认响应式 class）
	//   - 新建和编辑都走同一个 Dialog
	ed := mb.Editing("Title", "Status", "Priority", "Notes", "RelatedID")

	ed.Field("Title").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Input().
				Label(field.Label).
				Placeholder("请输入标题").
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	ed.Field("Status").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Select(
				shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("选择状态")),
				shadcn.SelectContent(
					shadcn.SelectItem(h.Text("激活")).Value("active"),
					shadcn.SelectItem(h.Text("待定")).Value("pending"),
					shadcn.SelectItem(h.Text("停用")).Value("inactive"),
				),
			).Label(field.Label).
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	ed.Field("Priority").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Select(
				shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("选择优先级")),
				shadcn.SelectContent(
					shadcn.SelectItem(h.Text("1 - 最低")).Value("1"),
					shadcn.SelectItem(h.Text("2 - 低")).Value("2"),
					shadcn.SelectItem(h.Text("3 - 中")).Value("3"),
					shadcn.SelectItem(h.Text("4 - 高")).Value("4"),
					shadcn.SelectItem(h.Text("5 - 最高")).Value("5"),
				),
			).Label(field.Label).
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	ed.Field("Notes").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Textarea().
				Label(field.Label).
				Placeholder("请输入备注...").
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	// RelatedID 字段：点击按钮打开 ListingCompo Dialog 选择关联记录
	ed.Field("RelatedID").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				shadcn.Input().
					Label(field.Label).
					Attr(web.VField(field.FormKey, field.Value(obj))...).
					Disabled(true),
				// 旧法：单建 InMenu(false) 的 selector 模型 + WrapCell 手挂行点击
				dialogDemoListingDialogTrigger(),
				// 新法（B）：OpenListing 按次配置——直接弹本模型列表，无需 selector 模型；
				// 尺寸/列/行菜单/选中事件全在调用点声明。OnSelectEvent 传已注册事件名，框架按固定
				// 模板构建行点击调用（id=$event.id），复用既有 eventSelectDialogDemo（弹提示并关弹窗）。
				presets.OpenListing(mb).
					Size(presets.DialogSizeMd).
					Columns("ID", "Title", "Status").
					HideRowMenu().
					HideNewButton().
					SearchOff().
					OnSelectEvent(eventSelectDialogDemo).
					Button(h.Text("选择关联记录（OpenListing 按次配置）")).
					Variant(shadcn.ButtonVariantSecondary),
			).Class("flex flex-col gap-2")
		})

	ed.ValidateFunc(func(obj any, ctx *web.EventContext) (err web.ValidationErrors) {
		demo := obj.(*models.DialogDemo)
		if demo.Title == "" {
			err.FieldError("Title", "标题不能为空")
		}
		return
	})

	// ========================================================================
	// 2. Listing 配置 + BulkAction Dialog + Action Dialog
	// ========================================================================
	lb := mb.Listing("ID", "Title", "Status", "Priority", "RelatedID", "UpdatedAt")
	lb.SearchColumns("title", "notes")
	lb.SelectableColumns(true) // 启用行选择（BulkAction 需要）
	lb.PerPage(10)

	lb.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {

		return []*shadcn.FilterItem{
			{
				Key:          "status",
				Label:        "Status",
				ItemType:     shadcn.FilterItemTypeSelect,
				SQLCondition: `status %s ?`,
				Options: []shadcn.FilterSelectOption{
					{Text: "Active", Value: "active"},
					{Text: "Pending", Value: "pending"},
					{Text: "Inactive", Value: "inactive"},
				},
			},
			{
				Key:          "priority",
				Label:        "Priority",
				ItemType:     shadcn.FilterItemTypeSelect,
				SQLCondition: `priority %s ?`,
				Options: []shadcn.FilterSelectOption{
					{Text: "1 - 最低", Value: "1"},
					{Text: "3 - 中", Value: "3"},
					{Text: "5 - 最高", Value: "5"},
				},
			},
		}
	})

	lb.Field("Status").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		demo := obj.(*models.DialogDemo)
		var variant shadcn.BadgeVariant
		switch demo.Status {
		case "active":
			variant = shadcn.BadgeVariantDefault
		case "pending":
			variant = shadcn.BadgeVariantSecondary
		case "inactive":
			variant = shadcn.BadgeVariantDestructive
		default:
			variant = shadcn.BadgeVariantOutline
		}
		return h.Td(shadcn.Badge(h.Text(demo.Status)).Variant(variant))
	})

	// ────────────────────────────────────────────────────────
	// 在「列表行内」打开 ListingDialog（按本行值预过滤）
	// ────────────────────────────────────────────────────────
	//
	// 对标「点某行的某字段 → 弹出按该行值过滤的子列表」（如点 UserAgentID → 弹该 UA 的访问统计）。
	// 关键：@click.stop 阻止冒泡到行点击（否则会同时触发行编辑）；Filter(key,val) 预设筛选（key
	// 须是 FilterDataFunc 注册项）；这里演示打开 DialogDemos 自身、按本行 Status 过滤。
	lb.Field("RelatedID").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		demo := obj.(*models.DialogDemo)
		return h.Div(
			h.Span(fmt.Sprint(demo.RelatedID)).Class("text-xs text-muted-foreground tabular-nums"),
			shadcn.Button(h.Text("同状态")).
				Variant(shadcn.ButtonVariantOutline).
				Size(shadcn.ButtonSizeSm).
				Attr("@click.stop", presets.OpenListing(mb).
					Size(presets.DialogSizeLg).
					Columns("ID", "Title", "Status", "Priority").
					HideRowMenu().HideNewButton().
					Filter("status", demo.Status). // 按本行状态预过滤
					Go()),
		).Class("flex items-center gap-2")
	})

	// ────────────────────────────────────────────────────────
	// 2a. BulkAction Dialog（批量操作对话框）
	// ────────────────────────────────────────────────────────
	//
	// 触发流程：勾选多行 → 点击 "批量修改状态" 按钮 → 弹出 Dialog
	// 回调签名：ComponentFunc(func(selectedIds []string, ctx) Component)
	// 执行回调：UpdateFunc(func(selectedIds []string, ctx, r) error)
	//
	// Dialog 内容由 ComponentFunc 定义，框架自动添加 Cancel/OK 按钮。
	// OK 点击后调用 UpdateFunc 执行实际操作。
	lb.BulkAction("ChangeStatus").
		ComponentFunc(func(selectedIds []string, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.P(h.Text(fmt.Sprintf("已选择 %d 条记录：%s", len(selectedIds), strings.Join(selectedIds, ", ")))).
					Class("text-sm text-muted-foreground mb-4"),
				shadcn.Select(
					shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("选择新状态")),
					shadcn.SelectContent(
						shadcn.SelectItem(h.Text("激活")).Value("active"),
						shadcn.SelectItem(h.Text("待定")).Value("pending"),
						shadcn.SelectItem(h.Text("停用")).Value("inactive"),
					),
				).Label("新状态").
					Attr(web.VField("NewStatus", "")...),
			)
		}).
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) error {
			newStatus := ctx.R.FormValue("NewStatus")
			if newStatus == "" {
				return fmt.Errorf("请选择状态")
			}
			if err := db.Model(&models.DialogDemo{}).
				Where("id IN (?)", selectedIds).
				Update("status", newStatus).Error; err != nil {
				return err
			}
			r.Emit(
				presets.NotifModelsUpdated(&models.DialogDemo{}),
				presets.PayloadModelsUpdated{Ids: selectedIds},
			)
			return nil
		}).
		DialogContentClass(presets.DialogSizeMd)

	// ────────────────────────────────────────────────────────
	// 2b. BulkAction Dialog #2（批量修改优先级 — 第二个示例）
	// ────────────────────────────────────────────────────────
	lb.BulkAction("ChangePriority").
		ComponentFunc(func(selectedIds []string, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.P(h.Text(fmt.Sprintf("将为 %d 条记录设置优先级", len(selectedIds)))).
					Class("text-sm text-muted-foreground mb-4"),
				shadcn.Select(
					shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("选择优先级")),
					shadcn.SelectContent(
						shadcn.SelectItem(h.Text("1 - 最低")).Value("1"),
						shadcn.SelectItem(h.Text("3 - 中")).Value("3"),
						shadcn.SelectItem(h.Text("5 - 最高")).Value("5"),
					),
				).Label("优先级").
					Attr(web.VField("NewPriority", "")...),
			)
		}).
		UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) error {
			priority := ctx.R.FormValue("NewPriority")
			if priority == "" {
				return fmt.Errorf("请选择优先级")
			}
			if err := db.Model(&models.DialogDemo{}).
				Where("id IN (?)", selectedIds).
				Update("priority", priority).Error; err != nil {
				return err
			}
			r.Emit(
				presets.NotifModelsUpdated(&models.DialogDemo{}),
				presets.PayloadModelsUpdated{Ids: selectedIds},
			)
			return nil
		})

	// ────────────────────────────────────────────────────────
	// 2c. Action Dialog（操作对话框 — 列表级别）
	// ────────────────────────────────────────────────────────
	//
	// 触发流程：点击列表顶部的操作按钮 → 弹出 Dialog
	// 回调签名：ComponentFunc(func(id string, ctx) Component)
	//   - 在列表级别，id 参数为空字符串 ""
	//   - 在详情级别，id 参数为记录 ID
	//
	// 与 BulkAction 的区别：
	//   - Action 不需要勾选行，始终可点击
	//   - Action 接收单个 id（列表级为空），BulkAction 接收 selectedIds 数组
	lb.Action("BatchImport").
		ComponentFunc(func(id string, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.P(h.Text("Action Dialog：id 参数 = \""+id+"\"（列表级别为空字符串）")).
					Class("text-sm text-muted-foreground mb-4 p-3 bg-muted rounded-md"),
				shadcn.Textarea().
					Label("导入数据（JSON）").
					Placeholder("粘贴 JSON 数据...").
					Attr(web.VField("ImportData", "")...),
			)
		}).
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			data := ctx.R.FormValue("ImportData")
			if data == "" {
				return fmt.Errorf("请输入导入数据")
			}
			// 模拟导入
			presets.ShowMessage(r, "导入成功（模拟）", "success")
			return nil
		}).
		DialogContentClass(presets.DialogSizeMd)

	// Action（无 Dialog）— 纯按钮，自定义 ButtonCompFunc
	lb.Action("ExportData").
		ButtonCompFunc(func(ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Button(h.Text("导出 CSV")).
				Variant(shadcn.ButtonVariantOutline).
				Attr("@click", web.Plaid().EventFunc("dialog_demo_exportCSV").Go())
		})

	// 注册导出事件
	b.GetWebBuilder().RegisterEventFunc("dialog_demo_exportCSV", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		presets.ShowMessage(&r, "CSV 导出成功（模拟）", "success")
		return
	})

	// ========================================================================
	// 3. ListingCompo Dialog（列表选择对话框）
	// ========================================================================
	//
	// 这是一个独立的 Model 配置，专门用于在 Dialog 中展示列表。
	// 核心步骤：
	//   a) 创建 Model 并设置 InMenu(false)（不在侧栏显示）
	//   b) 配置 Listing 并设置 DialogContentClass
	//   c) 通过 EventFunc(actions.OpenListingDialog) 触发打开
	//   d) 用户选择后通过自定义事件回传选中的 ID
	//
	// 与 Global Dialog 的区别：
	//   - Global Dialog 展示编辑表单
	//   - ListingCompo Dialog 展示完整列表（含搜索、筛选、分页）
	ConfigDialogDemoSelector(b, db)

	// 注册选择事件
	b.GetWebBuilder().RegisterEventFunc(eventSelectDialogDemo, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		selectedID := ctx.R.FormValue("id")
		presets.ShowMessage(&r, fmt.Sprintf("已选择记录 ID: %s", selectedID), "success")
		// 关闭列表对话框
		r.RunScript = "vars.presetsListingDialog = false"
		return
	})

	// ========================================================================
	// 4. Detailing（详情页）+ Detailing Action Dialog
	// ========================================================================
	//
	// Action 也可以挂在 DetailingBuilder 上。
	// 与 Listing Action 的区别：id 参数为当前记录的 ID（非空）。
	dt := mb.Detailing(
		&presets.FieldsSection{
			Title: "基本信息",
			Rows: [][]string{
				{"ID", "Title"},
				{"Status", "Priority"},
				{"Notes"},
				{"RelatedID"},
			},
		},
		"EventFlowGraph",
	).Drawer(true) // 详情页用 Drawer 打开（和 Dialog 对比）
	dt.EnableRefreshOnUpdate()

	dt.Field("Status").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		demo := obj.(*models.DialogDemo)
		var variant string
		switch demo.Status {
		case "active":
			variant = "default"
		case "pending":
			variant = "secondary"
		case "inactive":
			variant = "destructive"
		}
		return readonlyFieldWithChild(field.Label, shadcn.Badge(h.Text(demo.Status)).Variant(shadcn.BadgeVariant(variant)))
	})

	// 事件流向图：展示 NotifModelsUpdated 事件的流转机制
	dt.Field("EventFlowGraph").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(
			h.H3("事件流向").Class("text-lg font-semibold mb-2"),
			unovis.Graph().
				Direction(unovis.GraphDirectionLR).
				NodeSize(40).
				ShowLabels(true).
				LinkFlow(true).
				LinkFlowParticleSpeed(40).
				Config(unovis.GraphConfig{
					"action": {
						Label: "Action 提交",
						Color: "var(--chart-1)",
						Shape: unovis.GraphNodeShapeCircle,
					},
					"section": {
						Label: "Section 保存",
						Color: "var(--chart-2)",
						Shape: unovis.GraphNodeShapeCircle,
					},
					"emit": {
						Label: "Emit 事件",
						Color: "var(--chart-3)",
						Shape: unovis.GraphNodeShapeHexagon,
					},
					"listing": {
						Label: "Listing 列表",
						Color: "var(--chart-4)",
						Shape: unovis.GraphNodeShapeSquare,
					},
					"detailing": {
						Label: "Detailing 详情",
						Color: "var(--chart-5)",
						Shape: unovis.GraphNodeShapeSquare,
					},
				}).
				Data(unovis.GraphData{
					Nodes: []map[string]any{
						{"id": "action"},
						{"id": "section"},
						{"id": "emit"},
						{"id": "listing"},
						{"id": "detailing"},
					},
					Links: []map[string]any{
						{"source": "action", "target": "emit"},
						{"source": "section", "target": "emit"},
						{"source": "emit", "target": "listing"},
						{"source": "emit", "target": "detailing"},
					},
				}),
		).Class("mb-4")
	})

	dt.TopActionsFunc(func(obj any, ctx *web.EventContext) h.HTMLComponent {
		var actionBtns h.HTMLComponents

		// 确认支付
		actionBtns = append(actionBtns, shadcn.Button(h.Text("导出")))
		actionBtns = append(actionBtns, shadcn.Button(h.Text("导出")))
		return shadcn.ButtonGroup(actionBtns...)
	})

	// Detailing Action Dialog — id 参数为记录 ID
	dt.Action("EditStatus").
		ComponentFunc(func(id string, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.P(h.Text("Detailing Action Dialog：id 参数 = \""+id+"\"（当前记录 ID）")).
					Class("text-sm text-muted-foreground mb-4 p-3 bg-muted rounded-md"),
				shadcn.Select(
					shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("选择新状态")),
					shadcn.SelectContent(
						shadcn.SelectItem(h.Text("激活")).Value("active"),
						shadcn.SelectItem(h.Text("待定")).Value("pending"),
						shadcn.SelectItem(h.Text("停用")).Value("inactive"),
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

			// 更新数据库
			if err := db.Model(&models.DialogDemo{}).
				Where("id = ?", id).
				Update("status", newStatus).Error; err != nil {
				return err
			}

			// 显示成功提示
			presets.ShowMessage(r, "状态修改成功", "")

			// 关闭 Action Dialog
			web.AppendRunScripts(r, presets.CloseDialogVarScript)

			// 通知列表页更新
			r.Emit(
				presets.NotifModelsUpdated(&models.DialogDemo{}),
				presets.PayloadModelsUpdated{Ids: []string{id}},
			)
			return nil
		})

	dt.Action("AddNote").
		ComponentFunc(func(id string, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.P(h.Text(fmt.Sprintf("为记录 #%s 添加备注", id))).
					Class("text-sm text-muted-foreground mb-4"),
				shadcn.Textarea().
					Label("备注内容").
					Placeholder("输入备注...").
					Attr(web.VField("AppendNotes", "")...),
			)
		}).
		UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) error {
			notes := ctx.R.FormValue("AppendNotes")
			if notes == "" {
				return fmt.Errorf("请输入备注")
			}

			// 追加备注
			var demo models.DialogDemo
			if err := db.First(&demo, id).Error; err != nil {
				return err
			}
			if demo.Notes != "" {
				demo.Notes += "\n"
			}
			demo.Notes += notes
			if err := db.Save(&demo).Error; err != nil {
				return err
			}

			// 显示成功提示
			// presets.ShowMessage(r, "备注添加成功", "")

			dt.UpdateOverlayContent(ctx, r, &demo, "备注添加成功", nil)

			// 关闭 Action Dialog
			web.AppendRunScripts(r, presets.CloseDialogVarScript)

			// 通知列表页更新
			r.Emit(
				presets.NotifModelsUpdated(&models.DialogDemo{}),
				presets.PayloadModelsUpdated{Ids: []string{id}},
			)
			return nil
		})
}

// configDialogDemoSelector 配置用于 ListingCompo Dialog 的选择器
//
// 这是一个独立的 ModelBuilder，不在菜单中显示（InMenu(false)），
// 专门用于在 Dialog 中展示列表供用户选择。
func ConfigDialogDemoSelector(pb *presets.Builder, db *gorm.DB) {
	b := pb.Model(&models.DialogDemo{}).
		URIName(uriNameDialogDemoSelector).
		InMenu(false)

	lb := b.Listing("ID", "Title", "Status", "Priority").
		DialogContentClass(presets.DialogSizeLg).
		SearchColumns("title").
		PerPage(10).
		SelectableColumns(true)

	// 筛选器（用于测试 Dialog 中的 Filter 布局）
	lb.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		return []*shadcn.FilterItem{
			{
				Key:      "status",
				Label:    "Status",
				ItemType: shadcn.FilterItemTypeSelect,
				Options: []shadcn.FilterSelectOption{
					{Text: "Active", Value: "active"},
					{Text: "Pending", Value: "pending"},
					{Text: "Inactive", Value: "inactive"},
				},
			},
			{
				Key:      "priority",
				Label:    "Priority",
				ItemType: shadcn.FilterItemTypeSelect,
				Options: []shadcn.FilterSelectOption{
					{Text: "1 - 最低", Value: "1"},
					{Text: "3 - 中", Value: "3"},
					{Text: "5 - 最高", Value: "5"},
				},
			},
		}
	})

	// 隐藏新建按钮（选择对话框不需要）
	lb.NewButtonFunc(func(ctx *web.EventContext) h.HTMLComponent { return nil })
	// 隐藏行菜单（选择对话框不需要）
	lb.RowMenu().Empty()

	// 点击行触发选择事件
	lb.WrapCell(func(in presets.CellProcessor) presets.CellProcessor {
		return func(evCtx *web.EventContext, cell h.MutableAttrHTMLComponent, id string, obj any) (h.MutableAttrHTMLComponent, error) {
			cell.SetAttr("@click", web.Plaid().
				Query("id", id).
				EventFunc(eventSelectDialogDemo).
				Go(),
			)
			return in(evCtx, cell, id, obj)
		}
	})

	lb.Field("Status").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		demo := obj.(*models.DialogDemo)
		var variant shadcn.BadgeVariant
		switch demo.Status {
		case "active":
			variant = shadcn.BadgeVariantDefault
		case "pending":
			variant = shadcn.BadgeVariantSecondary
		case "inactive":
			variant = shadcn.BadgeVariantDestructive
		default:
			variant = shadcn.BadgeVariantOutline
		}
		return h.Td(shadcn.Badge(h.Text(demo.Status)).Variant(variant))
	})
}

// dialogDemoListingDialogTrigger 创建打开 ListingCompo Dialog 的触发按钮
//
// 使用方法：在编辑表单字段中调用，点击后打开 ListingCompo Dialog
// 内部通过 EventFunc(actions.OpenListingDialog) 触发
func dialogDemoListingDialogTrigger() h.HTMLComponent {
	return shadcn.Button(h.Text("选择关联记录（ListingCompo Dialog）")).
		Variant(shadcn.ButtonVariantOutline).
		Attr("@click",
			web.Plaid().
				URL("/"+uriNameDialogDemoSelector).
				EventFunc(actions.OpenListingDialog).
				Go(),
		)
}

// readonlyFieldWithChild 创建一个只读字段组件（标签+子内容）
func readonlyFieldWithChild(label string, child h.HTMLComponent) h.HTMLComponent {
	return h.Div(
		h.Div(
			h.Label(label).Class("text-sm font-medium mb-1 block"),
			child,
		).Class("mb-4"),
	)
}
