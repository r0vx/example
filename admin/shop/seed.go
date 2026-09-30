package shop

import (
	"context"
	"fmt"

	"example/admin/pagebuilder/containers"

	"github.com/r0vx/admin/l10n"
	"github.com/r0vx/admin/pagebuilder"
	"github.com/r0vx/admin/publish"
	"github.com/r0vx/commerce/catalog"
	"github.com/r0vx/commerce/inventory"
	"github.com/r0vx/commerce/market"
	"github.com/r0vx/commerce/storefront"
	"gorm.io/gorm"
)

// seed 市场表为空时建演示数据：US / EU / UK 三个市场与配送方式、四件商品（一件只在 US / UK 卖、一件草稿）、
// 英德两版首页（含精选商品）与法务页占位。全部经领域服务写入。
func seed(db *gorm.DB) error {
	var n int64
	if err := db.Model(&market.Market{}).Count(&n).Error; err != nil || n > 0 {
		return err
	}
	ctx := context.Background()
	return db.Transaction(func(tx *gorm.DB) error {
		s := &seeder{ctx: ctx, db: tx}
		s.markets()
		s.products()
		s.pages()
		return s.err
	})
}

// seeder 演示数据：第一个错误之后的步骤都跳过。
type seeder struct {
	ctx        context.Context
	db         *gorm.DB
	err        error
	us, eu, uk *market.Market
	featured   []uint
}

// do 记下第一个错误。
func (s *seeder) do(err error) {
	if s.err == nil && err != nil {
		s.err = err
	}
}

// markets 三个市场：US（美元、税另计）、EU（欧元、含 19% 增值税）、UK（英镑、含 20%）。
func (s *seeder) markets() {
	s.us = &market.Market{Code: "US", Name: "United States", Currency: "USD", DefaultLocale: "en", Locales: []string{"en"}, Enabled: true, Position: 0, Countries: []string{"US"}}
	s.eu = &market.Market{Code: "EU", Name: "Europe", Currency: "EUR", DefaultLocale: "de", Locales: []string{"de", "en"}, PricesIncludeTax: true, TaxRateBps: 1900, Enabled: true, Position: 1, Countries: []string{"DE", "AT", "FR", "NL"}}
	s.uk = &market.Market{Code: "UK", Name: "United Kingdom", Currency: "GBP", DefaultLocale: "en", Locales: []string{"en"}, PricesIncludeTax: true, TaxRateBps: 2000, Enabled: true, Position: 2, Countries: []string{"GB"}}
	for _, m := range []*market.Market{s.us, s.eu, s.uk} {
		s.do(market.Save(s.ctx, s.db, m))
	}
	free := func(v int64) *int64 { return &v }
	for _, sm := range []struct {
		m          *market.Market
		code       string
		price      int64
		freeOver   *int64
		en, de, di string
	}{
		{s.us, "standard", 599, free(7500), "Standard (3–5 business days)", "", ""},
		{s.eu, "dhl", 490, free(5000), "DHL Parcel", "DHL Paket", "2–4 Werktage"},
		{s.uk, "royal-mail", 395, nil, "Royal Mail Tracked 48", "", ""},
	} {
		if s.err != nil {
			return
		}
		x := &market.ShippingMethod{MarketID: sm.m.ID, Code: sm.code, Price: sm.price, FreeOverAmount: sm.freeOver, Enabled: true}
		s.do(market.SaveShippingMethod(s.ctx, s.db, x))
		s.do(market.SaveShippingContent(s.ctx, s.db, market.ShippingMethodContent{ID: x.ID, Locale: l10n.Locale{LocaleCode: "en"}, Name: sm.en}))
		if sm.de != "" {
			s.do(market.SaveShippingContent(s.ctx, s.db, market.ShippingMethodContent{ID: x.ID, Locale: l10n.Locale{LocaleCode: "de"}, Name: sm.de, Description: sm.di}))
		}
	}
}

// text 一种语言的商品文案。
type text struct{ title, handle, desc string }

// product 建商品（英文为默认语言）并加德文译文，返回商品。
func (s *seeder) product(en, de text) *catalog.Product {
	if s.err != nil {
		return nil
	}
	p, err := catalog.CreateProduct(s.ctx, s.db, "en", catalog.ProductContent{Title: en.title, Handle: &en.handle, Description: en.desc})
	s.do(err)
	if s.err == nil && de.title != "" {
		s.do(catalog.SaveContent(s.ctx, s.db, catalog.ProductContent{ID: p.ID, Locale: l10n.Locale{LocaleCode: "de"}, Title: de.title, Handle: &de.handle, Description: de.desc}))
	}
	return p
}

