package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"example/admin"
	"example/admin/pagebuilder/containers"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets/actions"
	"github.com/r0vx/admin/publish"
	"github.com/r0vx/admin/seo"
	. "github.com/r0vx/web/multipartestutils"
	"github.com/theplant/gofixtures"
	"gorm.io/gorm"
)

var pageBuilderData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.page_builder_pages (id, created_at, updated_at, deleted_at, title, slug, category_id, version, locale_code)
VALUES (10, '2024-01-01 00:00:00', '2024-01-01 00:00:00', null, 'Test Page', '/test', 0, '2024-01-01-v01', 'International');

INSERT INTO public.container_headers (id, color)
VALUES (1, 'black');

INSERT INTO public.container_headings (id, add_top_space, add_bottom_space, anchor_id, heading, font_color, background_color, link, link_text, link_display_option)
VALUES (2, false, false, '', 'Test Heading', 'black', 'white', '', '', '');

INSERT INTO public.page_builder_containers (id, created_at, updated_at, deleted_at, page_id, page_version, page_model_name, model_name, model_id, display_order, shared, hidden, display_name, locale_code)
VALUES
(1, '2024-01-01 00:00:00', '2024-01-01 00:00:00', null, 10, '2024-01-01-v01', 'Page', 'Header', 1, 1, false, false, 'Header', 'International'),
(2, '2024-01-01 00:00:00', '2024-01-01 00:00:00', null, 10, '2024-01-01-v01', 'Page', 'Heading', 2, 2, false, false, 'Heading', 'International');
`, []string{"page_builder_pages", "page_builder_containers", "page_builder_templates", "container_headers", "container_headings"}))

// newPageBuilderEventBuilder creates an event request with the page identity used by the real editor.
func newPageBuilderEventBuilder(pageID, pageVersion, event string) *Builder {
	return NewMultipartBuilder().
		PageURL(fmt.Sprintf("/page-builder-pages/%s_%s_International", pageID, pageVersion)).
		EventFunc(event).
		Query("pageID", pageID).
		Query("pageVersion", pageVersion).
		Query("pageModelName", "Page").
		Query("locale", "International")
}

// resetInNumbersDemo removes prior demo state so each drawer scenario is independent.
func resetInNumbersDemo(t *testing.T) {
	t.Helper()
	if err := TestDB.Unscoped().
		Where("model_name = ?", "InNumbers").
		Delete(&pagebuilder.DemoContainer{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := TestDB.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&containers.InNumbers{}).Error; err != nil {
		t.Fatal(err)
	}
}

// TestPageBuilder 页面构建器集成测试
func TestPageBuilder(t *testing.T) {
	h := admin.TestHandler(TestDB, nil)
	dbr, _ := TestDB.DB()

	cases := []TestCase{
		{
			Name: "PageBuilder Page List",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				return httptest.NewRequest("GET", "/page-builder-pages?locale=international", http.NoBody)
			},
			ExpectPageBodyContainsInOrder: []string{"Test Page"},
		},
		{
			Name: "PageBuilder Page Detail (Editor)",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				return httptest.NewRequest("GET", "/page-builder-pages/10_2024-01-01-v01_International", http.NoBody)
			},
			ExpectPageBodyContainsInOrder: []string{"Header", "Heading"},
		},
		{
			Name: "PageBuilder Add Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.AddContainerEvent).
					Query("modelName", "Footer").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var cons []pagebuilder.Container
				TestDB.Order("display_order asc").Find(&cons)
				if len(cons) != 3 {
					t.Fatalf("expected 3 containers, got %d", len(cons))
				}
				if cons[2].ModelName != "Footer" {
					t.Fatalf("expected last container to be Footer, got %s", cons[2].ModelName)
				}
			},
		},
		{
			Name: "PageBuilder Delete Container Confirmation",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.DeleteContainerConfirmationEvent).
					Query("containerID", "1").
					Query("containerName", "Header").
					BuildEventFuncRequest()
				return req
			},
			ExpectPortalUpdate0ContainsInOrder: []string{"Header"},
		},
		{
			Name: "PageBuilder Delete Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.DeleteContainerEvent).
					Query("containerID", "1").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var count int64
				TestDB.Model(&pagebuilder.Container{}).Where("deleted_at IS NULL").Count(&count)
				if count != 1 {
					t.Fatalf("expected 1 container after delete, got %d", count)
				}
			},
		},
		{
			Name: "PageBuilder Move Container Up",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.MoveUpDownContainerEvent).
					Query("containerID", "2").
					Query("moveDirection", "up").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var cons []pagebuilder.Container
				TestDB.Order("display_order asc").Find(&cons)
				if len(cons) != 2 {
					t.Fatalf("expected 2 containers, got %d", len(cons))
				}
				if cons[0].ModelName != "Heading" {
					t.Fatalf("expected Heading first after move up, got %s", cons[0].ModelName)
				}
			},
		},
		{
			Name: "PageBuilder Move Container Down",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.MoveUpDownContainerEvent).
					Query("containerID", "1").
					Query("moveDirection", "down").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var cons []pagebuilder.Container
				TestDB.Order("display_order asc").Find(&cons)
				if cons[0].ModelName != "Heading" {
					t.Fatalf("expected Heading first after move down, got %s", cons[0].ModelName)
				}
			},
		},
		{
			Name: "PageBuilder Toggle Container Visibility",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.ToggleContainerVisibilityEvent).
					Query("containerID", "1").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var c pagebuilder.Container
				TestDB.First(&c, 1)
				if !c.Hidden {
					t.Fatal("expected container to be hidden")
				}
			},
		},
		{
			Name: "PageBuilder Rename Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.RenameContainerEvent).
					Query("containerID", "1").
					Query("displayName", "My Header").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var c pagebuilder.Container
				TestDB.First(&c, 1)
				if c.DisplayName != "My Header" {
					t.Fatalf("expected display name 'My Header', got '%s'", c.DisplayName)
				}
			},
		},
		{
			Name: "PageBuilder Edit Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.EditContainerEvent).
					Query("containerID", "1").
					BuildEventFuncRequest()
				return req
			},
			ExpectRunScriptContainsInOrder: []string{
				`.url("/headers")`,
				`.query("id", "1")`,
				`.query("overlay", "content")`,
			},
		},
		{
			Name: "PageBuilder Update Header Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := NewMultipartBuilder().
					PageURL("/headers").
					EventFunc(actions.Update).
					Query("id", "1").
					AddField("Color", "white").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var header containers.WebHeader
				TestDB.First(&header, 1)
				if header.Color != "white" {
					t.Fatalf("expected color 'white', got '%s'", header.Color)
				}
			},
		},
		{
			Name: "PageBuilder Update Heading Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := NewMultipartBuilder().
					PageURL("/headings").
					EventFunc(actions.Update).
					Query("id", "2").
					AddField("Heading", "Updated Title").
					AddField("FontColor", "blue").
					AddField("BackgroundColor", "grey").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var heading containers.Heading
				TestDB.First(&heading, 2)
				if heading.Heading != "Updated Title" {
					t.Fatalf("expected heading 'Updated Title', got '%s'", heading.Heading)
				}
				if heading.FontColor != "blue" {
					t.Fatalf("expected font color 'blue', got '%s'", heading.FontColor)
				}
			},
		},
		{
			Name: "PageBuilder Replicate Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.ReplicateContainerEvent).
					Query("containerID", "1").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var cons []pagebuilder.Container
				TestDB.Order("display_order asc").Find(&cons)
				if len(cons) != 3 {
					t.Fatalf("expected 3 containers after replicate, got %d", len(cons))
				}
				// 复制的容器应该紧跟在原容器后面
				if cons[1].ModelName != "Header" {
					t.Fatalf("expected replicated Header at position 2, got %s", cons[1].ModelName)
				}
			},
		},
		{
			Name: "PageBuilder Mark As Shared Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.MarkAsSharedContainerEvent).
					Query("containerID", "1").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
				var c pagebuilder.Container
				TestDB.First(&c, 1)
				if !c.Shared {
					t.Fatal("expected container to be marked as shared")
				}
			},
		},
		{
			Name: "PageBuilder Demo Edit Uses Persisted ID",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				resetInNumbersDemo(t)
				return httptest.NewRequest(http.MethodGet, "/demo-containers?locale=international", http.NoBody)
			},
			ExpectPageBodyNotContains: []string{`query("id", "0")`, "id=0"},
			ResponseMatch: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				if recorder.Code != http.StatusOK {
					t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
				}
				var demo pagebuilder.DemoContainer
				if err := TestDB.Where("model_name = ?", "InNumbers").
					Order("id DESC").
					First(&demo).Error; err != nil {
					t.Fatal(err)
				}
				if demo.ModelID == 0 {
					t.Fatal("InNumbers demo model ID is zero")
				}
				var model containers.InNumbers
				if err := TestDB.First(&model, demo.ModelID).Error; err != nil {
					t.Fatalf("InNumbers model %d does not exist: %v", demo.ModelID, err)
				}
			},
		},
		{
			Name: "PageBuilder Demo Drawer Opens Persisted Model",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				resetInNumbersDemo(t)
				recorder := httptest.NewRecorder()
				h.ServeHTTP(
					recorder,
					httptest.NewRequest(
						http.MethodGet,
						"/demo-containers?locale=international",
						http.NoBody,
					),
				)
				if recorder.Code != http.StatusOK {
					t.Fatalf(
						"initialize demo containers: status = %d, body = %s",
						recorder.Code,
						recorder.Body.String(),
					)
				}
				var demo pagebuilder.DemoContainer
				if err := TestDB.Where("model_name = ?", "InNumbers").
					Order("id DESC").
					First(&demo).Error; err != nil {
					t.Fatal(err)
				}
				return NewMultipartBuilder().
					PageURL("/in-numbers").
					EventFunc(actions.Edit).
					Query("id", strconv.FormatUint(uint64(demo.ModelID), 10)).
					Query("overlay", actions.Dialog).
					BuildEventFuncRequest()
			},
			EventResponseMatch: func(t *testing.T, response *TestEventResponse) {
				rendered := response.Body
				for _, portal := range response.UpdatePortals {
					rendered += portal.Body
				}
				if !strings.Contains(rendered, "Heading") {
					t.Fatalf("demo edit drawer did not render the InNumbers form: %#v", response)
				}
			},
		},
		{
			Name: "PageBuilder Delete From Confirmation Action",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				return newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.DeleteContainerConfirmationEvent).
					Query("containerID", "1").
					Query("containerName", "Header").
					BuildEventFuncRequest()
			},
			EventResponseMatch: func(t *testing.T, response *TestEventResponse) {
				if len(response.UpdatePortals) == 0 ||
					!strings.Contains(response.UpdatePortals[0].Body, pagebuilder.DeleteContainerEvent) {
					t.Fatalf("confirmation action is missing delete event: %#v", response.UpdatePortals)
				}
				deleteRequest := newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.DeleteContainerEvent).
					Query("containerID", "1").
					BuildEventFuncRequest()
				recorder := httptest.NewRecorder()
				h.ServeHTTP(recorder, deleteRequest)
				if recorder.Code != http.StatusOK {
					t.Fatalf("delete status = %d, body = %s", recorder.Code, recorder.Body.String())
				}
				var count int64
				TestDB.Model(&pagebuilder.Container{}).Where("id = ?", 1).Count(&count)
				if count != 0 {
					t.Fatalf("container still exists after confirmation action")
				}
			},
		},
		{
			Name: "PageBuilder Drag Sort",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				return newPageBuilderEventBuilder("10", "2024-01-01-v01", pagebuilder.MoveContainerEvent).
					AddField(
						"moveResult",
						`[{"container_id":"2","locale":"International"},{"container_id":"1","locale":"International"}]`,
					).
					BuildEventFuncRequest()
			},
			EventResponseMatch: func(t *testing.T, _ *TestEventResponse) {
				var relationships []pagebuilder.Container
				if err := TestDB.Where(
					"page_id = ? AND page_version = ? AND locale_code = ?",
					10,
					"2024-01-01-v01",
					"International",
				).Order("display_order ASC").Find(&relationships).Error; err != nil {
					t.Fatal(err)
				}
				if len(relationships) != 2 || relationships[0].ID != 2 {
					t.Fatalf("drag order = %#v", relationships)
				}
			},
		},
		{
			Name: "PageBuilder Add Shared Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				if err := pagebuilder.MarkContainerAsShared(TestDB, 1, "International"); err != nil {
					t.Fatal(err)
				}
				page := pagebuilder.Page{Title: "Shared Destination", Slug: "/shared-destination"}
				page.ID = 11
				page.Version.Version = "v1"
				page.LocaleCode = "International"
				if err := TestDB.Create(&page).Error; err != nil {
					t.Fatal(err)
				}
				return newPageBuilderEventBuilder("11", "v1", pagebuilder.AddSharedContainerEvent).
					Query("sourceContainerID", "1").
					BuildEventFuncRequest()
			},
			EventResponseMatch: func(t *testing.T, _ *TestEventResponse) {
				var relationship pagebuilder.Container
				if err := TestDB.Where(
					"page_id = ? AND page_version = ? AND locale_code = ?",
					11,
					"v1",
					"International",
				).First(&relationship).Error; err != nil {
					t.Fatal(err)
				}
				if !relationship.Shared || relationship.ModelID != 1 || relationship.ModelName != "Header" {
					t.Fatalf("shared relationship = %#v", relationship)
				}
			},
		},
		{
			Name: "PageBuilder Apply Template",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				template := pagebuilder.Template{Name: "Landing"}
				template.ID = 20
				template.LocaleCode = "International"
				if err := TestDB.Create(&template).Error; err != nil {
					t.Fatal(err)
				}
				source := pagebuilder.Container{
					PageID:        template.ID,
					PageModelName: "Template",
					ModelName:     "Heading",
					ModelID:       2,
					DisplayOrder:  1,
					DisplayName:   "Template Heading",
				}
				source.LocaleCode = "International"
				if err := TestDB.Create(&source).Error; err != nil {
					t.Fatal(err)
				}
				page := pagebuilder.Page{Title: "Template Destination", Slug: "/template-destination"}
				page.ID = 12
				page.Version.Version = "v1"
				page.LocaleCode = "International"
				if err := TestDB.Create(&page).Error; err != nil {
					t.Fatal(err)
				}
				return newPageBuilderEventBuilder("12", "v1", pagebuilder.ApplyTemplateEvent).
					Query("templateID", "20").
					BuildEventFuncRequest()
			},
			EventResponseMatch: func(t *testing.T, _ *TestEventResponse) {
				var relationship pagebuilder.Container
				if err := TestDB.Where(
					"page_id = ? AND page_version = ? AND locale_code = ? AND page_model_name = ?",
					12,
					"v1",
					"International",
					"Page",
				).First(&relationship).Error; err != nil {
					t.Fatal(err)
				}
				if relationship.ModelName != "Heading" || relationship.ModelID == 0 || relationship.ModelID == 2 {
					t.Fatalf("template relationship = %#v", relationship)
				}
			},
		},
		{
			Name: "PageBuilder Duplicate Version Copies Containers",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				return newPageBuilderEventBuilder("10", "2024-01-01-v01", publish.EventDuplicateVersion).
					BuildEventFuncRequest()
			},
			EventResponseMatch: func(t *testing.T, _ *TestEventResponse) {
				var duplicated pagebuilder.Page
				if err := TestDB.Where(
					"id = ? AND parent_version = ? AND locale_code = ?",
					10,
					"2024-01-01-v01",
					"International",
				).First(&duplicated).Error; err != nil {
					t.Fatal(err)
				}
				var relationships []pagebuilder.Container
				if err := TestDB.Where(
					"page_id = ? AND page_version = ? AND locale_code = ?",
					10,
					duplicated.Version.Version,
					"International",
				).Order("display_order ASC").Find(&relationships).Error; err != nil {
					t.Fatal(err)
				}
				if len(relationships) != 2 {
					t.Fatalf("duplicated relationships = %#v", relationships)
				}
				for _, relationship := range relationships {
					if relationship.ModelID == 0 {
						t.Fatalf("duplicated relationship has zero model ID: %#v", relationship)
					}
					if relationship.ModelName == "Header" && relationship.ModelID == 1 {
						t.Fatalf("private Header model was not cloned: %#v", relationship)
					}
					if relationship.ModelName == "Heading" && relationship.ModelID == 2 {
						t.Fatalf("private Heading model was not cloned: %#v", relationship)
					}
				}
			},
		},
		{
			Name: "PageBuilder Preview Renders SEO And Style",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				if err := TestDB.Model(&pagebuilder.Page{}).
					Where("id = ? AND version = ? AND locale_code = ?", 10, "2024-01-01-v01", "International").
					Update("seo", seo.Setting{
						EnabledCustomize: true,
						Title:            "PageBuilder SEO",
						Description:      "PageBuilder preview metadata",
					}).Error; err != nil {
					t.Fatal(err)
				}
				return httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf(
						"/page_builder/preview?pageID=10&pageVersion=%s&locale=International&pageModelName=Page",
						"2024-01-01-v01",
					),
					http.NoBody,
				)
			},
			ResponseMatch: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				if recorder.Code != http.StatusOK {
					t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
				}
				for _, want := range []string{
					"data-pagebuilder-style",
					"<title>PageBuilder SEO</title>",
					"PageBuilder preview metadata",
					"rel='canonical'",
				} {
					if !strings.Contains(recorder.Body.String(), want) {
						t.Errorf("preview missing %q: %s", want, recorder.Body.String())
					}
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunCase(t, c, h)
		})
	}
}
