package shadcn_demo

import (
	"encoding/json"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	h "github.com/r0vx/htmlgo"
)

// Configure 注册所有 shadcn 组件演示到 admin 面板
func Configure(b *presets.Builder) {
	configBasicInputs(b)
	configSelections(b)
	configDialog(b)
	configTable(b)
	configDataTable(b)
	configLazyPortals(b)
	configGrid(b)
	configList(b)
	configPopoverMenu(b)
	configSheetDrawer(b)
	configVariantSubForm(b)
	configProgress(b)
	configSidebarDemo(b)
	configInvoiceList(b)
	configNewComponents(b)
	configRangePicker(b)
	configFormField(b)
	configAutocomplete(b)
	configDisplayComponents(b)
	configFilter(b)
	configTreeView(b)
	configCascader(b)
	configChart(b)
	configTimeline(b)
	configAdminDemo(b)
}

// emptySearchFunc 返回空数据的通用 SearchFunc，避免数据库查询
func emptySearchFunc[T any]() presets.SearchFunc {
	return func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
		totalCount := 0
		result = &presets.SearchResult{
			Nodes:      []T{},
			TotalCount: &totalCount,
		}
		return
	}
}

// injectDemoCSS 注入 demo 页面通用样式
func injectDemoCSS(ctx *web.EventContext) {
	ctx.Injector.HeadHTML(`<style>
		.demo-section { padding:16px; border:1px solid #e5e7eb; border-radius:8px; margin-bottom:24px; background:white; }
		.demo-section h2 { margin:0 0 16px 0; font-size:18px; }
		.demo-row { display:flex; gap:8px; flex-wrap:wrap; align-items:center; margin-bottom:12px; }
		.demo-row:last-child { margin-bottom:0; }
	</style>`)
}

// prettyFormAsJSON 显示 form 提交数据的 JSON（替代 examples.PrettyFormAsJSON）
func prettyFormAsJSON(ctx *web.EventContext) h.HTMLComponent {
	if ctx.R.MultipartForm == nil {
		return nil
	}
	formData, err := json.MarshalIndent(ctx.R.MultipartForm, "", "\t")
	if err != nil {
		panic(err)
	}
	return h.Pre(string(formData))
}
