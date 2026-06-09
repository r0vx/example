package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnCascaderDemo 虚拟模型
type ShadcnCascaderDemo struct{}

// configCascader 注册 Cascader demo
func configCascader(b *presets.Builder) {
	m := b.Model(&ShadcnCascaderDemo{}).
		Label("Cascader").
		URIName("shadcn-cascader")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnCascaderDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Cascader"
		r.Body = shadcnCascaderBody(ctx)
		return
	})
}

// shadcnCascaderBody 级联选择器演示页面
func shadcnCascaderBody(ctx *web.EventContext) h.HTMLComponent {
	// 省市区数据
	provinces := []CascaderItem{
		{ID: "zj", Name: "浙江省", ChildrenIDs: []string{"hz", "nb", "wz"}},
		{ID: "js", Name: "江苏省", ChildrenIDs: []string{"nj", "sz", "wx"}},
		{ID: "gd", Name: "广东省", ChildrenIDs: []string{"gz", "sz2", "dg"}},
	}

	cities := []CascaderItem{
		{ID: "hz", Name: "杭州市", ChildrenIDs: []string{"xh", "gs", "yh"}},
		{ID: "nb", Name: "宁波市", ChildrenIDs: []string{"hx", "zx", "yc"}},
		{ID: "wz", Name: "温州市", ChildrenIDs: []string{"lc", "oa", "yj"}},
		{ID: "nj", Name: "南京市", ChildrenIDs: []string{"xw", "jy", "qh"}},
		{ID: "sz", Name: "苏州市", ChildrenIDs: []string{"gc", "wz2", "kc"}},
		{ID: "wx", Name: "无锡市", ChildrenIDs: []string{"cy", "hb", "xs"}},
		{ID: "gz", Name: "广州市", ChildrenIDs: []string{"th", "yy", "hp"}},
		{ID: "sz2", Name: "深圳市", ChildrenIDs: []string{"ft", "ns", "bj"}},
		{ID: "dg", Name: "东莞市", ChildrenIDs: []string{"np", "cz", "hw"}},
	}

	districts := []CascaderItem{
		// 杭州
		{ID: "xh", Name: "西湖区"},
		{ID: "gs", Name: "拱墅区"},
		{ID: "yh", Name: "余杭区"},
		// 宁波
		{ID: "hx", Name: "海曙区"},
		{ID: "zx", Name: "镇海区"},
		{ID: "yc", Name: "鄞州区"},
		// 温州
		{ID: "lc", Name: "鹿城区"},
		{ID: "oa", Name: "瓯海区"},
		{ID: "yj", Name: "永嘉县"},
		// 南京
		{ID: "xw", Name: "玄武区"},
		{ID: "jy", Name: "江宁区"},
		{ID: "qh", Name: "秦淮区"},
		// 苏州
		{ID: "gc", Name: "姑苏区"},
		{ID: "wz2", Name: "吴中区"},
		{ID: "kc", Name: "昆山市"},
		// 无锡
		{ID: "cy", Name: "崇安区"},
		{ID: "hb", Name: "惠山区"},
		{ID: "xs", Name: "锡山区"},
		// 广州
		{ID: "th", Name: "天河区"},
		{ID: "yy", Name: "越秀区"},
		{ID: "hp", Name: "黄埔区"},
		// 深圳
		{ID: "ft", Name: "福田区"},
		{ID: "ns", Name: "南山区"},
		{ID: "bj", Name: "宝安区"},
		// 东莞
		{ID: "np", Name: "南城区"},
		{ID: "cz", Name: "长安镇"},
		{ID: "hw", Name: "虎门镇"},
	}

	return h.Div(
		h.H1("Cascader 级联选择器").Style("margin-bottom: 24px;"),

		// 基本用法
		h.Div(
			h.H2("基本用法"),
			h.P(h.Text("用于多级联动选择，如省市区选择")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					Cascader().
						Items(provinces, cities, districts).
						Labels("选择省份", "选择城市", "选择区县").
						Attr("v-model", "form.location").
						Class("max-w-md"),
				),
				h.Div(
					h.Span("选中: {{ form.location?.join(' / ') || '未选择' }}").Class("text-sm text-muted-foreground"),
				).Class("mt-2"),
			).VSlot("{ form }").FormInit(`{ "location": [] }`),
		).Class("demo-section"),

		// 水平排列
		h.Div(
			h.H2("水平排列"),
			h.P(h.Text("使用 Row 属性使选择器水平排列")).Class("text-muted-foreground mb-4"),
			web.Scope(
				Cascader().
					Items(provinces, cities, districts).
					Labels("省份", "城市", "区县").
					Row(true).
					Attr("v-model", "form.location"),
				h.Div(
					h.Span("选中: {{ form.location?.join(' / ') || '未选择' }}").Class("text-sm text-muted-foreground"),
				).Class("mt-2"),
			).VSlot("{ form }").FormInit(`{ "location": [] }`),
		).Class("demo-section"),

		// 带错误信息
		h.Div(
			h.H2("表单验证"),
			h.P(h.Text("显示各级的错误信息")).Class("text-muted-foreground mb-4"),
			Cascader().
				Items(provinces, cities, districts).
				Labels("选择省份", "选择城市", "选择区县").
				ErrorMessages("请选择省份", "", "").
				Class("max-w-md"),
		).Class("demo-section"),

		// 禁用状态
		h.Div(
			h.H2("禁用状态"),
			h.P(h.Text("禁用整个级联选择器")).Class("text-muted-foreground mb-4"),
			Cascader().
				Items(provinces, cities, districts).
				Labels("选择省份", "选择城市", "选择区县").
				ModelValue([]string{"zj", "hz", "xh"}).
				Disabled(true).
				Class("max-w-md"),
		).Class("demo-section"),

		// 允许跨级选择
		h.Div(
			h.H2("跨级选择"),
			h.P(h.Text("允许不按顺序选择")).Class("text-muted-foreground mb-4"),
			web.Scope(
				Cascader().
					Items(provinces, cities, districts).
					Labels("省份", "城市", "区县").
					SelectOutOfOrder(true).
					Row(true).
					Attr("v-model", "form.location"),
				h.Div(
					h.Span("选中: {{ form.location?.join(' / ') || '未选择' }}").Class("text-sm text-muted-foreground"),
				).Class("mt-2"),
			).VSlot("{ form }").FormInit(`{ "location": [] }`),
		).Class("demo-section"),
	).Style("max-width: 800px; margin: 0 auto;")
}
