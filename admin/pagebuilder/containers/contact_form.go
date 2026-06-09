package containers

import (
	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/web"
	. "github.com/r0vx/htmlgo"
)

// ContactForm 联系表单容器模型
type ContactForm struct {
	ID                 uint
	AddTopSpace        bool
	AddBottomSpace     bool
	AnchorID           string
	Heading            string
	Text               string
	SendButtonText     string
	FormButtonText     string
	MessagePlaceholder string
	NamePlaceholder    string
	EmailPlaceholder   string
	ThankyouMessage    string
	ActionUrl          string
	PrivacyPolicy      string
}

// TableName 联系表单表名
func (*ContactForm) TableName() string {
	return "container_contact_form"
}

// RegisterContactFormContainer 注册联系表单容器
func RegisterContactFormContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("ContactForm").
		RenderFunc(func(obj interface{}, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*ContactForm)
			return ContactFormBody(v, input)
		})
	vb.Model(&ContactForm{}).Editing(
		"AddTopSpace", "AddBottomSpace", "AnchorID",
		"Heading", "Text", "SendButtonText", "FormButtonText",
		"MessagePlaceholder", "NamePlaceholder", "EmailPlaceholder",
		"ThankyouMessage", "ActionUrl", "PrivacyPolicy",
	)
}

// ContactFormBody 联系表单渲染
func ContactFormBody(data *ContactForm, input *pagebuilder.RenderInput) (body HTMLComponent) {
	n := Div(
		Div(
			H2(data.Heading),
		).Class("container-contact_form-title"),
		Div(
			P(Text(data.Text)).Class("p-large"),
			A(Span(data.FormButtonText)).Class("container-contact_form-link button").Href("#"),
		).Class("container-contact_form-brief"),
		Form(
			Div(
				Div(Textarea("").Class("textarea").Name("message").Placeholder(data.MessagePlaceholder).Required(true)).Class("container-contact_form-message"),
				Div(
					Input("").Class("input").Type("text").Name("name").Placeholder(data.NamePlaceholder).Required(true),
					Input("").Class("input").Type("email").Name("email").Placeholder(data.EmailPlaceholder).Required(true),
				).Class("container-contact_form-contact"),
			).Class("container-contact_form-filled"),
			Div(
				Div(
					Label("").Children(
						Input("").Class("container-contact_form-policy-checkbox").Type("checkbox").Name("privacy_policy"),
						Div(RawHTML(data.PrivacyPolicy)).Class("container-contact_form-policy-text"),
					).Class("container-contact_form-policy"),
					Div(
						Div(Text(data.ThankyouMessage)).Class("container-contact_form-response-text"),
					).Class("container-contact_form-response"),
				).Class("container-contact_form-append"),
				Div(
					Button("").Class("button").Children(Span(data.SendButtonText)),
				).Class("container-contact_form-button"),
			).Class("container-contact_form-submit"),
		).Class("container-contact_form-form").Action(data.ActionUrl).Method("POST"),
	).Class("container-contact_form-inner")

	body = ContainerWrapper(
		data.AnchorID, "container-contact_form",
		"", "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(n).Class("container-wrapper"),
	)
	return
}
