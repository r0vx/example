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

// InNumbers 数字展示容器模型
type InNumbers struct {
	ID             uint
	AddTopSpace    bool
	AddBottomSpace bool
	AnchorID       string
	Heading        string
	Items          InNumbersItems
}

// InNumbersItem 数字项
type InNumbersItem struct {
	Heading string
	Text    string
}

// TableName 数字展示表名
func (*InNumbers) TableName() string {
	return "container_in_numbers"
}

// InNumbersItems JSON 序列化的数字项列表
type InNumbersItems []*InNumbersItem

// Value 实现 driver.Valuer 接口
func (items InNumbersItems) Value() (driver.Value, error) {
	return json.Marshal(items)
}

// Scan 实现 sql.Scanner 接口
func (items *InNumbersItems) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), items)
	case []byte:
		return json.Unmarshal(v, items)
	default:
		return errors.New("not supported")
	}
}

// RegisterInNumbersContainer 注册数字展示容器
func RegisterInNumbersContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("InNumbers").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*InNumbers)
			return InNumbersBody(v, input)
		})
	vb.Model(&InNumbers{}).Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "Heading", "Items")
	vb.ConfigureEditing(func(eb *presets.EditingBuilder) {
		eb.ValidateFunc(func(obj any, ctx *web.EventContext) (err web.ValidationErrors) {
			p := obj.(*InNumbers)
			for i, v := range p.Items {
				if v == nil {
					continue
				}
				if v.Heading == "" {
					err.FieldError(fmt.Sprintf("Items[%v].Heading", i), "Heading 不能为空")
				}
			}
			return
		})
		fb := pb.GetPresetsBuilder().NewFieldsBuilder(presets.WRITE).Model(&InNumbersItem{}).Only("Heading", "Text")
		eb.Field("Items").Nested(fb, &presets.DisplayFieldInSorter{Field: "Heading"})
	})
}

// InNumbersBody 数字展示渲染
func InNumbersBody(data *InNumbers, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper(
		data.AnchorID, "container-in_numbers container-corner",
		"", "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(
			H2(data.Heading).Class("container-in_numbers-heading"),
			InNumbersItemsBody(data.Items),
		).Class("container-wrapper"),
	)
	return
}

// InNumbersItemsBody 数字项列表渲染
func InNumbersItemsBody(items []*InNumbersItem) HTMLComponent {
	itemsDiv := Div().Class("container-in_numbers-grid")
	for _, i := range items {
		itemsDiv.AppendChildren(
			Div(
				Div(
					H2(i.Heading).Class("container-in_numbers-item-title"),
					Div(Text(i.Text)).Class("container-in_numbers-item-description"),
				).Class("container-in_numbers-inner"),
			).Class("container-in_numbers-item"),
		)
	}
	return itemsDiv
}
