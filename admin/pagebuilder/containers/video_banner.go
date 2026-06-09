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

// VideoBanner 视频横幅容器模型
type VideoBanner struct {
	ID             uint
	AddTopSpace    bool
	AddBottomSpace bool
	AnchorID       string
	VideoURL       string
	VideoCoverURL  string
	Heading        string
	Text           string
	LinkText       string
	Link           string
}

// TableName 视频横幅表名
func (*VideoBanner) TableName() string {
	return "container_video_banners"
}

// RegisterVideoBannerContainer 注册视频横幅容器
func RegisterVideoBannerContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("VideoBanner").Group("Content").
		RenderFunc(func(obj interface{}, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*VideoBanner)
			return VideoBannerBody(v, input)
		})
	vb.Model(&VideoBanner{}).Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "VideoURL", "VideoCoverURL", "Heading", "Text", "LinkText", "Link")
	vb.ConfigureEditing(func(eb *presets.EditingBuilder) {
		eb.Field("Heading").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return shadcn.Textarea().
				Label(field.Label).
				Attr(presets.ShadcnVFieldError(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)), field.Errors)...).
				Disabled(field.Disabled)
		})
		eb.Field("Text").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return shadcn.Textarea().
				Label(field.Label).
				Attr(presets.ShadcnVFieldError(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)), field.Errors)...).
				Disabled(field.Disabled)
		})
	})
}

// VideoBannerBody 视频横幅渲染
func VideoBannerBody(data *VideoBanner, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper(
		data.AnchorID, "container-video_banner",
		"", "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div().Class("container-video_banner-mask"),
		VideoBannerHeadBody(data),
		VideoBannerFootBody(data),
	)
	return
}

// VideoBannerHeadBody 视频横幅头部
func VideoBannerHeadBody(data *VideoBanner) HTMLComponent {
	return Div(
		Div().Class("container-video_banner-background container-video_banner-background-image"),
		If(data.VideoURL != "",
			Video(
				Source("").Src(data.VideoURL),
			).Class("container-video_banner-background container-video_banner-background-desktop").
				Attr("preload", "none").Attr("loop", "true").Attr("muted", "true").
				Attr("playsinline", "true").Attr("data-cover-image-url", data.VideoCoverURL),
		),
		Div(
			If(data.Heading != "", H1(data.Heading).Class("container-video_banner-heading")),
		).Class("container-video_banner-head-wrap container-wrapper"),
	).Class("container-video_banner-head")
}

// VideoBannerFootBody 视频横幅底部
func VideoBannerFootBody(data *VideoBanner) HTMLComponent {
	return Div(
		Div(
			P(Text(data.Text)).Class("container-video_banner-text p-large"),
			LinkTextWithArrow(data.LinkText, data.Link),
		).Class("container-wrapper"),
	).Class("container-video_banner-foot")
}
