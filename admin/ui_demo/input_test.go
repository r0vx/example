package ui_demo

import (
	"net/http/httptest"
	"testing"

	"example/models"

	"github.com/r0vx/web"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestCityTextForRequest_memoPerRequest 验证按请求 memo：
// 同一 ctx 只查库一次（读缓存，不受新增影响）；换 ctx 重查（实时）。
func TestCityTextForRequest_memoPerRequest(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.City{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&models.City{ID: 1, Name: "Tokyo"})

	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	if got := cityTextForRequest(db, ctx)["1"]; got != "Tokyo" {
		t.Fatalf("首次应查库得 Tokyo, got %q", got)
	}

	// 同一 ctx：库里新增也读缓存、不重查（证明 per-request memo，各行复用一次查询）
	db.Create(&models.City{ID: 2, Name: "Osaka"})
	if _, ok := cityTextForRequest(db, ctx)["2"]; ok {
		t.Fatal("同请求应读缓存，不该看到新增的 Osaka")
	}

	// 新 ctx：重查、看到 Osaka（实时，无进程内快照陈旧）
	ctx2 := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	if got := cityTextForRequest(db, ctx2)["2"]; got != "Osaka" {
		t.Fatalf("新请求应重查看到 Osaka, got %q", got)
	}
}
