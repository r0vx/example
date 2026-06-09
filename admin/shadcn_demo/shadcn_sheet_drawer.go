package shadcn_demo

import (
	"fmt"
	"time"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	. "github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// ShadcnSheetDrawerDemo 虚拟模型
type ShadcnSheetDrawerDemo struct{}

// configSheetDrawer 注册 Sheet Drawer demo
func configSheetDrawer(b *presets.Builder) {
	m := b.Model(&ShadcnSheetDrawerDemo{}).
		Label("Sheet Drawer").
		URIName("shadcn-sheet-drawer")
	m.Listing().SearchFunc(emptySearchFunc[ShadcnSheetDrawerDemo]())
	m.Editing().Only()

	m.RegisterEventFunc("showRemoteSheet", showRemoteSheet)
	m.RegisterEventFunc("sheetUpdateParentAndClose", sheetUpdateParentAndClose)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		injectDemoCSS(ctx)
		r.PageTitle = "Shadcn Sheet Drawer"
		r.Body = shadcnSheetDrawerBody(ctx)
		return
	})
}

// shadcnSheetDrawerBody Sheet 抽屉演示
func shadcnSheetDrawerBody(ctx *web.EventContext) h.HTMLComponent {
	return h.Div(
		h.H1("Shadcn Sheet Drawer").Style("margin-bottom: 24px;"),

		// Basic Sheet with close button
		h.Div(
			h.H2("Sheet with Close Button"),
			web.Scope(
				Sheet(
					SheetTrigger(
						Button(h.Text("Open Sheet")),
					),
					SheetContent(
						SheetHeader(
							SheetTitle(h.Text("Sheet Title")),
							SheetDescription(h.Text("This is a sheet that can be closed.")),
						),
						h.Div(
							h.P(h.Text("Sheet content goes here. You can add any content.")),
						).Class("py-4"),
						SheetFooter(
							SheetClose(
								Button(h.Text("Close")).Variant(ButtonVariantOutline),
							),
						),
					).Side(SheetSideRight),
				),
			).VSlot("{ locals }"),
		).Class("demo-section"),

		// Sheet from different sides
		h.Div(
			h.H2("Sheet from Different Sides"),
			h.Div(
				Sheet(
					SheetTrigger(
						Button(h.Text("Left")).Variant(ButtonVariantOutline),
					),
					SheetContent(
						SheetHeader(
							SheetTitle(h.Text("Left Sheet")),
							SheetDescription(h.Text("Slides in from the left")),
						),
						h.Div(h.Text("Left side content")).Class("py-4"),
					).Side(SheetSideLeft),
				),
				Sheet(
					SheetTrigger(
						Button(h.Text("Right")).Variant(ButtonVariantOutline),
					),
					SheetContent(
						SheetHeader(
							SheetTitle(h.Text("Right Sheet")),
							SheetDescription(h.Text("Slides in from the right")),
						),
						h.Div(h.Text("Right side content")).Class("py-4"),
					).Side(SheetSideRight),
				),
				Sheet(
					SheetTrigger(
						Button(h.Text("Top")).Variant(ButtonVariantOutline),
					),
					SheetContent(
						SheetHeader(
							SheetTitle(h.Text("Top Sheet")),
							SheetDescription(h.Text("Slides in from the top")),
						),
						h.Div(h.Text("Top side content")).Class("py-4"),
					).Side(SheetSideTop),
				),
				Sheet(
					SheetTrigger(
						Button(h.Text("Bottom")).Variant(ButtonVariantOutline),
					),
					SheetContent(
						SheetHeader(
							SheetTitle(h.Text("Bottom Sheet")),
							SheetDescription(h.Text("Slides in from the bottom")),
						),
						h.Div(h.Text("Bottom side content")).Class("py-4"),
					).Side(SheetSideBottom),
				),
			).Class("flex gap-2"),
		).Class("demo-section"),

		// Sheet loaded from remote
		h.Div(
			h.H2("Sheet Loaded from Remote"),
			h.P(h.Text("Click button to load sheet content from server")).Class("text-muted-foreground mb-4"),

			Button(h.Text("Show Remote Sheet")).On("click", web.POST().EventFunc("showRemoteSheet").Go()),

			web.Portal().Name("sheetDrawerUpdateContent"),
			web.Portal().Name("sheetDrawer"),
		).Class("demo-section"),
	).Style("max-width: 600px; margin: 0 auto;")
}

// showRemoteSheet 显示远程加载的 Sheet
func showRemoteSheet(ctx *web.EventContext) (er web.EventResponse, err error) {
	er.UpdatePortals = append(er.UpdatePortals,
		&web.PortalUpdate{
			Name: "sheetDrawer",
			Body: web.Scope(
				Sheet(
					SheetContent(
						SheetHeader(
							SheetTitle(h.Text("Remote Sheet")),
							SheetDescription(h.Text("Content loaded from server")),
						),
						h.Div(
							web.Portal(
								sheetInputField(""),
							).Name("SheetInputPortal"),
						).Class("py-4"),
						SheetFooter(
							Button(h.Text("Update Parent and Close")).
								On("click", web.POST().EventFunc("sheetUpdateParentAndClose").Go()),
						),
					).Side(SheetSideRight).Class("sm:max-w-lg"),
				).Attr("v-model:open", "locals.sheetOpen"),
			).VSlot("{ locals, form }").Init("{ sheetOpen: true }"),
		},
	)
	return
}

// sheetInputField 创建输入字段
func sheetInputField(value string, fieldErrors ...string) h.HTMLComponent {
	input := Input().
		Placeholder("Enter at least 10 characters").
		Attr("v-model", "form.SheetInput")

	if len(fieldErrors) > 0 {
		input = input.ErrorMessages(fieldErrors...)
	}

	return h.Div(
		Label(h.Text("Input")).Class("mb-2"),
		input,
	)
}

// sheetUpdateParentAndClose 更新父级并关闭
func sheetUpdateParentAndClose(ctx *web.EventContext) (er web.EventResponse, err error) {
	inputValue := ctx.R.FormValue("SheetInput")

	if len(inputValue) < 10 {
		er.UpdatePortals = append(er.UpdatePortals, &web.PortalUpdate{
			Name: "SheetInputPortal",
			Body: sheetInputField(inputValue, "Input must be at least 10 characters"),
		})
		return
	}

	er.UpdatePortals = append(er.UpdatePortals, &web.PortalUpdate{
		Name: "sheetDrawerUpdateContent",
		Body: Card(
			CardContent(
				h.Div(
					Badge(h.Text("Updated")).Variant(BadgeVariantSecondary),
					h.Text(fmt.Sprintf(" Content updated at %s", time.Now().Format("15:04:05"))),
				),
			).Class("pt-6"),
		).Class("mt-4"),
	})

	er.RunScript = "locals.sheetOpen = false"
	return
}
