package crud_demo

import (
	"fmt"

	"example/models"

	"github.com/r0vx/x/i18n"
	"github.com/r0vx/x/ui/shadcn"

	"github.com/r0vx/admin/activity"
	plogin "github.com/r0vx/admin/login"
	"github.com/r0vx/admin/publish"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/presets/gorm2op"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/perm"
	"gorm.io/gorm"
)

// ConfigUser 配置用户管理模块
func ConfigUser(b *presets.Builder, ab *activity.Builder, db *gorm.DB, publisher *publish.Builder, loginSessionBuilder *plogin.SessionBuilder) {
	mb := b.Model(&models.User{}).URIName("users")
	defer mb.Use(ab)

	lb := mb.Listing("ID", "Name", "Account", "Company", "Status", "Roles", "UpdatedAt").
		SearchColumns("name", "account")

	// Roles 是 m2m（User.Roles []perm.Role gorm:"many2many"），默认查询不带关联 → 列表角色列空。
	// 经 gorm2op hook 给列表查询注入 Preload("Roles")，使每行 user.GetRoles() 有值。
	lb.WrapSearchFunc(func(in presets.SearchFunc) presets.SearchFunc {
		return func(ctx *web.EventContext, params *presets.SearchParams) (*presets.SearchResult, error) {
			return in(gorm2op.EventContextWithHook(ctx, func(db *gorm.DB) *gorm.DB {
				return db.Preload("Roles")
			}), params)
		}
	})

	rmb := lb.RowMenu().InlineDefaultsInMenu(true)

	// test
	rmb.RowMenuItem("readmeDoc").ComponentFunc(func(obj any, id string, ctx *web.EventContext) h.HTMLComponent {
		cu := obj.(*models.User)

		return shadcn.RowMenuItem("test").SetOnclick(
			fmt.Sprintf("vars.__window.open('/readme?id=%d', '_blank')", cu.ID),
			//web.GET().URL("/readme").PushState(true).Query("id", fmt.Sprint(cu.ID)).Go(),
		)

	})

	lb.Field("Roles").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		user := obj.(*models.User)
		roles := user.GetRoles()
		badges := make([]h.HTMLComponent, 0, len(roles))
		for _, r := range roles {
			badges = append(badges, shadcn.Badge(h.Text(r)).Variant(shadcn.BadgeVariantSecondary))
		}
		return h.Div(badges...).Class("flex gap-1 flex-wrap")
	})

	lb.Field("UpdatedAt").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		user := obj.(*models.User)
		v := user.UpdatedAt.Local().Format("2006-01-02 15:04:05")
		return h.Text(v)
	})

	lb.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		statusOptions := []shadcn.FilterSelectOption{
			{Text: "Active", Value: "active"},
			{Text: "Inactive", Value: "inactive"},
		}
		return []*shadcn.FilterItem{
			{
				Key:          "created_at",
				Label:        "Create Time",
				ItemType:     shadcn.FilterItemTypeDatetimeRangePicker,
				SQLCondition: `created_at %s ?`,
			},
			{
				Key:          "status",
				Label:        "Status",
				ItemType:     shadcn.FilterItemTypeSelect,
				Options:      statusOptions,
				SQLCondition: `status %s ?`,
			},
		}
	})

	mb.Editing(
		&presets.FieldsSection{
			Title: "Basic Info",
			Rows:  [][]string{{"Name", "Company"}, {"Status", "Account"}, {"RegistrationDate"}, {"Roles"}},
		},
	).ValidateFunc(func(obj any, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.User)
		if u.Name == "" {
			err.FieldError("Name", "Name is required")
		}
		if u.Account == "" {
			err.FieldError("Account", "Account is required")
		}
		return
	})

	// loadUserRoles 把 m2m Roles 关联加载进 *User（Fetch 不走 gorm2op 的 preload hook，须手动）。
	loadUserRoles := func(obj any) {
		if u, ok := obj.(*models.User); ok && u.ID != 0 {
			_ = db.Model(u).Association("Roles").Find(&u.Roles)
		}
	}

	eb := mb.Editing()
	eb.Field("Roles").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		user := obj.(*models.User)
		roleItems := make([]shadcn.DefaultOptionItem, 0)
		for _, r := range models.DefaultRoles {
			roleItems = append(roleItems, shadcn.DefaultOptionItem{Text: r, Value: r})
		}
		currentRoles := user.GetRoles()
		return shadcn.Select().Items(roleItems).Label(field.Label).
			Attr(web.VField(field.FormKey, currentRoles)...).
			Multiple(true).
			ErrorMessages(field.Errors...)
	}).SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
		// 显式清空 Roles：presets 会先批量 reflectutils unmarshal 整个 form，把 "Manager" 塞进
		// []perm.Role 得到一个零值 Role；若不清掉，主 Save 的 Select("*").Updates 会连带 upsert 该零值
		// （INSERT roles name=''、user_role_join role_id=0）→ 外键违例。真正的 m2m 写入在 WrapSaveFunc。
		if u, ok := obj.(*models.User); ok {
			u.Roles = nil
		}
		return nil
	})

	// 编辑/新建打开表单时把当前角色加载进来，供 ComponentFunc 回显已选。
	eb.WrapFetchFunc(func(in presets.FetchFunc) presets.FetchFunc {
		return func(obj any, id string, ctx *web.EventContext) (any, error) {
			r, err := in(obj, id, ctx)
			if err == nil {
				loadUserRoles(r)
			}
			return r, err
		}
	})

	// 保存后用表单提交的角色名重置 m2m 关联（角色名 → perm.Role → Association.Replace）。
	eb.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj any, id string, ctx *web.EventContext) error {
			_ = ctx.R.ParseForm()
			names := ctx.R.Form["Roles"] // 须在 in() 前取（multipart 表单值）
			// 清空 obj.Roles 再存：杜绝主 Save 误 upsert 零值/旧关联（role_id=0 外键违例）；m2m 全交给下方 Replace。
			if u, ok := obj.(*models.User); ok {
				u.Roles = nil
			}
			if err := in(obj, id, ctx); err != nil {
				return err
			}
			u, ok := obj.(*models.User)
			if !ok {
				return nil
			}
			var roles []perm.Role
			if len(names) > 0 {
				if err := db.Where("name IN ?", names).Find(&roles).Error; err != nil {
					return err
				}
			}
			return db.Model(u).Association("Roles").Replace(roles)
		}
	})

	dp := mb.Detailing(
		&presets.FieldsSection{
			Title: "User Info",
			Rows:  [][]string{{"ID", "Name"}, {"Company", "Status"}, {"Account"}, {"Roles"}, {"RegistrationDate"}},
		},
	).Drawer(true)

	// 详情页同样需手动加载 m2m Roles（Fetch 不走 preload hook）。
	dp.WrapFetchFunc(func(in presets.FetchFunc) presets.FetchFunc {
		return func(obj any, id string, ctx *web.EventContext) (any, error) {
			r, err := in(obj, id, ctx)
			if err == nil {
				loadUserRoles(r)
			}
			return r, err
		}
	})

	dp.Field("Roles").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		user := obj.(*models.User)
		roles := user.GetRoles()
		badges := make([]h.HTMLComponent, 0, len(roles))
		for _, r := range roles {
			badges = append(badges, shadcn.Badge(h.Text(r)).Variant(shadcn.BadgeVariantSecondary))
		}
		return h.Div(
			h.Div(
				h.Label(i18n.PT(ctx.R, presets.ModelsI18nModuleKey, mb.Info().Label(), field.Label)).Class("text-sm font-medium mb-1 block"),
				h.Div(badges...).Class("flex gap-1 flex-wrap"),
			).Class("mb-4"),
		)
	})

	// FavorPost 选择对话框
	ConfigureFavorPostSelectDialog(db, mb, publisher)
}

// ConfigureFavorPostSelectDialog 配置收藏文章选择对话框
func ConfigureFavorPostSelectDialog(db *gorm.DB, mb *presets.ModelBuilder, publisher *publish.Builder) {
	mb.Editing().Field("FavorPostID").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		user := obj.(*models.User)
		var postTitle string
		if user.FavorPostID > 0 {
			var post models.Post
			if err := db.First(&post, user.FavorPostID).Error; err == nil {
				postTitle = post.Title
			}
		}
		return h.Div(
			shadcn.Input().Label(field.Label).Value(fmt.Sprint(user.FavorPostID)).
				Attr(web.VField(field.FormKey, fmt.Sprint(user.FavorPostID))...).
				Readonly(true).Class("hidden"),
			shadcn.Button(h.Text(postTitle)).
				Variant(shadcn.ButtonVariantOutline).
				On("click", "locals.showPostSelector = true"),
			h.Span(postTitle).Class("ml-2"),
		).Class("flex items-center gap-2")
	})
}
