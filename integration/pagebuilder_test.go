package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example/admin"
	"example/admin/pagebuilder/containers"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets/actions"
	. "github.com/r0vx/web/multipartestutils"
	"github.com/theplant/gofixtures"
)

var pageBuilderData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.page_builder_pages (id, created_at, updated_at, deleted_at, title, slug, category_id, version, locale_code)
VALUES (10, '2024-01-01 00:00:00', '2024-01-01 00:00:00', null, 'Test Page', '/test', 0, '2024-01-01-v01', '');

INSERT INTO public.container_headers (id, color)
VALUES (1, 'black');

INSERT INTO public.container_headings (id, add_top_space, add_bottom_space, anchor_id, heading, font_color, background_color, link, link_text, link_display_option)
VALUES (2, false, false, '', 'Test Heading', 'black', 'white', '', '', '');

INSERT INTO public.page_builder_containers (id, created_at, updated_at, deleted_at, page_id, page_version, page_model_name, model_name, model_id, display_order, shared, hidden, display_name, locale_code)
VALUES
(1, '2024-01-01 00:00:00', '2024-01-01 00:00:00', null, 10, '2024-01-01-v01', 'Page', 'Header', 1, 1, false, false, 'Header', ''),
(2, '2024-01-01 00:00:00', '2024-01-01 00:00:00', null, 10, '2024-01-01-v01', 'Page', 'Heading', 2, 2, false, false, 'Heading', '');
`, []string{"page_builder_pages", "page_builder_containers", "container_headers", "container_headings"}))

// TestPageBuilder 页面构建器集成测试
func TestPageBuilder(t *testing.T) {
	h := admin.TestHandler(TestDB, nil)
	dbr, _ := TestDB.DB()

	cases := []TestCase{
		{
			Name: "PageBuilder Page List",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				return httptest.NewRequest("GET", "/page-builder-pages", http.NoBody)
			},
			ExpectPageBodyContainsInOrder: []string{"Test Page"},
		},
		{
			Name: "PageBuilder Page Detail (Editor)",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				return httptest.NewRequest("GET", "/page-builder-pages/10_2024-01-01-v01_", http.NoBody)
			},
			ExpectPageBodyContainsInOrder: []string{"Header", "Heading"},
		},
		{
			Name: "PageBuilder Add Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.AddContainerEvent).
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.DeleteContainerConfirmationEvent).
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.DeleteContainerEvent).
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.MoveUpDownContainerEvent).
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.MoveUpDownContainerEvent).
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.ToggleContainerVisibilityEvent).
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.RenameContainerEvent).
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.EditContainerEvent).
					Query("containerID", "1").
					BuildEventFuncRequest()
				return req
			},
			ExpectPortalUpdate0ContainsInOrder: []string{"Color"},
		},
		{
			Name: "PageBuilder Update Header Container",
			ReqFunc: func() *http.Request {
				pageBuilderData.TruncatePut(dbr)
				req := NewMultipartBuilder().
					PageURL("/headers/1").
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
					PageURL("/headings/2").
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.ReplicateContainerEvent).
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
				req := NewMultipartBuilder().
					PageURL("/page-builder-pages/10_2024-01-01-v01_").
					EventFunc(pagebuilder.MarkAsSharedContainerEvent).
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
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunCase(t, c, h)
		})
	}
}
