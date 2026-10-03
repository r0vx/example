// Package shop example 的电商演示（r0vx/commerce）：后台「Shop」菜单 + 独立端口的公开商店。
//
// 公开商店用 demo 主题渲染：首页是 pagebuilder 页面（含精选商品容器），商品、购物车、结账、订单是应用页；
// 支付用假支付（本地模拟付款页），交易邮件打印到标准输出。
package shop

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/r0vx/admin/activity"
	"github.com/r0vx/admin/l10n"
	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/seo"
	"github.com/r0vx/commerce"
	"github.com/r0vx/commerce/payment"
	"github.com/r0vx/commerce/payment/antom"
	"github.com/r0vx/commerce/payment/fake"
	"github.com/r0vx/commerce/payment/stripe"
	"github.com/r0vx/web"
	"github.com/r0vx/x/i18n"
	"github.com/r0vx/x/mail"
	"github.com/theplant/osenv"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

var (
	// storeURL 公开商店的根地址（邮件链接、付款返回地址）
	storeURL = osenv.Get("STORE_URL", "public store base URL", "http://127.0.0.1:9510")
	// tokenSecret 订单令牌密钥（≥ 32 字节）；不设时每次启动随机生成（重启后旧订单链接失效，仅限开发）
	tokenSecret = os.Getenv("SHOP_ORDER_TOKEN_SECRET")
)

// Shop 装配好的商店：后台已注册进 presets，Public 是公开站 handler，Run 跑后台任务。
type Shop struct {
	*commerce.Builder
	i18n *i18n.Builder
	pb   *pagebuilder.Builder
}

// Configure 装配商店：迁移表、注册后台与精选商品容器、给主题接全站数据，表空时建演示数据。
// demoPassword 非空时顺带建三种角色的演示账号（见 DemoAccounts）。
func Configure(db *gorm.DB, b *presets.Builder, pb *pagebuilder.Builder, l *l10n.Builder, ab *activity.Builder, sb *seo.Builder, demoPassword string) *Shop {
	secret := []byte(tokenSecret)
	if len(secret) < 32 {
		log.Printf("shop: 未设置 SHOP_ORDER_TOKEN_SECRET（≥ 32 字节），本次启动随机生成（重启后旧订单链接失效）")
		secret = make([]byte, 32)
		_, _ = rand.Read(secret)
	}
	ib := i18n.New().SupportLanguages(language.English, language.German).
		RegisterForModule(language.English, pagebuilder.I18nThemeKeyPrefix+"demo", themeEN).
		RegisterForModule(language.German, pagebuilder.I18nThemeKeyPrefix+"demo", themeDE)

	shop := commerce.New(db).
		Payment(paymentProvider()).
		Mailer(mail.Log(os.Stdout)).
		OrderTokenSecret(secret).
		OrderNumberPrefix("R-").
		Storefront(pb.Renderer()).
		I18n(ib).
		PublicURL(storeURL).
		DefaultCountry("US").
		L10n(l).
		Activity(ab).
		SEO(sb).
		StripeTestMode(true)
	if err := shop.Migrate(); err != nil {
		log.Fatalf("shop: 迁移: %v", err)
	}
	if err := shop.Install(b); err != nil {
		log.Fatalf("shop: 注册后台: %v", err)
	}
	shop.ProductGrid(pb)
	pb.SiteFunc(func(ctx *web.EventContext) (any, error) { return shop.SiteData(ctx.R) })
	if err := seed(db); err != nil {
		log.Fatalf("shop: 演示数据: %v", err)
	}
	if err := seedDemoUsers(db, demoPassword); err != nil {
		log.Fatalf("shop: 演示账号: %v", err)
	}
	if err := shop.Validate(); err != nil {
		log.Fatalf("shop: 配置校验: %v", err)
	}
	return &Shop{Builder: shop, i18n: ib, pb: pb}
}

