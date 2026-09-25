package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/r0vx/web/multipartestutils"
	"github.com/theplant/gofixtures"

	"example/admin"
	"example/admin/ec_demo"

	"github.com/r0vx/admin/presets/actions"
)

var demoCaseData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.demo_cases (id, created_at, updated_at, deleted_at, name,field_data) VALUES (1, '2024-10-10 03:18:50.316417 +00:00', '2024-10-10 03:18:50.316417 +00:00', null, '12313','{"Text":"121231321\u0026\u0026","Textarea":"1231","TextValidate":"21312","TextareaValidate":"1😋11231"}');
`, []string{`demo_cases`}))

func TestDemoCase(t *testing.T) {
	h := admin.TestHandler(TestDB, nil)
	dbr, _ := TestDB.DB()

	cases := []TestCase{
		{
			Name:  "Index Demo Case",
			Debug: true,
			ReqFunc: func() *http.Request {
				demoCaseData.TruncatePut(dbr)
				return httptest.NewRequest("GET", "/demo-cases", http.NoBody)
			},
			ExpectPageBodyContainsInOrder: []string{"Name", "12313"},
		},
		{
			Name:  "Create Demo Case",
			Debug: true,
			ReqFunc: func() *http.Request {
				demoCaseData.TruncatePut(dbr)
				req := NewMultipartBuilder().
					PageURL("/demo-cases").
					EventFunc(actions.Update).
					AddField("Name", "test").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				m := ec_demo.DemoCase{}
				TestDB.Order("id desc").First(&m)
				if m.Name != "test" {
					t.Fatalf("Create Demo Case Failed: %v", m)
				}
			},
		},
		{
			Name:  "Demo Case Detail",
			Debug: true,
			ReqFunc: func() *http.Request {
				demoCaseData.TruncatePut(dbr)
				req := NewMultipartBuilder().
					PageURL("/demo-cases/1").
					BuildEventFuncRequest()
				return req
			},
			// 详情页现在是三个展示区块（admin/ec_demo/demo_case.go）；原先的字段保存 / 校验 / 日期校验区块已移除，对应用例一并删除
			ExpectPageBodyContainsInOrder: []string{
				"Input Components", "Text Input", "<shd-textarea", "<shd-checkbox",
				"Button Components", "Outline",
				"Dialog Components", "Open Dialog",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunCase(t, c, h)
		})
	}
}
