package models

// City 参考表：演示「DB 驱动 Select 选项 + 列表按 ID 查表显示名称」。
// InputDemo.Select1 存 City.ID（字符串），编辑下拉来自本表，列表显示 City.Name。
type City struct {
	ID   uint
	Name string
}
