package admin

import (
	"net/http"

	"example/models"

	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/admin/activity"
	plogin "github.com/r0vx/admin/login"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/role"
	"github.com/r0vx/x/login"
	"github.com/r0vx/x/login/provider/wechat"
	"github.com/theplant/osenv"
	"gorm.io/gorm"
)

var (
	loginSecret      = osenv.Get("LOGIN_SECRET", "Login secret use to sign session", "")
	baseURL          = osenv.Get("BASE_URL", "Base URL for Login", "http://localhost:9500")
	// 微信扫码登录（微信开放平台「网站应用」）：
	//   1) open.weixin.qq.com 注册「网站应用」→ 拿 AppID / AppSecret；
	//   2) 应用「授权回调域」填本服务域名（只填域名、不含 https:// 和路径，须与 BASE_URL 域名一致）；
	//   3) 用环境变量注入（可写进 dev_env）：export WECHAT_APPID=...  / export WECHAT_APPSECRET=...
	//   未配置（留空）则登录页不显示「微信登录」按钮；配置后自动启用。详见 x/login/README.md。
	wechatAppID     = osenv.Get("WECHAT_APPID", "WeChat Open Platform website-app AppID (scan login)", "")
	wechatAppSecret = osenv.Get("WECHAT_APPSECRET", "WeChat Open Platform website-app AppSecret", "")
	// Telegram 登录（Login Widget，非 OAuth2）：@BotFather 建 bot → 拿 token + 用户名；
	//   用 BotFather 的 /setdomain 把 bot 域名设为 BASE_URL 的域名。两者留空则不显示 Telegram 按钮。详见 x/login/README.md。
	telegramBotToken = osenv.Get("TELEGRAM_BOT_TOKEN", "Telegram bot token (@BotFather) for login", "")
	telegramBotName  = osenv.Get("TELEGRAM_BOT_NAME", "Telegram bot username (without @)", "")
	recaptchaSiteKey   = osenv.Get("RECAPTCHA_SITE_KEY", "Recaptcha site key for Login with Recaptcha", "")
	recaptchaSecret          = osenv.Get("RECAPTCHA_SECRET_KEY", "Recaptcha secret for Login with Recaptcha", "")
	loginInitialUserEmail    = osenv.Get("LOGIN_INITIAL_USER_EMAIL", "Initial user email for Login", "")
	loginInitialUserPassword = osenv.Get("LOGIN_INITIAL_USER_PASSWORD", "Initial user password for Login", "123")
)

// getCurrentUser 从请求中获取当前登录用户
func getCurrentUser(r *http.Request) (u *models.User) {
	u, ok := login.GetCurrentUser(r).(*models.User)
	if !ok {
		return nil
	}

	return u
}

