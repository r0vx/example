package ui_demo

import (
	"example/models"
	"fmt"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	shadcn "github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// ID → 中文名 映射表
var (
	brandNames = map[string]string{
		"1": "宝马",
		"2": "奔驰",
	}
	makerNames = map[string]string{
		"1": "华晨宝马",
		"2": "进口宝马",
		"3": "进口奔驰",
		"4": "北京奔驰",
	}
	seriesNames = map[string]string{
		"1": "3系", "2": "5系", "3": "X3",
		"4": "X5", "5": "7系",
		"6": "S级", "7": "GLE", "8": "G级",
		"9": "C级", "10": "E级", "11": "GLC",
	}
)

// configVehicleFilterDemo 配置车型级联筛选演示
// 展示 LinkageSelectItem 在品牌/厂商/车系三级联动场景下的筛选效果
func ConfigVehicleFilterDemo(b *presets.Builder, db *gorm.DB) {
	seedVehicleFilterDemoData(db)

	mb := b.Model(&models.VehicleFilterDemo{}).URIName("vehicle-filter-demos").Label("Vehicle Filter Demos")

	listing := mb.Listing("ID", "Title", "Brand", "Maker", "Series", "Price", "CreatedAt")
	listing.ResponsiveCards(false) // ponytail: 临时开表格视图测移动端 action 列横滚；测完删
	listing.SearchColumns("title")
	listing.PerPage(20)

	// 品牌 — ID 转中文名
	listing.Field("Brand").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		d := obj.(*models.VehicleFilterDemo)
		return h.Text(brandNames[d.Brand])
	})
	// 厂商 — ID 转中文名
	listing.Field("Maker").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		d := obj.(*models.VehicleFilterDemo)
		return h.Text(makerNames[d.Maker])
	})
	// 车系 — ID 转中文名
	listing.Field("Series").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		d := obj.(*models.VehicleFilterDemo)
		return h.Text(seriesNames[d.Series])
	})
	// 价格格式化
	listing.Field("Price").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		d := obj.(*models.VehicleFilterDemo)
		return h.Text(fmt.Sprintf("¥%.2f", d.Price))
	})

	// 品牌（第一级）
	brandLevel := []shadcn.FilterLinkageItem{
		{ID: "1", Name: "宝马", ChildrenIDs: []string{"1", "2"}},
		{ID: "2", Name: "奔驰", ChildrenIDs: []string{"3", "4"}},
	}
	// 厂商（第二级）
	makerLevel := []shadcn.FilterLinkageItem{
		{ID: "1", Name: "华晨宝马", ChildrenIDs: []string{"1", "2", "3"}},
		{ID: "2", Name: "进口宝马", ChildrenIDs: []string{"1", "2", "4", "5"}},
		{ID: "3", Name: "进口奔驰", ChildrenIDs: []string{"6", "7", "8"}},
		{ID: "4", Name: "北京奔驰", ChildrenIDs: []string{"9", "10", "11"}},
	}
	// 车系（第三级）— 同一车系可属于多个厂商
	seriesLevel := []shadcn.FilterLinkageItem{
		{ID: "1", Name: "3系"},
		{ID: "2", Name: "5系"},
		{ID: "3", Name: "X3"},
		{ID: "4", Name: "X5"},
		{ID: "5", Name: "7系"},
		{ID: "6", Name: "S级"},
		{ID: "7", Name: "GLE"},
		{ID: "8", Name: "G级"},
		{ID: "9", Name: "C级"},
		{ID: "10", Name: "E级"},
		{ID: "11", Name: "GLC"},
	}

	listing.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		return []*shadcn.FilterItem{
			{
				Key:           "vehicle",
				Label:         "车型（品牌/厂商/车系）",
				ItemType:      shadcn.FilterItemTypeLinkageSelect,
				LinkageItems:  [][]shadcn.FilterLinkageItem{brandLevel, makerLevel, seriesLevel},
				LinkageLabels: []string{"品牌", "厂商", "车系"},
				LinkageSelectData: shadcn.FilterLinkageSelectData{
					SQLConditions: []string{
						`brand = ?`,
						`maker = ?`,
						`series = ?`,
					},
				},
			},
		}
	})

	mb.Editing("Title", "Brand", "Maker", "Series", "Price")
}

// seedVehicleFilterDemoData 创建车型筛选种子数据
func seedVehicleFilterDemoData(db *gorm.DB) {
	var count int64
	db.Model(&models.VehicleFilterDemo{}).Count(&count)
	if count > 0 {
		return
	}

	demos := []models.VehicleFilterDemo{
		{Title: "前刹车片", Brand: "1", Maker: "1", Series: "1", Price: 580},     // 宝马/华晨宝马/3系
		{Title: "空气滤芯", Brand: "1", Maker: "1", Series: "2", Price: 120},     // 宝马/华晨宝马/5系
		{Title: "后减震器", Brand: "1", Maker: "1", Series: "3", Price: 1200},    // 宝马/华晨宝马/X3
		{Title: "机油滤清器", Brand: "1", Maker: "2", Series: "5", Price: 95},     // 宝马/进口宝马/7系
		{Title: "前大灯总成", Brand: "1", Maker: "2", Series: "4", Price: 3500},   // 宝马/进口宝马/X5
		{Title: "进口3系刹车盘", Brand: "1", Maker: "2", Series: "1", Price: 880},  // 宝马/进口宝马/3系
		{Title: "雨刮片", Brand: "2", Maker: "3", Series: "6", Price: 280},      // 奔驰/进口奔驰/S级
		{Title: "空调滤芯", Brand: "2", Maker: "3", Series: "7", Price: 150},     // 奔驰/进口奔驰/GLE
		{Title: "火花塞（4支装）", Brand: "2", Maker: "4", Series: "9", Price: 320}, // 奔驰/北京奔驰/C级
		{Title: "后视镜片", Brand: "2", Maker: "4", Series: "10", Price: 180},    // 奔驰/北京奔驰/E级
		{Title: "刹车盘", Brand: "2", Maker: "4", Series: "11", Price: 760},     // 奔驰/北京奔驰/GLC
	}
	db.Create(&demos)
}
