package containers

import (
	"fmt"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/sunfmin/reflectutils"
	. "github.com/r0vx/htmlgo"
)

// LinkDisplayOption 链接显示选项
const (
	LinkDisplayOptionDesktop = "desktop"
	LinkDisplayOptionMobile  = "mobile"
	LinkDisplayOptionAll     = "all"
)

// LinkDisplayOptions 链接显示选项列表
var LinkDisplayOptions = []string{LinkDisplayOptionAll, LinkDisplayOptionDesktop, LinkDisplayOptionMobile}

// Heading 标题容器模型
type Heading struct {
	ID                uint
	AddTopSpace       bool
	AddBottomSpace    bool
	AnchorID          string
	Heading           string
	FontColor         string
	BackgroundColor   string
	Link              string
	LinkText          string
	LinkDisplayOption string
	Text              string
}

// TableName 标题容器表名
func (*Heading) TableName() string {
	return "container_headings"
}

// RegisterHeadingContainer 注册标题容器
func RegisterHeadingContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("Heading").Group("Navigation").
		RenderFunc(func(obj interface{}, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*Heading)
			return HeadingBody(v, input)
		})
	vb.Model(&Heading{}).Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "Heading", "FontColor", "BackgroundColor", "Link", "LinkText", "LinkDisplayOption", "Text")
	vb.ConfigureEditing(func(eb *presets.EditingBuilder) {
		eb.Field("FontColor").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions(FontColors))
		})
		eb.Field("BackgroundColor").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions(BackgroundColors))
		})
		eb.Field("LinkDisplayOption").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions(LinkDisplayOptions))
		})
		eb.Field("Text").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return shadcn.Textarea().
				Label(field.Label).
				Attr(presets.ShadcnVFieldError(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)), field.Errors)...).
				Disabled(field.Disabled)
		})
	})
}

// HeadingBody 标题容器渲染
func HeadingBody(data *Heading, input *pagebuilder.RenderInput) (body HTMLComponent) {
	headingBody := Div(
		Div(
			If(data.Heading != "",
				If(data.Link != "",
					A(H2(data.Heading).Class("container-heading-title")).Class("container-heading-title-link").Href(data.Link),
				),
				If(data.Link == "",
					H2(data.Heading).Class("container-heading-title"),
				),
			),
			If(data.Text != "", Div(RawHTML(data.Text)).Class("container-heading-content")),
		).Class("container-heading-wrap"),
		If(data.LinkText != "" && data.Link != "",
			Div(
				LinkTextWithArrow(data.LinkText, data.Link),
			).Class("container-heading-link").Attr("data-display", data.LinkDisplayOption),
		),
	).Class("container-heading-inner")

	body = ContainerWrapper(
		data.AnchorID, "container-heading", data.BackgroundColor, "", data.FontColor,
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(headingBody).Class("container-wrapper"),
	)
	return
}
