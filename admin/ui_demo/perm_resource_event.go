package ui_demo

import (
	"time"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/perm"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

// permEventKey 自定义权限资源键（snake）。PermResource 声明、SnakeOn 校验、两处必须一致。
const permEventKey = "export_report"

// Report 演示模型（仅为让 listing 存在；事件本身与表数据无关）。
type Report struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	CreatedAt time.Time
}

// ConfigPermResourceEventDemo 演示「裸 RegisterEventFunc 事件 + PermResource 进权限树独立勾选」。
//
// 三件套：
//  1. lb.PermResource(key, label)      → 进角色权限树 Custom 类（fc_<key> 闸），可独立勾选
//  2. mb.RegisterEventFunc(id, handler) → 注册裸事件供前端调用（框架不自动校验）
//  3. handler 内 Verifier 自校          → 把树里勾的权限接到事件执行上，fail-closed
func ConfigPermResourceEventDemo(b *presets.Builder, db *gorm.DB) {
	if err := db.AutoMigrate(&Report{}); err != nil {
		panic(err)
	}

	mb := b.Model(&Report{}).URIName("perm-resource-event-demo")
	lb := mb.Listing("ID", "Name", "CreatedAt")

	// ① 声明自定义权限资源 → 自动进权限树（Custom 类，资源串 *:<uri>:*fc_export_report:*，动作 presets:list）
	lb.PermResource(permEventKey, "导出报表")

	// ② 注册裸事件
	mb.RegisterEventFunc("exportReport", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		// ③ 手校 fc_ 闸：与 PermResource 同 key。被拒 → 返回 PermissionDenied（框架弹错误 toast），fail-closed。
		if mb.Info().Verifier().Do(presets.PermList).SnakeOn("fc_"+permEventKey).WithReq(ctx.R).IsAllowed() != nil {
			return r, perm.PermissionDenied
		}

		// 业务逻辑（demo：仅提示成功）
		presets.ShowMessage(&r, "报表已导出", "success")
		return r, nil
	})

	// 工具栏触发按钮（「列设置」左侧）。点击调用裸事件；无权限时事件内被拒、弹错误。
	lb.ToolbarTrailing(func(ctx *web.EventContext) h.HTMLComponent {
		return shadcn.Button(h.Text("导出报表")).
			Variant(shadcn.ButtonVariantOutline).
			Size(shadcn.ButtonSizeSm).
			On("click", web.Plaid().EventFunc("exportReport").Go())
	})
}
