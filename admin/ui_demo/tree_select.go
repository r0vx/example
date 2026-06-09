package ui_demo

import (
	"encoding/json"
	"example/models"

	"github.com/r0vx/admin/presets"
	shadcn "github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/web"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// VehicleTreeNode 车型树形节点
type VehicleTreeNode struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Children []VehicleTreeNode `json:"children,omitempty"`
}

// vehicleTree 静态车型树数据（品牌 → 厂商 → 车系）
var vehicleTree = []VehicleTreeNode{
	{ID: "b1", Name: "宝马", Children: []VehicleTreeNode{
		{ID: "m1", Name: "华晨宝马", Children: []VehicleTreeNode{
			{ID: "1", Name: "3系"},
			{ID: "2", Name: "5系"},
			{ID: "3", Name: "X3"},
		}},
		{ID: "m2", Name: "进口宝马", Children: []VehicleTreeNode{
			{ID: "4", Name: "X5"},
			{ID: "5", Name: "7系"},
		}},
	}},
	{ID: "b2", Name: "奔驰", Children: []VehicleTreeNode{
		{ID: "m3", Name: "进口奔驰", Children: []VehicleTreeNode{
			{ID: "6", Name: "S级"},
			{ID: "7", Name: "GLE"},
			{ID: "8", Name: "G级"},
		}},
		{ID: "m4", Name: "北京奔驰", Children: []VehicleTreeNode{
			{ID: "9", Name: "C级"},
			{ID: "10", Name: "E级"},
			{ID: "11", Name: "GLC"},
		}},
	}},
}

// configTreeSelectDemo 配置树形选择演示
func ConfigTreeSelectDemo(b *presets.Builder, db *gorm.DB) {
	seedTreeSelectDemoData(db)

	mb := b.Model(&models.TreeSelectDemo{}).URIName("tree-select-demos").Label("Tree Select Demos")

	listing := mb.Listing("ID", "Name", "CreatedAt")
	listing.SearchColumns("name")

	// 编辑表单 — 自定义 SeriesIDs 字段
	editing := mb.Editing("Name", "SeriesIDs")

	// SeriesIDs 字段 — Checkbox TreeView
	editing.Field("SeriesIDs").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		demo := obj.(*models.TreeSelectDemo)

		// 从关联表读取已选的车系 ID
		var seriesIDs []string
		if demo.ID > 0 {
			db.Model(&models.TreeSelectDemoSeries{}).
				Where("tree_select_demo_id = ?", demo.ID).
				Pluck("series_id", &seriesIDs)
		}

		return h.Div(
			h.Label("经营车型").Class("text-sm font-medium"),
			web.Scope().VSlot("{locals}").Init(`{seriesIDs: `+h.JSONString(seriesIDs)+`}`).Children(
				shadcn.TreeView().
					Items(vehicleTree).
					ItemValue("id").
					ItemTitle("name").
					ItemChildren("children").
					Multiple(true).
					Checkbox(true).
					PropagateSelect(true).
					ShowTags(true).
					Attr("v-model", "locals.seriesIDs").
					Class("mt-2 border rounded-md p-2 max-h-80 overflow-y-auto"),
				// hidden input 传值到表单
				h.Input("").Type("hidden").
					Attr(":value", "JSON.stringify(locals.seriesIDs || [])").
					Attr("name", "SeriesIDs"),
			),
		).Class("mb-4")
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		// 从表单取 SeriesIDs JSON 字符串，解析后保存到关联表
		demo := obj.(*models.TreeSelectDemo)
		seriesJSON := ctx.R.FormValue("SeriesIDs")

		// 先保存主记录（确保有 ID）
		if demo.ID == 0 {
			if err = db.Save(demo).Error; err != nil {
				return
			}
		}

		// 删除旧关联
		db.Where("tree_select_demo_id = ?", demo.ID).Delete(&models.TreeSelectDemoSeries{})

		// 解析并插入新关联
		if seriesJSON != "" && seriesJSON != "[]" {
			var ids []string
			if e := json.Unmarshal([]byte(seriesJSON), &ids); e == nil {
				// 只保存叶子节点 ID（纯数字的是车系，b/m 前缀的是品牌/厂商）
				for _, id := range ids {
					if id != "" && id[0] >= '0' && id[0] <= '9' {
						db.Create(&models.TreeSelectDemoSeries{
							TreeSelectDemoID: demo.ID,
							SeriesID:         id,
						})
					}
				}
			}
		}
		return nil
	})
}

// seedTreeSelectDemoData 创建种子数据
func seedTreeSelectDemoData(db *gorm.DB) {
	var count int64
	db.Model(&models.TreeSelectDemo{}).Count(&count)
	if count > 0 {
		return
	}

	demos := []models.TreeSelectDemo{
		{Name: "张三汽配"},
		{Name: "李四配件"},
	}
	db.Create(&demos)

	// 张三汽配经营华晨宝马全系
	for _, sid := range []string{"1", "2", "3"} {
		db.Create(&models.TreeSelectDemoSeries{TreeSelectDemoID: demos[0].ID, SeriesID: sid})
	}
	// 李四配件经营北京奔驰 C级、E级
	for _, sid := range []string{"9", "10"} {
		db.Create(&models.TreeSelectDemoSeries{TreeSelectDemoID: demos[1].ID, SeriesID: sid})
	}
}
