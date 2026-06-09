package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnDataTableDemo 虚拟模型
type ShadcnDataTableDemo struct{}

// configDataTable 注册 Data Table demo
func configDataTable(b *presets.Builder) {
	m := b.Model(&ShadcnDataTableDemo{}).
		Label("Data Table").
		URIName("shadcn-data-table")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnDataTableDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "DataTable Demo"
		r.Body = shadcnDataTableBody(ctx)
		return
	})
}

// shadcnDataTableBody DataTable 演示页面
func shadcnDataTableBody(ctx *web.EventContext) h.HTMLComponent {
	// 示例数据
	users := []map[string]interface{}{
		{"id": "1", "name": "张三", "email": "zhangsan@example.com", "role": "管理员", "status": "active"},
		{"id": "2", "name": "李四", "email": "lisi@example.com", "role": "编辑", "status": "active"},
		{"id": "3", "name": "王五", "email": "wangwu@example.com", "role": "访客", "status": "inactive"},
		{"id": "4", "name": "赵六", "email": "zhaoliu@example.com", "role": "编辑", "status": "active"},
		{"id": "5", "name": "孙七", "email": "sunqi@example.com", "role": "管理员", "status": "active"},
		{"id": "6", "name": "周八", "email": "zhouba@example.com", "role": "访客", "status": "inactive"},
		{"id": "7", "name": "吴九", "email": "wujiu@example.com", "role": "编辑", "status": "active"},
		{"id": "8", "name": "郑十", "email": "zhengshi@example.com", "role": "访客", "status": "active"},
		{"id": "9", "name": "冯十一", "email": "fengshiyi@example.com", "role": "管理员", "status": "inactive"},
		{"id": "10", "name": "陈十二", "email": "chenshier@example.com", "role": "编辑", "status": "active"},
	}

	columns := []DataTableColumn{
		{Name: "name", Title: "姓名", Sortable: true},
		{Name: "email", Title: "邮箱", Sortable: true},
		{Name: "role", Title: "角色", Sortable: true},
		{Name: "status", Title: "状态", Sortable: true},
	}

	return h.Div(
		h.H1("DataTable Demo").Style("margin-bottom: 24px;"),
		h.P(h.Text("高级数据表格组件，支持排序、分页、行选择、行菜单等功能")).Class("text-muted-foreground mb-6"),

		// 基础表格
		h.Div(
			h.H2("Basic Table"),
			h.P(h.Text("基础数据表格")).Class("text-muted-foreground mb-4"),
			DataTable().
				Data(users[:5]).
				Columns(columns),
		).Class("demo-section"),

		// 可排序表格
		h.Div(
			h.H2("Sortable Table"),
			h.P(h.Text("点击列标题进行排序")).Class("text-muted-foreground mb-4"),
			DataTable().
				Data(users).
				Columns(columns).
				Hover(true),
		).Class("demo-section"),

		// 带分页的表格
		h.Div(
			h.H2("Table with Pagination"),
			h.P(h.Text("分页显示数据")).Class("text-muted-foreground mb-4"),
			DataTable().
				Data(users).
				Columns(columns).
				Pagination(true).
				PageSize(5),
		).Class("demo-section"),

		// 可选择行的表格
		h.Div(
			h.H2("Selectable Table"),
			h.P(h.Text("支持行选择和批量操作")).Class("text-muted-foreground mb-4"),
			web.Scope(
				h.Div(
					DataTable().
						Data(users).
						Columns(columns).
						Selectable(true).
						Pagination(true).
						PageSize(5).
						Attr("v-model:selected-ids", "form.selectedIds").
						Attr("@update:selected-ids", "form.selectedIds = $event"),
				).Class("mb-4"),
				h.Div(
					h.H4("Selected IDs:").Class("text-sm font-medium mb-2"),
					h.Pre("{{ JSON.stringify(form.selectedIds, null, 2) }}").Class("text-xs bg-muted p-2 rounded"),
				),
			).VSlot("{ form }").FormInit(`{ "selectedIds": [] }`),
		).Class("demo-section"),

		// 带行菜单的表格
		h.Div(
			h.H2("Table with Row Menu"),
			h.P(h.Text("每行带操作菜单")).Class("text-muted-foreground mb-4"),
			web.Scope(
				DataTable().
					Data(users).
					Columns(columns).
					Pagination(true).
					PageSize(5).
					Children(
						// 行菜单插槽
						h.Tag("template").Attr("v-slot:row-menu", "{ row }").Children(
							h.Tag("shd-dropdown-menu-item").
								Attr("@click", "console.log('Edit: ' + row.name)").
								Children(h.Text("编辑")),
							h.Tag("shd-dropdown-menu-item").
								Attr("@click", "console.log('Delete: ' + row.name)").
								Children(h.Text("删除")),
							h.Tag("shd-dropdown-menu-separator"),
							h.Tag("shd-dropdown-menu-item").
								Attr("@click", "console.log('View: ' + row.name)").
								Children(h.Text("查看详情")),
						),
					),
			).VSlot("{ form }").FormInit(`{}`),
		).Class("demo-section"),

		// 自定义单元格渲染
		h.Div(
			h.H2("Custom Cell Rendering"),
			h.P(h.Text("自定义单元格内容")).Class("text-muted-foreground mb-4"),
			DataTable().
				Data(users).
				Columns(columns).
				Pagination(true).
				PageSize(5).
				Children(
					// 自定义 status 列
					h.Tag("template").Attr("v-slot:cell-status", "{ value }").Children(
						h.Tag("shd-badge").
							Attr(":variant", "value === 'active' ? 'default' : 'secondary'").
							Children(h.RawHTML("{{ value === 'active' ? '活跃' : '未激活' }}")),
					),
					// 自定义 role 列
					h.Tag("template").Attr("v-slot:cell-role", "{ value }").Children(
						h.Tag("shd-badge").
							Attr("variant", "outline").
							Children(h.RawHTML("{{ value }}")),
					),
				),
		).Class("demo-section"),

		// 空数据状态
		h.Div(
			h.H2("Empty State"),
			h.P(h.Text("无数据时的显示")).Class("text-muted-foreground mb-4"),
			DataTable().
				Data([]map[string]interface{}{}).
				Columns(columns).
				EmptyText("暂无数据"),
		).Class("demo-section"),

		// 完整示例：用户管理表格
		h.Div(
			h.H2("Complete Example - User Management"),
			h.P(h.Text("完整的用户管理表格示例")).Class("text-muted-foreground mb-4"),
			web.Scope(
				Card(
					CardHeader(
						CardTitle(h.Text("用户列表")),
						CardDescription(h.Text("管理系统用户")),
					),
					CardContent(
						DataTable().
							Data(users).
							Columns(columns).
							Selectable(true).
							Pagination(true).
							PageSize(5).
							SelectedCountLabel("已选择 {count} 条记录").
							ClearSelectionLabel("清除选择").
							Attr("v-model:selected-ids", "form.selectedUsers").
							Children(
								// 状态列自定义渲染
								h.Tag("template").Attr("v-slot:cell-status", "{ value }").Children(
									h.Tag("shd-badge").
										Attr(":variant", "value === 'active' ? 'default' : 'secondary'").
										Children(h.RawHTML("{{ value === 'active' ? '活跃' : '未激活' }}")),
								),
								// 角色列自定义渲染
								h.Tag("template").Attr("v-slot:cell-role", "{ value }").Children(
									h.Tag("shd-badge").
										Attr("variant", "outline").
										Attr(":class", `{
											'bg-red-100 text-red-800': value === '管理员',
											'bg-blue-100 text-blue-800': value === '编辑',
											'bg-gray-100 text-gray-800': value === '访客'
										}`).
										Children(h.RawHTML("{{ value }}")),
								),
								// 行菜单
								h.Tag("template").Attr("v-slot:row-menu", "{ row }").Children(
									h.Tag("shd-dropdown-menu-item").
										Attr("@click", "console.log('编辑用户: ' + row.name)").
										Children(h.Text("✏️ 编辑")),
									h.Tag("shd-dropdown-menu-item").
										Attr("@click", "console.log('重置密码: ' + row.name)").
										Children(h.Text("🔑 重置密码")),
									h.Tag("shd-dropdown-menu-separator"),
									h.Tag("shd-dropdown-menu-item").
										Attr("@click", "console.log('删除用户: ' + row.name)").
										Attr("class", "text-red-600").
										Children(h.Text("🗑️ 删除")),
								),
							),
					),
					CardFooter(
						h.Div(
							Button(h.Text("批量删除")).
								Variant(ButtonVariantDestructive).
								Disabled(true).
								Attr(":disabled", "form.selectedUsers.length === 0").
								Attr("@click", "console.log('删除 ' + form.selectedUsers.length + ' 个用户')").
								Class("mr-2"),
							Button(h.Text("导出选中")).
								Variant(ButtonVariantOutline).
								Disabled(true).
								Attr(":disabled", "form.selectedUsers.length === 0").
								Attr("@click", "console.log('导出 ' + form.selectedUsers.length + ' 个用户')"),
						),
					),
				),
			).VSlot("{ form }").FormInit(`{ "selectedUsers": [] }`),
		).Class("demo-section"),
	).Style("max-width: 1200px; margin: 0 auto;")
}
