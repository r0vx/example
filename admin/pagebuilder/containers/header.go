package containers

import (
	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/htmlgo"
)

// WebHeader 页头容器模型
type WebHeader struct {
	ID    uint
	Color string
}

// TableName 页头表名
func (*WebHeader) TableName() string {
	return "container_headers"
}

// RegisterHeader 注册页头容器
func RegisterHeader(pb *pagebuilder.Builder) {
	header := pb.RegisterContainer("Header").Group("Navigation").
		RenderFunc(func(obj interface{}, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			header := obj.(*WebHeader)
			return HeaderTemplate(header, input)
		})

	header.Model(&WebHeader{}).Editing("Color")
	header.ConfigureEditing(func(eb *presets.EditingBuilder) {
		eb.Field("Color").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions([]string{"black", "white"}))
		})
	})
}

// HeaderTemplate 页头渲染模板
func HeaderTemplate(data *WebHeader, input *pagebuilder.RenderInput) (body HTMLComponent) {
	style := "color: #fff;background: #000;"
	if data.Color == "white" {
		style = "color: #000;background: #fff;"
	}

	body = ContainerWrapper(
		"", "container-header", "", "", "",
		"", false, false, style,
		Div(RawHTML(`
<a href="/" class="container-header-logo">Logo</a>
<ul data-list-unset="true" class="container-header-links">
<li><a href="/about/">About</a></li>
<li><a href="/services/">Services</a></li>
<li><a href="/projects/">Projects</a></li>
<li><a href="/contact/">Contact</a></li>
</ul>
<button class="container-header-menu">
<span class="container-header-menu-icon"></span>
</button>`)).Class("container-wrapper"),
	)
	return
}
