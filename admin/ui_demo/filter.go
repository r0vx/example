package ui_demo

import (
	"example/models"
	"fmt"

	"github.com/r0vx/admin/presets"
	shadcn "github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/web"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// configFilterDemo 配置 Filter 全类型演示页面
// 展示所有 12 种 FilterItemType 的筛选组件效果
func ConfigFilterDemo(b *presets.Builder, db *gorm.DB) {
	// 种子数据
	seedFilterDemoData(db)

	mb := b.Model(&models.FilterDemo{}).URIName("filter-demos").Label("Filter Demos")

	// 列表展示
	listing := mb.Listing("ID", "Title", "Amount", "Status", "Category", "IsActive", "Country", "Province", "City", "CreatedAt")
	listing.SearchColumns("title")
	listing.PerPage(20)

	// 金额格式化
	listing.Field("Amount").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		d := obj.(*models.FilterDemo)
		return h.Text(fmt.Sprintf("¥%.2f", d.Amount))
	})

	// ========== 所有 12 种 FilterItemType ==========
	listing.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		// 状态选项（Select / MultipleSelect 共用）
		statusOptions := []shadcn.FilterSelectOption{
			{Text: "草稿", Value: "draft"},
			{Text: "已发布", Value: "published"},
			{Text: "已归档", Value: "archived"},
		}

		// 分类选项（AutoComplete 使用）
		categoryOptions := []shadcn.FilterSelectOption{
			{Text: "电子产品", Value: "electronics"},
			{Text: "服装", Value: "clothing"},
			{Text: "食品", Value: "food"},
			{Text: "家居", Value: "home"},
			{Text: "运动", Value: "sports"},
		}

		// 级联选择数据（LinkageSelect 使用）
		// 第一级：国家
		level1 := []shadcn.FilterLinkageItem{
			{ID: "china", Name: "中国", ChildrenIDs: []string{"beijing", "shanghai", "guangdong"}},
			{ID: "usa", Name: "美国", ChildrenIDs: []string{"california", "newyork"}},
		}
		// 第二级：省/州
		level2 := []shadcn.FilterLinkageItem{
			{ID: "beijing", Name: "北京", ChildrenIDs: []string{"chaoyang", "haidian"}},
			{ID: "shanghai", Name: "上海", ChildrenIDs: []string{"pudong", "xuhui"}},
			{ID: "guangdong", Name: "广东", ChildrenIDs: []string{"guangzhou", "shenzhen"}},
			{ID: "california", Name: "加州", ChildrenIDs: []string{"losangeles", "sanfrancisco"}},
			{ID: "newyork", Name: "纽约州", ChildrenIDs: []string{"nyc", "buffalo"}},
		}
		// 第三级：城市
		level3 := []shadcn.FilterLinkageItem{
			{ID: "chaoyang", Name: "朝阳区"},
			{ID: "haidian", Name: "海淀区"},
			{ID: "pudong", Name: "浦东新区"},
			{ID: "xuhui", Name: "徐汇区"},
			{ID: "guangzhou", Name: "广州"},
			{ID: "shenzhen", Name: "深圳"},
			{ID: "losangeles", Name: "洛杉矶"},
			{ID: "sanfrancisco", Name: "旧金山"},
			{ID: "nyc", Name: "纽约市"},
			{ID: "buffalo", Name: "布法罗"},
		}

		return []*shadcn.FilterItem{
			// 1. StringItem — 字符串筛选（支持等于/包含）
			{
				Key:          "title",
				Label:        "1. String（标题）",
				ItemType:     shadcn.FilterItemTypeString,
				SQLCondition: `title %s ?`,
			},
			// 2. NumberItem — 数字筛选（支持等于/大于/小于/区间）
			{
				Key:          "amount",
				Label:        "2. Number（金额）",
				ItemType:     shadcn.FilterItemTypeNumber,
				SQLCondition: `amount %s ?`,
			},
			// 3. SelectItem — 单选下拉
			{
				Key:          "status",
				Label:        "3. Select（状态）",
				ItemType:     shadcn.FilterItemTypeSelect,
				Options:      statusOptions,
				SQLCondition: `status %s ?`,
			},
			// 4. MultipleSelectItem — 多选下拉
			{
				Key:          "multi_status",
				Label:        "4. MultipleSelect（状态多选）",
				ItemType:     shadcn.FilterItemTypeMultipleSelect,
				Options:      statusOptions,
				SQLCondition: `status %s ?`,
			},
			// 5. AutoCompleteItem — 自动完成搜索
			{
				Key:          "category",
				Label:        "5. AutoComplete（分类）",
				ItemType:     shadcn.FilterItemTypeAutoComplete,
				Options:      categoryOptions,
				SQLCondition: `category %s ?`,
			},
			// 5.5 BooleanItem — 布尔值筛选（Switch 开关）
			{
				Key:          "is_active",
				Label:        "5.5 Boolean（是否启用）",
				ItemType:     shadcn.FilterItemTypeBoolean,
				SQLCondition: `is_active %s ?`,
			},
			// 6. DateItem — 日期输入（手动输入格式）
			{
				Key:          "date",
				Label:        "6. Date（创建日期）",
				ItemType:     shadcn.FilterItemTypeDate,
				SQLCondition: `created_at %s ?`,
			},
			// 7. DateRangeItem — 日期范围输入
			{
				Key:          "date_range",
				Label:        "7. DateRange（日期范围）",
				ItemType:     shadcn.FilterItemTypeDateRange,
				SQLCondition: `created_at %s ?`,
			},
			// 8. DatePickerItem — 日期选择器（弹出日历）
			{
				Key:          "date_picker",
				Label:        "8. DatePicker（日期选择器）",
				ItemType:     shadcn.FilterItemTypeDatePicker,
				SQLCondition: `created_at %s ?`,
			},
			// 9. DateRangePickerItem — 日期范围选择器（弹出日历）
			{
				Key:          "date_range_picker",
				Label:        "9. DateRangePicker（日期范围选择器）",
				ItemType:     shadcn.FilterItemTypeDateRangePicker,
				SQLCondition: `created_at %s ?`,
			},
			// 10. DatetimeRangeItem — 日期时间范围（手动输入）
			{
				Key:          "datetime_range",
				Label:        "10. DatetimeRange（时间范围）",
				ItemType:     shadcn.FilterItemTypeDatetimeRange,
				SQLCondition: `created_at %s ?`,
			},
			// 11. DatetimeRangePickerItem — 日期时间范围选择器（弹出日历+时间）
			{
				Key:          "datetime_range_picker",
				Label:        "11. DatetimeRangePicker（时间范围选择器）",
				ItemType:     shadcn.FilterItemTypeDatetimeRangePicker,
				SQLCondition: `created_at %s ?`,
			},
			// 12. LinkageSelectItem — 级联选择（多级联动）
			{
				Key:          "region",
				Label:        "12. LinkageSelect（地区级联）",
				ItemType:     shadcn.FilterItemTypeLinkageSelect,
				LinkageItems:  [][]shadcn.FilterLinkageItem{level1, level2, level3},
				LinkageLabels: []string{"国家", "省/州", "城市"},
				LinkageSelectData: shadcn.FilterLinkageSelectData{
					SQLConditions: []string{
						`country = ?`,
						`province = ?`,
						`city = ?`,
					},
				},
			},
		}
	})

	// 编辑配置
	mb.Editing("Title", "Amount", "Status", "Category", "Country", "Province", "City")
}

