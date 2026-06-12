package ui_demo

import (
	"fmt"

	"example/models"

	"github.com/r0vx/admin/autocomplete"
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// ============================================================================
// admin/autocomplete 包功能演示
// ============================================================================
//
// ## 概述
//
// admin/autocomplete 包提供「数据库驱动的自动补全 JSON API」，核心功能：
//   - 自动为 GORM 模型创建 HTTP JSON 端点（如 /complete/products）
//   - 支持搜索（SQLCondition）、分页（Paging）、排序（OrderBy）
//   - 配合 shadcn.Autocomplete 组件实现远程搜索下拉
//   - 配合 FilterItemTypeAutoComplete 实现列表页筛选
//
// ## 包结构
//
//   autocomplete.New()           → 创建 Builder（配置 DB、Prefix 等）
//   builder.Model(&MyModel{})   → 注册模型，返回 ModelBuilder
//   modelBuilder.Columns(...)    → 指定 JSON 返回的字段列
//   modelBuilder.SQLCondition()  → 设置搜索 SQL（如 "name ILIKE ?"）
//   modelBuilder.OrderBy()       → 设置排序规则
//   modelBuilder.Paging(true)    → 启用分页（返回 pages/current 等字段）
//   modelBuilder.JsonHref()      → 获取 API 端点路径
//
// ## API 请求/响应格式
//
//   请求: GET /complete/products?search=手机&page=1&pageSize=5
//
//   响应:
//   {
//     "data": [
//       {"id": 1, "name": "iPhone 15", "code": "SKU-001"},
//       {"id": 2, "name": "华为手机", "code": "SKU-002"}
//     ],
//     "total": 42,     // 匹配的总记录数
//     "pages": 9,      // 总页数（启用 Paging 时）
//     "current": 5     // 当前已加载到的记录位置
//   }
//
// ## 两种使用方式
//
//   方式1: 编辑字段 — 用 shadcn.Autocomplete + Remote(true) + RemoteURL（声明式远程搜索）
//   方式2: 列表筛选 — 用 FilterItemTypeAutoComplete + AutocompleteDataSource
//
// ============================================================================

// AutocompletePrefix 自动补全 API 统一路径前缀
const AutocompletePrefix = "/complete"

// configAutocompleteDemo 配置 autocomplete 演示模块
//
// 做两件事：
//  1. 创建 autocomplete.Builder 并注册 Product/Category 模型的 API 端点
//  2. 创建 AutocompleteDemo 模型，演示编辑字段和列表筛选的 autocomplete 集成
//
// 返回 autocomplete.Builder（实现 http.Handler），需要在 Router 中挂载到 /complete/ 路径
func ConfigAutocompleteDemo(b *presets.Builder, db *gorm.DB) *autocomplete.Builder {
	// ========================================================================
	// Part 1: 创建 autocomplete API 端点
	// ========================================================================

	// 创建 Builder，配置数据库连接和 URL 前缀
	// Prefix 决定所有端点的根路径，如 Prefix("/complete") → /complete/products
	ab := autocomplete.New().
		DB(db).
		Prefix(AutocompletePrefix).
		AllowCrossOrigin(true) // 允许跨域（开发环境常用）

	// 注册 Product 模型的自动补全端点
	// 端点路径: GET /complete/products?search=xxx&page=1&pageSize=5
	// 返回字段: id, name, code（由 Columns 指定）
	// 搜索条件: name ILIKE '%xxx%'（由 SQLCondition 指定，? 会被替换为 %search%）
	// 排序规则: 按 name 升序
	// 分页: 启用（返回 pages/current 字段）
	ab.Model(&models.Product{}).
		Columns("id", "name", "code"). // 指定返回的数据库列（SELECT id, name, code FROM products）
		SQLCondition("name ILIKE ?").  // 搜索 SQL 条件，search 参数会被包装成 %search%
		OrderBy("name asc").           // 排序规则
		Paging(true)                   // 启用分页，返回 pages/current/total

	// 注册 Category 模型的自动补全端点
	// 端点路径: GET /complete/categories?search=xxx
	// 返回字段: id, name
	// 搜索条件: name ILIKE '%xxx%'
	// 排序规则: 按 id 降序
	ab.Model(&models.Category{}).
		Columns("id", "name").
		SQLCondition("name ILIKE ?").
		OrderBy("id desc").
		Paging(true)

	// ========================================================================
	// Part 2: 创建 AutocompleteDemo 模型的 CRUD 配置
	// ========================================================================

	// 数据库迁移
	db.AutoMigrate(&models.AutocompleteDemo{})

	mb := b.Model(&models.AutocompleteDemo{})

	// ------------------------------------------------------------------------
	// 列表配置
	// ------------------------------------------------------------------------
	listing := mb.Listing("ID", "Title", "ProductName", "CategoryName", "UpdatedAt").
		SearchColumns("title").
		PerPage(20)

	// 列表筛选 — 使用 FilterItemTypeAutoComplete + AutocompleteDataSource
	//
	// NewDefaultAutocompleteDataSource 创建一个标准数据源配置，字段映射如下：
	//   RemoteURL   → API 端点 URL（如 /complete/products）
	//   IsPaging    → 启用分页加载
	//   ItemTitle   → 响应中作为显示文本的字段名（默认 "title"，此处需改为 "name"）
	//   ItemValue   → 响应中作为值的字段名（默认 "id"）
	//   Separator   → 值分隔符（默认 "__"），存储格式为 "title__value"
	//   SearchField → 搜索参数名（默认 "search"）
	//   PageField   → 分页参数名（默认 "page"）
	//   ItemsField  → 响应中数据数组的字段名（默认 "data"）
	//   TotalField  → 响应中总数的字段名（默认 "total"）
	listing.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		// 创建 Product 的筛选数据源
		productDS := autocomplete.NewDefaultAutocompleteDataSource(AutocompletePrefix + "/products")
		productDS.ItemTitle = "name" // 默认是 "title"，Product 模型用 "name" 字段显示

		// 创建 Category 的筛选数据源
		categoryDS := autocomplete.NewDefaultAutocompleteDataSource(AutocompletePrefix + "/categories")
		categoryDS.ItemTitle = "name"

		return []*shadcn.FilterItem{
			{
				Key:      "title",
				Label:    "标题",
				ItemType: shadcn.FilterItemTypeString,
				// SQL 条件中的 %s 会被替换为比较运算符（= / LIKE 等）
				SQLCondition: `title %s ?`,
			},
			{
				Key:      "product_id",
				Label:    "关联产品（远程搜索）",
				ItemType: shadcn.FilterItemTypeAutoComplete,
				// SQLCondition 中的 %s 会被替换为 =
				// 筛选值是 AutocompleteDataSource.Separator 分隔的 "title__id"
				// 框架会自动解析出 id 部分用于 SQL 查询
				SQLCondition:           `product_id %s ?`,
				AutocompleteDataSource: productDS,
			},
			{
				Key:                    "category_id",
				Label:                  "关联分类（远程搜索）",
				ItemType:               shadcn.FilterItemTypeAutoComplete,
				SQLCondition:           `category_id %s ?`,
				AutocompleteDataSource: categoryDS,
			},
		}
	})

	// 列表自定义字段：显示关联的产品/分类名称
	listing.Field("ProductName").Label("产品").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			demo := obj.(*models.AutocompleteDemo)
			return h.Td(h.Text(demo.ProductName))
		})

	listing.Field("CategoryName").Label("分类").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			demo := obj.(*models.AutocompleteDemo)
			return h.Td(h.Text(demo.CategoryName))
		})

	// ------------------------------------------------------------------------
	// 编辑配置
	// ------------------------------------------------------------------------
	ed := mb.Editing("Title", "AssigneeID", "ProductID", "ProductIDs", "CategoryID")

	// Title — 基本输入框
	ed.Field("Title").Label("标题").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Input().
				Label(field.Label).
				Placeholder("请输入标题").
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	// AssigneeID — 静态 items + icon 头像演示
	// 展示 DefaultOptionItem 的 Icon 字段，每个选项前显示圆形头像
	ed.Field("AssigneeID").Label("负责人").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			// 静态用户列表，带头像 URL
			items := []shadcn.DefaultOptionItem{
				{Text: "张三", Value: "zhangsan", Icon: "https://api.dicebear.com/9.x/avataaars/svg?seed=zhangsan"},
				{Text: "李四", Value: "lisi", Icon: "https://api.dicebear.com/9.x/avataaars/svg?seed=lisi"},
				{Text: "王五", Value: "wangwu", Icon: "https://api.dicebear.com/9.x/avataaars/svg?seed=wangwu"},
				{Text: "赵六", Value: "zhaoliu", Icon: "https://api.dicebear.com/9.x/avataaars/svg?seed=zhaoliu"},
			}

			return shadcn.Autocomplete().
				Label(field.Label).
				Placeholder("选择负责人...").
				Items(items).
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	// ProductID — 通过 autocomplete HTTP API + RemoteURL 实现远程搜索
	//
	// 使用声明式 RemoteURL 模式：
	//   只需指定 API 地址和字段映射，Vue 组件内部自动处理 fetch + 数据转换。
	//   相比 RemoteMethod（需要传递 JS 函数表达式），RemoteURL 模式更适合服务端渲染。
	//
	// 数据流：
	//   用户输入 → Vue 组件自动 fetch(remoteUrl + ?search=keyword)
	//   → 后端 autocomplete.ModelBuilder.ServeHTTP 查询数据库
	//   → 返回 JSON {data: [{id, name, code}, ...], total, pages, current}
	//   → Vue 组件根据 remoteItemValue/remoteItemText 映射为选项
	//
	// API 端点由 Part 1 注册的 ab.Model(&models.Product{}) 自动生成。
	ed.Field("ProductID").Label("关联产品").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			demo := obj.(*models.AutocompleteDemo)

			// 编辑时回显：查询已选中的产品信息，作为初始 items
			var selectedItems []shadcn.DefaultOptionItem
			if demo.ProductID > 0 {
				var p models.Product
				if err := db.Select("id, name").First(&p, demo.ProductID).Error; err == nil {
					selectedItems = append(selectedItems, shadcn.DefaultOptionItem{
						Text:  p.Name,
						Value: fmt.Sprint(p.ID),
					})
				}
			}

			return shadcn.Autocomplete().
				Label(field.Label).
				Placeholder("搜索产品名称...").
				Remote(true).                                         // 启用远程搜索模式
				RemoteURL(AutocompletePrefix + "/products").          // 声明式：API 地址
				RemoteItemValue("id").                                // 响应中 value 取 id 字段
				RemoteItemText("name").                               // 响应中 text 取 name 字段
				Items(selectedItems).                                 // 编辑回显用
				Attr(web.VField(field.FormKey, field.Value(obj))...). // v-model 绑定
				ErrorMessages(field.Errors...)
		})

	// ProductIDs — 多选模式演示（Multiple(true)）
	// 值以 JSON 数组字符串存储，如 ["1","3"]
	ed.Field("ProductIDs").Label("关联产品（多选）").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			demo := obj.(*models.AutocompleteDemo)

			// 解析已存储的 JSON 数组，查询回显数据
			ac := shadcn.Autocomplete()
			ids := ac.ParseValue(demo.ProductIDs)
			var selectedItems []shadcn.DefaultOptionItem
			if len(ids) > 0 {
				var products []models.Product
				db.Select("id, name").Where("id IN ?", ids).Find(&products)
				for _, p := range products {
					selectedItems = append(selectedItems, shadcn.DefaultOptionItem{
						Text:  p.Name,
						Value: fmt.Sprint(p.ID),
					})
				}
			}

			return ac.
				Label(field.Label).
				Placeholder("搜索并选择多个产品...").
				Remote(true).
				RemoteURL(AutocompletePrefix + "/products").
				RemoteItemValue("id").
				RemoteItemText("name").
				Items(selectedItems).
				Multiple(true).                          // 启用多选
				Attr(web.VField(field.FormKey, ids)...). // 传 []string，Vue 端接收数组
				ErrorMessages(field.Errors...)
		}).SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		// Plaid 将数组值作为重复 form field 提交（ProductIDs=1&ProductIDs=3）
		// 用 ctx.R.Form 获取所有值，再用 FormatValue 转为 JSON 数组字符串存储
		demo := obj.(*models.AutocompleteDemo)
		values := ctx.R.Form[field.Name]
		demo.ProductIDs = shadcn.Autocomplete().FormatValue(values)
		return
	})

	// CategoryID — 同理，使用 RemoteURL 声明式远程搜索
	ed.Field("CategoryID").Label("关联分类").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			demo := obj.(*models.AutocompleteDemo)

			var selectedItems []shadcn.DefaultOptionItem
			if demo.CategoryID > 0 {
				var c models.Category
				if err := db.Select("id, name").First(&c, demo.CategoryID).Error; err == nil {
					selectedItems = append(selectedItems, shadcn.DefaultOptionItem{
						Text:  c.Name,
						Value: fmt.Sprint(c.ID),
					})
				}
			}

			return shadcn.Autocomplete().
				Label(field.Label).
				Placeholder("搜索分类名称...").
				Remote(true).
				RemoteURL(AutocompletePrefix + "/categories").
				RemoteItemValue("id").
				RemoteItemText("name").
				Items(selectedItems).
				Attr(web.VField(field.FormKey, field.Value(obj))...).
				ErrorMessages(field.Errors...)
		})

	// 保存前钩子：写入冗余名称字段（方便列表展示）
	ed.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj any, id string, ctx *web.EventContext) error {
			demo := obj.(*models.AutocompleteDemo)
			if demo.ProductID > 0 {
				var p models.Product
				if db.Select("name").First(&p, demo.ProductID).Error == nil {
					demo.ProductName = p.Name
				}
			}
			if demo.CategoryID > 0 {
				var c models.Category
				if db.Select("name").First(&c, demo.CategoryID).Error == nil {
					demo.CategoryName = c.Name
				}
			}
			return in(obj, id, ctx)
		}
	})

	// 验证
	ed.ValidateFunc(func(obj any, ctx *web.EventContext) (err web.ValidationErrors) {
		demo := obj.(*models.AutocompleteDemo)
		if demo.Title == "" {
			err.FieldError("Title", "标题不能为空")
		}
		return
	})

	return ab
}

// ============================================================================
// 备注：三种远程搜索模式对比
// ============================================================================
//
// 1. RemoteURL 模式（本 demo 使用，推荐）
//    Go 端声明式配置 URL + 字段映射，Vue 组件内部自动 fetch。
//    适用于 admin/autocomplete 包生成的标准 JSON API。
//    示例: Remote(true).RemoteURL("/complete/products").RemoteItemValue("id").RemoteItemText("name")
//
// 2. RemoteMethod 模式
//    传递 JS 函数表达式给 Vue 组件的 :remote-method 属性。
//    需要注意 HTML 属性中的引号转义问题。
//    示例: Remote(true).RemoteMethod("async (q) => { ... }")
//
// 3. EventFunc 模式
//    使用 Plaid EventFunc 注册服务端事件处理搜索。
//    参考 category_config.go 中 productsSelector 的实现。
// ============================================================================
