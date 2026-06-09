package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/worker"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	h "github.com/r0vx/htmlgo"
)

// addJobs 添加后台任务
func addJobs(w *worker.Builder) {
	w.NewJob("noArgJob").
		Handler(func(ctx context.Context, job worker.GoJobInterface) error {
			job.AddLog("hoho1")
			job.AddLog("hoho2")
			job.AddLog("hoho3")
			return nil
		})
	w.NewJob("progressTextJob").
		Handler(func(ctx context.Context, job worker.GoJobInterface) error {
			for i := 1; i <= 10; i++ {
				select {
				case <-ctx.Done():
					job.AddLog("job aborted")
					return nil
				default:
					job.SetProgress(uint(i * 10))
					job.AddLog(fmt.Sprintf("Processing step %d", i))
					time.Sleep(time.Second)
				}
			}
			job.SetProgressText(`<a href="https://www.google.com">Download users</a>`)
			return nil
		})
	type ArgJobResource struct {
		F1 string
		F2 int
		F3 bool
	}
	ajb := w.NewJob("argJob").
		Resource(&ArgJobResource{}).
		Handler(func(ctx context.Context, job worker.GoJobInterface) error {
			jobInfo, _ := job.GetJobInfo()
			job.AddLog(fmt.Sprintf("Argument %#+v", jobInfo.Argument))
			job.AddLog(fmt.Sprintf("Context %#+v", jobInfo.Context))
			return nil
		})
	ajb.GetResourceBuilder().Editing().Field("F1").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		var vErr web.ValidationErrors
		if ve, ok := ctx.Flash.(*web.ValidationErrors); ok {
			vErr = *ve
		}
		return shadcn.Input().
			Attr(web.VField(field.Name, field.Value(obj))...).
			Label(field.Label).
			ErrorMessages(vErr.GetFieldErrors(field.Name)...)
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		v := ctx.R.FormValue("F1")
		obj.(*ArgJobResource).F1 = v

		if v == "aaa" {
			return errors.New("cannot be aaa")
		}
		return nil
	})

	type ScheduleJobResource struct {
		F1 string
		worker.Schedule
	}
	w.NewJob("scheduleJob").
		Resource(&ScheduleJobResource{}).
		Handler(func(ctx context.Context, job worker.GoJobInterface) error {
			jobInfo, _ := job.GetJobInfo()
			job.AddLog(fmt.Sprintf("%#+v", jobInfo.Argument))
			return nil
		})

	w.NewJob("errorJob").
		Handler(func(ctx context.Context, job worker.GoJobInterface) error {
			job.AddLog("=====perform error job")
			return errors.New("imError")
		})

	w.NewJob("panicJob").
		Handler(func(ctx context.Context, job worker.GoJobInterface) error {
			job.AddLog("=====perform panic job")
			panic("letsPanic")
		})
}
