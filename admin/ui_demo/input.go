package ui_demo

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/pkg/errors"
	"github.com/r0vx/admin/activity"
	"github.com/r0vx/admin/media"
	"github.com/r0vx/admin/media/base"
	"github.com/r0vx/admin/presets/actions"
	"github.com/r0vx/admin/worker"
	"github.com/samber/lo"

	"example/models"

	"github.com/r0vx/admin/media/media_library"
	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/amap"
	"github.com/r0vx/x/ui/codemirror"
	"github.com/r0vx/x/ui/shadcn"
	"gorm.io/gorm"
)

func ConfigInputDemo(b *presets.Builder, db *gorm.DB, ab *activity.Builder, wb *worker.Builder) {
	inputDemo := b.Model(&models.InputDemo{})
	// MenuIcon("view_quilt")

	defer func() {
		ab.RegisterModel(inputDemo)
	}()

	// 注册 Switch 切换事件
	b.GetWebBuilder().RegisterEventFunc("eventToggleSwitch", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id := ctx.R.FormValue("id")
		var demo models.InputDemo
		if err = db.First(&demo, id).Error; err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
			return
		}
		demo.Switch1 = !demo.Switch1
		if err = db.Save(&demo).Error; err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
			return
		}
		presets.ShowMessage(&r, "状态已更新", "success")
		return
	})

	// 列表配置
	cl := inputDemo.Listing("ID", "Switch1", "Slider1", "Select1", "UpdatedAt").
		SelectableColumns(true).
		SearchColumns("corp_id", "provider_name").
		PerPage(20)

	cl.OrderableFields("Slider1", "UpdatedAt")

	// ID 列 - 固定列宽 60px
	cl.Field("ID").CellClass("w-[60px]")

	// Switch1 列表字段 - 带事件的开关，固定列宽 80px
	cl.Field("Switch1").
		CellClass("w-[80px]").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			info := obj.(*models.InputDemo)
			onclick := web.Plaid().
				EventFunc("eventToggleSwitch").
				Query("id", fmt.Sprint(info.ID)).Go()
			return shadcn.Switch().
				Checked(info.Switch1).
				Disabled(field.Disabled).
				OnChange(onclick).
				Attr("@click.stop", true)
		})

	// Slider1 列使用 Progress 进度条展示
	cl.Field("Slider1").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		info := obj.(*models.InputDemo)
		return shadcn.Progress().ModelValue(info.Slider1)
	})

	// 快捷筛选标签
	cl.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		return []*presets.FilterTab{
			{
				ID:    "all",
				Label: "全部",
				Query: url.Values{"all": []string{""}},
			},
			{
				// ft_ 权限 demo：给 tab 显式 ID 后，即可在 admin/perm.go 用 *:input_demos:ft_enabled:*
				// 控制其可见性（与 fl_/fm_ 同构，框架自动剔除无权限 tab，无需在本 func 里写 if）。
				ID:    "enabled",
				Label: "启用",
				Query: url.Values{"enabled": []string{"true"}},
			},
		}
	})

	// 排序白名单示例：cl.OrderableFields("ID", "TextField1", "Switch1")

	// 编辑配置
	ed := inputDemo.Editing(
		&presets.FieldsSection{
			Title: "BasicInfo",
			Rows: [][]string{
				{"Location"},
				{"TextField1"},
				{"TextArea1"},
				{"Switch1"},
				{"Slider1"},
				{"Select1"},
				{"Radio1"},
				{"FileInput1"},
				{"Combobox1"},
				{"Checkbox1"},
				{"Autocomplete1"},
				{"ButtonGroup1"},
				{"Badge", "BadgeSelect"},
				{"ItemGroup1"},
				{"ListItemGroup1"},
				{"SlideGroup1"},
				{"ColorPicker1"},
				{"DatePicker1"},
				{"DatePickerMonth1"},
				{"TimePicker1"},
				{"CodeMirror1"},
				{"MediaLibrary1"},

				{"SelectedCustomers"},
			},
		},
	)

	// TextField1       string
	// TextArea1        string
	// Switch1          bool
	// Slider1          int
	// Select1          string
	// RangeSlider1     string
	// Radio1           string
	// FileInput1       string
	// Combobox1        string
	// Checkbox1        string
	// Autocomplete1    string
	// ButtonGroup1     string
	// ChipGroup1       string
	// ItemGroup1       string
	// ListItemGroup1   string
	// SlideGroup1      string
	// ColorPicker1     string
	// DatePicker1      string
	// DatePickerMonth1 string
	// TimePicker1      string
	// MediaLibrary1    media_library.MediaBox
	// TextField1 - 使用 shadcn Input
	ed.Field("TextField1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.Input().Label(field.Label).Attr(web.VField(field.Name, field.Value(obj))...).
				ErrorMessages(field.Errors...).Tips("提示信息")
		})

	// // TextArea1 - 使用 shadcn Textarea
	// ed.Field("TextArea1").
	// 	ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	// 		return shadcn.Textarea().Label(field.Label).Attr(web.VField(field.Name, field.Value(obj))...).ErrorMessages(field.Errors...)
	// 	})

	// Switch1 - 使用 shadcn Switch
	ed.Field("Switch1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			checked, _ := field.Value(obj).(bool)
			return shadcn.Switch().Label(field.Label).Checked(checked).Disabled(field.Disabled).Attr(web.VField(field.Name, field.Value(obj))...)
		})

	// Slider1 - 使用 shadcn Slider
	ed.Field("Slider1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(int)
			return shadcn.Slider().Label(field.Label).ModelValue(val).Disabled(field.Disabled).Attr(web.VField(field.Name, val)...).ErrorMessages(field.Errors...)
		})

	// Select1 - 使用 shadcn Select
	ed.Field("Select1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			// val, _ := field.Value(obj).(string)

			var items = []shadcn.DefaultOptionItem{}

			items = append(items, shadcn.DefaultOptionItem{
				Text:  "Tokyo",
				Value: "1",
			})
			items = append(items, shadcn.DefaultOptionItem{
				Text:  "Canberra",
				Value: "2",
			})
			items = append(items, shadcn.DefaultOptionItem{
				Text:  "Hangzhou",
				Value: "3",
			})

			return shadcn.Select().
				Items(items).
				Placeholder("选择城市"). // 通过方法设置
				Label(field.Label).
				Attr(web.VField(field.Name, field.Value(obj))...).ErrorMessages(field.Errors...)

			// return shadcn.Select(
			// 	shadcn.SelectTrigger(
			// 		shadcn.SelectValue().Placeholder("Select a city2"),
			// 	),
			// ).Items(items).Label(field.Label).Attr(web.VField(field.Name, val)...).ErrorMessages(field.Errors...)
		})

	// Radio1 - 使用 shadcn RadioGroup
	ed.Field("Radio1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)
			return h.Div(
				h.Label(field.Label).Class("text-sm font-medium"),
				shadcn.RadioGroup(
					h.Div(
						shadcn.RadioGroupItem().Value("1").Id("radio-tokyo"),
						h.Label("Tokyo").Attr("for", "radio-tokyo").Class("ml-2"),
					).Class("flex items-center space-x-2"),
					h.Div(
						shadcn.RadioGroupItem().Value("2").Id("radio-canberra"),
						h.Label("Canberra").Attr("for", "radio-canberra").Class("ml-2"),
					).Class("flex items-center space-x-2"),
					h.Div(
						shadcn.RadioGroupItem().Value("3").Id("radio-hangzhou"),
						h.Label("Hangzhou").Attr("for", "radio-hangzhou").Class("ml-2"),
					).Class("flex items-center space-x-2"),
				).DefaultValue(val).Attr(web.VField(field.Name, val)...).Class("mt-2").ErrorMessages(field.Errors...),
			).Class("space-y-2")
		})

	// FileInput1 - 使用 shadcn FileInput
	ed.Field("FileInput1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return shadcn.FileInput().Attr("name", field.Name).Class("mt-2").ErrorMessages(field.Errors...).Label(field.Label)

		}).
		SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			fs := ctx.R.MultipartForm.File[field.Name]
			if len(fs) == 0 {
				return
			}
			f, err := fs[0].Open()
			if err != nil {
				return
			}
			defer f.Close()

			b, err := io.ReadAll(f)
			if err != nil {
				return
			}
			obj.(*models.InputDemo).FileInput1 = fmt.Sprint(len(b))

			return
		})

	// Combobox1 - 使用 shadcn Combobox
	ed.Field("Combobox1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)
			return shadcn.Combobox(
				shadcn.ComboboxAnchor(
					shadcn.ComboboxInput().Placeholder("Search..."),
					shadcn.ComboboxTrigger(),
				).Class("w-full flex"),
				shadcn.ComboboxList(
					shadcn.ComboboxEmpty(h.Text("No results found.")),
					shadcn.ComboboxGroup(
						shadcn.ComboboxItem(h.Text("Tokyo")).Value("Tokyo"),
						shadcn.ComboboxItem(h.Text("Canberra")).Value("Canberra"),
						shadcn.ComboboxItem(h.Text("Hangzhou")).Value("Hangzhou"),
					),
				),
			).Attr(web.VField(field.Name, val)...).Class("mt-2").ErrorMessages(field.Errors...).Label(field.Label)

		})

		// Checkbox1 - 使用 shadcn Checkbox
		// ed.Field("Checkbox1").
		// 	ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		// 		checked, _ := field.Value(obj).(string)
		// 		return shadcn.Checkbox().Label(field.Label).Checked(checked == "true" || checked == "1").
		// 			Attr(web.VField(field.Name, field.Value(obj))...)
		// 	})

		// Autocomplete1 - 使用 shadcn Autocomplete（JSON 格式存储）
	autocomplete1 := shadcn.Autocomplete().
		Items([]shadcn.DefaultOptionItem{
			{Text: "Tokyo", Value: "Tokyo"},
			{Text: "Canberra", Value: "Canberra"},
			{Text: "Hangzhou", Value: "Hangzhou"},
		}).
		Multiple(true)
	ed.Field("Autocomplete1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return autocomplete1.Label(field.Label).
				Attr(web.VField(field.FormKey, autocomplete1.ParseValue(field.Value(obj)))...).
				// Autocomplete1 选择变化时，同步选中项 name 到 TextField1
				Attr("@update:modelValue", `(val) => { if (Array.isArray(val)) { form['TextField1'] = val.join(', ') } else if (val) { form['TextField1'] = val } }`).
				ErrorMessages(field.Errors...)
		}).
		SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			values := ctx.R.Form[field.FormKey]
			obj.(*models.InputDemo).Autocomplete1 = autocomplete1.FormatValue(values)
			return nil
		})

	// ButtonGroup1 - 使用 shadcn Button 组合
	// ButtonGroup1 - 使用 shadcn ButtonGroup
	ed.Field("ButtonGroup1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)
			return shadcn.ButtonGroup(
				shadcn.Button(h.Text("Left")).
					Variant(shadcn.ButtonVariantOutline).
					ClassIf("bg-primary text-primary-foreground", val == "left").
					Attr("@click", fmt.Sprintf("$refs.%s.value = 'left'", field.Name)),
				shadcn.Button(h.Text("Center")).
					Variant(shadcn.ButtonVariantOutline).
					ClassIf("bg-primary text-primary-foreground", val == "center").
					Attr("@click", fmt.Sprintf("$refs.%s.value = 'center'", field.Name)),
				shadcn.Button(h.Text("Right")).
					Variant(shadcn.ButtonVariantOutline).
					ClassIf("bg-primary text-primary-foreground", val == "right").
					Attr("@click", fmt.Sprintf("$refs.%s.value = 'right'", field.Name)),
				shadcn.Button(h.Text("Justify")).
					Variant(shadcn.ButtonVariantOutline).
					ClassIf("bg-primary text-primary-foreground", val == "justify").
					Attr("@click", fmt.Sprintf("$refs.%s.value = 'justify'", field.Name)),
			)
			// h.Input("").Type("hidden").Attr("ref", field.Name).Attr(web.VField(field.Name, val)...),

		})

		// Badge - 展示各种 Badge 样式
	ed.Field("Badge").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(
				h.Label(field.Label).Class("text-sm font-medium"),
				h.Div(
					// Default
					shadcn.Badge(h.Text("Default")),
					// Secondary
					shadcn.Badge(h.Text("Secondary")).Variant(shadcn.BadgeVariantSecondary),
					// Outline
					shadcn.Badge(h.Text("Outline")).Variant(shadcn.BadgeVariantOutline),
					// Destructive
					shadcn.Badge(h.Text("Destructive")).Variant(shadcn.BadgeVariantDestructive),
				).Class("flex flex-wrap gap-2 mt-2"),
			).Class("space-y-2")
		})
	// BadgeSelect - 可选择的 Badge 组
	ed.Field("BadgeSelect").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)

			// 辅助函数：根据是否选中返回 variant
			variantFor := func(option string) shadcn.BadgeVariant {
				if val == option {
					return shadcn.BadgeVariantDefault
				}
				return shadcn.BadgeVariantOutline
			}

			options := []string{"left", "center", "right", "justify"}
			var badges []h.HTMLComponent
			for _, opt := range options {
				badges = append(badges,
					shadcn.Badge(h.Text(cases.Title(language.English).String(opt))).
						Variant(variantFor(opt)).
						Class("cursor-pointer").
						Attr("@click", fmt.Sprintf("$refs.%s.value = '%s'", field.Name, opt)),
				)
			}
			badges = append(badges,
				h.Input("").Type("hidden").Attr("ref", field.Name).Attr(web.VField(field.Name, val)...),
			)

			return h.Div(
				h.Label(field.Label).Class("text-sm font-medium"),
				h.Div(badges...).Class("flex flex-wrap gap-2 mt-2"),
			).Class("space-y-2")
		})

	// DatePicker1 - 使用 shadcn DatePicker
	ed.Field("DatePicker1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)
			return shadcn.DatePicker().Label(field.Label).ModelValue(val).Attr(web.VField(field.Name, val)...).Class("mt-2").ErrorMessages(field.Errors...)
		})

	// DatePickerMonth1 - 月份选择器（shadcn DatePicker 暂不支持月份模式，使用普通 DatePicker）
	ed.Field("DatePickerMonth1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)
			return shadcn.DatePicker().Label(field.Label).ModelValue(val).Attr(web.VField(field.Name, val)...).Class("mt-2").ErrorMessages(field.Errors...).Disabled(field.Disabled)

		})

	// TimePicker1 - 使用 shadcn TimePicker
	ed.Field("TimePicker1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)
			return shadcn.TimePicker().Label(field.Label).ModelValue(val).Attr(web.VField(field.Name, val)...).Class("mt-2").ErrorMessages(field.Errors...).Disabled(field.Disabled)

		})

	// CodeMirror1 - 代码编辑器
	ed.Field("CodeMirror1").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)
			// 临时换行验证：空值时用一段含超长无空格 query 串的 HTTP 报文作默认，测试软换行
			if val == "" {
				val = "# ERROR\nPost \"https://www.myshop2.com/v3/pay/UnifiedOrder\": tls: failed to verify certificate: x509: certificate has expired or is not yet valid\n\nPOST /v3/pay/UnifiedOrder HTTP/1.1\nHost: www.myshop2.com\nContent-Type: application/x-www-form-urlencoded\n\namount=1000&code=8888&mchOrderNo=TEST1178124797807448400&mchid=10000&notifyUrl=http%3A%2F%2Fwww.myshop2.com%2Fv3%2Fpay%2Fnotify%2FTEST1178&sign=C2FF991CC869F1F5FAA48A8443A0BFC6"
			}
			return codemirror.CodeMirror().
				Label(field.Label).
				ModelValue(val).
				Language(codemirror.LangText).
				Theme(codemirror.ThemeDark).
				Height("300px").
				Placeholder("请输入...").
				Attr(web.VField(field.Name, val)...).
				ErrorMessages(field.Errors...)
		})

	// 添加必填验证
	ed.ValidateFunc(func(obj any, ctx *web.EventContext) (err web.ValidationErrors) {
		demo := obj.(*models.InputDemo)
		if demo.TextField1 == "" {
			err.FieldError("TextField1", "TextField1 is required")
		}
		if demo.TextArea1 == "" {
			err.FieldError("TextArea1", "TextArea1 is required")
		}

		if demo.Slider1 == 0 {
			err.FieldError("Slider1", "Slider1 is required")
		}

		if demo.Select1 == "" {
			err.FieldError("Select1", "Select1 is required")
		}
		if demo.Radio1 == "" {
			err.FieldError("Radio1", "Radio1 is required")
		}

		if len(demo.Autocomplete1) == 0 {
			err.FieldError("Autocomplete1", "Autocomplete1 is required")
		}
		if demo.DatePicker1 == "" {
			err.FieldError("DatePicker1", "DatePicker1 is required")
		}
		if demo.TimePicker1 == "" {
			err.FieldError("TimePicker1", "TimePicker1 is required")
		}
		return
	})

	// Location - 高德地图选点组件
	ed.Field("Location").
		ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			val, _ := field.Value(obj).(string)
			return h.Div(
				h.Label("Location").Class("text-sm font-medium"),
				h.Div(
					amap.AmapPicker().
						ApiKey("d806a2f74b0016c8190c71640d44b98d").
						SecurityJsCode("2c8ed246088985a69ed378d820e77602").
						Attr("v-model", fmt.Sprintf("form[%q]", field.Name)),
					h.Input("").Type("hidden").Attr(web.VField(field.Name, val)...),
				),
			).Class("space-y-2")
		})

	ed.Field("MediaLibrary1").
		WithContextValue(
			media.MediaBoxConfig,
			&media_library.MediaBoxConfig{
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
			})

	// Detailing 配置（只读详情页）
	dt := inputDemo.Detailing("Location")
	dt.Field("Location").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		val, _ := field.Value(obj).(string)
		return amap.AmapDisplay().
			ApiKey("d806a2f74b0016c8190c71640d44b98d").
			SecurityJsCode("2c8ed246088985a69ed378d820e77602").
			Value(val).
			Width(500).Height(250)
	})

	// 添加 Worker ActionJob，并在列表页添加触发按钮
	if wb != nil {
		inputDemoJob := wb.ActionJob(
			"Input Demo Job",
			inputDemo,
			func(ctx context.Context, job worker.GoJobInterface) error {
				for i := 1; i <= 5; i++ {
					select {
					case <-ctx.Done():
						job.AddLog("job aborted")
						return nil
					default:
						job.SetProgress(uint(i * 20))
						job.AddLog(fmt.Sprintf("Processing step %d", i))
						time.Sleep(time.Second)
					}
				}
				return nil
			},
		).Description("Input Demo 示例任务")

		cl.BulkAction("Input Demo Job").
			ButtonCompFunc(func(ctx *web.EventContext) h.HTMLComponent {
				return shadcn.Button(h.Text("Input Demo Job")).Size(shadcn.ButtonSizeSm).
					Attr("@click", inputDemoJob.URL())
			})
	}

	// ListingDialog Selector（列表弹窗选择器）
	ConfigureDialogCustomerSelector(db, b)
	ed.Field("SelectedCustomers").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		// 从模型读取已保存的 selectedIds
		demo := obj.(*models.InputDemo)
		var selectedIds []string
		if demo.SelectedCustomers != "" {
			selectedIds = strings.Split(demo.SelectedCustomers, ",")
		}
		compo, err := customerSelector(ctx, db, selectedIds)
		if err != nil {
			panic(err)
		}
		return web.Portal(compo).Name(portalCustomerSelector)
	})
}

