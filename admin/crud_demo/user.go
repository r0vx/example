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
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"gorm.io/gorm"
)

// ConfigUser 配置用户管理模块
func ConfigUser(b *presets.Builder, ab *activity.Builder, db *gorm.DB, publisher *publish.Builder, loginSessionBuilder *plogin.SessionBuilder) {
	mb := b.Model(&models.User{}).URIName("users")
	defer mb.Use(ab)

	lb := mb.Listing("ID", "Name", "Account", "Company", "Status", "Roles", "UpdatedAt").
		SearchColumns("name", "account")

	rmb := lb.RowMenu().InlineDefaultsInMenu(true)

	// test
	rmb.RowMenuItem("readmeDoc").ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
		cu := obj.(*models.User)

		return shadcn.RowMenuItem("test").SetOnclick(
			fmt.Sprintf("vars.__window.open('/readme?id=%d', '_blank')", cu.ID),
			//web.GET().URL("/readme").PushState(true).Query("id", fmt.Sprint(cu.ID)).Go(),
		)

	})

	lb.Field("Roles").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		user := obj.(*models.User)
		roles := user.GetRoles()
		badges := make([]h.HTMLComponent, 0, len(roles))
		for _, r := range roles {
			badges = append(badges, shadcn.Badge(h.Text(r)).Variant(shadcn.BadgeVariantSecondary))
		}
		return h.Div(badges...).Class("flex gap-1 flex-wrap")
	})

	lb.Field("UpdatedAt").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
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
	).ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.User)
		if u.Name == "" {
			err.FieldError("Name", "Name is required")
		}
		if u.Account == "" {
			err.FieldError("Account", "Account is required")
		}
		return
	})

	eb := mb.Editing()
	eb.Field("Roles").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		user := obj.(*models.User)
		roleItems := make([]shadcn.DefaultOptionItem, 0)
		for _, r := range models.DefaultRoles {
			roleItems = append(roleItems, shadcn.DefaultOptionItem{Text: r, Value: r})
		}
		currentRoles := user.GetRoles()
		return shadcn.Select().Items(roleItems).Label(field.Label).
			Attr(web.VField(field.FormKey, currentRoles)...).
			Attr("multiple", "true").
			ErrorMessages(field.Errors...)
	})

	dp := mb.Detailing(
		&presets.FieldsSection{
			Title: "User Info",
			Rows:  [][]string{{"ID", "Name"}, {"Company", "Status"}, {"Account"}, {"Roles"}, {"RegistrationDate"}},
		},
	).Drawer(true)

	dp.Field("Roles").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
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
	mb.Editing().Field("FavorPostID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
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
				Attr("readonly", "true").Class("hidden"),
			shadcn.Button(h.Text(postTitle)).
				Variant(shadcn.ButtonVariantOutline).
				On("click", "locals.showPostSelector = true"),
			h.Span(postTitle).Class("ml-2"),
		).Class("flex items-center gap-2")
	})
}