// seedFilterDemoData 创建种子数据
func seedFilterDemoData(db *gorm.DB) {
	var count int64
	db.Model(&models.FilterDemo{}).Count(&count)
	if count > 0 {
		return
	}

	demos := []models.FilterDemo{
		{Title: "iPhone 15 Pro", Amount: 8999, Status: "published", Category: "electronics", Country: "china", Province: "guangdong", City: "shenzhen"},
		{Title: "MacBook Air M3", Amount: 9499, Status: "published", Category: "electronics", Country: "china", Province: "shanghai", City: "pudong"},
		{Title: "Nike Air Max", Amount: 899, Status: "draft", Category: "sports", Country: "china", Province: "guangdong", City: "guangzhou"},
		{Title: "有机大米 5kg", Amount: 68, Status: "published", Category: "food", Country: "china", Province: "beijing", City: "haidian"},
		{Title: "纯棉 T 恤", Amount: 129, Status: "archived", Category: "clothing", Country: "usa", Province: "newyork", City: "nyc"},
		{Title: "智能台灯", Amount: 299, Status: "draft", Category: "home", Country: "usa", Province: "california", City: "losangeles"},
		{Title: "运动水壶", Amount: 59, Status: "published", Category: "sports", Country: "china", Province: "beijing", City: "chaoyang"},
		{Title: "羊毛围巾", Amount: 399, Status: "archived", Category: "clothing", Country: "china", Province: "shanghai", City: "xuhui"},
	}
	db.Create(&demos)
}
