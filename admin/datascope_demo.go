package admin

import (
	"slices"
	"time"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// ── 数据隔离 Demo：同一角色在不同表用不同隔离字段 ───────────────────────────────
// 场景：非 admin 用户登录后，三张表各按「自己的字段」把数据隔离到属于当前用户的行；
// admin 看全部。
//
// 关键点：隔离字段名是 per-model 的固定知识（表A 永远 agent_id、表B 永远 user_id…），
// 全局 DataScopeDynamicDefault 拿不到「当前是哪个模型」，硬返回某个 field 时碰到没有该
// 字段的表会退化成「看全部」(数据泄露)。所以正确做法是每张表各自声明
// DataScope("字段").Resolver(...)——字段每模型固定声明，属主值由各自 resolver 给。

// ScopeAgentDeal 表A：按 AgentID（经纪人）隔离
type ScopeAgentDeal struct {
	ID        uint `gorm:"primarykey"`
	Title     string
	AgentID   uint // 隔离列 agent_id
	CreatedAt time.Time
}

// ScopeUserNote 表B：按 UserID（用户）隔离
type ScopeUserNote struct {
	ID        uint `gorm:"primarykey"`
	Content   string
	UserID    uint // 隔离列 user_id
	CreatedAt time.Time
}

// ScopeOrgDoc 表C：按 ParentID（归属节点）隔离
type ScopeOrgDoc struct {
	ID        uint `gorm:"primarykey"`
	Title     string
	ParentID  uint // 隔离列 parent_id
	CreatedAt time.Time
}

// scopeOwnerOf 生成 demo 的属主解析器：admin → bypass 看全部、未登录 → 隔离到空，
// 其余角色按 getVal(u) 取「该表那一维」的属主值。getVal 让三张字段不同的表各取各的值
// （真实场景如 u.AgentID / u.ID / u.OrgID），而 admin/未登录的公共逻辑只写一遍。
func scopeOwnerOf(getVal func(u *models.User) any) presets.DataScopeResolverFunc {
	return func(ctx *web.EventContext) (any, bool) {
		u := getCurrentUser(ctx.R)
		if u == nil {
			return nil, false // 未登录 → 隔离到空结果
		}
		if slices.Contains(u.GetRoles(), models.RoleAdmin) {
			return nil, true // admin → bypass 看全部
		}
		return getVal(u), false // uint，匹配各表的 uint 隔离列（类型须与列兼容）
	}
}

// ConfigDataScopeDemo 注册「同角色跨表不同隔离字段」演示：三张表各用 DataScope("字段").Resolver(...)。
func ConfigDataScopeDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&ScopeAgentDeal{}, &ScopeUserNote{}, &ScopeOrgDoc{}); err != nil {
		panic(err)
	}
	seedScopeDemo(db)

	// 表A：按 agent_id 隔离
	mbA := b.Model(&ScopeAgentDeal{}).URIName("scope-agent-deal")
	lbA := mbA.Listing("ID", "Title", "AgentID")
	mbA.Editing("Title", "AgentID")
	mbA.DataScope("AgentID").Resolver(scopeOwnerOf(func(u *models.User) any { return u.GetID() }))
	// 行菜单：弹中央 Dialog 只编辑 Title 字段（复用 mb.Editing 配置 + partial-safe，AgentID 不受影响）。须在 Editing 之后调。
	// DialogSizeSm：单字段用小尺寸（不设则继承标准 Dialog 半屏默认，单字段偏宽）。
	lbA.RowMenu().RowMenuItem("快速改标题").Icon("pencil").
		EditInDialog("Title").DialogContentClass(presets.DialogSizeSm)
	// 筛选项（用于测 SSE 推送时不冲掉正在编辑的筛选）：AgentID 多选（可勾多项再应用，最契合测试）+ CreatedAt 范围。
	// 整表 reload 路径（无 RowLevelRefresh），正是 guard 保护的场景。
	// 新写法：构造器 + presets.Filters 适配；SQLCondition 自动推导（agent_id/created_at 即列名）。
	lbA.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		return presets.Filters(
			shadcn.MultipleSelectFilterItem("agent_id", "Agent", []shadcn.FilterSelectOption{
				{Value: "1", Text: "Agent 1"},
				{Value: "2", Text: "Agent 2"},
				{Value: "3", Text: "Agent 3"},
			}),
			shadcn.DatetimeRangeFilterItem("created_at", "Created At"),
		)
	})

	// 表B：按 user_id 隔离
	mbB := b.Model(&ScopeUserNote{}).URIName("scope-user-note")
	mbB.Listing("ID", "Content", "UserID")
	mbB.Editing("Content", "UserID")
	mbB.DataScope("UserID").Resolver(scopeOwnerOf(func(u *models.User) any { return u.GetID() }))

	// 表C：按 parent_id 隔离
	mbC := b.Model(&ScopeOrgDoc{}).URIName("scope-org-doc")
	mbC.Listing("ID", "Title", "ParentID")
	mbC.Editing("Title", "ParentID")
	mbC.DataScope("ParentID").Resolver(scopeOwnerOf(func(u *models.User) any { return u.GetID() }))
}

// seedScopeDemo 首次启动插入演示数据（owner=1/2/3 各几条，便于切换登录用户看隔离效果）。
func seedScopeDemo(db *gorm.DB) {
	var count int64
	db.Model(&ScopeAgentDeal{}).Count(&count)
	if count > 0 {
		return
	}
	db.Create(&[]ScopeAgentDeal{
		{Title: "经纪人1的单A", AgentID: 1}, {Title: "经纪人1的单B", AgentID: 1},
		{Title: "经纪人2的单", AgentID: 2}, {Title: "经纪人3的单", AgentID: 3},
	})
	db.Create(&[]ScopeUserNote{
		{Content: "用户1的笔记", UserID: 1}, {Content: "用户2的笔记A", UserID: 2},
		{Content: "用户2的笔记B", UserID: 2}, {Content: "用户3的笔记", UserID: 3},
	})
	db.Create(&[]ScopeOrgDoc{
		{Title: "归属1的文档A", ParentID: 1}, {Title: "归属1的文档B", ParentID: 1},
		{Title: "归属2的文档", ParentID: 2}, {Title: "归属3的文档", ParentID: 3},
	})
}
