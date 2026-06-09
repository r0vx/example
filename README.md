# Example - r0vx 完整示例应用

这是一个基于 r0vx 框架的完整示例应用，展示了如何使用 r0vx Admin、Web 和 X 模块构建功能齐全的企业级后台管理系统。

## 🎯 项目简介

本项目是一个电商管理系统示例，集成了 r0vx 框架的所有核心功能，包括：

- 用户管理和权限控制
- 产品和订单管理
- 媒体文件管理
- 内容发布系统
- 多语言支持
- 操作日志审计
- 后台任务调度
- SEO 管理
- OAuth 认证

## ✨ 功能特性

### 核心功能

- ✅ **用户系统**：用户注册、登录、角色权限管理、多语言角色名称
- ✅ **产品管理**：产品 CRUD、图片上传、发布/下线
- ✅ **订单管理**：订单列表、详情、状态跟踪
- ✅ **分类管理**：产品分类的树形结构管理
- ✅ **媒体库**：统一的文件和图片管理
- ✅ **内容发布**：草稿、发布、版本控制
- ✅ **操作日志**：自动记录所有 CRUD 操作
- ✅ **后台任务**：异步任务和定时任务
- ✅ **多语言**：内容的多语言版本管理、界面国际化
- ✅ **权限控制**：基于角色的细粒度权限控制

### UI 组件

- 使用 **shadcn/ui** 组件库
- 响应式设计
- 暗色模式支持
- 现代化交互体验

## 🚀 快速开始

### 环境要求

- Go 1.25+
- PostgreSQL 14+（或 SQLite 用于开发）
- Node.js 18+ 和 pnpm（可选，用于前端开发）

### 安装依赖

```bash
# 克隆仓库（如果需要）
cd example

# 安装 Go 依赖
go mod download
```

### 配置数据库

创建 `.env` 文件或设置环境变量：

```bash
# PostgreSQL 配置
DB_PARAMS="host=localhost user=postgres password=postgres dbname=example_dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"

# 或使用 SQLite（开发环境）
DB_PARAMS="example.db"

# 服务配置
HOST="0.0.0.0"
PORT="9500"

# 是否重置并导入初始数据
RESET_AND_IMPORT_INITIAL_DATA="true"
```

### 启动服务

```bash
# 开发模式运行
go run main.go

# 或使用脚本
./dev.sh
```

访问 http://localhost:9500 即可看到管理后台。

### 默认账号

初始化数据后，可以使用以下账号登录：

| 角色 | 邮箱 | 密码 | 权限 |
|------|------|------|------|
| 管理员 | admin@example.com | admin123 | 所有权限 |
| 经理 | manager@example.com | manager123 | 大部分权限 |
| 编辑 | editor@example.com | editor123 | 编辑权限 |
| 查看者 | viewer@example.com | viewer123 | 只读权限 |

## 📁 项目结构

```
example/
├── main.go                    # 应用入口
├── admin/                     # Admin 配置
│   ├── config.go             # Admin Builder 配置
│   ├── auth.go               # 认证和会话管理
│   ├── db.go                 # 数据库连接
│   ├── router.go             # 路由配置
│   ├── perm.go               # 权限配置
│   ├── data_init.go          # 初始数据
│   ├── product_config.go     # 产品模型配置
│   ├── user_config.go        # 用户模型配置
│   ├── order_config.go       # 订单模型配置
│   ├── category_config.go    # 分类模型配置
│   ├── input_demo_config.go  # UI 组件示例
│   └── messages.go           # 国际化消息
├── models/                    # 数据模型
│   ├── user.go               # 用户模型
│   ├── product.go            # 产品模型
│   ├── order.go              # 订单模型
│   ├── category.go           # 分类模型
│   ├── post.go               # 文章模型
│   ├── l10n_model.go         # 多语言模型
│   └── list_model.go         # 列表示例模型
├── cmd/                       # 命令行工具
│   ├── seed-users/           # 用户数据生成
│   ├── publisher/            # 发布任务
│   └── data-resetor/         # 数据重置
├── integration/               # 集成测试
│   ├── user_test.go
│   ├── product_test.go
│   ├── order_test.go
│   └── ...
├── pages/                     # 自定义页面
├── public/                    # 静态资源
└── assets/                    # 嵌入资源
```

## 🌍 国际化 (I18n)

Example 项目完整支持界面国际化，包括：

### 角色名称国际化

角色名称会根据用户语言自动翻译：

| 英文 | 中文 | 数据库存储 |
|------|------|-----------|
| Admin | 管理员 | `Admin` |
| Manager | 经理 | `Manager` |
| Editor | 编辑员 | `Editor` |
| Viewer | 查看者 | `Viewer` |

