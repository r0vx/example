package wizard_demo

import (
	"fmt"
	"strings"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/wizard"
	"github.com/r0vx/web"
	"gorm.io/gorm"
)

// configWizardFullPageDemo 演示 Wizard FullPage（独立路由）模式
//
// 与 configWizardDeclarativeDemo（Dialog 模式）对照：
//  1. 入口：URL `/wizard-fullpage` 直接进入，不依赖按钮
//  2. 整页布局：使用 presets.NewCustomPage layout
//  3. 浏览器刷新可恢复（需配 SessionStore + PersistStepInQueryString）
//  4. 提交成功跳转到 SubmitRedirect，取消跳转到 CancelRedirect
func ConfigWizardFullPageDemo(b *presets.Builder, db *gorm.DB) {
	db.AutoMigrate(&models.WizardDemo{})

	w := wizard.New("CreateTenantFullPage").
		AddStep(
			wizard.Step("基本信息").
				Description("租户基础资料").
				Icon("building-2").
				Field("Name").Label("租户名称").Required(true).End().
				Field("Industry").Label("行业").End().
				ValidateFunc(func(ctx *web.EventContext) error {
					if strings.TrimSpace(ctx.R.FormValue("Name")) == "" {
						return wizard.NewFieldError("Name", "租户名称不能为空")
					}
					return nil
				}),
		).
		AddStep(
			wizard.Step("联系方式").
				Description("接收通知").
				Icon("phone").
				Field("Phone").Label("联系电话").Required(true).End().
				Field("Address").Label("地址").End(),
		).
		AddStep(
			wizard.Step("确认").
				Description("核对后提交").
				Icon("check-circle").
				Summary(),
		).
		Skippable(true).
		PersistStepInQueryString().
		SessionStore(wizardSessionStore).
		IdempotencyStore(wizardIdempotencyStore).
		NextLabel("继续").
		SubmitLabel("确认创建").
		CancelLabel("放弃").
		SubmitRedirect("/wizard-declarative"). // 提交成功回 dialog demo 列表
		CancelRedirect("/wizard-declarative")  // 取消同上

	w.OnSubmit(func(ctx *web.EventContext, r *web.EventResponse) error {
		data := w.Data(ctx)
		// FullPage demo：仅做最小演示，不真落库；toast 提示成功
		web.AppendRunScripts(r,
			fmt.Sprintf(`window.plaid.toast.success("租户已创建：%s")`, data["Name"]),
		)
		return nil
	})

	w.InstallFullPage(b, "/wizard-fullpage", "新建租户（FullPage 模式）")
}
