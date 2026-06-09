package main

import (
	"log"
	"net/http"

	"example/presets/examples"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/theplant/osenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var dbParamsString = osenv.Get("DB_PARAMS", "presets example database connection string", "")

func main() {
	db, err := gorm.Open(postgres.Open(dbParamsString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Logger.LogMode(gormLogger.Info)

	p := examples.Preset1(db)

	mux := http.NewServeMux()
	mux.Handle("/",
		middleware.RequestID(
			middleware.Logger(
				middleware.Recoverer(p),
			),
		),
	)

	log.Println("serving on :7001")
	log.Fatal(http.ListenAndServe(":7001", mux))
}
