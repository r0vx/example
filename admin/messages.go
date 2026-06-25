package admin

import (
	"example/models"

	"github.com/r0vx/x/i18n"
)

const I18nExampleKey i18n.ModuleKey = "I18nExampleKey"

type Messages struct {
	FilterTabsAll                  string
	FilterTabsHasUnreadNotes       string
	FilterTabsActive               string
	DemoTips                       string
	DemoUsernameLabel              string
	DemoPasswordLabel              string
	LoginProviderGoogleText        string
	LoginProviderMicrosoftText     string
	LoginProviderGithubText        string
	LoginProviderWechatText        string
	OAuthCompleteInfoTitle         string
	OAuthCompleteInfoPositionLabel string
	OAuthCompleteInfoAgreeLabel    string
	OAuthCompleteInfoBackLabel     string
	Demo                           string
	DBResetTipLabel                string
	Name                           string
	Email                          string
	Company                        string
	Role                           string
	Status                         string
	ChangePassword                 string
	LoginSessions                  string
	RoleAdmin                      string
	RoleManager                    string
	RoleEditor                     string
	RoleViewer                     string
	RoleUnknown                    string
	PasswordChangePublicUserError  string
	PasswordLengthError            string
}

var Messages_en_US = &Messages{
	FilterTabsAll:                  "All",
	FilterTabsHasUnreadNotes:       "Has Unread Notes",
	FilterTabsActive:               "Active",
	DemoTips:                       "Please note that the database would be reset every even hour.",
	DemoUsernameLabel:              "Demo Username: ",
	DemoPasswordLabel:              "Demo Password: ",
	LoginProviderGoogleText:        "Login with Google",
	LoginProviderMicrosoftText:     "Login with Microsoft",
	LoginProviderGithubText:        "Login with Github",
	LoginProviderWechatText:        "Login with WeChat",
	OAuthCompleteInfoTitle:         "Complete your information",
	OAuthCompleteInfoPositionLabel: "Position(Optional)",
	OAuthCompleteInfoAgreeLabel:    "Subscribe to R0VX newsletter(Optional)",
	OAuthCompleteInfoBackLabel:     "Back to login",
	Demo:                           "DEMO",
	DBResetTipLabel:                "Database reset countdown",
	Name:                           "Name",
	Email:                          "Email",
	Company:                        "Company",
	Role:                           "Role",
	Status:                         "Status",
	ChangePassword:                 "Change Password",
	LoginSessions:                  "Login Sessions",
	RoleAdmin:                      "Admin",
	RoleManager:                    "Manager",
	RoleEditor:                     "Editor",
	RoleViewer:                     "Viewer",
	RoleUnknown:                    "Unknown Role",
	PasswordChangePublicUserError:  "Cannot change password for public user",
	PasswordLengthError:            "Password must be 6 to 20 characters",
}

var Messages_zh_CN = &Messages{
	FilterTabsAll:                  "全部",
	FilterTabsHasUnreadNotes:       "未读备注",
	FilterTabsActive:               "有效",
	DemoTips:                       "请注意，数据库将每隔偶数小时重置一次。",
	DemoUsernameLabel:              "演示账户：",
	DemoPasswordLabel:              "演示密码：",
	LoginProviderGoogleText:        "使用Google登录",
	LoginProviderMicrosoftText:     "使用Microsoft登录",
	LoginProviderGithubText:        "使用Github登录",
	LoginProviderWechatText:        "微信登录",
	OAuthCompleteInfoTitle:         "请填写您的信息",
	OAuthCompleteInfoPositionLabel: "职位（可选）",
	OAuthCompleteInfoAgreeLabel:    "订阅R0VX新闻（可选）",
	OAuthCompleteInfoBackLabel:     "返回登录",
	Demo:                           "演示",
	DBResetTipLabel:                "数据库重置倒计时",
	Name:                           "姓名",
	Email:                          "邮箱",
	Company:                        "公司",
	Role:                           "角色",
	Status:                         "状态",
	ChangePassword:                 "修改密码",
	LoginSessions:                  "登录会话",
	RoleAdmin:                      "管理员",
	RoleManager:                    "经理",
	RoleEditor:                     "编辑员",
	RoleViewer:                     "查看者",
	RoleUnknown:                    "未知角色",
	PasswordChangePublicUserError:  "无法修改公共用户的密码",
	PasswordLengthError:            "密码长度需为 6-20 个字符",
}

