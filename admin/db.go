package admin

import (
	"time"

	"github.com/theplant/osenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbParamsString = osenv.Get("DB_PARAMS", "admin example database connection string", "user=docs password=docs dbname=docs sslmode=disable host=localhost port=6532 TimeZone=Asia/Tokyo connect_timeout=300 statement_timeout=300000")

func ConnectDB() (db *gorm.DB) {
	var err error
	db, err = gorm.Open(postgres.Open(dbParamsString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Logger = db.Logger.LogMode(logger.Info)

	// Set database connection pool settings for handling large form data
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// Configure connection pool for large data operations
	// 注意：goque worker 会占用约 10 个连接，需要留足给其他请求使用
	sqlDB.SetMaxOpenConns(100)          // Maximum open connections
	sqlDB.SetMaxIdleConns(100)          // Maximum idle connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Connection max lifetime

	return db
}