**实现方式**：
1. 数据库存储英文 key
2. 通过 `Messages` 结构定义翻译
3. UI 渲染时使用 `GetRoleName()` 函数自动翻译

**显示位置**：
- 用户列表的"角色"列
- 用户编辑页面的角色选择器
- 任何显示角色的地方

### 添加新的翻译

在 [admin/messages.go](admin/messages.go) 中添加：

```go
// 1. 在 Messages 结构中添加字段
type Messages struct {
    // ...
    NewField string
}

// 2. 添加英文翻译
var Messages_en_US = &Messages{
    // ...
    NewField: "New Field",
}

// 3. 添加中文翻译
var Messages_zh_CN = &Messages{
    // ...
    NewField: "新字段",
}

// 4. 在组件中使用
msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)
label := msgr.NewField
```

### 语言切换

框架会根据以下优先级确定用户语言：
1. 用户偏好设置
2. Cookie 中的语言设置
3. HTTP Accept-Language 头
4. 默认语言（en-US）

## 📚 核心模块说明

### 1. 用户管理 (User)

**文件**: [models/user.go](models/user.go), [admin/user_config.go](admin/user_config.go)

**功能**：
- 用户 CRUD 操作
- 角色分配（Admin、Manager、Editor、Viewer）
- **角色国际化**：支持多语言角色名称显示（中文/英文自动切换）
- 密码管理和重置
- OAuth 登录支持（Google、GitHub、Microsoft）
- 会话管理和 TOTP 双因素认证
- 账户锁定/解锁

**字段**：
- 基本信息：姓名、公司、状态、注册日期
- 认证：邮箱、密码（加密）、TOTP 密钥
- OAuth：提供商、标识符
- 角色：多对多关联（支持 i18n 显示）
- 会话：安全会话管理

**国际化角色名称**：

```go
// 角色在列表和编辑页面自动翻译
// 英文：Admin, Manager, Editor, Viewer
// 中文：管理员, 经理, 编辑员, 查看者

// 实现方式（admin/messages.go）
func (m *Messages) GetRoleName(roleKey string) string {
    switch roleKey {
    case models.RoleAdmin:
        return m.RoleAdmin  // "管理员" 或 "Admin"
    case models.RoleManager:
        return m.RoleManager
    // ...
    }
}
```

**示例代码**：

```go
// 配置用户模型
user := b.Model(&models.User{}).
    Listing("ID", "Name", "Account", "Roles", "Status", "CreatedAt").
    Editing("Name", "Company", "Roles", "Status")

// 角色字段自动使用 i18n 翻译
cl.Field("Roles").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
    msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)

    // 显示翻译后的角色名称
    for _, r := range user.Roles {
        roleNames = append(roleNames, h.Text(msgr.GetRoleName(r.Name)))
    }
    return h.Td(roleNames...)
})

// 角色选择器（编辑页面）
ed.Field("Roles").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
    msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)

    return shadcn.Autocomplete().
        Label(field.Label).
        Multiple(true).
        Items(h.JSONString(lo.Map(roles, func(r role.Role, _ int) map[string]string {
            return map[string]string{
                "label": msgr.GetRoleName(r.Name),  // 翻译后的名称
                "value": fmt.Sprint(r.ID),
            }
        }))).
        Attr(presets.VFieldError(field.Name, values, field.Errors)...)
})
```

**顶部操作按钮**：
- 发送密码重置邮件
- 解锁账户
- 撤销 TOTP

### 2. 产品管理 (Product)

**文件**: [models/product.go](models/product.go), [admin/product_config.go](admin/product_config.go)

**功能**：
- 产品 CRUD
- 图片上传和管理
- 发布/下线控制
- 版本管理
- 定时发布

**集成插件**：
- `media` - 图片上传
- `publish` - 发布系统

**示例代码**：

```go
type Product struct {
    gorm.Model
    Code     string
    Name     string
    Price    int
    Image    media_library.MediaBox `sql:"type:text;"`
    publish.Status
    publish.Schedule
    publish.Version
}

// 配置产品模型
product := b.Model(&Product{}).
    Listing("ID", "Code", "Name", "Price", "Status").
    Editing("Code", "Name", "Price", "Image")
```

### 3. 订单管理 (Order)

**文件**: [models/order.go](models/order.go), [admin/order_config.go](admin/order_config.go)

**功能**：
- 订单列表和详情
- 订单状态管理
- 关联用户和产品
- 地址信息

**示例代码**：

