package containers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/htmlgo"
)

// ListContent 列表内容容器模型
type ListContent struct {
	ID                uint
	AddTopSpace       bool
	AddBottomSpace    bool
	AnchorID          string
	Items             ListItems `sql:"type:text;"`
	BackgroundColor   string
	Link              string
	LinkText          string
	LinkDisplayOption string
}

// ListItem 列表项
type ListItem struct {
	Heading  string
	Text     string
	Link     string
	LinkText string
}

// TableName 列表内容表名
func (*ListContent) TableName() string {
	return "container_list_content"
}

// ListItems JSON 序列化的列表项
type ListItems []*ListItem

// Value 实现 driver.Valuer 接口
func (items ListItems) Value() (driver.Value, error) {
	return json.Marshal(items)
}

// Scan 实现 sql.Scanner 接口
func (items *ListItems) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), items)
	case []byte:
		return json.Unmarshal(v, items)
	default:
		return errors.New("not supported")
	}
}

// RegisterListContentContainer 注册列表内容容器
func RegisterListContentContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("ListContent").Group("Content").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*ListContent)
			return ListContentBody(v, input)
		})
	vb.Model(&ListContent{}).Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "BackgroundColor", "Items", "Link", "LinkText", "LinkDisplayOption")
	vb.ConfigureEditing(func(eb *presets.EditingBuilder) {
		eb.Field("BackgroundColor").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions([]string{White, Grey}))
		})
		eb.Field("LinkDisplayOption").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions(LinkDisplayOptions))
		})
		fb := pb.GetPresetsBuilder().NewFieldsBuilder(presets.WRITE).Model(&ListItem{}).Only("Heading", "Text", "Link", "LinkText")
		eb.Field("Items").Nested(fb, &presets.DisplayFieldInSorter{Field: "Heading"})
	})
}

// ListContentBody 列表内容渲染
func ListContentBody(data *ListContent, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper(
		data.AnchorID, "container-list_content",
		data.BackgroundColor, "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(
			ListItemsBody(data.Items),
			If(data.LinkText != "" && data.Link != "",
				Div(
					LinkTextWithArrow(data.LinkText, data.Link),
				).Class("container-list_content-link").Attr("data-display", data.LinkDisplayOption),
			),
		).Class("container-wrapper"),
	)
	return
}

// ListItemsBody 列表项渲染
func ListItemsBody(items []*ListItem) HTMLComponent {
	itemsDiv := Div().Class("container-list_content-grid")
	for _, i := range items {
		itemsDiv.AppendChildren(
			Div(
				Div(
					If(i.Link != "",
						A(H3(i.Heading)).Class("container-list_content-heading").Href(i.Link),
					),
					If(i.Link == "",
						Div(H3(i.Heading)).Class("container-list_content-heading"),
					),
					Div(
						P(Text(i.Text)),
						LinkTextWithArrow(i.LinkText, i.Link),
					).Class("container-list_content-content"),
				).Class("container-list_content-inner"),
			).Class("container-list_content-item"),
		)
	}
	return itemsDiv
}
