package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example/admin"

	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/seo"
)

// TestPageBuilderThemePreview 示例应用装了 demo 主题后，公开预览整页走主题：
// 挂了 .View() 的容器用主题模板、页面壳用主题 layout，SEO 标题只出现一次且优先
func TestPageBuilderThemePreview(t *testing.T) {
	h := admin.TestHandler(TestDB, nil)
	dbr, _ := TestDB.DB()
	pageBuilderData.TruncatePut(dbr)
	if err := TestDB.Model(&pagebuilder.Page{}).
		Where("id = ? AND version = ? AND locale_code = ?", 10, "2024-01-01-v01", "International").
		Update("seo", seo.Setting{EnabledCustomize: true, Title: "Themed SEO", Description: "themed preview"}).Error; err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/page_builder/preview?pageID=10&pageVersion=2024-01-01-v01&locale=International&pageModelName=Page", http.NoBody))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{
		`href="/theme-assets/theme.css"`, // 主题 layout 接管了页面壳
		`demo-header demo-header--black`, // Header 容器走主题模板
		`class="demo-heading"`,           // Heading 容器走主题模板
		"Test Heading",                   // 容器数据
		"<title>Themed SEO</title>",      // SEO 标题优先
		"themed preview",                 // SEO 其余 meta 经 .Head 注入
		"data-pagebuilder-style",         // PageStyle 经 .Head 注入
	} {
		if !strings.Contains(body, want) {
			t.Errorf("预览缺 %q:\n%s", want, body)
		}
	}
	if n := strings.Count(body, "<title>"); n != 1 {
		t.Errorf("应只有 1 个 <title>，得到 %d 个", n)
	}
}