const (
	eventSelectCustomer           = "eventSelectCustomer"
	portalCustomerSelector        = "portalCustomerSelector"
	uriNameDialogCustomerSelector = "dialog-customer-selector"
)

func customerSelector(_ *web.EventContext, db *gorm.DB, selectedIds []string) (h.HTMLComponent, error) {
	var items []*models.Customer
	if len(selectedIds) > 0 {
		if err := db.Where("id IN (?)", selectedIds).Find(&items).Error; err != nil {
			return nil, errors.Wrap(err, "find customers")
		}
	}

	// 将 selectedIds 拼接为逗号分隔字符串，通过隐藏 input 提交到表单
	value := strings.Join(selectedIds, ",")

	plaidCall := web.Plaid().URL("/" + uriNameDialogCustomerSelector).
		EventFunc(actions.OpenListingDialog)
	// 空切片序列化为 null 会导致 JS 错误，仅在有值时传 Query
	if len(selectedIds) > 0 {
		plaidCall.Query("selected_ids", selectedIds)
	}

	children := []h.HTMLComponent{
		// 隐藏 input，绑定到表单字段 SelectedCustomers
		h.Input("").Type("hidden").Attr(web.VField("SelectedCustomers", value)...),
		shadcn.Button(h.Text("Select Customers")).Class("mb-2").Attr("@click",
			plaidCall.Go(),
		),
	}
	children = append(children, lo.Map(items, func(v *models.Customer, _ int) h.HTMLComponent {
		return h.Div().Children(
			h.Text(v.Name),
		)
	})...)

	return h.Div().Class("flex flex-col gap-2 items-start").Children(children...), nil
}