// prices 给变体在各市场标价（nil 表示该市场不卖）。
func (s *seeder) prices(variantID uint, us, eu, uk *int64, euCompare *int64) {
	for _, p := range []struct {
		m       *market.Market
		amount  *int64
		compare *int64
	}{{s.us, us, nil}, {s.eu, eu, euCompare}, {s.uk, uk, nil}} {
		if p.amount != nil && s.err == nil {
			s.do(catalog.SetPrice(s.ctx, s.db, variantID, p.m.ID, *p.amount, p.compare))
		}
	}
}

// products 四件商品：T 恤（颜色 × 尺码，一个组合缺货）、马克杯（不跟踪库存）、帆布袋（欧盟不卖）、连帽衫（草稿）。
func (s *seeder) products() {
	amt := func(v int64) *int64 { return &v }

	tee := s.product(
		text{"Classic T-Shirt", "classic-tee", "<p>Heavyweight organic cotton, relaxed fit. Printed in Berlin.</p>"},
		text{"Klassisches T-Shirt", "klassisches-t-shirt", "<p>Schwere Bio-Baumwolle, lockere Passform. Gedruckt in Berlin.</p>"})
	if s.err != nil {
		return
	}
	color, err := catalog.AddOption(s.ctx, s.db, tee.ID)
	s.do(err)
	size, err := catalog.AddOption(s.ctx, s.db, tee.ID)
	s.do(err)
	if s.err != nil {
		return
	}
	en := catalog.ProductContent{OptionNames: map[uint]string{color.ID: "Color", size.ID: "Size"}, OptionValueLabels: map[uint]string{}}
	de := catalog.ProductContent{OptionNames: map[uint]string{color.ID: "Farbe", size.ID: "Größe"}, OptionValueLabels: map[uint]string{}}
	var colors, sizes []uint
	for _, c := range [][2]string{{"Black", "Schwarz"}, {"White", "Weiß"}} {
		v, err := catalog.AddOptionValue(s.ctx, s.db, color.ID)
		s.do(err)
		if s.err != nil {
			return
		}
		colors = append(colors, v.ID)
		en.OptionValueLabels[v.ID], de.OptionValueLabels[v.ID] = c[0], c[1]
	}
	for _, z := range []string{"S", "M", "L"} {
		v, err := catalog.AddOptionValue(s.ctx, s.db, size.ID)
		s.do(err)
		if s.err != nil {
			return
		}
		sizes = append(sizes, v.ID)
		en.OptionValueLabels[v.ID], de.OptionValueLabels[v.ID] = z, z
	}
	s.saveOptionNames(tee.ID, "en", en)
	s.saveOptionNames(tee.ID, "de", de)
	i := 0
	for ci, c := range colors {
		for zi, z := range sizes {
			sku := fmt.Sprintf("TEE-%s-%s", []string{"BLK", "WHT"}[ci], []string{"S", "M", "L"}[zi])
			v, err := catalog.CreateVariant(s.ctx, s.db, catalog.VariantInput{ProductID: tee.ID, SKU: sku, OptionValueIDs: []uint{c, z}, TrackInventory: true, Position: i})
			s.do(err)
			if s.err != nil {
				return
			}
			s.prices(v.ID, amt(2900), amt(2499), amt(2200), amt(2999))
			if sku != "TEE-WHT-L" { // 白色 L 码缺货
				s.do(inventory.Adjust(s.ctx, s.db, v.ID, 25))
			}
			i++
		}
	}
	s.do(catalog.SetStatus(s.ctx, s.db, tee.ID, catalog.StatusActive))

	mug := s.product(
		text{"Ceramic Mug", "ceramic-mug", "<p>350 ml stoneware mug, dishwasher safe.</p>"},
		text{"Keramiktasse", "keramiktasse", "<p>Steinzeugtasse, 350 ml, spülmaschinenfest.</p>"})
	if s.err != nil {
		return
	}
	mv, err := catalog.CreateVariant(s.ctx, s.db, catalog.VariantInput{ProductID: mug.ID, SKU: "MUG-350"})
	s.do(err)
	if s.err != nil {
		return
	}
	s.prices(mv.ID, amt(1200), amt(1190), amt(1000), nil)
	s.do(catalog.SetStatus(s.ctx, s.db, mug.ID, catalog.StatusActive))

	tote := s.product(text{"Canvas Tote", "canvas-tote", "<p>Sturdy 12 oz canvas. Currently shipping to the US and UK only.</p>"}, text{})
	if s.err != nil {
		return
	}
	tv, err := catalog.CreateVariant(s.ctx, s.db, catalog.VariantInput{ProductID: tote.ID, SKU: "TOTE"})
	s.do(err)
	if s.err != nil {
		return
	}
	s.prices(tv.ID, amt(1800), nil, amt(1500), nil)
	s.do(catalog.SetStatus(s.ctx, s.db, tote.ID, catalog.StatusActive))

	hoodie := s.product(text{"Hoodie (coming soon)", "hoodie", "<p>Coming soon.</p>"}, text{})
	if s.err != nil {
		return
	}
	hv, err := catalog.CreateVariant(s.ctx, s.db, catalog.VariantInput{ProductID: hoodie.ID, SKU: "HOODIE"})
	s.do(err)
	if s.err == nil {
		s.prices(hv.ID, amt(5900), amt(5490), amt(4900), nil)
	}
	s.featured = []uint{tee.ID, mug.ID, tote.ID, hoodie.ID}
}

