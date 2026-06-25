package admin

import (
	"net/http"

	"example/models"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/x/perm"
	"gorm.io/gorm"
)

func initPermission(b *presets.Builder, db *gorm.DB) {
	perm.Verbose = true
	b.Permission(
		perm.New().Policies(
			//perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
			// Admin 超管放行（静态兜底，与 DB 策略叠加；Denied 仍可覆盖 Allow）。
			// 注意：seo 编辑权限闸 editIsAllowed 接线后（移植回归修复），启用 perm 的项目
			// 必须有能匹配 `:seo:seo_settings:` + `perm_seo_edit` 的 allow 策略，否则 SEO 不可编辑。
			perm.PolicyFor(models.RoleAdmin).WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
			// RowLevelRefresh demo：放行所有人访问（演示用，新模型默认无授权策略会被拒）
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On("*:row_refresh_demo:*"),
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On("*:relay_pagination_demo:*"),
			// 功能角色 demo（同角色不同用户、不同权限）：在 Editor 等「基础角色」之上，给用户附加
			// ProductManager / UserManager 这种「功能角色」按模块细分能力。subject 字符串须与角色 Name 一致
			//（User.GetRoles() 返回 role.Name）。用户 A=[Editor,ProductManager]→能管产品；B=[Editor,UserManager]→能管用户。
			// ⚠️ 多角色是权限并集（OR）：要互斥，须保证基础角色 Editor 不含 products/users 管理权，否则会从 Editor 漏出。
			// ⚠️ 要让此 demo 可验证：还需在 roles 后台（或 models.DefaultRoles）建同名角色并分配给测试用户。
			perm.PolicyFor("ProductManager").WhoAre(perm.Allowed).ToDo(perm.Anything).On("*:products:*"),
			perm.PolicyFor("UserManager").WhoAre(perm.Allowed).ToDo(perm.Anything).On("*:users:*"),
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermCreate).On("*:orders:*"),
			perm.PolicyFor(
				models.RoleViewer,
				models.RoleEditor,
				models.RoleManager,
			).WhoAre(perm.Denied).ToDo(presets.PermCreate, presets.PermUpdate, presets.PermDelete).On("*:roles:*", "*:users:*"),
			// 【已注释】原本：Viewer 全局禁 增/改/删（On Anything）。这是静态兜底策略，
			// ladon 中 Deny 恒覆盖 Allow，会盖掉 roles 后台（权限树）给 Viewer 授的任何 DB 编辑权，
			// 导致「明明在树里设了编辑权限却仍被拒」。注释后 Viewer 的增改删完全由 DB 策略（权限树）说了算。
			// 如需恢复「Viewer 全局只读」演示，取消下一行注释即可。
			// perm.PolicyFor(models.RoleViewer).WhoAre(perm.Denied).ToDo(presets.PermCreate, presets.PermUpdate, presets.PermDelete).On(perm.Anything),
			// Filter 项权限控制示例：Viewer 角色看不到用户列表的 name 筛选器
			perm.PolicyFor(models.RoleViewer).WhoAre(perm.Denied).ToDo(presets.PermList).On("*:users:fl_name:*"),
			// RowMenuItem 权限控制示例（fm_，与筛选项 fl_ 同构）：Viewer 角色看不到、也无法操作
			// action-enhance-demo 列表的「设置费率」行操作按钮（资源 fm_set_user_poundage、action presets:list）。
			// 隐藏按钮 + 服务端拒绝事件（防绕过 UI）；日志 Resource:"*:action_enhance_demo:fm_set_user_poundage:"。
			perm.PolicyFor(models.RoleAdmin).WhoAre(perm.Denied).ToDo(presets.PermList).On("*:action_enhance_demo:fm_set_user_poundage:*"),

			perm.PolicyFor(models.RoleAdmin).WhoAre(perm.Denied).ToDo(presets.PermUpdate).On("*:action_enhance_demo:fm_set_user_poundage:*"),

			// Action / BulkAction 权限控制示例（fa_ 资源闸，与 fm_/fl_ 同构）：
			// Viewer 看不到也无法执行 action-enhance-demo 列表的「刷新」行动作与「批量归档」批量操作。
			// 资源串格式：*:<uri>:*fa_<snake(name)>:*，action = presets:list。
			//   Action("Refresh")        → *:action_enhance_demo:*fa_refresh:*
			//   BulkAction("BulkArchive")→ *:action_enhance_demo:*fa_bulk_archive:*
			// 双层拦截：渲染隐藏按钮（listing_compo.go actionsComponent / floatingBulkBar）
			//           + 执行拒绝（fetchAction / fetchBulkAction）。
			// 开 perm.Verbose 时日志里 Resource:"*:action_enhance_demo:*fa_refresh:*" 可直接看到真实串。
			perm.PolicyFor(models.RoleViewer).WhoAre(perm.Denied).ToDo(presets.PermList).On("*:action_enhance_demo:*fa_refresh:*"),
			perm.PolicyFor(models.RoleViewer).WhoAre(perm.Denied).ToDo(presets.PermList).On("*:action_enhance_demo:*fa_bulk_archive:*"),

			// FilterTab 权限控制示例（ft_，与筛选项 fl_ 同构）：Viewer 看不到 input-demos 列表的「启用」筛选标签。
			// 资源 *:input_demos:ft_enabled:*、action presets:list；框架渲染时自动剔除无权限 tab
			//（listing_compo.go 的 compactTabsFilter），无需在 FilterTabsFunc 里写 if。tab 须显式给 ID（这里 "enabled"）。
			perm.PolicyFor(models.RoleViewer).WhoAre(perm.Denied).ToDo(presets.PermList).On("*:input_demos:ft_enabled:*"),

			// OptionsFunc 按用户 ID 过滤示例：Viewer 角色在订单 Source 筛选器中看不到 ID=1007 的用户
			perm.PolicyFor(models.RoleViewer).WhoAre(perm.Denied).ToDo(presets.PermList).On(":presets:users:users:1001:"),
			// 字段权限示例：所有人对 input_demos 的 Switch1 字段 PermUpdate 只读 →
			// 列表里那个可交互开关禁用（cell honor field.Disabled，与列可见性同款 f_<name> 资源、无 ObjectOn）。
			// 注意资源名是 SnakeOn 后的形式：字段 Switch1 → f_switch_1（数字前带下划线），不是 f_switch1。
			// 如需连编辑表单的 Switch1 也禁用，用更宽的 "*:f_switch_1:*"（表单路径经 ObjectOn 多一段记录 id）。
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermUpdate).On("*:input_demos:f_switch_1:*"),
		).SubjectsFunc(func(r *http.Request) []string {
			u := getCurrentUser(r)
			if u == nil {
				return nil
			}
			return u.GetRoles()
		}).DBPolicy(perm.NewDBPolicy(db)),
	)
}