// paymentProvider 选支付服务商（演示只允许沙箱 / 测试模式）：
//   - 设了 ANTOM_CLIENT_ID：Antom（国际支付宝）沙箱，另需 ANTOM_MERCHANT_PRIVATE_KEY、ANTOM_PUBLIC_KEY（Antom Dashboard 里的
//     不带头 base64 或 PEM），可选 ANTOM_GATEWAY、ANTOM_SETTLEMENT_CURRENCY；Antom 回调到 STORE_URL/shop/webhooks/antom（须公网可达）。
//   - 设了 STRIPE_SECRET_KEY：Stripe 测试模式，Webhook 密钥来自 STRIPE_WEBHOOK_SECRET，
//     本地用 `stripe listen --forward-to 127.0.0.1:9510/shop/webhooks/stripe` 转发。
//   - 都没设：带模拟付款页的假支付。
func paymentProvider() payment.Provider {
	if id := os.Getenv("ANTOM_CLIENT_ID"); id != "" {
		if !strings.HasPrefix(id, "SANDBOX_") {
			log.Fatal("shop: 演示只允许 Antom 沙箱（client-id 以 SANDBOX_ 开头）")
		}
		var opts []antom.Option
		if g := os.Getenv("ANTOM_GATEWAY"); g != "" {
			opts = append(opts, antom.WithBaseURL(g))
		}
		if c := os.Getenv("ANTOM_SETTLEMENT_CURRENCY"); c != "" {
			opts = append(opts, antom.WithSettlementCurrency(c))
		}
		p, err := antom.New(id, os.Getenv("ANTOM_MERCHANT_PRIVATE_KEY"), os.Getenv("ANTOM_PUBLIC_KEY"), storeURL+"/shop/webhooks/antom", opts...)
		if err != nil {
			log.Fatalf("shop: %v", err)
		}
		log.Printf("shop: 支付用 Antom 沙箱")
		return p
	}
	sk := os.Getenv("STRIPE_SECRET_KEY")
	if sk == "" {
		return fake.New(storeURL + "/shop/fake-pay/")
	}
	whsec := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if whsec == "" {
		log.Fatal("shop: 设了 STRIPE_SECRET_KEY 还需要 STRIPE_WEBHOOK_SECRET")
	}
	p := stripe.New(sk, whsec)
	if err := p.Validate(false); err != nil { // 演示只允许测试模式
		log.Fatalf("shop: %v", err)
	}
	log.Printf("shop: 支付用 Stripe 测试模式")
	return p
}

// Public 公开商店 handler：主题资源、商品页、/shop/*，其余交给 pagebuilder 页面（首页、法务页）。
// 整站包路径前缀语言（/de/...）与商店中间件；主题资源不经过它们。
func (s *Shop) Public() http.Handler {
	site := http.NewServeMux()
	s.Mount(site, "/shop")
	site.Handle("GET /products", s.ProductsHandler())
	site.Handle("GET /products/{handle}", s.ProductHandler())
	site.Handle("/", s.pb.PageHandler())

	root := http.NewServeMux()
	root.Handle("/theme-assets/", s.pb.ThemeAssetHandler("/theme-assets/"))
	root.Handle("/", s.i18n.EnsureLanguageFromPath(s.Wrap(site)))
	return root
}

// Start 在 addr 上起公开商店，并跑商店后台任务（事件分发、过期结账清理），直到 ctx 取消。
func (s *Shop) Start(ctx context.Context, addr string) {
	go func() {
		if err := s.Run(ctx); err != nil {
			log.Printf("shop: 后台任务退出: %v", err)
		}
	}()
	go func() {
		log.Printf("shop: 公开商店 %s（%s）", storeURL, addr)
		if err := http.ListenAndServe(addr, s.Public()); err != nil { //nolint:gosec // 演示用
			log.Printf("shop: 公开商店退出: %v", err)
		}
	}()
}
