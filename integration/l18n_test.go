package integration_test

import (
	"net/http"
	"testing"

	"example/admin"

	"github.com/r0vx/web/multipartestutils"
)

func TestL18n(t *testing.T) {
	h, _ := admin.TestL18nHandler(TestDB)

	dbr, _ := TestDB.DB()
	profileData.TruncatePut(dbr)

	cases := []multipartestutils.TestCase{
		{
			Name:  "view by zh",
			Debug: true,
			ReqFunc: func() *http.Request {
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/auth/login").
					BuildEventFuncRequest()
				req.Header.Add("accept-language", "zh")
				return req
			},
			ExpectPageBodyContainsInOrder: []string{`邮箱`},
		},
		{
			Name:  "view by ja",
			Debug: true,
			ReqFunc: func() *http.Request {
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/auth/login").
					BuildEventFuncRequest()
				req.Header.Add("accept-language", "ja")
				return req
			},
			// 登录模块的日文翻译已于 2025-07 移除（x/login 只注册 en / zh），ja 回退英文
			ExpectPageBodyContainsInOrder: []string{`Email`},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			multipartestutils.RunCase(t, c, h)
		})
	}
}
