package examples

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/r0vx/admin/presets"
	"github.com/r0vx/admin/presets/actions"
	"github.com/r0vx/admin/presets/gorm2op"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

type Thumb struct {
	Name string
}

type Customer struct {
	ID              int
	Name            string
	Email           string
	Description     string
	Thumb1          *Thumb `gorm:"-"`
	CompanyID       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ApprovedAt      *time.Time
	TermAgreedAt    *time.Time
	ApprovalComment string
	LanguageCode    string
	Events          []*Event `gorm:"-"`
}

func (c *Customer) PageTitle() string {
	return c.Name
}

type Note struct {
	ID         int
	SourceType string
	SourceID   int
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreditCard struct {
	ID              int
	CustomerID      int
	Number          string
	ExpireYearMonth string
	Name            string
	Type            string
	Phone           string
	Email           string
}

type Payment struct {
	ID                   int
	CustomerID           int
	CurrencyCode         string
	Amount               int
	PaymentMethodID      int
	StatementDescription string
	Description          string
	AuthorizeOnly        bool
	CreatedAt            time.Time
}

type Event struct {
	ID          int
	SourceType  string // Payment, Customer
	SourceID    int
	CreatedAt   time.Time
	Type        string
	Description string
}

type Language struct {
	Code string `gorm:"unique;not null"`
	Name string
}

func (l *Language) PrimarySlug() string {
	return l.Code
}

func (l *Language) PrimaryColumnValuesBySlug(slug string) map[string]string {
	return map[string]string{
		"code": slug,
	}
}

type Company struct {
	ID   int
	Name string
}

type Product struct {
	ID        int
	Name      string
	OwnerName string
}

func (*Product) TableName() string {
	return "preset_products"
}

// addListener 创建一个监听器组件
func addListener(v any) h.HTMLComponent {
	simpleReload := web.Plaid().MergeQuery(true).Go()
	return web.Listen(
		presets.NotifModelsCreated(v), simpleReload,
		presets.NotifModelsUpdated(v), simpleReload,
		presets.NotifModelsDeleted(v), simpleReload,
	)
}

// =====================================
// 辅助组件函数 (替代 vuetifyx 组件)
// =====================================

// actionCard 创建带标题和操作按钮的卡片
func actionCard(title string, content h.HTMLComponent, actions ...h.HTMLComponent) h.HTMLComponent {
	return shadcn.Card(
		shadcn.CardHeader(
			h.Div(
				shadcn.CardTitle(h.Text(title)),
				h.Div(actions...).Class("flex gap-2"),
			).Class("flex justify-between items-center"),
		),
		shadcn.CardContent(content),
	).Class("mb-4")
}

// detailInfo 创建详情信息布局
func detailInfo(columns ...h.HTMLComponent) h.HTMLComponent {
	return h.Div(columns...).Class("grid grid-cols-1 md:grid-cols-2 gap-6 p-4")
}

// detailColumn 创建详情列
func detailColumn(header string, fields ...h.HTMLComponent) h.HTMLComponent {
	return h.Div(
		h.If(header != "", h.H5(header).Class("text-sm font-medium text-muted-foreground pb-2")),
		h.Div(fields...).Class("space-y-2"),
	)
}

// detailField 创建详情字段
func detailField(label string, value h.HTMLComponent) h.HTMLComponent {
	return h.Div(
		h.Label(label).Class("text-sm text-muted-foreground min-w-44 inline-block"),
		value,
	).Class("flex pb-2")
}

// optionalText 创建可选文本，空值时显示零值标签
func optionalText(text, zeroLabel string) h.HTMLComponent {
	if text != "" {
		return h.Span(text)
	}
	return h.Span(zeroLabel).Class("text-muted-foreground")
}

// Preset1 创建 Admin 预设配置
func Preset1(db *gorm.DB) (r *presets.Builder) {
	err := db.AutoMigrate(
		&Customer{},
		&Note{},
		&CreditCard{},
		&Payment{},
		&Event{},
		&Company{},
		&Product{},
		&Language{},
	)
	if err != nil {
		panic(err)
	}

	p := presets.New().URIPrefix("/admin")

	p.BrandFunc(func(ctx *web.EventContext) h.HTMLComponent {
		return h.Components(
			shadcn.Icon("ship").Class("pr-2"),
			h.Span("My Admin").Class("font-semibold"),
		)
	})
	// .BrandTitle("My Admin")

	writeFieldDefaults := p.FieldDefaults(presets.WRITE)
	writeFieldDefaults.FieldType(&Thumb{}).ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		i, err := reflectutils.Get(obj, field.Name)
		if err != nil {
			panic(err)
		}
		return h.Text(i.(*Thumb).Name)
	})

	p.FieldDefaults(presets.LIST).FieldType(&Thumb{}).ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		i, err := reflectutils.Get(obj, field.Name)
		if err != nil {
			panic(err)
		}
		return h.Text(i.(*Thumb).Name)
	})

	p.FieldDefaults(presets.DETAIL).FieldType([]*Event{}).ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		events := reflectutils.MustGet(obj, field.Name).([]*Event)
		typeName := reflect.ValueOf(obj).Elem().Type().Name()
		objId := fmt.Sprint(reflectutils.MustGet(obj, "ID"))

		// 构建 DataTable 数据
		var data []map[string]any
		for _, e := range events {
			data = append(data, map[string]any{
				"id":          fmt.Sprint(e.ID),
				"Type":        e.Type,
				"Description": e.Description,
			})
		}
		cols := []shadcn.DataTableColumn{
			{Name: "Type", Title: "Type"},
			{Name: "Description", Title: "Description"},
		}
		dt := shadcn.DataTable().
			Data(data).
			Columns(cols).
			WithoutHeaders(true).
			Hover(false)

		return actionCard(
			field.Label,
			dt,
			addListener(&Event{}),
			shadcn.Button(h.Text("Add Event")).
				Attr("@click",
					web.Plaid().EventFunc(actions.New).
						Query("model", typeName).
						Query("model_id", objId).
						URL("/admin/events").
						Go(),
				),
		)
	})

	p.DataOperator(gorm2op.DataOperator(db))

	p.MenuGroup("Customer Management").Icon("group").SubItems("my_customers", "company")
	mp := p.Model(&Product{}).MenuIcon("laptop")
	mp.Listing().PerPage(3)

	m := p.Model(&Customer{}).URIName("my_customers")
	p.Model(&Company{})
	m.Labels(
		"Name", "名字",
		"Bool1", "性别",
		"Float1", "体重",
		"CompanyID", "公司",
	).Placeholders(
		"Name", "请输入你的名字",
	)

	l := m.Listing("Name", "CompanyID", "ApprovalComment").SearchColumns("name", "email", "description").PerPage(5).SelectableColumns(true)
	l.Field("Name").Label("列表的名字")
	l.Field("CompanyID").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*Customer)
		var comp Company
		err := db.Find(&comp, u.CompanyID).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return h.Td(
			h.A().Text(comp.Name).
				Attr("@click.stop",
					web.Plaid().URL("/admin/companies").
						EventFunc(actions.Edit).
						Query(presets.ParamID, fmt.Sprint(comp.ID)).
						Go()),
		)
	})

	l.BulkAction("Approve").Label("Approve").UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) (err error) {
		comment := ctx.R.FormValue("ApprovalComment")
		if len(comment) < 10 {
			ctx.Flash = "comment should larger than 10"
			return
		}
		err = db.Model(&Customer{}).
			Where("id IN (?)", selectedIds).
			Updates(map[string]any{"approved_at": time.Now(), "approval_comment": comment}).Error
		if err != nil {
			ctx.Flash = err.Error()
		} else {
			r.Emit(
				presets.NotifModelsUpdated(&Customer{}),
				presets.PayloadModelsUpdated{Ids: selectedIds},
			)
		}
		return
	}).ComponentFunc(func(selectedIds []string, ctx *web.EventContext) h.HTMLComponent {
		comment := ctx.R.FormValue("ApprovalComment")
		errorMessage := ""
		if ctx.Flash != nil {
			errorMessage = ctx.Flash.(string)
		}
		return shadcn.Input().
			Attr(web.VField("ApprovalComment", comment)...).
			Label("Content").
			ErrorMessages(errorMessage)
	})

	l.BulkAction("Delete").Label("Delete").UpdateFunc(func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) (err error) {
		err = db.Where("id IN (?)", selectedIds).Delete(&Customer{}).Error
		if err == nil {
			r.Emit(
				presets.NotifModelsDeleted(&Customer{}),
				presets.PayloadModelsDeleted{Ids: selectedIds},
			)
		}
		return
	}).ComponentFunc(func(selectedIds []string, ctx *web.EventContext) h.HTMLComponent {
		return h.Div().Text(fmt.Sprintf("Are you sure you want to delete %s ?", selectedIds)).Class("text-destructive font-medium")
	})

	l.FilterDataFunc(func(ctx *web.EventContext) shadcn.FilterData {
		var companyOptions []shadcn.FilterSelectOption
		err := db.Model(&Company{}).Select("name as text, id as value").Scan(&companyOptions).Error
		if err != nil {
			panic(err)
		}

		return []*shadcn.FilterItem{
			{
				Key:          "created",
				Label:        "Created",
				Folded:       true,
				ItemType:     shadcn.FilterItemTypeDatetimeRange,
				SQLCondition: `extract(epoch from created_at) %s ?`,
			},
			{
				Key:          "approved",
				Label:        "Approved",
				ItemType:     shadcn.FilterItemTypeDatetimeRange,
				SQLCondition: `extract(epoch from approved_at) %s ?`,
			},
			{
				Key:          "name",
				Label:        "Name",
				Folded:       true,
				ItemType:     shadcn.FilterItemTypeString,
				SQLCondition: `name %s ?`,
			},
			{
				Key:          "company",
				Label:        "Company",
				ItemType:     shadcn.FilterItemTypeSelect,
				SQLCondition: `company_id %s ?`,
				Options:      companyOptions,
			},
		}
	})

	l.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		var c Company
		db.First(&c)
		return []*presets.FilterTab{
			{
				Label: "All",
				Query: url.Values{"all": []string{"1"}},
			},
			{
				Label: "Felix",
				Query: url.Values{"name.ilike": []string{"felix"}},
			},
			{
				Label: "The Plant",
				Query: url.Values{"company": []string{fmt.Sprint(c.ID)}},
			},
			{
				Label: "Approved",
				Query: url.Values{"approved.gt": []string{fmt.Sprint(1)}},
			},
		}
	})

	ef := m.Editing("Name", "CompanyID", "LanguageCode").
		ValidateFunc(func(obj any, ctx *web.EventContext) (err web.ValidationErrors) {
			cu := obj.(*Customer)
			if len(cu.Name) < 5 {
				err.FieldError("Name", "input more than 5 chars")
				err.GlobalError("there are some errors")
			}
			return
		})
	ef.Field("LanguageCode").Label("语言").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*Customer)
		var langs []Language
		err := db.Find(&langs).Error
		if err != nil {
			panic(err)
		}

		// 构建选项
		var selectItems []h.HTMLComponent
		for _, lang := range langs {
			selectItems = append(selectItems,
				shadcn.SelectItem(h.Text(lang.Name)).Value(lang.Code),
			)
		}

		return shadcn.Select(
			shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("Select language")),
			shadcn.SelectContent(selectItems...),
		).Attr(web.VField(field.Name, u.LanguageCode)...).
			Label(field.Label)
	})

	ef.Field("CompanyID").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u := obj.(*Customer)
		var companies []*Company
		err := db.Find(&companies).Error
		if err != nil {
			panic(err)
		}

		// 构建选项
		var selectItems []h.HTMLComponent
		for _, comp := range companies {
			selectItems = append(selectItems,
				shadcn.SelectItem(h.Text(comp.Name)).Value(fmt.Sprint(comp.ID)),
			)
		}

		return shadcn.Select(
			shadcn.SelectTrigger(shadcn.SelectValue().Placeholder("Select company")),
			shadcn.SelectContent(selectItems...),
		).Attr(web.VField("CompanyID", fmt.Sprint(u.CompanyID))...).
			Label(field.Label)
	})

	dp := m.Detailing("MainInfo", "Details", "Cards", "Events")

	dp.FetchFunc(func(obj any, id string, ctx *web.EventContext) (r any, err error) {
		cus := &Customer{}
		err = db.Find(cus, id).Error
		if err != nil {
			return
		}

		var events []*Event
		err = db.Where("source_type = ? AND source_id = ?", "Customer", id).Find(&events).Error
		if err != nil {
			return
		}
		cus.Events = events
		r = cus
		return
	})

	dp.Field("MainInfo").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		cu := obj.(*Customer)

		title := cu.Name
		if len(title) == 0 {
			title = cu.Description
		}

		var notes []*Note
		err := db.Where("source_type = 'Customer' AND source_id = ?", cu.ID).
			Order("id DESC").
			Find(&notes).Error
		if err != nil {
			panic(err)
		}

		cusID := fmt.Sprint(cu.ID)

		// 构建 DataTable 数据
		var data []map[string]any
		for _, n := range notes {
			data = append(data, map[string]any{
				"id":      fmt.Sprint(n.ID),
				"Content": n.Content,
				"_html_Content": fmt.Sprintf(`<div class="py-3"><div class="text-sm"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-blue-500 mr-2 inline-block"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>%s</div><div class="text-muted-foreground pl-7 text-xs">%s by Felix Sun</div></div>`,
					n.Content, n.CreatedAt.Format("Jan 02,15:04 PM")),
			})
		}
		cols := []shadcn.DataTableColumn{
			{Name: "Content", Title: "Content"},
		}
		dt := shadcn.DataTable().
			Data(data).
			Columns(cols).
			WithoutHeaders(true).
			CellHtmlMode(true).
			Hover(false)

		return actionCard(
			title,
			dt,
			addListener(&Note{}),
			shadcn.Button(h.Text("Add Note")).
				Attr("@click",
					web.Plaid().EventFunc(actions.New).
						Query("model", "Customer").
						Query("model_id", cusID).
						URL("/admin/notes").
						Go(),
				),
		)
	})

	dp.Field("Details").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		cu := obj.(*Customer)
		cusID := fmt.Sprint(cu.ID)

		var lang Language
		db.Where("code = ?", cu.LanguageCode).First(&lang)

		var termAgreed string
		if cu.TermAgreedAt != nil {
			termAgreed = cu.TermAgreedAt.Format("Jan 02,15:04 PM")
		}

		detail := detailInfo(
			detailColumn("ACCOUNT INFORMATION",
				detailField("Name", optionalText(cu.Name, "No Name")),
				detailField("Email", optionalText(cu.Email, "No Email")),
				detailField("Description", optionalText(cu.Description, "No Description")),
				detailField("ID", optionalText(cusID, "No ID")),
				detailField("Created", optionalText(cu.CreatedAt.Format("Jan 02,15:04 PM"), "")),
				detailField("Terms Agreed", optionalText(termAgreed, "Not Agreed Yet")),
				detailField("Language", optionalText(lang.Name, "No Language")),
			),
			detailColumn("BILLING INFORMATION"),
		)

		return actionCard(
			"Details",
			detail,
			web.Listen(
				m.NotifModelsUpdated(), web.Plaid().MergeQuery(true).Go(),
			),
			shadcn.Button(h.Text("Agree Terms")).
				Class("mr-2").
				Attr("@click", web.Plaid().
					EventFunc(actions.Action).
					Query(presets.ParamAction, "AgreeTerms").
					Query("customerID", cusID).
					Go()),
			shadcn.Button(h.Text("Update details")).
				Attr("@click", web.Plaid().
					EventFunc(actions.Edit).
					Query("customerID", cusID).
					URL("/admin/customers").Go()),
		)
	})

	dp.Field("Cards").ComponentFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		cu := obj.(*Customer)
		cusID := fmt.Sprint(cu.ID)

		var cards []*CreditCard
		err := db.Where("customer_id = ?", cu.ID).Order("id ASC").Find(&cards).Error
		if err != nil {
			panic(err)
		}

		// 构建 DataTable 数据
		var data []map[string]any
		for _, c := range cards {
			data = append(data, map[string]any{
				"id":              fmt.Sprint(c.ID),
				"Type":            c.Type,
				"Number":          c.Number,
				"ExpireYearMonth": c.ExpireYearMonth,
			})
		}
		cols := []shadcn.DataTableColumn{
			{Name: "Type", Title: "Type"},
			{Name: "Number", Title: "Number"},
			{Name: "ExpireYearMonth", Title: "Expire"},
		}
		dt := shadcn.DataTable().
			Data(data).
			Columns(cols).
			WithoutHeaders(true).
			Hover(false)

		return actionCard(
			"Cards",
			dt,
			addListener(&CreditCard{}),
			shadcn.Button(h.Text("Add Card")).
				Attr("@click",
					web.Plaid().EventFunc(
						actions.New).Query("customerID", cusID).
						URL("/admin/credit-cards").
						Go()),
		)
	})

	dp.Action("AgreeTerms").UpdateFunc(func(id string, ctx *web.EventContext, r *web.EventResponse) (err error) {
		if ctx.R.FormValue("Agree") != "true" {
			ve := &web.ValidationErrors{}
			ve.GlobalError("You must agree the terms")
			err = ve
			return
		}

		err = db.Model(&Customer{}).Where("id = ?", id).
			Updates(map[string]any{"term_agreed_at": time.Now()}).Error
		if err == nil {
			r.Emit(
				presets.NotifModelsUpdated(&Customer{}),
				presets.PayloadModelsUpdated{Ids: []string{id}},
			)
		}
		return
	}).ComponentFunc(func(id string, ctx *web.EventContext) h.HTMLComponent {
		var alert h.HTMLComponent

		if ve, ok := ctx.Flash.(*web.ValidationErrors); ok {
			alert = shadcn.Alert(
				shadcn.AlertDescription(h.Text(ve.GetGlobalError())),
			).Variant(shadcn.AlertVariantDestructive)
		}

		return h.Components(
			alert,
			shadcn.Checkbox().
				Attr(web.VField("Agree", ctx.R.FormValue("Agree"))...).
				Label("Agree the terms"),
		)
	})

	p.Model(&Note{}).
		InMenu(false).
		Editing("Content").
		SetterFunc(func(obj any, ctx *web.EventContext) {
			note := obj.(*Note)
			note.SourceID = ctx.ParamAsInt("model_id")
			note.SourceType = ctx.R.FormValue("model")
		})

	p.Model(&Event{}).
		InMenu(false).
		Editing("Type", "Description").
		SetterFunc(func(obj any, ctx *web.EventContext) {
			note := obj.(*Event)
			note.SourceID = ctx.ParamAsInt("model_id")
			note.SourceType = ctx.R.FormValue("model")
		})

	cc := p.Model(&CreditCard{}).
		InMenu(false)

	ccedit := cc.Editing("ExpireYearMonth", "Phone", "Email").
		SetterFunc(func(obj any, ctx *web.EventContext) {
			card := obj.(*CreditCard)
			card.CustomerID = ctx.ParamAsInt("customerID")
		})

	ccedit.Creating("Number")

	p.Model(&Language{}).PrimaryField("Code")

	return p
}