func initLoginSessionBuilder(db *gorm.DB, pb *presets.Builder, ab *activity.Builder) *plogin.SessionBuilder {
	// 按需装配 OAuth providers：仅当配置了对应凭据才启用（未配置则登录页不显示该按钮）。
	var oauthProviders []*login.Provider
	if wechatAppID != "" && wechatAppSecret != "" {
		oauthProviders = append(oauthProviders, &login.Provider{
			Goth: wechat.New(wechatAppID, wechatAppSecret,
				baseURL+"/auth/callback?provider=wechat"),
			Key:  "wechat",
			Text: "微信登录",
			Logo: h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor"><path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-2.187 6.112-2.187.202 0 .399.013.6.023C17.253 4.82 13.3 2.188 8.691 2.188zm-2.5 5.085a1.04 1.04 0 0 1-1.042-1.042 1.04 1.04 0 0 1 1.043-1.043 1.04 1.04 0 0 1 1.042 1.043 1.04 1.04 0 0 1-1.042 1.042zm5.56 0a1.04 1.04 0 0 1-1.043-1.042 1.04 1.04 0 0 1 1.042-1.043 1.04 1.04 0 0 1 1.042 1.043 1.04 1.04 0 0 1-1.042 1.042zm3.378 6.312c0 3.428 3.209 6.208 7.167 6.208.86 0 1.687-.12 2.444-.34a.72.72 0 0 1 .59.08l1.56.912a.264.264 0 0 0 .136.044c.133 0 .24-.107.24-.24 0-.06-.023-.12-.038-.177l-.327-1.233a.49.49 0 0 1 .176-.548C23.5 17.447 24.5 15.77 24.5 13.89c0-3.428-3.209-6.21-7.167-6.21s-7.163 2.783-7.163 6.21v-.105zm4.653-1.735a.869.869 0 0 1-.869-.869.869.869 0 0 1 .869-.869.869.869 0 0 1 .869.869.869.869 0 0 1-.87.869zm5.026 0a.869.869 0 0 1-.869-.869.869.869 0 0 1 .87-.869.869.869 0 0 1 .868.869.869.869 0 0 1-.869.869z"/></svg>`),
		})
	}

	loginBuilder := plogin.New(pb).
		DB(db).
		UserModel(&models.User{}).
		Secret(loginSecret).
		OAuthProviders(oauthProviders...).
		HomeURLFunc(func(r *http.Request, user interface{}) string {
			return "/"
		}).
		MaxRetryCount(5).
		// TODO online  to set  true
		Recaptcha(false, login.RecaptchaConfig{
			SiteKey:   recaptchaSiteKey,
			SecretKey: recaptchaSecret,
		}).
		WrapBeforeSetPassword(func(in login.HookFunc) login.HookFunc {
			return func(r *http.Request, user interface{}, extraVals ...interface{}) error {
				if err := in(r, user, extraVals...); err != nil {
					return err
				}
				u := user.(*models.User)
				if u.GetAccountName() == loginInitialUserEmail {
					return &login.NoticeError{
						Level:   login.NoticeLevel_Error,
						Message: "Cannot change password for public user",
					}
				}
				password := extraVals[0].(string)
				if len(password) < 12 {
					return &login.NoticeError{
						Level:   login.NoticeLevel_Error,
						Message: "Password cannot be less than 12 characters",
					}
				}
				return nil
			}
		}).
		// 移除OAuth完成后的钩子函数
		// 两步验证：可选模式——用户在资料菜单自助开/关
		TOTPMode(login.TOTPOptional, login.TOTPConfig{Issuer: "r0vx"}).
		MaxRetryCount(6)
	// 使用默认登录页面（shadcn-vue 风格）
	// loginBuilder.LoginPageFunc(plogin.NewAdvancedLoginPage(...)) 注释掉后使用默认样式
	// 注：微信自助绑定的回调处理已内置于框架（x/login，识别绑定 state 自动绑定），无需在此注册钩子。

	// Telegram 登录（Login Widget，非 OAuth2）：留空 token/name 时为 no-op、登录页不显示按钮
	loginBuilder.TelegramLogin(telegramBotToken, telegramBotName)

	genInitialUser(db)

	return plogin.NewSessionBuilder(loginBuilder, db).
		Activity(ab.RegisterModel(&models.User{})).
		IsPublicUser(func(u interface{}) bool {
			user, ok := u.(*models.User)
			if !ok {
				return false
			}
			return user.GetAccountName() == loginInitialUserEmail
		}).
		TablePrefix("cms_").
		// WithSessionTableHook(func(next plogin.SessionTableFunc) plogin.SessionTableFunc {
		// 	return func(ctx context.Context, input *plogin.SessionTableInput) (*plogin.SessionTableOutput, error) {
		// 		output, err := next(ctx, input)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		output.Component = h.Components(
		// 			output.Component,
		// 			h.Div().Class("text-caption pt-2 text-warning").Text("Customized Bottom Text"),
		// 		)
		// 		return output, nil
		// 	}
		// }).
		// ParseIPFunc(func(ctx context.Context, lang language.Tag, addr string) (string, error) {
		// 	city, err := locationDB.GetCity(ctx, addr)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	return location.GeneralLocalizedCountryCity(city, lang, language.English), nil
		// }).
		AutoMigrate()
}

func genInitialUser(db *gorm.DB) {
	email := loginInitialUserEmail
	password := loginInitialUserPassword
	if email == "" || password == "" {
		return
	}

	var count int64
	if err := db.Model(&models.User{}).Where("account = ?", email).Count(&count).Error; err != nil {
		panic(err)
	}

	if count > 0 {
		return
	}
	if err := initDefaultRoles(db); err != nil {
		panic(err)
	}

	user := &models.User{
		Name:   email,
		Status: models.StatusActive,
		UserPass: login.UserPass{
			Account:  email,
			Password: password,
		},
	}
	user.EncryptPassword()
	if err := db.Create(user).Error; err != nil {
		panic(err)
	}
	if err := grantUserRole(db, user.ID, models.RoleManager); err != nil {
		panic(err)
	}
}

func grantUserRole(db *gorm.DB, userID uint, roleName string) error {
	var roleID int
	if err := db.Table("roles").Where("name = ?", roleName).Pluck("id", &roleID).Error; err != nil {
		panic(err)
	}
	return db.Table("user_role_join").Create(
		&map[string]interface{}{
			"user_id": userID,
			"role_id": roleID,
		}).Error
}

func initDefaultRoles(db *gorm.DB) error {
	var cnt int64
	if err := db.Model(&role.Role{}).Count(&cnt).Error; err != nil {
		return err
	}

	if cnt == 0 {
		var roles []*role.Role
		for _, r := range models.DefaultRoles {
			roles = append(roles, &role.Role{
				Name: r,
			})
		}

		if err := db.Create(roles).Error; err != nil {
			return err
		}
	}

	return nil
}
