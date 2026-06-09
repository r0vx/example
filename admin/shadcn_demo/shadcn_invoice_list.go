package shadcn_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// DemoInvoice 发票数据
type DemoInvoice struct {
	ID        string
	Date      string
	Recipient string
	Email     string
	Status    string
	Progress  int
	Amount    string
}

// demoInvoices 模拟数据
var demoInvoices = []DemoInvoice{
	{"INV001", "2024-01-13", "Alice Chen", "alice@company.com", "Paid", 100, "$2,300.00"},
	{"INV002", "2024-01-13", "Bob Wang", "bob@tech.com", "Unpaid", 50, "$780.00"},
	{"INV003", "2024-01-13", "Carol Liu", "carol@startup.com", "Paid", 100, "$567.00"},
	{"INV004", "2024-01-13", "David Zhang", "david@corp.com", "Unpaid", 25, "$1,500.00"},
	{"INV005", "2024-01-13", "Eva Li", "eva@enterprise.com", "Draft", 10, "$900.00"},
	{"INV006", "2024-07-03", "Frank Wu", "frank@digital.com", "Paid", 100, "$4,750.00"},
}

// invoiceStatusBadge 根据状态返回对应的 Badge
func invoiceStatusBadge(status string) h.HTMLComponent {
	switch status {
	case "Paid":
		return Badge(h.Text(status)).Class("bg-emerald-500 hover:bg-emerald-600")
	case "Unpaid":
		return Badge(h.Text(status)).Class("bg-rose-500 hover:bg-rose-600")
	case "Draft":
		return Badge(h.Text(status)).Variant(BadgeVariantSecondary)
	default:
		return Badge(h.Text(status)).Variant(BadgeVariantOutline)
	}
}

// invoiceActionButtons 操作按钮组
func invoiceActionButtons() h.HTMLComponent {
	return h.Div(
		Button(h.Text("Edit")).Variant(ButtonVariantOutline).Size(ButtonSizeSm),
		DropdownMenu(
			DropdownMenuTrigger(
				Button(h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg>`)).Variant(ButtonVariantGhost).Size(ButtonSizeSm).Class("px-2"),
			),
			DropdownMenuContent(
				DropdownMenuItem(h.Text("View")),
				DropdownMenuItem(h.Text("Download")),
				DropdownMenuSeparator(),
				DropdownMenuItem(h.Text("Delete")).Class("text-destructive"),
			).SideOffset(4),
		),
	).Class("flex items-center gap-1")
}

// ShadcnInvoiceListDemo 虚拟模型
type ShadcnInvoiceListDemo struct{}

// configInvoiceList 注册 Invoice List demo
func configInvoiceList(b *presets.Builder) {
	m := b.Model(&ShadcnInvoiceListDemo{}).
		Label("Invoice List").
		URIName("shadcn-invoice-list")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnInvoiceListDemo]())
	m.Editing().Only()

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Invoice List"
		r.Body = shadcnInvoiceListBody(ctx)
		return
	})
}

// shadcnInvoiceListBody 发票列表演示
func shadcnInvoiceListBody(ctx *web.EventContext) h.HTMLComponent {
	// 表格行
	var tableRows []h.HTMLComponent
	for _, inv := range demoInvoices {
		tableRows = append(tableRows, TableRow(
			TableCell(Checkbox()).Class("w-10"),
			TableCell(h.A(h.Text(inv.ID)).Href("#").Class("text-primary hover:underline")),
			TableCell(h.Text(inv.Date)),
			TableCell(
				h.Div(
					h.Div(h.Text(inv.Recipient)).Class("font-medium"),
					h.Div(h.Text(inv.Email)).Class("text-sm text-muted-foreground"),
				),
			),
			TableCell(invoiceStatusBadge(inv.Status)),
			TableCell(
				h.Div(
					Progress().ModelValue(inv.Progress).Class("h-2 w-16"),
					h.Span(inv.Status).Class("ml-2 text-xs text-muted-foreground"),
				).Class("flex items-center"),
			),
			TableCell(h.Text(inv.Amount)).Class("text-right"),
			TableCell(invoiceActionButtons()),
		))
	}

	return h.Div(
		h.H1("Invoice List Demo").Style("margin-bottom: 24px;"),

		// 工具栏
		h.Div(
			h.Div(
				Input().Placeholder("Search invoices...").Type("search").Class("w-64"),
				Select(
					SelectTrigger(
						SelectValue().Placeholder("Status"),
					).Class("w-32"),
					SelectContent(
						SelectItem(h.Text("All")).Value("all"),
						SelectItem(h.Text("Paid")).Value("paid"),
						SelectItem(h.Text("Unpaid")).Value("unpaid"),
						SelectItem(h.Text("Draft")).Value("draft"),
					),
				),
			).Class("flex items-center gap-2"),
			h.Div(
				DropdownMenu(
					DropdownMenuTrigger(
						Button(h.Text("Export"), h.RawHTML(`<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="ml-2"><polyline points="6 9 12 15 18 9"/></svg>`)).Variant(ButtonVariantOutline).Size(ButtonSizeSm),
					),
					DropdownMenuContent(
						DropdownMenuItem(h.Text("Export CSV")),
						DropdownMenuItem(h.Text("Export Excel")),
						DropdownMenuItem(h.Text("Export PDF")),
					),
				),
				Button(h.Text("New Invoice")),
			).Class("flex items-center gap-2"),
		).Class("flex items-center justify-between mb-4"),

		// 表格
		Card(
			Table(
				TableHeader(
					TableRow(
						TableHead(Checkbox()).Class("w-10"),
						TableHead(h.Text("Invoice")),
						TableHead(h.Text("Date")),
						TableHead(h.Text("Recipient")),
						TableHead(h.Text("Status")),
						TableHead(h.Text("Progress")),
						TableHead(h.Text("Amount")).Class("text-right"),
						TableHead(h.Text("Actions")),
					),
				),
				TableBody(tableRows...),
			),
		),

		// 分页
		h.Div(
			h.Span("1 - 6 of 6 items").Class("text-sm text-muted-foreground"),
			Pagination(
				PaginationContent(
					PaginationItem(PaginationPrevious()),
					PaginationItem(Button(h.Text("1")).Size(ButtonSizeSm)),
					PaginationItem(Button(h.Text("2")).Variant(ButtonVariantOutline).Size(ButtonSizeSm)),
					PaginationItem(PaginationNext()),
				),
			).Total(100).ItemsPerPage(10),
		).Class("flex items-center justify-between mt-4"),
	).Style("max-width: 1000px; margin: 0 auto;")
}
