package ui_demo

import (
	"fmt"

	"example/models"

	"github.com/r0vx/admin/presets"
	"gorm.io/gorm"
)

// configDisableRowClickDemo 演示 ListingBuilder.DisableRowClick + 列宽 API
//
// URI: /readonly-list
//
// 特点：
//   - 整行点击不触发任何动作（不弹 Drawer / 不跳详情）
//   - 演示 ColumnWidth 一行锁定列宽
//   - 演示 HeaderClass / CellClass 分别设置
//   - 演示 ResizableColumns 拖动表头边缘调整列宽（双击手柄回默认，localStorage 记住）
//   - 演示 ReorderableColumns 拖动 ⋮⋮ 手柄改变列顺序（localStorage 记住）
//   - 演示 StickyHeader 数据超一屏时表头钉在表格区顶部
func ConfigDisableRowClickDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&models.WizardDemo{}); err != nil {
		panic(err)
	}

	// 填充测试数据（数据超一屏方便测试 StickyHeader 等），幂等：不足 50 条才补
	var count int64
	db.Model(&models.WizardDemo{}).Count(&count)
	if count < 50 {
		industries := []string{"液化石油气", "天然气", "工业气体", "医用气体", "特种气体"}
		statuses := []string{"draft", "active", "pending", "closed"}
		rows := make([]models.WizardDemo, 0, 50)
		for i := 1; i <= 50; i++ {
			rows = append(rows, models.WizardDemo{
				Name:     fmt.Sprintf("测试用户%02d", i),
				Industry: industries[i%len(industries)],
				Phone:    fmt.Sprintf("13%09d", i),
				Address:  fmt.Sprintf("测试地址 %d 号", i),
				Status:   statuses[i%len(statuses)],
			})
		}
		db.Create(&rows)
	}

	mb := b.Model(&models.WizardDemo{}).
		URIName("readonly-list").
		Label("只读列表演示")

	lb := mb.Listing("ID", "Name", "Phone", "Industry", "Status", "UpdatedAt").
		DisableRowClick(true).    // 整行点击不响应
		ResizableColumns(true).   // 列宽可拖动（localStorage 持久化）
		ReorderableColumns(true). // 列可拖动排序（localStorage 持久化）
		StickyHeader(true)        // ← 数据超出时表头吸附钉顶

	// 列宽演示（Field 挂在 ListingBuilder 上）
	lb.Field("ID").ColumnWidth("w-16")                                  // 64px
	lb.Field("Name").ColumnWidth("w-48")                                // 192px
	lb.Field("Phone").HeaderClass("text-center").CellClass("font-mono") // 表头居中 + 单元格等宽字体
	lb.Field("Industry").ColumnWidth("w-32")                            // 128px
	lb.Field("Status").ColumnWidth("w-24")                              // 96px
	lb.Field("UpdatedAt").ColumnWidth("w-44")                           // 176px

	// 不挂 Editing/Detailing，纯只读
}