// saveOptionNames 把规格名、规格值名写进某语言文案（保留该语言已有的标题等）。
func (s *seeder) saveOptionNames(productID uint, locale string, names catalog.ProductContent) {
	if s.err != nil {
		return
	}
	var c catalog.ProductContent
	s.do(s.db.Where("id = ? AND locale_code = ?", productID, locale).First(&c).Error)
	if s.err != nil {
		return
	}
	c.OptionNames, c.OptionValueLabels = names.OptionNames, names.OptionValueLabels
	s.do(catalog.SaveContent(s.ctx, s.db, c))
}

// pageVersion 演示页面的版本号。
const pageVersion = "2026-09-30-v01"

// page 建一个已上线的页面（英德各一版，同一 ID），每版挂 build 返回的容器。
func (s *seeder) page(slug string, titles [2]string, build func(locale string) []pagebuilder.Container) {
	if s.err != nil {
		return
	}
	var id uint
	for i, locale := range []string{"en", "de"} {
		p := &pagebuilder.Page{Title: titles[i], Slug: slug}
		p.ID = id
		p.Version.Version = pageVersion
		p.Status.Status = publish.StatusOnline
		p.LocaleCode = locale
		s.do(s.db.Create(p).Error)
		if s.err != nil {
			return
		}
		id = p.ID
		cs := build(locale)
		for j := range cs {
			cs[j].PageID, cs[j].PageVersion, cs[j].PageModelName, cs[j].DisplayOrder = p.ID, pageVersion, "Page", float64(j+1)
			cs[j].LocaleCode = locale
		}
		if len(cs) > 0 {
			s.do(s.db.Create(&cs).Error)
		}
	}
}

// heading 建一个标题容器。
func (s *seeder) heading(title, body string) pagebuilder.Container {
	h := &containers.Heading{Heading: title, Text: body}
	s.do(s.db.Create(h).Error)
	return pagebuilder.Container{ModelName: "Heading", ModelID: h.ID, DisplayName: "Heading"}
}

// pages 首页（欢迎语 + 精选商品）与五个法务页占位（spec §8.0 #10，跨境欧盟必备）。
func (s *seeder) pages() {
	s.page("/", [2]string{"Home", "Startseite"}, func(locale string) []pagebuilder.Container {
		title, body, featured := "Made to last", "Small-batch basics, shipped to the US, EU and UK.", "Featured"
		if locale == "de" {
			title, body, featured = "Gemacht für lange", "Basics in kleinen Auflagen – Versand in die USA, EU und nach UK.", "Empfohlen"
		}
		grid := &storefront.ProductGrid{Title: featured, ProductIDs: s.featured}
		s.do(s.db.Create(grid).Error)
		return []pagebuilder.Container{s.heading(title, body), {ModelName: "ProductGrid", ModelID: grid.ID, DisplayName: "ProductGrid"}}
	})
	for _, lp := range []struct{ slug, en, de string }{
		{"/legal/terms", "Terms & Conditions", "Allgemeine Geschäftsbedingungen"},
		{"/legal/privacy", "Privacy Policy", "Datenschutzerklärung"},
		{"/legal/returns", "Returns & Refunds", "Widerrufsbelehrung"},
		{"/legal/shipping", "Shipping", "Versand"},
		{"/legal/imprint", "Imprint", "Impressum"},
	} {
		s.page(lp.slug, [2]string{lp.en, lp.de}, func(locale string) []pagebuilder.Container {
			title, body := lp.en, "Placeholder — replace with your legal text before going live."
			if locale == "de" {
				title, body = lp.de, "Platzhalter – vor dem Livegang durch den Rechtstext ersetzen."
			}
			return []pagebuilder.Container{s.heading(title, body)}
		})
	}
}
