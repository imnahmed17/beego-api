package main

import (
	"bookingApi/db"
	// "bookingApi/models"
	_ "bookingApi/routers"
	// "fmt"

	// "github.com/lib/pq"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"

	beego "github.com/beego/beego/v2/server/web"
	// "gorm.io/gorm/logger"
)

// var Db *gorm.DB
// var err error

func main() {
	db.Connect()

	beego.Run()
}
