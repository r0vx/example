package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/r0vx/web"
	"github.com/r0vx/web/multipartestutils"
	"gorm.io/gorm"

	"example/admin"

	"example/models"

	plogin "github.com/r0vx/admin/login"
	"github.com/r0vx/admin/role"
)

func TestLogin(t *testing.T) {
	h, _ := admin.TestL18nHandler(TestDB)

	dbr, _ := TestDB.DB()
	profileData.TruncatePut(dbr)

	cases := []multipartestutils.TestCase{
		{
			Name:  "view by en",
			Debug: false,
			ReqFunc: func() *http.Request {
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/auth/login").
					BuildEventFuncRequest()
				req.Header.Add("accept-language", "en")
				return req
			},
			ExpectPageBodyContainsInOrder: []string{`MATERIALPRO`, `邮箱`, `密码`, `忘记密码？`, `登录`},
		},
		{
			Name:  "view by zh",
			Debug: false,
			ReqFunc: func() *http.Request {
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/auth/login").
					BuildEventFuncRequest()
				req.Header.Add("accept-language", "zh")
				return req
			},
			ExpectPageBodyContainsInOrder: []string{`MATERIALPRO`, `Email`, `Password`, `Forget your password?`, `Sign in`},
		},
		{
			Name:  "view by ja",
			Debug: false,
			ReqFunc: func() *http.Request {
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/auth/login").
					BuildEventFuncRequest()
				req.Header.Add("accept-language", "ja")
				return req
			},
			ExpectPageBodyContainsInOrder: []string{`MATERIALPRO`, `Email`, `Password`, `Forget your password?`, `Sign in`},
		},
		{
			Name:  "view by en (customized)",
			Debug: false,
			HandlerMaker: func() http.Handler {
				mux := http.NewServeMux()
				c := admin.NewConfig(TestDB, false)
				loginSessionBuilder := c.GetLoginSessionBuilder()
				loginBuilder := loginSessionBuilder.GetLoginBuilder()
				loginBuilder.LoginPageFunc(plogin.NewAdvancedLoginPage(func(ctx *web.EventContext, config *plogin.AdvancedLoginPageConfig) (*plogin.AdvancedLoginPageConfig, error) {
					config.SignInButtonLabel = "Hello"
					return config, nil
				})(loginBuilder.ViewHelper(), c.GetPresetsBuilder()))
				loginSessionBuilder.Secret("test")
				loginSessionBuilder.Mount(mux)
				mux.Handle("/", c.GetPresetsBuilder())
				return mux
			},
			ReqFunc: func() *http.Request {
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/auth/login").
					BuildEventFuncRequest()
				req.Header.Add("accept-language", "en")
				return req
			},
			ExpectPageBodyContainsInOrder: []string{`MATERIALPRO`, `Email`, `Password`, `Forget your password?`, `Hello`},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			multipartestutils.RunCase(t, c, h)
		})
	}
}

func TestChangePassword(t *testing.T) {
	h, c := admin.TestHandlerComplex(TestDB, &models.User{
		Model: gorm.Model{ID: 888},
		Name:  "viwer@theplant.jp",
		Roles: []role.Role{
			{
				Name: models.RoleEditor,
			},
		},
	}, false)
	// dbr, _ := TestDB.DB()
	sb := c.GetLoginSessionBuilder()
	sb.Secret("test")
	sb.Mount(h.(*http.ServeMux))

	cases := []multipartestutils.TestCase{
		{
			Name:  "View Change Password Page",
			Debug: true,
			ReqFunc: func() *http.Request {
				return httptest.NewRequest("GET", "/auth/change-password", http.NoBody)
			},
			ExpectPageBodyContainsInOrder: []string{"Change your password", "Old password", "New password", "zxcvbn.js", "Re-enter new password"},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			multipartestutils.RunCase(t, c, h)
		})
	}
}
