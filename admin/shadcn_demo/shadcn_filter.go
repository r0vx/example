package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnFilterDemo 虚拟模型
type ShadcnFilterDemo struct{}

// configFilter 注册 Filter demo
func configFilter(b *presets.Builder) {
	m := b.Model(&ShadcnFilterDemo{}).
		Label("Filter").
		URIName("shadcn-filter")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnFilterDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Filter"
		r.Body = shadcnFilterBody(ctx)
		return
	})
}

// shadcnFilterBody Filter 组件演示
func shadcnFilterBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		// 页面标题
		h.Div(
			h.H1("Filter 组件").Style("font-size: 28px; font-weight: bold; margin-bottom: 8px;"),
			h.P(h.Text("强大的筛选组件,支持多种数据类型和筛选方式")).Style("color: #666; margin-bottom: 24px;"),
		),

		// 完整示例
		h.Div(
			h.H2("完整示例").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("展示所有筛选类型:文本、数字、单选、多选、自动完成、日期、日期范围")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),

			web.Scope(
				Filter().Items(
					// 文本筛选
					StringFilterItem("name", "Name"),

					// 数字筛选
					NumberFilterItem("price", "Price"),

					// 单选筛选
					SelectFilterItem("status", "Status", []FilterSelectOption{
						{Text: "Active", Value: "active"},
						{Text: "Inactive", Value: "inactive"},
						{Text: "Pending", Value: "pending"},
					}),

					// 多选筛选
					MultipleSelectFilterItem("tags", "Tags", []FilterSelectOption{
						{Text: "Featured", Value: "featured"},
						{Text: "New", Value: "new"},
						{Text: "Sale", Value: "sale"},
						{Text: "Popular", Value: "popular"},
					}),

					// 自动完成筛选
					AutoCompleteFilterItem("category", "Category", []FilterSelectOption{
						{Text: "Electronics", Value: "electronics"},
						{Text: "Clothing", Value: "clothing"},
						{Text: "Books", Value: "books"},
						{Text: "Home & Garden", Value: "home"},
						{Text: "Sports", Value: "sports"},
					}),

					// 日期筛选
					DateFilterItem("created_at", "Created Date"),

					// 日期范围筛选(默认折叠)
					func() FilterItem {
						item := DateRangeFilterItem("updated_at", "Updated Range")
						item.Folded = true
						return item
					}(),
				).On("change", "console.log('Filter changed:', $event)"),

				// 显示当前筛选值
				h.Div(
					h.Div(
						h.H3("当前筛选条件:").Style("font-weight: 600; margin-bottom: 8px;"),
						h.Tag("pre").Children(
							h.Text("{{ JSON.stringify(locals.filterValue || {}, null, 2) }}"),
						).Style("background: #f3f4f6; padding: 12px; border-radius: 4px; font-size: 12px; overflow-x: auto;"),
					).Class("mt-4"),
				),
			).Init(`{filterValue: null}`).VSlot("{ locals }"),
		).Class("demo-section"),

		// 基础文本筛选
		h.Div(
			h.H2("基础文本筛选").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("简单的文本输入筛选")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),

			Filter().Items(
				StringFilterItem("search", "Search"),
			).On("change", "console.log('Search:', $event)"),
		).Class("demo-section"),

		// 数字范围筛选
		h.Div(
			h.H2("数字范围筛选").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("支持最小值和最大值")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),

			Filter().Items(
				NumberFilterItem("min_price", "Min Price"),
				NumberFilterItem("max_price", "Max Price"),
			).On("change", "console.log('Price range:', $event)"),
		).Class("demo-section"),

		// 选择筛选
		h.Div(
			h.H2("选择筛选").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("单选和多选下拉列表")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),

			Filter().Items(
				SelectFilterItem("department", "Department", []FilterSelectOption{
					{Text: "Engineering", Value: "eng"},
					{Text: "Sales", Value: "sales"},
					{Text: "Marketing", Value: "marketing"},
					{Text: "HR", Value: "hr"},
				}),
				MultipleSelectFilterItem("skills", "Skills", []FilterSelectOption{
					{Text: "Go", Value: "go"},
					{Text: "Vue", Value: "vue"},
					{Text: "TypeScript", Value: "ts"},
					{Text: "Python", Value: "python"},
				}),
			).On("change", "console.log('Selection:', $event)"),
		).Class("demo-section"),

		// 日期筛选
		h.Div(
			h.H2("日期筛选").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("单日期和日期范围选择")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),

			Filter().Items(
				DateFilterItem("start_date", "Start Date"),
				DateRangeFilterItem("date_range", "Date Range"),
			).On("change", "console.log('Date:', $event)"),
		).Class("demo-section"),

		// 日期选择器（带弹出日历）
		h.Div(
			h.H2("日期选择器（带弹出日历）").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("带有弹出日历面板的日期选择器")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),

			Filter().Items(
				DatePickerFilterItem("birth_date", "Birth Date"),
				DateRangePickerFilterItem("vacation", "Vacation Period"),
			).On("change", "console.log('DatePicker:', $event)"),
		).Class("demo-section"),

		// 日期时间范围
		h.Div(
			h.H2("日期时间范围").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("选择日期和时间范围")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),

			Filter().Items(
				DatetimeRangeFilterItem("event_time", "Event Time"),
			).On("change", "console.log('DatetimeRange:', $event)"),
		).Class("demo-section"),

		// 级联选择
		h.Div(
			h.H2("级联选择").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.P(h.Text("多级联动选择，适用于省市区等层级数据")).Style("color: #666; font-size: 14px; margin-bottom: 16px;"),

			Filter().Items(
				LinkageSelectFilterItem("location", "Location",
					[][]FilterLinkageItem{
						// 第一级：省份
						{
							{ID: "zj", Name: "浙江省", ChildrenIDs: []string{"hz", "nb"}},
							{ID: "js", Name: "江苏省", ChildrenIDs: []string{"nj", "sz"}},
						},
						// 第二级：城市
						{
							{ID: "hz", Name: "杭州市", ChildrenIDs: []string{"xihu", "binjiang"}},
							{ID: "nb", Name: "宁波市", ChildrenIDs: []string{"haishu", "jiangbei"}},
							{ID: "nj", Name: "南京市", ChildrenIDs: []string{"xuanwu", "gulou"}},
							{ID: "sz", Name: "苏州市", ChildrenIDs: []string{"gusu", "wuzhong"}},
						},
						// 第三级：区县
						{
							{ID: "xihu", Name: "西湖区"},
							{ID: "binjiang", Name: "滨江区"},
							{ID: "haishu", Name: "海曙区"},
							{ID: "jiangbei", Name: "江北区"},
							{ID: "xuanwu", Name: "玄武区"},
							{ID: "gulou", Name: "鼓楼区"},
							{ID: "gusu", Name: "姑苏区"},
							{ID: "wuzhong", Name: "吴中区"},
						},
					},
					[]string{"省份", "城市", "区县"},
				),
			).On("change", "console.log('Linkage:', $event)"),
		).Class("demo-section"),

		// 功能特性说明
		h.Div(
			h.H2("功能特性").Style("font-size: 18px; font-weight: 600; margin-bottom: 12px;"),
			h.Ul(
				h.Li(h.Text("StringFilterItem - 文本输入筛选")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("NumberFilterItem - 数字输入筛选")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("SelectFilterItem - 单选下拉筛选")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("MultipleSelectFilterItem - 多选下拉筛选")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("AutoCompleteFilterItem - 自动完成筛选")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("DateFilterItem - 单日期选择")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("DateRangeFilterItem - 日期范围选择")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("DatePickerFilterItem - 日期选择器（带弹出日历）")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("DateRangePickerFilterItem - 日期范围选择器（带弹出日历）")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("DatetimeRangeFilterItem - 日期时间范围选择")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("LinkageSelectFilterItem - 级联选择")).Style("padding: 8px 0; border-bottom: 1px solid #f3f4f6;"),
				h.Li(h.Text("Folded - 支持默认折叠状态")).Style("padding: 8px 0;"),
			).Style("list-style: none; padding: 0; margin: 0;"),
		).Class("demo-section").Style("background: #eff6ff;"),
	)
}
