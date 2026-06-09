package pages

import (
	"log"

	"github.com/r0vx/admin/media"
	"github.com/r0vx/admin/media/base"
	"github.com/r0vx/admin/media/media_library"

	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/cropper"
	"github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// Settings 设置页面
func Settings(db *gorm.DB) web.PageFunc {
	return func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.PageTitle = "Settings"

		r.Body = h.Div(
			h.Div(
				h.Div(
					h.Div(
						h.H1("Example of use QMediaBox in any page").Class("text-xl font-medium pt-4 pl-2"),
						media.QMediaBox(db).
							FieldName("test").
							Value(&media_library.MediaBox{}).
							Config(&media_library.MediaBoxConfig{
								AllowType: "image",
								Sizes: map[string]*base.Size{
									"thumb": {
										Width:  400,
										Height: 300,
									},
									"main": {
										Width:  800,
										Height: 500,
									},
								},
							}),
					).Class("md:w-1/2"),
				).Class("flex flex-wrap"),

				h.Div(
					h.Div(
						shadcn.Textarea().
							Label("Body").
							Value(`Could you do an actual logo instead of a font I cant pay you? Can we try some other colors maybe? I cant pay you. You might wanna give it another shot, so make it pop and this is just a 5 minutes job the target audience makes and families aged zero and up will royalties in the company do instead of cash.

Jazz it up a little I was wondering if my cat could be placed over the logo in the flyer I have printed it out, but the animated gif is not moving I have printed it out, but the animated gif is not moving make it original. Can you make it stand out more? Make it original.`).
							Rows(6),
					).Class("w-full"),
				).Class("flex flex-wrap mt-4"),

				h.Div(
					cropper.Cropper().
						Src("https://agontuk.github.io/assets/images/berserk.jpg").
						ModelValue(cropper.Value{X: 1141, Y: 540, Width: 713, Height: 466}).
						AspectRatio(713, 466).
						Attr("@input", web.Plaid().
							FieldValue("CropperEvent", web.Var("JSON.stringify($event)")).EventFunc(LogInfoEvent).Go()),
				).Class("flex flex-wrap mt-4"),
			).Class("container mx-auto px-4"),
		)
		return
	}
}

const LogInfoEvent = "logInfo"

// LogInfo 打印日志
func LogInfo(ctx *web.EventContext) (r web.EventResponse, err error) {
	log.Println("CropperEvent", ctx.R.FormValue("CropperEvent"))
	return
}
