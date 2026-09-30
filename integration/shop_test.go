package integration_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"example/admin"
	"example/models"

	"github.com/r0vx/commerce/catalog"
	"github.com/r0vx/commerce/events"
	"github.com/r0vx/commerce/market"
	"github.com/r0vx/commerce/order"
	"github.com/r0vx/x/perm"
	"gorm.io/gorm"
)

// shopClient 公开商店的一个顾客（带 cookie，不自动跟跳转）
type shopClient struct {
	t    *testing.T
	base string
	c    *http.Client
	hdr  map[string]string
}

// newShopClient 新顾客；hdr 是每个请求都带的头（如 CF-IPCountry）
func newShopClient(t *testing.T, base string, hdr map[string]string) *shopClient {
	jar, _ := cookiejar.New(nil)
	return &shopClient{t: t, base: base, hdr: hdr, c: &http.Client{Jar: jar,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}

// do 发请求，返回状态码、正文、跳转地址
func (s *shopClient) do(method, path string, form url.Values) (int, string, string) {
	s.t.Helper()
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, _ := http.NewRequest(method, s.base+path, body)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for k, v := range s.hdr {
		req.Header.Set(k, v)
	}
	res, err := s.c.Do(req)
	if err != nil {
		s.t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b), res.Header.Get("Location")
}

// csrfRe 页面里的 CSRF 令牌
var csrfRe = regexp.MustCompile(`name="csrf" value="([0-9a-f]+)"`)

// post 带 CSRF 令牌提交（令牌取自购物车页）
func (s *shopClient) post(path string, form url.Values) (int, string, string) {
	s.t.Helper()
	_, page, _ := s.do(http.MethodGet, "/en/shop/cart", nil)
	m := csrfRe.FindStringSubmatch(page)
	if m == nil {
		s.t.Fatalf("没有 CSRF 令牌")
	}
	form.Set("csrf", m[1])
	return s.do(http.MethodPost, path, form)
}

// buy 加购第一个 T 恤变体 2 件，按 country 的地址结账、假付款，返回订单与订单页路径（经付款返回页跳转得到）
func (s *shopClient) buy(db *gorm.DB, lang, country string) (*order.Order, string) {
	s.t.Helper()
	var v catalog.Variant
	if err := db.Where("sku = ?", "TEE-BLK-M").First(&v).Error; err != nil {
		s.t.Fatal(err)
	}
	if code, body, _ := s.post("/"+lang+"/shop/cart/lines", url.Values{"variant_id": {fmt.Sprint(v.ID)}, "quantity": {"2"}}); code != http.StatusSeeOther {
		s.t.Fatalf("加购 %d\n%s", code, body)
	}
	m, err := market.ForCountry(context.Background(), db, country)
	if err != nil {
		s.t.Fatal(err)
	}
	var ship market.ShippingMethod
	if err := db.Where("market_id = ?", m.ID).First(&ship).Error; err != nil {
		s.t.Fatal(err)
	}
	f := url.Values{"email": {"buyer@example.com"}, "name": {"Buyer"}, "line1": {"Street 1"}, "city": {"City"},
		"postal_code": {"10115"}, "country": {country}, "shipping_method_id": {fmt.Sprint(ship.ID)}}
	code, body, loc := s.post("/"+lang+"/shop/checkout", f)
	if code != http.StatusSeeOther || !strings.Contains(loc, "/shop/fake-pay/") {
		s.t.Fatalf("结账应跳假付款页: %d %s\n%s", code, loc, body)
	}
	u, _ := url.Parse(loc)
	code, _, back := s.post(u.Path, url.Values{"action": {"pay"}})
	if code != http.StatusSeeOther || !strings.Contains(back, "/"+lang+"/shop/checkout/return") {
		s.t.Fatalf("付款后应回下单语言的返回页: %d %s", code, back)
	}
	ru, _ := url.Parse(back)
	code, _, orderPath := s.do(http.MethodGet, ru.RequestURI(), nil)
	if code != http.StatusSeeOther || !strings.HasPrefix(orderPath, "/"+lang+"/shop/orders/") {
		s.t.Fatalf("返回页应跳订单页: %d %s", code, orderPath)
	}
	var o order.Order
	if err := db.Order("id DESC").First(&o).Error; err != nil {
		s.t.Fatal(err)
	}
	return &o, orderPath
}

// TestShopEndToEnd 验收（spec §12.2）：德国顾客（德语、欧元、含 VAT）与美国顾客（英语、美元、税另计）都能下单；
// 首页精选商品按市场显示价格、当前市场不卖的不出现；后台发货后订单已发货、交易邮件投递成功。
func TestShopEndToEnd(t *testing.T) {
	h, c := admin.TestHandlerComplex(TestDB, nil, false)
	srv := httptest.NewServer(c.Shop().Public())
	defer srv.Close()

	de := newShopClient(t, srv.URL, map[string]string{"CF-IPCountry": "DE"})
	_, home, _ := de.do(http.MethodGet, "/de/", nil)
	for _, want := range []string{"Klassisches T-Shirt", "24,99 €", "Keramiktasse"} {
		if !strings.Contains(home, want) {
			t.Errorf("德国首页缺 %q", want)
		}
	}
	if strings.Contains(home, "Canvas Tote") || strings.Contains(home, "Hoodie") {
		t.Error("德国首页不该出现欧盟不卖的商品与草稿")
	}
	us := newShopClient(t, srv.URL, map[string]string{"CF-IPCountry": "US"})
	_, home, _ = us.do(http.MethodGet, "/en/", nil)
	for _, want := range []string{"Classic T-Shirt", "$29.00", "Canvas Tote"} {
		if !strings.Contains(home, want) {
			t.Errorf("美国首页缺 %q", want)
		}
	}

	deOrder, dePath := de.buy(TestDB, "de", "DE")
	if deOrder.Currency != "EUR" || deOrder.Locale != "de" || deOrder.Total != 2*2499+490 {
		t.Errorf("德国订单: %s %s %d", deOrder.Currency, deOrder.Locale, deOrder.Total)
	}
	_, page, _ := de.do(http.MethodGet, dePath, nil)
	if !strings.Contains(page, "inkl. MwSt. 19%") || !strings.Contains(page, "54,88 €") {
		t.Errorf("德国订单页应含税说明与欧元合计:\n%s", page)
	}
	usOrder, _ := us.buy(TestDB, "en", "US")
	if usOrder.Currency != "USD" || usOrder.TaxTotal != 0 || usOrder.Total != 2*2900+599 {
		t.Errorf("美国订单: %s tax=%d total=%d", usOrder.Currency, usOrder.TaxTotal, usOrder.Total)
	}

	// 后台发货（管理员）
	ls, err := order.Lines(context.Background(), TestDB, deOrder.ID)
	if err != nil {
		t.Fatal(err)
	}
	form := url.Values{"id": {fmt.Sprint(deOrder.ID)}, "action": {"Ship"}, fmt.Sprintf("Qty_%d", ls[0].ID): {"2"}, "Carrier": {"DHL"}, "TrackingNumber": {"JJD1"}}
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/commerce-orders/%d?__execute_event__=presets_DoAction", deOrder.ID), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(httptest.NewRecorder(), req)
	var shipped order.Order
	TestDB.First(&shipped, deOrder.ID)
	if shipped.FulfillmentStatus != order.Fulfilled {
		t.Fatalf("后台发货后应已发货: %s", shipped.FulfillmentStatus)
	}

	// 事件分发：下单、发货邮件投递成功（邮件打到标准输出）
	d := c.Shop().Bus().Dispatcher(TestDB)
	for {
		n, err := d.RunOnce(context.Background(), time.Now())
		if err != nil {
			t.Fatal(err)
		}
		if n == 0 {
			break
		}
	}
	var failed int64
	TestDB.Model(&events.Delivery{}).Where("status <> ?", events.StatusSucceeded).Count(&failed)
	if failed != 0 {
		t.Errorf("还有 %d 个投递没成功", failed)
	}
}

// TestShopRoles 三种角色（spec §9）：运营管商品、客服管订单、译者只能在德语下改文案
func TestShopRoles(t *testing.T) {
	get := func(role, path string) string {
		u := &models.User{Model: gorm.Model{ID: 901}, Roles: []perm.Role{{Name: role}}}
		h, _ := admin.TestHandlerComplex(TestDB, u, false)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec.Body.String()
	}
	denied := func(body string) bool { return strings.Contains(body, "Access Denied") }
	var tee catalog.Product
	if err := TestDB.Order("id").First(&tee).Error; err != nil {
		t.Fatal(err)
	}
	detail := fmt.Sprintf("/commerce-products/%d?locale=en", tee.ID)

	if denied(get(models.RoleShopOperator, "/commerce-products")) || !denied(get(models.RoleShopOperator, "/commerce-orders")) {
		t.Error("运营：能进商品、不能进订单")
	}
	if denied(get(models.RoleShopSupport, "/commerce-orders")) || !denied(get(models.RoleShopSupport, "/commerce-products")) {
		t.Error("客服：能进订单、不能进商品")
	}
	if body := get(models.RoleShopOperator, detail); !strings.Contains(body, `"AddOption"`) {
		t.Error("运营在英语下应能改结构")
	}
	body := get(models.RoleShopTranslator, detail)
	if denied(body) || strings.Contains(body, `"AddOption"`) || !strings.Contains(body, "Edit content") {
		t.Error("译者：被限定在德语，只能编辑文案、看不到结构操作")
	}
}
