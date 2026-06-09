package wizard_demo

import (
	"strings"
	"time"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/wizard"
	"github.com/r0vx/web"
	"gorm.io/gorm"
)

// wizardSessionStore 演示用：进程级 wizard 会话存储（生产环境推荐用 Redis 实现）
var wizardSessionStore = wizard.NewMemorySessionStore()

// wizardIdempotencyStore 演示用：进程级 wizard 幂等占位（生产环境推荐用 Redis SETNX）
var wizardIdempotencyStore = wizard.NewMemoryIdempotencyStore()

// configWizardDeclarativeDemo 演示 Wizard 声明式 Field + Summary + StepError 模式
//
// 与 configWizardDemo 的差别：
//  1. 用 StepBuilder.Field() 替代 ComponentFunc 手写表单
//  2. 用 StepBuilder.Summary() 让最后一步自动汇总
//  3. 用 FieldError/StepError 进行结构化错误返回（避免 err.Error() 外露）
//  4. 用 wiz.Data(ctx) 一次性读取所有累积字段值，无需逐字段 FormValue
func ConfigWizardDeclarativeDemo(b *presets.Builder, db *gorm.DB) {
	db.AutoMigrate(&models.WizardDemo{})

	mb := b.Model(&models.WizardDemo{}).URIName("wizard-declarative")

	mb.Listing("ID", "Name", "Phone", "Industry", "Status")
	mb.Editing("Name", "Phone", "Industry", "Address", "Status")

	w := wizard.New("CreateDeclarative").
		AddStep(
			wizard.Step("基本信息").
				Description("商户基础资料").
				Icon("building-2").
				Field("Name").Label("商户名称").Required(true).Placeholder("请输入名称").End().
				Field("Industry").Label("行业类型").Placeholder("如：电商").End().
				ValidateFunc(func(ctx *web.EventContext) error {
					if strings.TrimSpace(ctx.R.FormValue("Name")) == "" {
						return wizard.NewFieldError("Name", "商户名称不能为空")
					}
					return nil
				}),
		).
		AddStep(
			wizard.Step("联系方式").
				Description("接收通知的渠道").
				Icon("phone").
				Field("Phone").Label("联系电话").Required(true).Type("tel").End().
				Field("Address").Label("地址").HelpText("可填详细地址或省市区").End().
				ValidateFunc(func(ctx *web.EventContext) error {
					phone := ctx.R.FormValue("Phone")
					if phone == "" {
						return wizard.NewFieldError("Phone", "联系电话不能为空")
					}
					if len(phone) < 7 {
						return wizard.NewFieldError("Phone", "电话号码长度不足")
					}
					return nil
				}),
		).
		AddStep(
			wizard.Step("确认提交").
				Description("核对信息后提交").
				Icon("check-circle").
				Summary(),
		).
		Skippable(true).
		PersistStepInQueryString().
		NextLabel("继续").
		PrevLabel("返回").
		SubmitLabel("确认提交").
		SessionStore(wizardSessionStore).
		SessionTTL(15 * time.Minute).
		DialogContentClass("max-w-xl")

	w.OnSubmit(func(ctx *web.EventContext, r *web.EventResponse) error {
		data := w.Data(ctx)

		// 业务侧二次校验（前端 hidden field 不可信）
		if data["Name"] == "" {
			return wizard.NewStepError(0, "商户名称为空，请回到第一步重新填写")
		}

		demo := &models.WizardDemo{
			Name:     data["Name"],
			Industry: data["Industry"],
			Phone:    data["Phone"],
			Address:  data["Address"],
			Status:   "active",
		}
		if err := db.Create(demo).Error; err != nil {
			// 不暴露 SQL 细节；StepError 跳回首步让用户改名
			return wizard.NewStepError(0, "商户名称已存在，请更换")
		}
		presets.ShowMessage(r, "商户创建成功", "success")
		r.Emit(presets.NotifModelsUpdated(&models.WizardDemo{}), presets.PayloadModelsUpdated{})
		return nil
	})

	// 顶级 Action 挂法：使用 BindToAction 显式绑定
	w.Install(mb, "声明式创建（Field+Summary）")
}
