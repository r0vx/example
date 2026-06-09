package shadcn_demo

import (
	"fmt"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// Invoice 发票数据结构
type Invoice struct {
	ID     string
	Status string
	Method string
	Amount float64
}

var invoices = []Invoice{
	{ID: "INV001", Status: "Paid", Method: "Credit Card", Amount: 250.00},
	{ID: "INV002", Status: "Pending", Method: "PayPal", Amount: 150.00},
	{ID: "INV003", Status: "Paid", Method: "Bank Transfer", Amount: 350.00},
	{ID: "INV004", Status: "Failed", Method: "Credit Card", Amount: 450.00},
	{ID: "INV005", Status: "Paid", Method: "PayPal", Amount: 550.00},
}

// ShadcnTableDemo 虚拟模型
type ShadcnTableDemo struct{}

// configTable 注册 Table demo
func configTable(b *presets.Builder) {
	m := b.Model(&ShadcnTableDemo{}).
		Label("Table").
		URIName("shadcn-table")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnTableDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("viewAll", viewAll)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Table"
		r.Body = shadcnTableBody(ctx)
		return
	})
}

// shadcnTableBody 表格组件演示
func shadcnTableBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("Shadcn Table").Style("margin-bottom: 24px;"),

		prettyFormAsJSON(ctx),

		// Basic Table
		h.Div(
			h.H2("Basic Table"),
			Table(
				TableHeader(
					TableRow(
						TableHead(h.Text("Invoice")),
						TableHead(h.Text("Status")),
						TableHead(h.Text("Method")),
						TableHead(h.Text("Amount")).Class("text-right"),
					),
				),
				TableBody(
					h.RawHTML(buildTableRows(invoices)),
				),
				TableFooter(
					TableRow(
						h.Td(h.Text("Total")).Attr("colspan", "3"),
						TableCell(h.Text(fmt.Sprintf("$%.2f", calculateTotal(invoices)))).Class("text-right font-bold"),
					),
				),
			),
		).Class("demo-section"),

		// Table with Actions
		h.Div(
			h.H2("Table with Actions"),
			Table(
				TableHeader(
					TableRow(
						TableHead(h.Text("Invoice")),
						TableHead(h.Text("Status")),
						TableHead(h.Text("Amount")).Class("text-right"),
						TableHead(h.Text("Actions")).Class("text-right"),
					),
				),
				TableBody(
					h.RawHTML(buildTableRowsWithActions(invoices)),
				),
			),
		).Class("demo-section"),

		// Card with Table
		h.Div(
			h.H2("Card with Table"),
			Card(
				CardHeader(
					CardTitle(h.Text("Recent Invoices")),
					CardDescription(h.Text("A list of your recent invoices.")),
				),
				CardContent(
					Table(
						TableHeader(
							TableRow(
								TableHead(h.Text("Invoice")),
								TableHead(h.Text("Status")),
								TableHead(h.Text("Amount")).Class("text-right"),
							),
						),
						TableBody(
							h.RawHTML(buildSimpleTableRows(invoices[:3])),
						),
					),
				),
				CardFooter(
					Button(h.Text("View All")).Variant(ButtonVariantOutline).On("click", web.POST().EventFunc("viewAll").Go()),
				),
			).Class("max-w-lg"),
		).Class("demo-section"),
	).Style("max-width: 800px; margin: 0 auto;")
}

// buildTableRows 构建表格行
func buildTableRows(invoices []Invoice) string {
	var rows string
	for _, inv := range invoices {
		statusBadge := Badge(h.Text(inv.Status))
		switch inv.Status {
		case "Paid":
			statusBadge = statusBadge.Variant(BadgeVariantDefault)
		case "Pending":
			statusBadge = statusBadge.Variant(BadgeVariantSecondary)
		case "Failed":
			statusBadge = statusBadge.Variant(BadgeVariantDestructive)
		}

		row := TableRow(
			TableCell(h.Text(inv.ID)).Class("font-medium"),
			TableCell(statusBadge),
			TableCell(h.Text(inv.Method)),
			TableCell(h.Text(fmt.Sprintf("$%.2f", inv.Amount))).Class("text-right"),
		)
		rows += h.MustString(row, nil)
	}
	return rows
}

// buildTableRowsWithActions 构建带操作的表格行
func buildTableRowsWithActions(invoices []Invoice) string {
	var rows string
	for _, inv := range invoices {
		statusBadge := Badge(h.Text(inv.Status))
		switch inv.Status {
		case "Paid":
			statusBadge = statusBadge.Variant(BadgeVariantDefault)
		case "Pending":
			statusBadge = statusBadge.Variant(BadgeVariantSecondary)
		case "Failed":
			statusBadge = statusBadge.Variant(BadgeVariantDestructive)
		}

		row := TableRow(
			TableCell(h.Text(inv.ID)).Class("font-medium"),
			TableCell(statusBadge),
			TableCell(h.Text(fmt.Sprintf("$%.2f", inv.Amount))).Class("text-right"),
			TableCell(
				DropdownMenu(
					DropdownMenuTrigger(
						Button(h.Text("...")).Variant(ButtonVariantGhost).Size(ButtonSizeSm),
					),
					DropdownMenuContent(
						DropdownMenuItem(h.Text("View")),
						DropdownMenuItem(h.Text("Edit")),
						DropdownMenuSeparator(),
						DropdownMenuItem(h.Text("Delete")).Class("text-red-600"),
					),
				),
			).Class("text-right"),
		)
		rows += h.MustString(row, nil)
	}
	return rows
}

// buildSimpleTableRows 构建简单表格行
func buildSimpleTableRows(invoices []Invoice) string {
	var rows string
	for _, inv := range invoices {
		row := TableRow(
			TableCell(h.Text(inv.ID)).Class("font-medium"),
			TableCell(h.Text(inv.Status)),
			TableCell(h.Text(fmt.Sprintf("$%.2f", inv.Amount))).Class("text-right"),
		)
		rows += h.MustString(row, nil)
	}
	return rows
}

// calculateTotal 计算总金额
func calculateTotal(invoices []Invoice) float64 {
	var total float64
	for _, inv := range invoices {
		total += inv.Amount
	}
	return total
}

// viewAll 查看全部
func viewAll(ctx *web.EventContext) (r web.EventResponse, err error) {
	r.Reload = true
	return
}
