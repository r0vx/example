package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // 添加 pprof 支持

	"example/admin"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/theplant/osenv"
)

func main() {
	h := admin.Router(admin.ConnectDB())

	host := osenv.Get("HOST", "The host to serve the admin on", "127.0.0.1")
	port := osenv.Get("PORT", "The port to serve the admin on", "9500")
	addr := host + ":" + port

	fmt.Println("Served at http://" + addr)

	mux := http.NewServeMux()
	// 添加 pprof 端点用于调试（不经过应用中间件）
	mux.Handle("/debug/", http.DefaultServeMux)
	mux.Handle("/",
		middleware.RequestID(
			middleware.Logger(
				middleware.Recoverer(h),
			),
		),
	)

	// 同时在另一个端口启动 pprof 调试服务
	go func() {
		fmt.Println("pprof available at http://localhost:6060/debug/pprof/")
		http.ListenAndServe(":6060", nil)
	}()

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
