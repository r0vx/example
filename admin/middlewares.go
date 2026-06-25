package admin

import (
	"net/http"

	"github.com/r0vx/x/perm"
	"gorm.io/gorm"
)

func withRoles(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := getCurrentUser(r)
			if u == nil {
				next.ServeHTTP(w, r)
				return
			}

			var roleIDs []uint
			if err := db.Table("user_role_join").Select("role_id").Where("user_id=?", u.ID).Scan(&roleIDs).Error; err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if len(roleIDs) > 0 {
				var roles []perm.Role
				if err := db.Where("id in (?)", roleIDs).Find(&roles).Error; err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				u.Roles = roles
			}
			next.ServeHTTP(w, r)
		})
	}
}

func securityMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			w.Header().Add("Cache-control", "no-cache, no-store, max-age=0, must-revalidate")
			w.Header().Add("Pragma", "no-cache")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			// SAMEORIGIN（非 DENY）：仍防跨源点击劫持，但放行同源 iframe 内嵌——
			// 主题编辑器 /theme-editor 的预览用 <iframe> 同源加载 admin，DENY 会拦死。
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")

			next.ServeHTTP(w, req)
		})
	}
}
