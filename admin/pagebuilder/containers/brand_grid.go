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

// BrandGrid 品牌网格容器模型
type BrandGrid struct {
	ID             uint
	AddTopSpace    bool
	AddBottomSpace bool
	AnchorID       string
	Brands         Brands `sql:"type:text;"`
}

// Brand 品牌项
type Brand struct {
	ImageURL string
	Name     string
}

// TableName 品牌网格表名
func (*BrandGrid) TableName() string {
	return "container_brand_grids"
}

// Brands JSON 序列化的品牌列表
type Brands []*Brand

// Value 实现 driver.Valuer 接口
func (b Brands) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan 实现 sql.Scanner 接口
func (b *Brands) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), b)
	case []byte:
		return json.Unmarshal(v, b)
	default:
		return errors.New("not supported")
	}
}

// RegisterBrandGridContainer 注册品牌网格容器
func RegisterBrandGridContainer(pb *pagebuilder.Builder) {
	vb := pb.RegisterContainer("BrandGrid").Group("Content").
		RenderFunc(func(obj any, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*BrandGrid)
			return BrandGridBody(v, input)
		})
	vb.Model(&BrandGrid{}).Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "Brands")
	vb.ConfigureEditing(func(eb *presets.EditingBuilder) {
		fb := pb.GetPresetsBuilder().NewFieldsBuilder(presets.WRITE).Model(&Brand{}).Only("ImageURL", "Name")
		eb.Field("Brands").Nested(fb).SorterField("Name")
	})
}

// BrandGridBody 品牌网格渲染
func BrandGridBody(data *BrandGrid, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper(data.AnchorID, "container-brand_grid",
		"", "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(
			BrandsBody(data.Brands),
		).Class("container-wrapper"),
	)
	return
}

// BrandsBody 品牌列表渲染
func BrandsBody(brands []*Brand) HTMLComponent {
	brandsDiv := Div().Class("container-brand_grid-wrap")
	for _, b := range brands {
		brandsDiv.AppendChildren(
			Div(
				If(b.ImageURL != "", ImageHtml(b.ImageURL, b.Name)),
				If(b.ImageURL == "", Span(b.Name)),
			).Class("container-brand_grid-item"),
		)
	}
	return brandsDiv
}
