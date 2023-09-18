package db

import (
	"bookingApi/models"
	// _ "bookingApi/routers"
	"fmt"
	// "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	beego "github.com/beego/beego/v2/server/web"
	// "gorm.io/gorm/logger"
)

var Db *gorm.DB
var err error

func Connect() {
	dbHost, _ := beego.AppConfig.String("dbhost")
	dbPort, _ := beego.AppConfig.String("dbport")
	dbUser, _ := beego.AppConfig.String("dbuser")
	dbPass, _ := beego.AppConfig.String("dbpass")
	dbName, _ := beego.AppConfig.String("dbname")
	dbSslMode, _ := beego.AppConfig.String("dbsslmode")
	dbTimeZone, _ := beego.AppConfig.String("dbtimezone")

	dsn := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " sslmode=" + dbSslMode + " TimeZone=" + dbTimeZone
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config {
		// Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		fmt.Println("Database connection failed: ", err)
		return
	}

	Db.AutoMigrate(&models.Hotel_Locations{}, &models.Hotel_Lists{}, &models.Hotel_Details{})
}