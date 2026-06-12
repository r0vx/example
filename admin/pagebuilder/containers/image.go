package containers

import (
	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	. "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
)

// ImageContainer 图片容器模型
type ImageContainer struct {
	ID                        uint
	AddTopSpace               bool
	AddBottomSpace            bool
	AnchorID                  string
	ImageURL                  string
	ImageAlt                  string
	BackgroundColor           string
	TransitionBackgroundColor string
}

// TableName 图片容器表名
func (*ImageContainer) TableName() string {
	return "container_images"
}

// RegisterImageContainer 注册图片容器
func RegisterImageContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("Image").Group("Content").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*ImageContainer)
			return ImageContainerBody(v, input)
		})
	vb.Model(&ImageContainer{}).Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "BackgroundColor", "TransitionBackgroundColor", "ImageURL", "ImageAlt")
	vb.ConfigureEditing(func(eb *presets.EditingBuilder) {
		eb.Field("BackgroundColor").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions(BackgroundColors))
		})
		eb.Field("TransitionBackgroundColor").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions(BackgroundColors))
		})
	})
}

// ImageContainerBody 图片容器渲染
func ImageContainerBody(data *ImageContainer, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper(
		data.AnchorID, "container-image",
		data.BackgroundColor, data.TransitionBackgroundColor, "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(
			ImageHtml(data.ImageURL, data.ImageAlt),
			Div().Class("container-image-corner"),
		).Class("container-wrapper"),
	)
	return
}
