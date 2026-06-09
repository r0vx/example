package containers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/htmlgo"
)

// ListContentWithImage 图文列表容器模型
type ListContentWithImage struct {
	ID             uint
	AddTopSpace    bool
	AddBottomSpace bool
	AnchorID       string
	Items          ImageListItems `sql:"type:text;"`
}

// ImageListItem 图文列表项
type ImageListItem struct {
	ImageURL   string
	Link       string
	Heading    string
	Subheading string
	Text       string
}

// TableName 图文列表表名
func (*ListContentWithImage) TableName() string {
	return "container_list_content_with_image"
}

// ImageListItems JSON 序列化的图文列表项
type ImageListItems []*ImageListItem

// Value 实现 driver.Valuer 接口
func (items ImageListItems) Value() (driver.Value, error) {
	return json.Marshal(items)
}

// Scan 实现 sql.Scanner 接口
func (items *ImageListItems) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), items)
	case []byte:
		return json.Unmarshal(v, items)
	default:
		return errors.New("not supported")
	}
}

// RegisterListContentWithImageContainer 注册图文列表容器
func RegisterListContentWithImageContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("ListContentWithImage").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*ListContentWithImage)
			return ListContentWithImageBody(v, input)
		})
	vb.Model(&ListContentWithImage{}).Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "Items")
	vb.ConfigureEditing(func(eb *presets.EditingBuilder) {
		fb := pb.GetPresetsBuilder().NewFieldsBuilder(presets.WRITE).Model(&ImageListItem{}).
			Only("ImageURL", "Link", "Heading", "Subheading", "Text")
		eb.Field("Items").Nested(fb)
	})
}

// ListContentWithImageBody 图文列表渲染
func ListContentWithImageBody(data *ListContentWithImage, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper(
		data.AnchorID, "container-list_content_with_image",
		"", "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(
			ImageListItemsBody(data.Items),
		).Class("container-wrapper"),
	)
	return
}

// ImageListItemsBody 图文列表项渲染
func ImageListItemsBody(items []*ImageListItem) HTMLComponent {
	listItemsDiv := Div().Class("container-list_content_with_image-inner")
	for _, i := range items {
		listItemsDiv.AppendChildren(
			Div(
				If(i.Link != "", A().Class("container-list_content_with_image-link").Href(i.Link)),
				If(i.ImageURL != "",
					Div().Class("container-list_content_with_image-image").
						Style(fmt.Sprintf("background-image: url(%s)", i.ImageURL)),
				),
				If(i.Heading != "" || i.Subheading != "" || i.Text != "",
					Div(
						If(i.Heading != "", H3(i.Heading).Class("container-list_content_with_image-heading")),
						If(i.Subheading != "", Div(Text(i.Subheading)).Class("container-list_content_with_image-subheading h5")),
						If(i.Text != "", P(Text(i.Text)).Class("container-list_content_with_image-text")),
					).Class("container-list_content_with_image-content"),
				),
			).Class("container-list_content_with_image-item"),
		)
	}
	return listItemsDiv
}
