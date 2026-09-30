package shop

// themeMessages demo 主题的文案（模板里 {{t "Key"}}；键经 ToCamel 匹配字段，"Status_open" → StatusOpen）。
type themeMessages struct {
	ShopName, NavProducts, NavCart, Country, NoProducts, SoldOut, NotAvailableHere, Quantity, AddToCart    string
	CartEmpty, ContinueShopping, Update, Checkout, BackToCart, Subtotal, Shipping, Tax, InclTax, Total     string
	Email, Name, Address, Address2, PostalCode, City, Phone, ShippingMethod, Free, PayNow, MarketSwitched  string
	OnlyLeft, Processing, ProcessingHint, PaymentExpired, OrderNumber, ShipTo, Shipments, Track            string
	StatusOpen, StatusCancelled, PaymentPaid, PaymentRefunded                                              string
	FulfillmentUnfulfilled, FulfillmentPartial, FulfillmentFulfilled                                       string
	PageNotFound, SomethingWrong, BackHome                                                                 string
	LegalTerms, LegalPrivacy, LegalReturns, LegalShipping, LegalImprint                                    string
	OrderPlacedSubject, OrderShippedSubject, OrderCancelledSubject                                         string
	OrderPlacedIntro, OrderShippedIntro, OrderCancelledIntro, ViewOrder                                    string
	ErrRequired, ErrInvalid, ErrLimit, ErrNotSellable, ErrUnsupportedCountry, ErrUnsellable, ErrOutOfStock string
	ErrTooMany, ErrPaymentUnavailable, ErrCartEmpty                                                        string
}

// themeEN 英文。
var themeEN = &themeMessages{
	ShopName: "r0vx Store", NavProducts: "Products", NavCart: "Cart", Country: "Ship to", NoProducts: "No products available in your country yet.",
	SoldOut: "Sold out", NotAvailableHere: "Not available in your country", Quantity: "Quantity", AddToCart: "Add to cart",
	CartEmpty: "Your cart is empty.", ContinueShopping: "Continue shopping", Update: "Update", Checkout: "Checkout", BackToCart: "Back to cart",
	Subtotal: "Subtotal", Shipping: "Shipping", Tax: "Tax", InclTax: "incl. VAT", Total: "Total",
	Email: "Email", Name: "Full name", Address: "Address", Address2: "Apartment, suite (optional)", PostalCode: "Postal code", City: "City",
	Phone: "Phone (optional)", ShippingMethod: "Shipping method", Free: "Free", PayNow: "Continue to payment",
	MarketSwitched: "Your shipping country uses a different currency. Prices were updated — please choose a shipping method again.",
	OnlyLeft:       "only available:", Processing: "Confirming your payment…", ProcessingHint: "This page refreshes automatically.",
	PaymentExpired: "Your payment session has expired", OrderNumber: "Order", ShipTo: "Shipping address", Shipments: "Shipments", Track: "Track",
	StatusOpen: "Open", StatusCancelled: "Cancelled", PaymentPaid: "Paid", PaymentRefunded: "Refunded",
	FulfillmentUnfulfilled: "Not shipped yet", FulfillmentPartial: "Partially shipped", FulfillmentFulfilled: "Shipped",
	PageNotFound: "Page not found", SomethingWrong: "Something went wrong", BackHome: "Back to the home page",
	LegalTerms: "Terms", LegalPrivacy: "Privacy", LegalReturns: "Returns", LegalShipping: "Shipping", LegalImprint: "Imprint",
	OrderPlacedSubject: "Order confirmed", OrderShippedSubject: "Your order has shipped", OrderCancelledSubject: "Order cancelled",
	OrderPlacedIntro: "Thanks for your order! We'll let you know when it ships.", OrderShippedIntro: "Good news — your order is on its way.",
	OrderCancelledIntro: "Your order was cancelled and the full amount refunded.", ViewOrder: "View order",
	ErrRequired: "Required", ErrInvalid: "Please check this field", ErrLimit: "Maximum quantity reached", ErrNotSellable: "This item is not available in your country",
	ErrUnsupportedCountry: "We don't ship to this country yet", ErrUnsellable: "Some items aren't available in your country — please remove them",
	ErrOutOfStock: "Some items are out of stock", ErrTooMany: "Too many attempts — please wait a minute", ErrPaymentUnavailable: "Payment is temporarily unavailable, please try again",
	ErrCartEmpty: "Your cart is empty",
}

