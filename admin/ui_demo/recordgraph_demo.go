package ui_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/recordgraph"
	"gorm.io/gorm"
)

// RGUser 记录关系图演示用户模型（has-many RGOrder）
type RGUser struct {
	ID     uint      `gorm:"primarykey"`
	Name   string
	Orders []RGOrder `gorm:"foreignKey:UserID"`
}

// RGOrder 记录关系图演示订单模型（belongs-to RGUser）
// UserID 用 *uint（可空）——可选关联，断开时置 NULL（非空 uint 的必填 FK 无法断开）。
type RGOrder struct {
	ID     uint `gorm:"primarykey"`
	Amount int
	UserID *uint
	User   *RGUser
}

// ConfigRecordGraphDemo 注册记录关系图演示（用户→订单有关联结构体字段，让 schema.Parse 能探出关联）
func ConfigRecordGraphDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&RGUser{}, &RGOrder{}); err != nil {
		panic(err)
	}
	seedRecordGraph(db)

	userMB := b.Model(&RGUser{}).URIName("rg-demo-users").Label("RG 用户")
	userMB.Listing("ID", "Name")
	userMB.Editing("Name")

	orderMB := b.Model(&RGOrder{}).URIName("rg-demo-orders").Label("RG 订单")
	orderMB.Listing("ID", "Amount", "UserID")
	orderMB.Editing("Amount", "UserID")

	rg := recordgraph.New(b, db)
	rg.EnableFor(userMB)
	rg.EnableFor(orderMB)
	rg.Editable(true)
	rg.Install()
}

// seedRecordGraph 幂等插入演示数据（2 用户各含 2-3 订单）
func seedRecordGraph(db *gorm.DB) {
	var n int64
	db.Model(&RGUser{}).Count(&n)
	if n > 0 {
		return
	}
	alice := RGUser{Name: "Alice"}
	bob := RGUser{Name: "Bob"}
	db.Create(&alice)
	db.Create(&bob)
	uptr := func(v uint) *uint { return &v }
	db.Create(&[]RGOrder{
		{Amount: 100, UserID: uptr(alice.ID)},
		{Amount: 200, UserID: uptr(alice.ID)},
		{Amount: 300, UserID: uptr(alice.ID)},
		{Amount: 50, UserID: uptr(bob.ID)},
		{Amount: 150, UserID: uptr(bob.ID)},
	})
}
