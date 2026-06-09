package models

import (
	"time"

	"github.com/lib/pq"
	"github.com/r0vx/admin/media/media_library"
)

type InputDemo struct {
	ID               uint
	TextField1       string
	TextArea1        string
	Switch1          bool
	Slider1          int
	Select1          string
	RangeSlider1     pq.Int64Array `gorm:"type:integer[]"`
	Radio1           string
	FileInput1       string
	Combobox1        string
	Checkbox1        bool
	Autocomplete1    string `gorm:"type:jsonb"` // JSON 格式 ["a","b","c"]
	ButtonGroup1     string
	Badge            string
	BadgeSelect      string
	ItemGroup1       string
	ListItemGroup1   string
	SlideGroup1      string
	ColorPicker1     string
	DatePicker1      string
	DatePickerMonth1 string
	TimePicker1      string
	CodeMirror1      string                 `gorm:"type:text"` // 代码编辑器内容
	MediaLibrary1    media_library.MediaBox `sql:"type:text;"`
	SelectedCustomers string                `gorm:"type:text"` // 逗号分隔的 Customer ID 列表
	Location          string                `gorm:"type:text"` // 高德地图选点 JSON（经纬度+地址）
	UpdatedAt        time.Time
	CreatedAt        time.Time
}
