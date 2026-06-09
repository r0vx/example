package admin

import (
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"example/models"

	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
)

type DataTableHeader struct {
	Text     string `json:"text"`
	Value    string `json:"value"`
	Width    string `json:"width"`
	Sortable bool   `json:"sortable"`
}

func getStringHash(v string, len int) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(v)))[:len]
}

func ip(r *http.Request) string {
	if r == nil {
		return ""
	}

	ips := proxy(r)
	if len(ips) > 0 && ips[0] != "" {
		rip, _, err := net.SplitHostPort(ips[0])
		if err != nil {
			rip = ips[0]
		}
		return rip
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}

func proxy(r *http.Request) []string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}

	return nil
}

func exportOrders(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var orders []*models.Order

		if err := db.Model(&models.Order{}).Limit(10000).Find(&orders).Error; err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		name := time.Now().Format("20060102150405")
		w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="orders-%v.csv"`, name))
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")

		if err := gocsv.Marshal(orders, w); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})
}