type Messages_ModelsI18nModuleKey struct {
	R0VXExample string
	Roles       string
	Users       string

	// action-enhance-demo（模型 label = "WizardDemos"）RowMenuItem i18n（D2 演示）
	WizardDemos                     string // 模型显示名
	WizardDemosUpgradeTooltip       string
	WizardDemosUpgradeConfirmTitle  string
	WizardDemosUpgradeConfirmPrompt string
	WizardDemosResetTooltip         string

	Posts          string
	PostsID        string
	PostsTitle     string
	PostsHeroImage string
	PostsBody      string
	Example        string
	Settings       string
	Post           string
	PostsBodyImage string

	SeoPost             string
	SeoVariableTitle    string
	SeoVariableSiteName string

	PageBuilder              string
	Pages                    string
	SharedContainers         string
	DemoContainers           string
	Templates                string
	PageCategories           string
	ECManagement             string
	ECDashboard              string
	Orders                   string
	InputDemos               string
	Products                 string
	NestedFieldDemos         string
	SiteManagement           string
	SEO                      string
	UserManagement           string
	Profile                  string
	FeaturedModelsManagement string
	Customers                string
	ListModels               string
	MicrositeModels          string
	Workers                  string
	MediaLibrary             string

	// 角色名翻译（role 模块用 i18n.T(ModelsI18nModuleKey, roleKey) 查；未注册角色回退原始名）
	Admin   string
	Manager string
	Editor  string
	Viewer  string

	PagesID         string
	PagesTitle      string
	PagesSlug       string
	PagesLocale     string
	PagesNotes      string
	PagesDraftCount string
	PagesPath       string
	PagesOnline     string
	PagesVersion    string
	PagesVersions   string
	PagesStartAt    string
	PagesEndAt      string
	PagesOption     string
	PagesLive       string

	Page                   string
	PagesStatus            string
	PagesSchedule          string
	PagesCategoryID        string
	PagesTemplateSelection string
	PagesEditContainer     string

	WebHeader       string
	WebHeadersColor string
	Header          string
	Navigation      string
	Content         string

	WebFooter             string
	WebFootersEnglishUrl  string
	WebFootersJapaneseUrl string
	Footer                string

	VideoBanner                       string
	VideoBannersAddTopSpace           string
	VideoBannersAddBottomSpace        string
	VideoBannersAnchorID              string
	VideoBannersVideo                 string
	VideoBannersBackgroundVideo       string
	VideoBannersMobileBackgroundVideo string
	VideoBannersVideoCover            string
	VideoBannersMobileVideoCover      string
	VideoBannersHeading               string
	VideoBannersPopupText             string
	VideoBannersText                  string
	VideoBannersLinkText              string
	VideoBannersLink                  string

	Heading                   string
	HeadingsAddTopSpace       string
	HeadingsAddBottomSpace    string
	HeadingsAnchorID          string
	HeadingsHeading           string
	HeadingsFontColor         string
	HeadingsBackgroundColor   string
	HeadingsLink              string
	HeadingsLinkText          string
	HeadingsLinkDisplayOption string
	HeadingsText              string

	BrandGrid                string
	BrandGridsAddTopSpace    string
	BrandGridsAddBottomSpace string
	BrandGridsAnchorID       string
	BrandGridsBrands         string

	ListContent                   string
	ListContentsAddTopSpace       string
	ListContentsAddBottomSpace    string
	ListContentsAnchorID          string
	ListContentsBackgroundColor   string
	ListContentsItems             string
	ListContentsLink              string
	ListContentsLinkText          string
	ListContentsLinkDisplayOption string

	ImageContainer                           string
	ImageContainersAddTopSpace               string
	ImageContainersAddBottomSpace            string
	ImageContainersAnchorID                  string
	ImageContainersBackgroundColor           string
	ImageContainersTransitionBackgroundColor string
	ImageContainersImage                     string
	Image                                    string

	InNumber                string
	InNumbersAddTopSpace    string
	InNumbersAddBottomSpace string
	InNumbersAnchorID       string
	InNumbersHeading        string
	InNumbersItems          string
	InNumbers               string

	ContactForm                    string
	ContactFormsAddTopSpace        string
	ContactFormsAddBottomSpace     string
	ContactFormsAnchorID           string
	ContactFormsHeading            string
	ContactFormsText               string
	ContactFormsSendButtonText     string
	ContactFormsFormButtonText     string
	ContactFormsMessagePlaceholder string
	ContactFormsNamePlaceholder    string
	ContactFormsEmailPlaceholder   string
	ContactFormsThankyouMessage    string
	ContactFormsActionUrl          string
	ContactFormsPrivacyPolicy      string

	ActivityActionLogIn           string
	ActivityActionExtendSession   string
	ActivityActionAddContainer    string
	ActivityActionDeleteContainer string

	PagesPage string

	// Demo 模块（Dialog / EditingActions / Notif）的 field + action i18n：
	// name 用英文标识符，中文显示在此注册。key = ToCamel(modelLabel + " " + name)。
	DialogDemosTitle          string
	DialogDemosStatus         string
	DialogDemosPriority       string
	DialogDemosNotes          string
	DialogDemosRelatedId      string
	DialogDemosChangeStatus   string
	DialogDemosChangePriority string
	DialogDemosBatchImport    string
	DialogDemosExportData     string
	DialogDemosEditStatus     string
	DialogDemosAddNote        string

	EditingActionsDemosTitle             string
	EditingActionsDemosStatus            string
	EditingActionsDemosContent           string
	EditingActionsDemosQuickChangeStatus string
	EditingActionsDemosAddNote           string

	NotifDemosBulkActivate string
}

