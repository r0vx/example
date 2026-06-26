package containers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	. "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
)

// ListContentLite 简约列表内容容器模型
type ListContentLite struct {
	ID              uint
	AddTopSpace     bool
	AddBottomSpace  bool
	AnchorID        string
	Items           ListItemLites `sql:"type:text;"`
	BackgroundColor string
}

// ListItemLite 简约列表项
type ListItemLite struct {
	Heading string
	Text    string
}

// TableName 简约列表内容表名
func (*ListContentLite) TableName() string {
	return "container_list_content_lite"
}

// ListItemLites JSON 序列化的简约列表项
type ListItemLites []*ListItemLite

// Value 实现 driver.Valuer 接口
func (items ListItemLites) Value() (driver.Value, error) {
	return json.Marshal(items)
}

// Scan 实现 sql.Scanner 接口
func (items *ListItemLites) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), items)
	case []byte:
		return json.Unmarshal(v, items)
	default:
		return errors.New("not supported")
	}
}

// RegisterListContentLiteContainer 注册简约列表内容容器
func RegisterListContentLiteContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("ListContentLite").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*ListContentLite)
			return ListContentLiteBody(v, input)
		})
	vb.Model(&ListContentLite{}).Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "Items", "BackgroundColor")
	vb.ConfigureEditing(func(eb *presets.EditingBuilder) {
		eb.Field("BackgroundColor").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			return presets.SelectField(obj, field, ctx).Items(StringsToOptions([]string{White, Grey}))
		})
		fb := pb.GetPresetsBuilder().NewFieldsBuilder(presets.WRITE).Model(&ListItemLite{}).Only("Heading", "Text")
		eb.Field("Items").Nested(fb).SorterField("Heading")
	})
}

// ListContentLiteBody 简约列表内容渲染
func ListContentLiteBody(data *ListContentLite, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper(
		data.AnchorID, "container-list_content_lite",
		data.BackgroundColor, "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(LiteItemsBody(data.Items)).Class("container-wrapper"),
	)
	return
}

// LiteItemsBody 简约列表项渲染
func LiteItemsBody(items []*ListItemLite) HTMLComponent {
	itemsDiv := Div().Class("container-list_content_lite-grid")
	for _, i := range items {
		itemsDiv.AppendChildren(
			Div(
				H3(i.Heading).Class("container-list_content_lite-heading"),
				Div(RawHTML(i.Text)).Class("container-list_content_lite-text"),
			).Class("container-list_content_lite-item"),
		)
	}
	return itemsDiv
}
