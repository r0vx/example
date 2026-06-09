package models

import "time"

// FilterDemo Filter 组件全类型演示模型
//
// 用于展示 admin 框架所有 12 种 FilterItemType 的效果：
//   - StringItem / NumberItem / SelectItem / MultipleSelectItem
//   - AutoCompleteItem / DateItem / DateRangeItem / DatePickerItem
//   - DateRangePickerItem / DatetimeRangeItem / DatetimeRangePickerItem
//   - LinkageSelectItem
type FilterDemo struct {
	ID        uint      `gorm:"primarykey"`
	Title     string    // 标题（StringItem）
	Amount    float64   // 金额（NumberItem）
	Status    string    // 状态（SelectItem / MultipleSelectItem）
	Category  string    // 分类（AutoCompleteItem）
	IsActive  bool      // 是否启用（BooleanItem）
	Country   string    // 国家（LinkageSelectItem 第一级）
	Province  string    // 省/州（LinkageSelectItem 第二级）
	City      string    // 城市（LinkageSelectItem 第三级）
	UpdatedAt time.Time // 更新时间
	CreatedAt time.Time // 创建时间（Date/DateRange/DatePicker/DatetimeRange 等）
}