```go
// 自定义订单详情页
order.Detailing().
    AppendTabsPanelFunc(func(obj interface{}, ctx *web.EventContext) (string, h.HTMLComponent) {
        o := obj.(*models.Order)
        return "订单明细", shadcn.Card().
            Title("订单信息").
            Content(renderOrderDetails(o))
    })
```

### 4. 分类管理 (Category)

**文件**: [models/category.go](models/category.go), [admin/category_config.go](admin/category_config.go)

**功能**：
- 树形结构分类
- 父子关系管理
- 拖拽排序

### 5. 媒体库 (Media)

使用 `github.com/r0vx/admin/media` 模块：

```go
// 配置媒体库
mb := media.New(db).
    Storage(filesystem.New("./uploads"))

b.Use(mb)

// 在模型中使用
type Product struct {
    Image media_library.MediaBox `sql:"type:text;"`
}
```

### 6. 操作日志 (Activity)

自动记录所有数据变更：

```go
ab := activity.New(db).
    AutoMigrate()

b.Use(ab)

// 查看日志
// 访问 /admin/activity-logs
```

### 7. 后台任务 (Worker)

异步任务和定时任务：

```go
wb := worker.New(db)

// 添加任务
wb.AddJob("send-email", func(ctx context.Context, job worker.Job) error {
    // 发送邮件逻辑
    return nil
})

// 定时任务
wb.AddCronJob("cleanup", "0 0 * * *", func(ctx context.Context) error {
    // 清理逻辑
    return nil
})

b.Use(wb)
```

### 8. 多语言 (I18n & L10n)

#### 界面国际化 (I18n)

**文件**: [admin/messages.go](admin/messages.go)

支持界面元素的多语言翻译：

```go
// 定义消息结构
type Messages struct {
    FilterTabsAll    string
    FilterTabsActive string
    RoleAdmin        string  // 角色翻译
    RoleManager      string
    RoleEditor       string
    RoleViewer       string
    // ...
}

// 英文翻译
var Messages_en_US = &Messages{
    FilterTabsAll:    "All",
    FilterTabsActive: "Active",
    RoleAdmin:        "Admin",
    RoleManager:      "Manager",
    // ...
}

// 中文翻译
var Messages_zh_CN = &Messages{
    FilterTabsAll:    "全部",
    FilterTabsActive: "有效",
    RoleAdmin:        "管理员",
    RoleManager:      "经理",
    // ...
}

// 在组件中使用
msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)
label := msgr.RoleAdmin  // 根据用户语言返回 "Admin" 或 "管理员"
```

**支持的翻译**：
- 菜单和导航
- 表单标签
- 按钮文本
- 角色名称（Admin → 管理员）
- 状态文本
- 验证错误消息

#### 内容多语言 (L10n)

内容数据的多语言版本管理：

```go
// 定义多语言模型
type Article struct {
    ID      uint
    l10n.Locale
    Title   string
    Content string
}

lb := l10n.New().
    SupportLocales("en", "zh", "ja")

b.Use(lb)
```

### 9. 权限控制 (Permission)

基于角色的访问控制：

```go
// 定义权限
permBuilder := perm.New()
permBuilder.
    Role(models.RoleAdmin).
        Policies(perm.PolicyEverything).
    Role(models.RoleEditor).
        Model(&models.Product{}).
        Perm(perm.ActionRead, perm.ActionUpdate).
    Role(models.RoleViewer).
        Model(&models.Product{}).
        Perm(perm.ActionRead)

b.Permission(permBuilder)
```

## 🎨 UI 组件示例

**文件**: [admin/input_demo_config.go](admin/input_demo_config.go)

项目包含完整的 shadcn/ui 组件使用示例，访问 `/admin/input-demos` 可查看：

### 基础组件
- Input、Textarea - 文本输入
- Select、MultiSelect - 下拉选择
- Autocomplete - 自动完成（支持多选，用于角色选择）
- Checkbox、Switch - 复选框和开关
- Radio - 单选按钮

### 日期时间
- DatePicker - 日期选择
- DateRangePicker - 日期范围
- DatetimeRangePicker - 日期时间范围

### 媒体上传
- FileUpload - 文件上传
- ImageCropper - 图片裁剪

### 富文本
- RichTextEditor - 富文本编辑器
- ColorPicker - 颜色选择器

### 数据展示
- DataTable - 数据表格（带分页、排序、筛选）
- Badge - 徽章（用于状态显示）
- Card - 卡片容器

### 交互组件
- Dialog - 对话框
- Sheet - 侧边抽屉
- Alert - 警告提示
- Toast - 消息通知

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test -v ./integration -run TestUser