var Messages_zh_CN_ModelsI18nModuleKey = &Messages_ModelsI18nModuleKey{
	// action-enhance-demo RowMenuItem i18n（D2 演示）
	WizardDemos:                     "Action 增强演示",
	WizardDemosUpgradeTooltip:       "升级该商户",
	WizardDemosUpgradeConfirmTitle:  "确认升级",
	WizardDemosUpgradeConfirmPrompt: "将该商户升级为已发布状态，继续吗？",
	WizardDemosResetTooltip:         "重置为草稿",

	Posts:          "帖子 示例",
	PostsID:        "ID",
	PostsTitle:     "标题",
	PostsHeroImage: "主图",
	PostsBody:      "内容",
	Example:        "R0VX演示",
	Settings:       "SEO 设置",
	Post:           "帖子",
	PostsBodyImage: "内容图片",

	SeoPost:             "帖子",
	SeoVariableTitle:    "标题",
	SeoVariableSiteName: "站点名称",

	R0VXExample: "R0VX 示例",
	Roles:       "权限管理",
	Users:       "用户管理",

	PageBuilder:              "页面管理菜单",
	Pages:                    "页面管理",
	SharedContainers:         "公用组件",
	DemoContainers:           "示例组件",
	Templates:                "模板页面",
	PageCategories:           "目录管理",
	ECManagement:             "电子商务管理",
	ECDashboard:              "电子商务仪表盘",
	Orders:                   "订单管理",
	InputDemos:               "表单 示例",
	Products:                 "产品管理",
	NestedFieldDemos:         "嵌套表单 示例",
	SiteManagement:           "站点管理菜单",
	SEO:                      "SEO 管理",
	UserManagement:           "用户管理菜单",
	Profile:                  "个人页面",
	FeaturedModelsManagement: "特色模块管理菜单",
	Customers:                "Customers 示例",
	ListModels:               "发布带排序及分页模块 示例",
	MicrositeModels:          "Microsite 示例",
	Workers:                  "后台工作进程管理",
	MediaLibrary:             "媒体库",

	Admin:   "管理员",
	Manager: "经理",
	Editor:  "编辑员",
	Viewer:  "查看者",

	PagesID:         "ID",
	PagesTitle:      "标题",
	PagesSlug:       "Slug",
	PagesLocale:     "地区",
	PagesNotes:      "备注",
	PagesDraftCount: "草稿数",
	PagesPath:       "路径",
	PagesOnline:     "在线",
	PagesVersion:    "版本",
	PagesVersions:   "版本",
	PagesStartAt:    "开始时间",
	PagesEndAt:      "结束时间",
	PagesOption:     "选项",
	PagesLive:       "发布状态",

	Page:                   "Page",
	PagesStatus:            "状态",
	PagesSchedule:          "PagesSchedule",
	PagesCategoryID:        "PagesCategoryID",
	PagesTemplateSelection: "PagesTemplateSelection",
	PagesEditContainer:     "PagesEditContainer",

	WebHeader:       "WebHeader",
	WebHeadersColor: "WebHeadersColor",
	Header:          "Header",
	Navigation:      "Navigation",
	Content:         "Content",

	WebFooter:             "WebFooter",
	WebFootersEnglishUrl:  "WebFootersEnglishUrl",
	WebFootersJapaneseUrl: "WebFootersJapaneseUrl",
	Footer:                "Footer",

	VideoBanner:                       "VideoBanner",
	VideoBannersAddTopSpace:           "VideoBannersAddTopSpace",
	VideoBannersAddBottomSpace:        "VideoBannersAddBottomSpace",
	VideoBannersAnchorID:              "VideoBannersAnchorID",
	VideoBannersVideo:                 "VideoBannersVideo",
	VideoBannersBackgroundVideo:       "VideoBannersBackgroundVideo",
	VideoBannersMobileBackgroundVideo: "VideoBannersMobileBackgroundVideo",
	VideoBannersVideoCover:            "VideoBannersVideoCover",
	VideoBannersMobileVideoCover:      "VideoBannersMobileVideoCover",
	VideoBannersHeading:               "VideoBannersHeading",
	VideoBannersPopupText:             "VideoBannersPopupText",
	VideoBannersText:                  "VideoBannersText",
	VideoBannersLinkText:              "VideoBannersLinkText",
	VideoBannersLink:                  "VideoBannersLink",

	Heading:                   "Heading",
	HeadingsAddTopSpace:       "HeadingsAddTopSpace",
	HeadingsAddBottomSpace:    "HeadingsAddBottomSpace",
	HeadingsAnchorID:          "HeadingsAnchorID",
	HeadingsHeading:           "HeadingsHeading",
	HeadingsFontColor:         "HeadingsFontColor",
	HeadingsBackgroundColor:   "HeadingsBackgroundColor",
	HeadingsLink:              "HeadingsLink",
	HeadingsLinkText:          "HeadingsLinkText",
	HeadingsLinkDisplayOption: "HeadingsLinkDisplayOption",
	HeadingsText:              "HeadingsText",

	BrandGrid:                "BrandGrid",
	BrandGridsAddTopSpace:    "BrandGridsAddTopSpace",
	BrandGridsAddBottomSpace: "BrandGridsAddBottomSpace",
	BrandGridsAnchorID:       "BrandGridsAnchorID",
	BrandGridsBrands:         "BrandGridsBrands",

	ListContent:                   "ListContent",
	ListContentsAddTopSpace:       "ListContentsAddTopSpace",
	ListContentsAddBottomSpace:    "ListContentsAddBottomSpace",
	ListContentsAnchorID:          "ListContentsAnchorID",
	ListContentsBackgroundColor:   "ListContentsBackgroundColor",
	ListContentsItems:             "ListContentsItems",
	ListContentsLink:              "ListContentsLink",
	ListContentsLinkText:          "ListContentsLinkText",
	ListContentsLinkDisplayOption: "ListContentsLinkDisplayOption",

	ImageContainer:                           "ImageContainer",
	ImageContainersAddTopSpace:               "ImageContainersAddTopSpace",
	ImageContainersAddBottomSpace:            "ImageContainersAddBottomSpace",
	ImageContainersAnchorID:                  "ImageContainersAnchorID",
	ImageContainersBackgroundColor:           "ImageContainersBackgroundColor",
	ImageContainersTransitionBackgroundColor: "ImageContainersTransitionBackgroundColor",
	ImageContainersImage:                     "ImageContainersImage",
	Image:                                    "Image",

	InNumber:                "InNumber",
	InNumbersAddTopSpace:    "InNumbersAddTopSpace",
	InNumbersAddBottomSpace: "InNumbersAddBottomSpace",
	InNumbersAnchorID:       "InNumbersAnchorID",
	InNumbersHeading:        "InNumbersHeading",
	InNumbersItems:          "InNumbersItems",
	InNumbers:               "InNumbers",

	ContactForm:                    "ContactForm",
	ContactFormsAddTopSpace:        "ContactFormsAddTopSpace",
	ContactFormsAddBottomSpace:     "ContactFormsAddBottomSpace",
	ContactFormsAnchorID:           "ContactFormsAnchorID",
	ContactFormsHeading:            "ContactFormsHeading",
	ContactFormsText:               "ContactFormsText",
	ContactFormsSendButtonText:     "ContactFormsSendButtonText",
	ContactFormsFormButtonText:     "ContactFormsFormButtonText",
	ContactFormsMessagePlaceholder: "ContactFormsMessagePlaceholder",
	ContactFormsNamePlaceholder:    "ContactFormsNamePlaceholder",
	ContactFormsEmailPlaceholder:   "ContactFormsEmailPlaceholder",
	ContactFormsThankyouMessage:    "ContactFormsThankyouMessage",
	ContactFormsActionUrl:          "ContactFormsActionUrl",
	ContactFormsPrivacyPolicy:      "ContactFormsPrivacyPolicy",

	ActivityActionLogIn:           "登录",
	ActivityActionExtendSession:   "延长会话",
	ActivityActionAddContainer:    "添加容器",
	ActivityActionDeleteContainer: "删除容器",

	PagesPage: "Page",

	DialogDemosTitle:          "标题",
	DialogDemosStatus:         "状态",
	DialogDemosPriority:       "优先级",
	DialogDemosNotes:          "备注",
	DialogDemosRelatedId:      "关联记录",
	DialogDemosChangeStatus:   "批量修改状态",
	DialogDemosChangePriority: "批量修改优先级",
	DialogDemosBatchImport:    "批量导入",
	DialogDemosExportData:     "导出 CSV",
	DialogDemosEditStatus:     "修改状态",
	DialogDemosAddNote:        "添加备注",

	EditingActionsDemosTitle:             "标题",
	EditingActionsDemosStatus:            "状态",
	EditingActionsDemosContent:           "内容",
	EditingActionsDemosQuickChangeStatus: "快速修改状态",
	EditingActionsDemosAddNote:           "添加备注",

	NotifDemosBulkActivate: "批量激活",
}

// GetRoleName 根据角色 key 获取翻译后的角色名称
func (m *Messages) GetRoleName(roleKey string) string {
	switch roleKey {
	case models.RoleAdmin:
		return m.RoleAdmin
	case models.RoleManager:
		return m.RoleManager
	case models.RoleEditor:
		return m.RoleEditor
	case models.RoleViewer:
		return m.RoleViewer
	default:
		return roleKey // 未知角色返回原始名称
	}
}
