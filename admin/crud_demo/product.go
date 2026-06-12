package crud_demo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/r0vx/admin/media"
	"github.com/r0vx/admin/media/base"
	"github.com/r0vx/admin/publish"

	"example/models"

	"github.com/r0vx/admin/media/media_library"
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"

	"github.com/r0vx/admin/worker"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// ConfigProduct 配置商品管理模块
func ConfigProduct(b *presets.Builder, _ *gorm.DB, wb *worker.Builder, publisher *publish.Builder) *presets.ModelBuilder {
	p := b.Model(&models.Product{}).Use(publisher)
	eb := p.Editing("StatusBar", "ScheduleBar", "Code", "Name", "Price", "Image")
	listing := p.Listing("Code", "Name", "Price", "Image").SearchColumns("Code", "Name").
		SelectableColumns(true)

	// listing.ActionsAsMenu(true)

	noParametersJob := wb.ActionJob(
		"No parameters",
		p,
		func(ctx context.Context, job worker.GoJobInterface) error {
			for i := 1; i <= 10; i++ {
				select {
				case <-ctx.Done():
					job.AddLog("job aborted")
					return nil
				default:
					job.SetProgress(uint(i * 10))
					time.Sleep(time.Second)
				}
			}
			job.SetProgressText(`<a href="https://r0vx-test.s3.ap-northeast-1.amazonaws.com/system/media_libraries/37/file.@r0vx_preview.png">Please download this file</a>`)
			return nil
		},
	).Description("This test demo is used to show that an no parameter job can be executed").
		DisplayLog(true)

	type JobResource struct {
		Name string
	}
	parametersBoxJob := wb.ActionJob(
		"Parameter input box",
		p,
		func(ctx context.Context, job worker.GoJobInterface) error {
			for i := 1; i <= 10; i++ {
				select {
				case <-ctx.Done():
					job.AddLog("job aborted")
					return nil
				default:
					job.SetProgress(uint(i * 10))
					time.Sleep(time.Second)
				}
			}
			job.SetProgressText(`<a href="https://r0vx-test.s3.ap-northeast-1.amazonaws.com/system/media_libraries/37/file.@r0vx_preview.png">Please download this file</a>`)
			return nil
		},
	).Description("This test demo is used to show an input box when there are parameters").
		Params(&JobResource{})

	if rmb := parametersBoxJob.GetParamsModelBuilder(); rmb != nil {
		rmb.Editing("Name").Field("Name").Label("Name").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			if obj == nil {
				obj = &JobResource{}
			}
			jobRes, ok := obj.(*JobResource)
			if !ok {
				jobRes = &JobResource{}
			}
			return shadcn.Input().
				Label(field.Label).
				Attr(web.VField(field.Name, jobRes.Name)...)
		})
	}

	displayLogJob := wb.ActionJob(
		"Display log",
		p,
		func(ctx context.Context, job worker.GoJobInterface) error {
			for i := 1; i <= 10; i++ {
				select {
				case <-ctx.Done():
					job.AddLog("job aborted")
					return nil
				default:
					job.SetProgress(uint(i * 10))
					job.AddLog(fmt.Sprintf("%v", i))
					time.Sleep(time.Second)
				}
			}
			job.SetProgressText(`<a href="https://r0vx-test.s3.ap-northeast-1.amazonaws.com/system/media_libraries/37/file.@r0vx_preview.png">Please download this file</a>`)
			return nil
		},
	).Description("This test demo is used to show the log section of this job").
		Params(&struct{ Name string }{}).
		DisplayLog(true).
		ProgressingInterval(4000)

	getArgsJob := wb.ActionJob(
		"Get Args",
		p,
		func(ctx context.Context, job worker.GoJobInterface) error {
			jobInfo, err := job.GetJobInfo()
			if err != nil {
				return err
			}

			job.AddLog(fmt.Sprintf("Action Params Name is  %#+v", jobInfo.Argument.(*struct{ Name string }).Name))
			job.AddLog(fmt.Sprintf("Origina Context AuthInfo is  %#+v", jobInfo.Context["AuthInfo"]))
			job.AddLog(fmt.Sprintf("Origina Context URL is  %#+v", jobInfo.Context["URL"]))

			for i := 1; i <= 10; i++ {
				select {
				case <-ctx.Done():
					return nil
				default:
					job.SetProgress(uint(i * 10))
					time.Sleep(time.Second)
				}
			}
			job.SetProgressText(`<a href="https://r0vx-test.s3.ap-northeast-1.amazonaws.com/system/media_libraries/37/file.@r0vx_preview.png">Please download this file</a>`)
			return nil
		},
	).Description("This test demo is used to show how to get the action's arguments and original page context").
		Params(&struct{ Name string }{}).
		DisplayLog(true).
		ContextHandler(func(ctx *web.EventContext) map[string]any {
			auth, err := ctx.R.Cookie("auth")
			if err == nil {
				return map[string]any{"AuthInfo": auth.Value}
			}
			return nil
		})

	// 走框架统一的紧凑图标渲染：无 Icon 时自动取 Label 首字母（N/P/D/G），hover 显示完整 label。
	// OnClickFunc 延迟到 render 期才求 job.URL()（job 的 ModelBuilder 配置期未挂载，eager 调用会 nil panic）。
	listing.BulkAction("Action Job - No parameters").
		Label("No parameters").
		OnClickFunc(func(ctx *web.EventContext) string { return noParametersJob.URL() })

	listing.BulkAction("Action Job - Parameter input box").
		Label("Parameter input box").
		OnClickFunc(func(ctx *web.EventContext) string { return parametersBoxJob.URL() })

	listing.BulkAction("Action Job - Display log").
		Label("Display log").
		OnClickFunc(func(ctx *web.EventContext) string { return displayLogJob.URL() })

	listing.BulkAction("Action Job - Get Args").
		Label("Get Args").
		OnClickFunc(func(ctx *web.EventContext) string { return getArgsJob.URL() })

	eb.ValidateFunc(func(obj any, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.Product)
		if u.Code == "" {
			err.FieldError("Code", "Code is required")
		}
		if u.Name == "" {
			err.FieldError("Name", "Name is required")
		}
		return
	})

	eb.Field("Image").
		WithContextValue(
			media.MediaBoxConfig,
			&media_library.MediaBoxConfig{
				AllowType: "image",
				Sizes: map[string]*base.Size{
					"thumb": {
						Width:  100,
						Height: 100,
					},
				},
			})

	return p
}

type productItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func productsSelector(db *gorm.DB) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var ps []models.Product
		var items []productItem
		searchKey := ctx.R.FormValue("keyword")
		sql := db.Order("created_at desc").Limit(10)
		if searchKey != "" {
			key := fmt.Sprintf("%%%s%%", searchKey)
			sql = sql.Where("name ILIKE ? or code ILIKE ?", key, key)
		}
		sql.Find(&ps)
		for _, p := range ps {
			items = append(items, productItem{
				ID:    strconv.Itoa(int(p.ID)),
				Name:  p.Name,
				Image: p.Image.URL("thumb"),
			})
		}
		r.Data = items
		return
	}
}
