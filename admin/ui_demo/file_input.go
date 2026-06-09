package ui_demo

// ============================================================================
// FileInput 文件上传演示（模拟进度 + 真实 XHR 上传）
// ============================================================================
//
// ## API 用法
//
//	// 模拟进度（纯 UI）
//	shadcn.FileInput().
//		ShowProgress(true).
//		Attr(":upload-progress", "locals.progress")
//
//	// 真实上传（内置 XHR）
//	shadcn.FileInput().
//		UploadURL("/file-input-upload").
//		AutoUpload(false).
//		Multiple(true)
//
// ============================================================================

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/r0vx/admin/presets"
	h "github.com/r0vx/htmlgo"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
)

const FileInputUploadPath = "/file-input-upload"

// FileInputDemo 虚拟模型（无数据库）
type FileInputDemo struct{}

// emptyFileInputSearchFunc 返回空数据，避免数据库查询
func emptyFileInputSearchFunc() presets.SearchFunc {
	return func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
		totalCount := 0
		result = &presets.SearchResult{
			Nodes:      []FileInputDemo{},
			TotalCount: &totalCount,
		}
		return
	}
}

// configFileInputDemo 配置 FileInput 文件上传演示页面（纯 UI，无数据库）
func ConfigFileInputDemo(b *presets.Builder) {
	m := b.Model(&FileInputDemo{}).
		Label("File Input").
		URIName("file-input-demo")
	m.Listing().SearchFunc(emptyFileInputSearchFunc())
	m.Editing().Only()

	// 注册模拟上传事件（通过 RunScript 在全局作用域执行，绕过 Vue 模板沙箱）
	b.GetWebBuilder().RegisterEventFunc("fileInputSimulateUpload", fileInputSimulateUpload)
	b.GetWebBuilder().RegisterEventFunc("fileInputResetUpload", fileInputResetUpload)

	m.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.PageTitle = "FileInput Demo"
		r.Body = fileInputDemoBody()
		return
	})
}

// FileInputUploadHandler 处理 FileInput 组件的文件上传（示例端点）
func FileInputUploadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"error":   "解析表单失败: " + err.Error(),
			})
			return
		}

		files := r.MultipartForm.File["NewFiles"]
		results := make([]map[string]any, 0, len(files))
		for _, fh := range files {
			results = append(results, map[string]any{
				"name": fh.Filename,
				"size": fh.Size,
				"type": fh.Header.Get("Content-Type"),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"files":   results,
			"message": fmt.Sprintf("成功接收 %d 个文件", len(files)),
		})
	}
}

// fileInputSimulateUpload 模拟上传进度（通过 RunScript 执行定时器）
func fileInputSimulateUpload(ctx *web.EventContext) (er web.EventResponse, err error) {
	er.RunScript = `
		locals.uploading = true;
		locals.progress = 0;
		if (locals.timer) clearInterval(locals.timer);
		locals.timer = setInterval(() => {
			locals.progress += Math.random() * 15;
			if (locals.progress >= 100) {
				locals.progress = 100;
				clearInterval(locals.timer);
				locals.timer = null;
				locals.uploading = false;
			}
		}, 300);
	`
	return
}

// fileInputResetUpload 重置上传状态
func fileInputResetUpload(ctx *web.EventContext) (er web.EventResponse, err error) {
	er.RunScript = `
		if (locals.timer) { clearInterval(locals.timer); locals.timer = null; }
		locals.progress = 0;
		locals.uploading = false;
		locals.hasFiles = false;
	`
	return
}

// fileInputDemoBody FileInput 演示页面主体
func fileInputDemoBody() h.HTMLComponent {
	return h.Div(
		h.H1("FileInput Upload Demo").Style("margin-bottom: 24px;"),

		// ── 模拟进度 ──
		h.Div(
			h.H2("模拟进度条"),
			h.P(h.Text("通过 ShowProgress + uploadProgress 控制进度（纯 UI，无真实上传）")).Class("text-muted-foreground mb-4"),
			web.Scope(
				shadcn.FileInput().
					Label("模拟上传").
					ShowProgress(true).
					Placeholder("选择文件后点击上传按钮").
					Attr(":upload-progress", "locals.progress").
					On("change", `locals.hasFiles = true; locals.progress = 0;`),
				h.Div(
					shadcn.Button(h.Text("模拟上传")).
						Size(shadcn.ButtonSizeSm).
						Attr(":disabled", "!locals.hasFiles || locals.uploading").
						On("click", web.Plaid().EventFunc("fileInputSimulateUpload").Go()),
					shadcn.Button(h.Text("重置")).
						Variant(shadcn.ButtonVariantOutline).
						Size(shadcn.ButtonSizeSm).
						On("click", web.Plaid().EventFunc("fileInputResetUpload").Go()),
				).Class("flex items-center gap-2 mt-3"),
			).VSlot("{ locals }").Init(`{ progress: 0, timer: null, hasFiles: false, uploading: false }`),
		).Class("demo-section"),

		// ── 真实 XHR 上传（手动触发） ──
		h.Div(
			h.H2("真实 XHR 上传"),
			h.P(h.Text("设置 UploadURL 后，组件通过 XHR 发送文件并追踪真实进度")).Class("text-muted-foreground mb-4"),
			web.Scope(
				shadcn.FileInput().
					Label("XHR 上传").
					Multiple(true).
					UploadURL(FileInputUploadPath).
					Placeholder("选择文件后点击上传").
					Accept("image/*,.pdf,.doc,.docx,.txt").
					MaxSize(10*1024*1024).
					Attr("ref", "xhrUploader").
					On("upload-success", `locals.result = JSON.stringify($event, null, 2)`).
					On("upload-error", `locals.result = '上传失败: ' + $event`),
				h.Div(
					shadcn.Button(h.Text("开始上传")).
						Size(shadcn.ButtonSizeSm).
						Attr("@click", "$refs.xhrUploader.startUpload()"),
				).Class("flex items-center gap-2 mt-3"),
				h.Pre("").
					Attr("v-if", "locals.result").
					Attr("v-text", "locals.result").
					Class("mt-3 p-3 bg-muted rounded text-xs"),
			).VSlot("{ locals }").Init(`{ result: '' }`),
		).Class("demo-section"),

		// ── 自动上传 ──
		h.Div(
			h.H2("自动上传"),
			h.P(h.Text("设置 AutoUpload 后，选择文件即自动上传")).Class("text-muted-foreground mb-4"),
			web.Scope(
				shadcn.FileInput().
					Label("自动上传").
					UploadURL(FileInputUploadPath).
					AutoUpload(true).
					Placeholder("拖拽或选择文件，自动开始上传").
					On("upload-success", `locals.result = '上传成功!'`).
					On("upload-error", `locals.result = '上传失败: ' + $event`),
				h.Span("").
					Attr("v-if", "locals.result").
					Attr("v-text", "locals.result").
					Class("text-sm text-green-600 mt-2 block"),
			).VSlot("{ locals }").Init(`{ result: '' }`),
		).Class("demo-section"),
	).Style("max-width: 600px; margin: 0 auto;")
}
