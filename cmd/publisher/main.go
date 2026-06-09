package main

import (
	"context"

	"example/admin"

	"github.com/r0vx/admin/publish"
)

func main() {
	db := admin.ConnectDB()
	config := admin.NewConfig(db, false)
	storage := admin.PublishStorage
	publish.RunPublisher(context.Background(), db, storage, config.Publisher)
	select {}
}
