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
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermCreate).On("*:orders:*"),
			perm.PolicyFor(
				models.RoleViewer,
				models.RoleEditor,
				models.RoleManager,
			).WhoAre(perm.Denied).ToDo(presets.PermCreate, presets.PermUpdate, presets.PermDelete).On("*:roles:*", "*:users:*"),
			perm.PolicyFor(models.RoleViewer).WhoAre(perm.Denied).ToDo(presets.PermCreate, presets.PermUpdate, presets.PermDelete).On(perm.Anything),
			// Filter 项权限控制示例：Viewer 角色看不到用户列表的 name 筛选器
			perm.PolicyFor(models.RoleViewer).WhoAre(perm.Denied).ToDo(presets.PermList).On("*:users:fl_name:*"),
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