# 运行集成测试
cd integration && go test -v
```

### 测试覆盖

项目包含完整的集成测试：

- `user_test.go` - 用户 CRUD 测试
- `product_test.go` - 产品管理测试
- `order_test.go` - 订单管理测试
- `login_test.go` - 登录认证测试
- `activity_logs_test.go` - 操作日志测试
- `worker_test.go` - 后台任务测试

## 🔧 开发指南

### 添加新模型

1. **定义模型** (`models/new_model.go`)：

```go
package models

type NewModel struct {
    gorm.Model
    Name        string
    Description string
}
```

2. **配置 Admin** (`admin/new_model_config.go`)：

```go
func configNewModel(b *presets.Builder, db *gorm.DB) {
    model := b.Model(&models.NewModel{})

    model.Listing("ID", "Name", "CreatedAt")
    model.Editing("Name", "Description")
}
```

3. **注册模型** (`admin/config.go`)：

```go
func NewConfig(db *gorm.DB) Config {
    // ...
    configNewModel(b, db)
    // ...
}
```

### 自定义字段组件

使用 shadcn/ui 组件自定义字段渲染：

```go
model.Field("Status").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
    return shadcn.Select().
        Label("状态").
        Options([]*shadcn.SelectOption{
            {Value: "active", Label: "激活"},
            {Value: "inactive", Label: "停用"},
        }).
        Attr(web.VField(field.FormKey, field.Value(obj))...)
})
```

### 添加自定义页面

```go
// 在 admin/router.go 中添加
b.GetWebBuilder().Page("/custom-page", func(ctx *web.EventContext) (pr web.PageResponse, err error) {
    pr.Body = Div(
        H1("自定义页面"),
        shadcn.Card().
            Title("数据统计").
            Content(/* ... */),
    )
    return
})

// 添加到菜单
b.MenuGroup("报表").
    Item("/custom-page").
    Label("自定义报表").
    Icon("chart")
```

## 🛠️ 命令行工具

### 数据初始化

```bash
# 重置数据库并导入初始数据
RESET_AND_IMPORT_INITIAL_DATA=true go run main.go

# 或使用数据重置工具
go run cmd/data-resetor/main.go
```

### 生成用户数据

```bash
# 生成测试用户
go run cmd/seed-users/main.go --count 100
```

### 发布任务

```bash
# 手动触发发布
go run cmd/publisher/main.go
```

## 📊 性能监控

项目集成了 pprof 性能分析工具：

- 主服务: http://localhost:9500
- pprof 分析: http://localhost:6060/debug/pprof/

可用的 pprof 端点：
- `/debug/pprof/` - 概览
- `/debug/pprof/heap` - 内存堆分析
- `/debug/pprof/goroutine` - 协程分析
- `/debug/pprof/profile` - CPU 分析

## 🚢 部署

### 构建

```bash
# 构建二进制文件
go build -o example main.go

# 交叉编译（Linux）
GOOS=linux GOARCH=amd64 go build -o example-linux main.go
```

### Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o example main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/example .
COPY --from=builder /app/public ./public

EXPOSE 9500
CMD ["./example"]
```

构建和运行：

```bash
docker build -t example-app .
docker run -p 9500:9500 \
  -e DB_PARAMS="host=db user=postgres password=postgres dbname=example" \
  example-app
```

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `HOST` | 服务监听地址 | `127.0.0.1` |
| `PORT` | 服务端口 | `9500` |
| `DB_PARAMS` | 数据库连接字符串 | - |
| `RESET_AND_IMPORT_INITIAL_DATA` | 是否重置数据 | `false` |
| `PUBLISH_URL` | 发布 URL | - |

## 📖 学习资源

本项目是学习 r0vx 框架的最佳起点，推荐阅读顺序：

1. [main.go](main.go) - 了解应用入口和服务启动
2. [admin/config.go](admin/config.go) - 了解 Admin 配置
3. [admin/user_config.go](admin/user_config.go) - 学习模型配置
4. [admin/product_config.go](admin/product_config.go) - 学习发布系统集成
5. [admin/input_demo_config.go](admin/input_demo_config.go) - 学习 UI 组件使用
6. [integration/](integration/) - 学习集成测试

## 🔗 相关链接

- [r0vx 框架文档](https://docs.r0vx.com)
- [r0vx/admin](../admin) - Admin 模块
- [r0vx/web](../web) - Web 模块
- [r0vx/x](../x) - UI 组件和扩展
- [shadcn/ui](https://ui.shadcn.com/) - UI 组件库

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📝 许可证

查看根目录的 LICENSE 文件。

---

*这是一个完整的、生产就绪的示例应用，展示了 r0vx 框架的最佳实践。*