// themeDE 德文。
var themeDE = &themeMessages{
	ShopName: "r0vx Store", NavProducts: "Produkte", NavCart: "Warenkorb", Country: "Lieferland", NoProducts: "In deinem Land sind noch keine Produkte verfügbar.",
	SoldOut: "Ausverkauft", NotAvailableHere: "In deinem Land nicht erhältlich", Quantity: "Menge", AddToCart: "In den Warenkorb",
	CartEmpty: "Dein Warenkorb ist leer.", ContinueShopping: "Weiter einkaufen", Update: "Aktualisieren", Checkout: "Zur Kasse", BackToCart: "Zurück zum Warenkorb",
	Subtotal: "Zwischensumme", Shipping: "Versand", Tax: "Steuer", InclTax: "inkl. MwSt.", Total: "Gesamt",
	Email: "E-Mail", Name: "Vollständiger Name", Address: "Adresse", Address2: "Wohnung, Etage (optional)", PostalCode: "PLZ", City: "Stadt",
	Phone: "Telefon (optional)", ShippingMethod: "Versandart", Free: "Kostenlos", PayNow: "Weiter zur Zahlung",
	MarketSwitched: "Für dein Lieferland gilt eine andere Währung. Die Preise wurden aktualisiert – bitte wähle die Versandart erneut.",
	OnlyLeft:       "nur noch verfügbar:", Processing: "Zahlung wird bestätigt …", ProcessingHint: "Diese Seite aktualisiert sich automatisch.",
	PaymentExpired: "Deine Zahlungssitzung ist abgelaufen", OrderNumber: "Bestellung", ShipTo: "Lieferadresse", Shipments: "Sendungen", Track: "Verfolgen",
	StatusOpen: "Offen", StatusCancelled: "Storniert", PaymentPaid: "Bezahlt", PaymentRefunded: "Erstattet",
	FulfillmentUnfulfilled: "Noch nicht versendet", FulfillmentPartial: "Teilweise versendet", FulfillmentFulfilled: "Versendet",
	PageNotFound: "Seite nicht gefunden", SomethingWrong: "Etwas ist schiefgelaufen", BackHome: "Zur Startseite",
	LegalTerms: "AGB", LegalPrivacy: "Datenschutz", LegalReturns: "Widerruf", LegalShipping: "Versand", LegalImprint: "Impressum",
	OrderPlacedSubject: "Bestellung bestätigt", OrderShippedSubject: "Deine Bestellung wurde versendet", OrderCancelledSubject: "Bestellung storniert",
	OrderPlacedIntro: "Danke für deine Bestellung! Wir melden uns, sobald sie versendet wird.", OrderShippedIntro: "Gute Nachricht – deine Bestellung ist unterwegs.",
	OrderCancelledIntro: "Deine Bestellung wurde storniert und der volle Betrag erstattet.", ViewOrder: "Bestellung ansehen",
	ErrRequired: "Pflichtfeld", ErrInvalid: "Bitte prüfe dieses Feld", ErrLimit: "Höchstmenge erreicht", ErrNotSellable: "Dieser Artikel ist in deinem Land nicht erhältlich",
	ErrUnsupportedCountry: "In dieses Land liefern wir noch nicht", ErrUnsellable: "Einige Artikel sind in deinem Land nicht erhältlich – bitte entferne sie",
	ErrOutOfStock: "Einige Artikel sind nicht vorrätig", ErrTooMany: "Zu viele Versuche – bitte warte eine Minute", ErrPaymentUnavailable: "Zahlung vorübergehend nicht möglich, bitte später erneut versuchen",
	ErrCartEmpty: "Dein Warenkorb ist leer",
}
