package wizard_demo

import (
	"fmt"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/wizard"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// configWizardDemo 配置 Wizard 多步向导演示
func ConfigWizardDemo(b *presets.Builder, db *gorm.DB) {
	db.AutoMigrate(&models.WizardDemo{})

	mb := b.Model(&models.WizardDemo{}).URIName("wizard-demos")

	// Listing
	mb.Listing("ID", "Name", "Phone", "Industry", "Status")

	// Editing
	mb.Editing("Name", "Phone", "Industry", "Address", "Status")

	// Wizard Action：多步创建商户
	w := wizard.New("CreateWizardDemo").
		AddStep(
			wizard.Step("基本信息").
				ComponentFunc(func(ctx *web.EventContext) h.HTMLComponent {
					return h.Div(
						h.Div(
							h.Label("商户名称").Class("text-sm font-medium"),
							shadcn.Input().
								Placeholder("请输入商户名称").
								Attr(web.VField("Name", ctx.R.FormValue("Name"))...),
						).Class("space-y-2"),
						h.Div(
							h.Label("行业类型").Class("text-sm font-medium"),
							shadcn.Select().
								Items([]shadcn.DefaultOptionItem{
									{Text: "餐饮", Value: "餐饮"},
									{Text: "零售", Value: "零售"},
									{Text: "服务", Value: "服务"},
									{Text: "制造", Value: "制造"},
								}).
								Placeholder("请选择行业").
								Attr(web.VField("Industry", ctx.R.FormValue("Industry"))...),
						).Class("space-y-2"),
					).Class("space-y-4")
				}).
				ValidateFunc(func(ctx *web.EventContext) error {
					if ctx.R.FormValue("Name") == "" {
						return fmt.Errorf("商户名称不能为空")
					}
					if ctx.R.FormValue("Industry") == "" {
						return fmt.Errorf("请选择行业类型")
					}
					return nil
				}),
		).
		AddStep(
			wizard.Step("联系方式").
				ComponentFunc(func(ctx *web.EventContext) h.HTMLComponent {
					return h.Div(
						h.Div(
							h.Label("联系电话").Class("text-sm font-medium"),
							shadcn.Input().
								Placeholder("请输入联系电话").
								Attr(web.VField("Phone", ctx.R.FormValue("Phone"))...),
						).Class("space-y-2"),
						h.Div(
							h.Label("地址").Class("text-sm font-medium"),
							shadcn.Input().
								Placeholder("请输入地址").
								Attr(web.VField("Address", ctx.R.FormValue("Address"))...),
						).Class("space-y-2"),
					).Class("space-y-4")
				}).
				ValidateFunc(func(ctx *web.EventContext) error {
					if ctx.R.FormValue("Phone") == "" {
						return fmt.Errorf("联系电话不能为空")
					}
					return nil
				}),
		).
		AddStep(
			wizard.Step("确认提交").
				ComponentFunc(func(ctx *web.EventContext) h.HTMLComponent {
					return h.Div(
						h.Div(h.Text("请确认以下信息：")).Class("text-sm font-medium mb-3"),
						summaryItem("商户名称", ctx.R.FormValue("Name")),
						summaryItem("行业类型", ctx.R.FormValue("Industry")),
						summaryItem("联系电话", ctx.R.FormValue("Phone")),
						summaryItem("地址", ctx.R.FormValue("Address")),
					).Class("space-y-2")
				}),
		).
		OnSubmit(func(ctx *web.EventContext, r *web.EventResponse) error {
			demo := &models.WizardDemo{
				Name:     ctx.R.FormValue("Name"),
				Industry: ctx.R.FormValue("Industry"),
				Phone:    ctx.R.FormValue("Phone"),
				Address:  ctx.R.FormValue("Address"),
				Status:   "active",
			}
			if err := db.Create(demo).Error; err != nil {
				return fmt.Errorf("创建失败: %w", err)
			}
			presets.ShowMessage(r, "商户创建成功", "success")
			r.Emit(presets.NotifModelsUpdated(&models.WizardDemo{}), presets.PayloadModelsUpdated{})
			return nil
		}).
		DialogContentClass("max-w-lg")

	w.Install(mb, "新建商户向导")
}

// summaryItem 渲染确认页的单行信息
func summaryItem(label, value string) h.HTMLComponent {
	if value == "" {
		value = "—"
	}
	return h.Div(
		h.Span(label+"：").Class("text-sm text-muted-foreground"),
		h.Span(value).Class("text-sm font-medium"),
	).Class("flex gap-2")
}
