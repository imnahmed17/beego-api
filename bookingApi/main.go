package main

import (
	"bookingApi/db"
	_ "bookingApi/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	db.Connect()

	beego.Run()
}
