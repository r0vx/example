package ui_demo

import (
	"fmt"

	"example/models"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/theplant/relay"
	"gorm.io/gorm"
)

// ============================================================================
// ListingBuilder Wrap* 方法演示
// ============================================================================
//
// 演示 ListingBuilder 的 6 个 Wrap* 包装方法：
//   - WrapColumns:        动态修改列配置（表头文字、样式等）
//   - WrapCell:           包装单元格（添加点击事件、条件样式等）
//   - WrapRow:            包装行（条件高亮、行级事件等）
//   - WrapFilterDataFunc: 动态追加或修改筛选项
//   - WrapNewButtonFunc:  修改「新建」按钮（隐藏、替换、添加权限控制等）
//   - WrapSearchFunc:     包装搜索逻辑（追加条件、权限过滤、默认排序等）
//
// 所有 Wrap* 方法都采用装饰器模式，支持链式调用多次，按注册顺序依次执行。
// ============================================================================

// configListingWrapDemo 配置 Wrap* 方法演示模块
func ConfigListingWrapDemo(b *presets.Builder, db *gorm.DB) {
	db.AutoMigrate(&models.ListingWrapDemo{})
	seedListingWrapDemo(db)

	mb := b.Model(&models.ListingWrapDemo{}).URIName("listing-wrap-demo").MenuIcon("layers")
	lb := mb.Listing("ID", "Title", "Status", "Priority", "Assignee", "Category", "UpdatedAt").
		SearchColumns("title", "assignee").
		PerPage(20)

	// ========================================================================
	// 1. WrapColumns — 动态修改列配置
	// ========================================================================
	//
	// 用途：修改列标签、表头样式、列顺序等。
	// 框架提供两个辅助函数：CustomizeColumnLabel 和 CustomizeColumnHeader。
	//
	// 支持链式多次调用，每次装饰上一次的结果。

	// 1a. CustomizeColumnLabel — 批量修改列标签文本
	// 常用于 i18n 国际化，mapper 返回 {字段名: 显示名} 映射
	lb.WrapColumns(presets.CustomizeColumnLabel(func(evCtx *web.EventContext) (map[string]string, error) {
		return map[string]string{
			"ID":        "编号",
			"Title":     "标题",
			"Status":    "状态",
			"Priority":  "优先级",
			"Assignee":  "负责人",
			"Category":  "分类",
			"UpdatedAt": "更新时间",
		}, nil
	}))

	// 1b. CustomizeColumnHeader — 自定义指定列的表头渲染
	// 第二个参数指定作用的列名，省略则作用于所有列
	lb.WrapColumns(presets.CustomizeColumnHeader(func(evCtx *web.EventContext, col *presets.Column, th h.MutableAttrHTMLComponent) (h.MutableAttrHTMLComponent, error) {
		// 给「标题」列设置最小宽度
		th.SetAttr(":class", "'min-w-40'")
		return th, nil
	}, "Title"))

	// 1c. 自定义 ColumnsProcessor — 完全自定义列处理逻辑
	// 示例：动态隐藏某些列（如根据权限控制列可见性）
	lb.WrapColumns(func(in presets.ColumnsProcessor) presets.ColumnsProcessor {
		return func(evCtx *web.EventContext, columns []*presets.Column) ([]*presets.Column, error) {
			columns, err := in(evCtx, columns)
			if err != nil {
				return nil, err
			}
			// 示例：可以根据 evCtx 中的用户信息过滤列
			// 这里演示给「更新时间」列也加最小宽度
			for _, col := range columns {
				if col.Name == "UpdatedAt" {
					w := col.WrapHeader
					col.WrapHeader = func(evCtx *web.EventContext, col *presets.Column, th h.MutableAttrHTMLComponent) (h.MutableAttrHTMLComponent, error) {
						if w != nil {
							var err error
							th, err = w(evCtx, col, th)
							if err != nil {
								return nil, err
							}
						}
						th.SetAttr(":class", "'min-w-36'")
						return th, nil
					}
				}
			}
			return columns, nil
		}
	})

	// ========================================================================
	// 2. WrapCell — 包装单元格
	// ========================================================================
	//
	// 用途：给单元格添加点击事件、条件样式、工具提示等。
	// 参数：(evCtx, cell, id, obj) → cell
	//   - cell: 当前单元格组件（<td>）
	//   - id:   行 ID（主键）
	//   - obj:  行数据对象
	//
	// 常见场景：
	//   - 点击单元格触发事件（如对话框列表中的选择）
	//   - 根据数据内容给单元格加条件样式

	lb.WrapCell(func(in presets.CellProcessor) presets.CellProcessor {
		return func(evCtx *web.EventContext, cell h.MutableAttrHTMLComponent, id string, obj any) (h.MutableAttrHTMLComponent, error) {
			demo := obj.(*models.ListingWrapDemo)
			// 根据状态给单元格添加不同样式
			switch demo.Status {
			case "archived":
				cell.SetAttr(":class", "'opacity-50'")
			}
			return in(evCtx, cell, id, obj)
		}
	})

	// ========================================================================
	// 3. WrapRow — 包装行
	// ========================================================================
	//
	// 用途：给整行添加条件高亮、点击事件、CSS 类等。
	// 参数与 WrapCell 相同：(evCtx, row, id, obj) → row
	//
	// 常见场景：
	//   - 高优先级行高亮显示
	//   - 行点击跳转详情
	//   - 根据状态显示不同背景色

	lb.WrapRow(func(in presets.RowProcessor) presets.RowProcessor {
		return func(evCtx *web.EventContext, row h.MutableAttrHTMLComponent, id string, obj any) (h.MutableAttrHTMLComponent, error) {
			demo := obj.(*models.ListingWrapDemo)
			// 高优先级行添加左边框高亮
			if demo.Priority == 3 {
				row.SetAttr(":class", "'border-l-2 border-l-destructive'")
			}
			return in(evCtx, row, id, obj)
		}
	})

	// ========================================================================
	// 4. WrapFilterDataFunc — 动态追加或修改筛选项
	// ========================================================================
	//
	// 用途：在已有筛选项基础上追加新筛选项，或修改现有筛选项的条件。
	// 适用于模块化场景：基础模块定义筛选项，插件模块追加额外筛选。
	//
	// 注意：需要先通过 FilterDataFunc 设置基础筛选项，
	// 然后用 WrapFilterDataFunc 在其基础上追加。

	// 先设置基础筛选项
	lb.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		return []*shadcn.FilterItem{
			{
				Key:          "status",
				Label:        "状态",
				ItemType:     shadcn.FilterItemTypeSelect,
				SQLCondition: `status %s ?`,
				Options: []shadcn.FilterSelectOption{
					{Text: "草稿", Value: "draft"},
					{Text: "活跃", Value: "active"},
					{Text: "归档", Value: "archived"},
				},
			},
			{
				Key:          "category",
				Label:        "分类",
				ItemType:     shadcn.FilterItemTypeSelect,
				SQLCondition: `category %s ?`,
				Options: []shadcn.FilterSelectOption{
					{Text: "技术", Value: "tech"},
					{Text: "产品", Value: "product"},
					{Text: "运营", Value: "ops"},
				},
			},
		}
	})

	// 用 WrapFilterDataFunc 追加优先级筛选项
	lb.WrapFilterDataFunc(func(in presets.FilterDataFunc) presets.FilterDataFunc {
		return func(ctx *web.EventContext) shadcn.FilterData {
			// 先获取已有筛选项
			fd := in(ctx)
			// 追加新的筛选项
			fd = append(fd, &shadcn.FilterItem{
				Key:          "priority",
				Label:        "优先级",
				ItemType:     shadcn.FilterItemTypeSelect,
				SQLCondition: `priority %s ?`,
				Options: []shadcn.FilterSelectOption{
					{Text: "低", Value: "1"},
					{Text: "中", Value: "2"},
					{Text: "高", Value: "3"},
				},
			})
			return fd
		}
	})

	// ========================================================================
	// 5. WrapNewButtonFunc — 修改「新建」按钮
	// ========================================================================
	//
	// 用途：修改、隐藏或替换默认的「新建」按钮。
	// 常见场景：
	//   - 根据权限隐藏新建按钮
	//   - 替换为自定义按钮（如批量导入）
	//   - 在新建按钮旁添加额外操作

	lb.WrapNewButtonFunc(func(in presets.ComponentFunc) presets.ComponentFunc {
		return func(ctx *web.EventContext) h.HTMLComponent {
			// 获取原始新建按钮
			originalBtn := in(ctx)
			if originalBtn == nil {
				return nil
			}
			// 在原始按钮旁添加一个导入按钮
			return h.Div(
				originalBtn,
				shadcn.Button(h.Text("导入")).
					Variant(shadcn.ButtonVariantOutline).
					Size(shadcn.ButtonSizeSm).
					Class("ml-2"),
			).Class("flex items-center")
		}
	})

	// ========================================================================
	// 6. WrapSearchFunc — 包装搜索逻辑
	// ========================================================================
	//
	// 用途：在默认搜索前后添加自定义逻辑。
	// 常见场景：
	//   - 追加默认排序（如按优先级降序）
	//   - 数据隔离（只显示当前用户的数据）
	//   - 追加 SQL 条件（如软删除过滤）
	//   - 搜索后处理（如聚合统计）

	// 6a. 追加默认排序：按优先级降序
	lb.WrapSearchFunc(func(in presets.SearchFunc) presets.SearchFunc {
		return func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
			// 在搜索前追加排序条件
			params.OrderBys = append(params.OrderBys, relay.OrderBy{
				Field: "Priority",
				Desc:  true,
			})
			return in(ctx, params)
		}
	})

	// 6b. 隐藏已归档数据（除非筛选器明确选择了 archived）
	lb.WrapSearchFunc(func(in presets.SearchFunc) presets.SearchFunc {
		return func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
			// 检查筛选器是否选择了 status
			hasStatusFilter := false
			for _, cond := range params.SQLConditions {
				if cond.Query == "status = ?" {
					hasStatusFilter = true
					break
				}
			}
			// 如果没有明确筛选状态，默认排除 archived
			if !hasStatusFilter {
				params.SQLConditions = append(params.SQLConditions, &presets.SQLCondition{
					Query: "status != ?",
					Args:  []any{"archived"},
				})
			}
			return in(ctx, params)
		}
	})

	// ========================================================================
	// 编辑配置（简单配置即可，重点是列表）
	// ========================================================================
	ed := mb.Editing("Title", "Status", "Priority", "Assignee", "Category")

	ed.Field("Status").Label("状态").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return shadcn.Select().
			Items([]shadcn.DefaultOptionItem{
				{Text: "草稿", Value: "draft"},
				{Text: "活跃", Value: "active"},
				{Text: "归档", Value: "archived"},
			}).
			Label(field.Label).
			Attr(web.VField(field.FormKey, field.Value(obj))...).
			ErrorMessages(field.Errors...)
	})

	ed.Field("Priority").Label("优先级").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return shadcn.Select().
			Items([]shadcn.DefaultOptionItem{
				{Text: "低", Value: "1"},
				{Text: "中", Value: "2"},
				{Text: "高", Value: "3"},
			}).
			Label(field.Label).
			Attr(web.VField(field.FormKey, field.Value(obj))...).
			ErrorMessages(field.Errors...)
	})

	ed.Field("Category").Label("分类").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return shadcn.Select().
			Items([]shadcn.DefaultOptionItem{
				{Text: "技术", Value: "tech"},
				{Text: "产品", Value: "product"},
				{Text: "运营", Value: "ops"},
			}).
			Label(field.Label).
			Attr(web.VField(field.FormKey, field.Value(obj))...).
			ErrorMessages(field.Errors...)
	})

	// 列表自定义字段：状态显示为 Badge
	lb.Field("Status").Label("状态").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			demo := obj.(*models.ListingWrapDemo)
			var variant shadcn.BadgeVariant
			var text string
			switch demo.Status {
			case "draft":
				variant = shadcn.BadgeVariantSecondary
				text = "草稿"
			case "active":
				variant = shadcn.BadgeVariantDefault
				text = "活跃"
			case "archived":
				variant = shadcn.BadgeVariantOutline
				text = "归档"
			default:
				variant = shadcn.BadgeVariantSecondary
				text = demo.Status
			}
			return h.Td(shadcn.Badge(h.Text(text)).Variant(variant))
		})

	// 列表自定义字段：优先级显示为彩色文字
	lb.Field("Priority").Label("优先级").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			demo := obj.(*models.ListingWrapDemo)
			var text, class string
			switch demo.Priority {
			case 1:
				text, class = "低", "text-muted-foreground"
			case 2:
				text, class = "中", "text-foreground"
			case 3:
				text, class = "高", "text-destructive font-medium"
			default:
				text = fmt.Sprint(demo.Priority)
			}
			return h.Td(h.Tag("span").Class(class).Children(h.Text(text)))
		})
}