const paramSelectedIds = "selected_ids"

func ConfigureDialogCustomerSelector(db *gorm.DB, pb *presets.Builder) {
	//b := pb.Model(&models.Customer{}).URIName(uriNameDialogCustomerSelector).InMenu(false)

	b := pb.Model(&models.Customer{}).URIName(uriNameDialogCustomerSelector).
		InMenu(true)
	registerEventSelectCustomer(db, pb)

	lb := b.Listing().DialogContentClass(presets.DialogSizeLg).
		PerPage(20).
		SearchColumns("name").
		SelectableColumns(true).
		// PopupExcludeFilters("region").
		PopupExcludeColumns("Titile").
		PopupSearchOff(false).
		PopupPerPage(10).
		PopupHideNewButton(true).
		PopupHideRowMenu(true)

	lb.WrapCell(func(in presets.CellProcessor) presets.CellProcessor {
		return in
	})

	lb.BulkAction("Confirm").ButtonCompFunc(func(ctx *web.EventContext) h.HTMLComponent {
		return shadcn.Button(h.Text("Confirm")).Attr("@click", web.Plaid().
			EventFunc(eventSelectCustomer).
			Query(paramSelectedIds, web.Var("locals.selected_ids")).
			MergeQuery(true).
			Go(),
		)
	}).PopupOnly(true)

	// ========== 所有 12 种 FilterItemType ==========
	lb.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		// 级联选择数据（LinkageSelect 使用）
		// 第一级：国家
		level1 := []shadcn.FilterLinkageItem{
			{ID: "china", Name: "中国", ChildrenIDs: []string{"beijing", "shanghai", "guangdong"}},
			{ID: "usa", Name: "美国", ChildrenIDs: []string{"california", "newyork"}},
		}
		// 第二级：省/州
		level2 := []shadcn.FilterLinkageItem{
			{ID: "beijing", Name: "北京", ChildrenIDs: []string{"chaoyang", "haidian"}},
			{ID: "shanghai", Name: "上海", ChildrenIDs: []string{"pudong", "xuhui"}},
			{ID: "guangdong", Name: "广东", ChildrenIDs: []string{"guangzhou", "shenzhen"}},
			{ID: "california", Name: "加州", ChildrenIDs: []string{"losangeles", "sanfrancisco"}},
			{ID: "newyork", Name: "纽约州", ChildrenIDs: []string{"nyc", "buffalo"}},
		}
		// 第三级：城市
		level3 := []shadcn.FilterLinkageItem{
			{ID: "chaoyang", Name: "朝阳区"},
			{ID: "haidian", Name: "海淀区"},
			{ID: "pudong", Name: "浦东新区"},
			{ID: "xuhui", Name: "徐汇区"},
			{ID: "guangzhou", Name: "广州"},
			{ID: "shenzhen", Name: "深圳"},
			{ID: "losangeles", Name: "洛杉矶"},
			{ID: "sanfrancisco", Name: "旧金山"},
			{ID: "nyc", Name: "纽约市"},
			{ID: "buffalo", Name: "布法罗"},
		}

		return []*shadcn.FilterItem{
			{
				Key:           "region",
				Label:         "12. LinkageSelect（地区级联）",
				ItemType:      shadcn.FilterItemTypeLinkageSelect,
				LinkageItems:  [][]shadcn.FilterLinkageItem{level1, level2, level3},
				LinkageLabels: []string{"国家", "省/州", "城市"},
				LinkageSelectData: shadcn.FilterLinkageSelectData{
					SQLConditions: []string{
						`country = ?`,
						`province = ?`,
						`city = ?`,
					},
				},
			},
		}
	})
}

func registerEventSelectCustomer(db *gorm.DB, b *presets.Builder) {
	b.GetWebBuilder().RegisterEventFunc(eventSelectCustomer, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		selectedIds := strings.Split(strings.TrimSpace(ctx.R.FormValue(paramSelectedIds)), ",")
		selectedIds = lo.Filter(selectedIds, func(v string, _ int) bool { return v != "" })
		compo, err := customerSelector(ctx, db, selectedIds)
		if err != nil {
			presets.ShowMessage(&r, err.Error(), "error")
			return
		}
		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: portalCustomerSelector,
			Body: compo,
		})
		web.AppendRunScripts(&r, presets.CloseListingDialogVarScript)
		return
	})
}
