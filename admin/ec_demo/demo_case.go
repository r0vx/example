package ec_demo

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

type (
	// DemoCase 演示用例模型
	DemoCase struct {
		gorm.Model
		Name              string
		FieldData         FieldData         `gorm:"type:json"`
		FieldTextareaData FieldTextareaData `gorm:"type:json"`
		FieldPasswordData FieldPasswordData `gorm:"type:json"`
		FieldNumberData   FieldNumberData   `gorm:"type:json"`
		SelectData        SelectData        `gorm:"type:json"`
		CheckboxData      CheckboxData      `gorm:"type:json"`
		DatepickerData    DatepickerData    `gorm:"type:json"`
		PaginatorData     PaginationData    `gorm:"type:json"`
		TabsData          TabsData          `gorm:"type:json"`
	}
	FieldData struct {
		Text         string
		TextValidate string
	}
	FieldTextareaData struct {
		Textarea         string
		TextareaValidate string
	}
	FieldPasswordData struct {
		Password        string
		PasswordDefault string
	}
	FieldNumberData struct {
		Number         int
		NumberValidate int
	}
	SelectData struct {
		AutoComplete []int
		NormalSelect int
	}
	CheckboxData struct {
		Checkbox bool
	}

	PaginationData struct {
		Current int
	}

	TabsData struct {
		Tab []string
	}

	DemoSelectItem struct {
		ID   int
		Name string
	}
	DatepickerData struct {
		Date                 int64
		DateTime             int64
		DateRange            []int64
		DateRangeNeedConfirm []int64
	}
)

func (c *FieldData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *FieldData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

func (c *FieldTextareaData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *FieldTextareaData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

func (c *FieldPasswordData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

func (c *FieldPasswordData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *FieldNumberData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *FieldNumberData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

func (c *SelectData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *SelectData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

func (c *CheckboxData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *CheckboxData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

func (c *DatepickerData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *DatepickerData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

func (c *PaginationData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *PaginationData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

func (c *TabsData) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		return sonic.Unmarshal(bytes, c)
	}
	return nil
}

func (c *TabsData) Value() (driver.Value, error) {
	return sonic.Marshal(c)
}

// configureDemoCase 配置演示用例模块 (Shadcn 版本)
func ConfigureDemoCase(b *presets.Builder, db *gorm.DB) {
	err := db.AutoMigrate(&DemoCase{})
	if err != nil {
		panic(err)
	}
	mb := b.Model(&DemoCase{})
	mb.Editing("Name").ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		p := obj.(*DemoCase)
		if p.Name == "" {
			err.FieldError("Name", "Name Can`t Empty")
		}
		return
	})
	mb.Listing("ID", "Name")

	// Shadcn 组件演示 - 简化版
	detailing := mb.Detailing(
		"InputSection",
		"ButtonSection",
		"DialogSection",
	)

	// Input 组件演示
	inputSection := presets.NewSectionBuilder(mb, "InputSection").
		Label("Input Components").
		ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.H2("Shadcn Input Components").Class("text-xl font-bold mb-4"),
				h.Div(
					shadcn.Input().Label("Text Input").Placeholder("Enter text..."),
					shadcn.Input().Label("Disabled Input").Disabled(true).Value("Disabled"),
					shadcn.Textarea().Label("Textarea").Placeholder("Enter long text...").Rows(4),
					shadcn.Checkbox().Label("Checkbox"),
				).Class("space-y-4"),
			).Class("p-4")
		})
	detailing.Section(inputSection)

	// Button 组件演示
	buttonSection := presets.NewSectionBuilder(mb, "ButtonSection").
		Label("Button Components").
		ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.H2("Shadcn Button Components").Class("text-xl font-bold mb-4"),
				h.Div(
					shadcn.Button(h.Text("Default")),
					shadcn.Button(h.Text("Secondary")),
					shadcn.Button(h.Text("Outline")),
					shadcn.Button(h.Text("Ghost")),
					shadcn.Button(h.Text("Destructive")),
					shadcn.Button(h.Text("Disabled")).Disabled(true),
				).Class("flex gap-2 flex-wrap"),
			).Class("p-4")
		})
	detailing.Section(buttonSection)

	// Dialog 组件演示
	dialogSection := presets.NewSectionBuilder(mb, "DialogSection").
		Label("Dialog Components").
		ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.H2("Shadcn Dialog Components").Class("text-xl font-bold mb-4"),
				shadcn.AlertDialog(
					shadcn.AlertDialogTrigger(
						shadcn.Button(h.Text("Open Dialog")),
					),
					shadcn.AlertDialogContent(
						shadcn.AlertDialogHeader(
							shadcn.AlertDialogTitle(h.Text("Confirm Action")),
							shadcn.AlertDialogDescription(h.Text("This is a confirmation dialog.")),
						),
						shadcn.AlertDialogFooter(
							shadcn.AlertDialogCancel(h.Text("Cancel")),
							shadcn.AlertDialogAction(h.Text("Continue")),
						),
					),
				),
			).Class("p-4")
		})
	detailing.Section(dialogSection)
}