// seedListingWrapDemo 填充演示数据
func seedListingWrapDemo(db *gorm.DB) {
	var count int64
	db.Model(&models.ListingWrapDemo{}).Count(&count)
	if count > 0 {
		return
	}

	demos := []models.ListingWrapDemo{
		{Title: "重构用户认证模块", Status: "active", Priority: 3, Assignee: "张三", Category: "tech"},
		{Title: "优化首页加载速度", Status: "active", Priority: 2, Assignee: "李四", Category: "tech"},
		{Title: "设计新版产品页面", Status: "draft", Priority: 2, Assignee: "王五", Category: "product"},
		{Title: "编写 Q4 运营报告", Status: "draft", Priority: 1, Assignee: "赵六", Category: "ops"},
		{Title: "修复支付回调 bug", Status: "active", Priority: 3, Assignee: "张三", Category: "tech"},
		{Title: "上线促销活动页", Status: "archived", Priority: 2, Assignee: "李四", Category: "ops"},
		{Title: "数据库索引优化", Status: "active", Priority: 1, Assignee: "王五", Category: "tech"},
		{Title: "产品需求文档更新", Status: "draft", Priority: 1, Assignee: "赵六", Category: "product"},
		{Title: "迁移旧版 API", Status: "archived", Priority: 3, Assignee: "张三", Category: "tech"},
		{Title: "客户反馈分析", Status: "active", Priority: 2, Assignee: "李四", Category: "ops"},
	}
	db.Create(&demos)
}
