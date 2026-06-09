package containers

import (
	"fmt"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/web"
	. "github.com/r0vx/htmlgo"
)

// PageTitle 页面标题容器模型
type PageTitle struct {
	ID                 uint
	AddTopSpace        bool
	AddBottomSpace     bool
	AnchorID           string
	HeroImageURL       string
	NavigationLink     string
	NavigationLinkText string
	Heading            string
	Text               string
}

// TableName 页面标题表名
func (*PageTitle) TableName() string {
	return "container_page_title"
}

// RegisterPageTitleContainer 注册页面标题容器
func RegisterPageTitleContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("PageTitle").Group("Navigation").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*PageTitle)
			return PageTitleBody(v, input)
		})
	vb.Model(&PageTitle{}).Editing(
		"AddTopSpace", "AddBottomSpace", "AnchorID",
		"HeroImageURL", "NavigationLink", "NavigationLinkText",
		"Heading", "Text",
	)
}

// PageTitleBody 页面标题渲染
func PageTitleBody(data *PageTitle, input *pagebuilder.RenderInput) (body HTMLComponent) {
	image := Div().Class("container-page_title-background").
		Style(fmt.Sprintf("background-image: url(%s)", data.HeroImageURL))

	wrapper := Div(
		Div().Class("container-page_title-corner"),
		Div(
			Div(
				Div(
					If(data.NavigationLinkText != "", A(
						RawHTML(`<svg height=".72em" viewBox="0 0 12 15" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M10 2L3 7.5L10 13" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/></svg>`),
						Span(data.NavigationLinkText),
					).Class("container-page_title-navigation").AttrIf("href", data.NavigationLink, data.NavigationLink != "")),
					Div(
						H1(data.Heading),
					).Class("container-page_title-title"),
					If(data.Text != "", P(Text(data.Text)).Class("container-page_title-content p-large")),
				).Class("container-page_title-heading"),
			).Class("container-page_title-inner").
				AttrIf("data-has-navigation", "true", data.NavigationLinkText != ""),
		).Class("container-wrapper"),
	).Class("container-page_title-wrap")

	body = ContainerWrapper(
		data.AnchorID, "container-page_title",
		"", "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		image, wrapper,
	)
	return
}
