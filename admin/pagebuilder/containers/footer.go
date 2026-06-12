package containers

import (
	"fmt"

	"github.com/r0vx/admin/pagebuilder"
	. "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
)

// WebFooter 页脚容器模型
type WebFooter struct {
	ID          uint
	EnglishUrl  string
	JapaneseUrl string
}

// TableName 页脚表名
func (*WebFooter) TableName() string {
	return "container_footers"
}

// RegisterFooter 注册页脚容器
func RegisterFooter(pb *pagebuilder.Builder) {
	footer := pb.RegisterContainer("Footer").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			footer := obj.(*WebFooter)
			return FooterTemplate(footer, input)
		})

	footer.Model(&WebFooter{}).Editing("EnglishUrl", "JapaneseUrl")
}

// FooterTemplate 页脚渲染模板
func FooterTemplate(data *WebFooter, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper("", "container-footer", "", "", "",
		"", false, false, "",
		Div(RawHTML(fmt.Sprintf(`
<div class='container-footer-main'>
<div class='container-footer-primary'>
<div class='container-footer-links'>
<div class='container-footer-links-group'>
<div class='container-footer-links-title'><a href='/about/'>About</a></div>
<ul data-list-unset='true' class='container-footer-links-list'>
<li class='container-footer-links-item'><a href='/about/team/'>Team</a></li>
<li class='container-footer-links-item'><a href='/about/culture/'>Culture</a></li>
</ul>
</div>
<div class='container-footer-links-group'>
<div class='container-footer-links-title'><a href='/services/'>Services</a></div>
<ul data-list-unset='true' class='container-footer-links-list'>
<li class='container-footer-links-item'><a href='/services/development/'>Development</a></li>
<li class='container-footer-links-item'><a href='/services/design/'>Design</a></li>
</ul>
</div>
</div>
</div>
<div class='container-footer-secondary'>
<ul data-list-unset='true' class='container-footer-language'>
<li><a href='%s'>English</a></li>
<li><a href='%s'>日本語</a></li>
</ul>
</div>
</div>`, data.EnglishUrl, data.JapaneseUrl))).Class("container-wrapper"))
	return
}
