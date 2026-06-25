package models

import "time"

// SCMember 非 ID 主键演示模型：主键是 ServiceUserID（uint），结构体**没有** ID 字段。
//
// 用于验证框架对非标准主键的支持：
//   - schema 自动解析主键（无需手写 PrimaryField）
//   - ObjectID 回退 schema 主键取值（行 id / 选择框 v-model / 编辑链接不再为空）
//   - 列表分页排序用真实主键（不再因找不到 "ID" 报错）
type SCMember struct {
	ServiceUserID uint `gorm:"primaryKey"` // 主键（非 ID）
	Name          string
	Phone         string
	Balance       float64
	CreatedAt     time.Time
}
