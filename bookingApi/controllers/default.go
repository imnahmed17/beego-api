package controllers

import (
	// "bookingApi/main"
	"bookingApi/db"
	"bookingApi/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	// "gorm.io/gorm"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	location := c.GetString("location")
	checkIn := c.GetString("t-start")
	checkOut := c.GetString("t-end")
	// db := main.Db 
	dbIns := db.Db

	if location == "" || checkIn == "" || checkOut == "" {
		c.Data["Error"] = "Please Fill the all Required Field"
	} else {
		url := "https://booking-com13.p.rapidapi.com/stays/properties/list-v2" +
		"?location=" + location + 
		"&checkin_date=" + checkIn + 
		"&checkout_date=" + checkOut + 
		"&language_code=en-us&currency_code=USD"

		fmt.Println(url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.Data["Error"] = "Error creating request"
			return
		}

		req.Header.Add("X-RapidAPI-Key", "04d45596a9mshafcf88d1434dc85p1fc8acjsnc24ebc76b973")
		req.Header.Add("X-RapidAPI-Host", "booking-com13.p.rapidapi.com")

		hotelDataChan := make(chan models.HotelData)

		go func() {
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				c.Data["Error"] = "Error making the request"
				hotelDataChan <- models.HotelData{}
				return
			}
			
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				c.Data["Error"] = "Error reading the response"
				hotelDataChan <- models.HotelData{}
				return
			}
			
			var allHotels models.HotelData
			if err = json.Unmarshal(body, &allHotels); err != nil {
				c.Data["Error"] = "Error parsing JSON response"
				hotelDataChan <- models.HotelData{}
				return
			}

			hotelDataChan <- allHotels
		}()

		extractedData := <- hotelDataChan
		hotels :=  extractedData.Data
		c.Data["Hotels"] = hotels
		fmt.Println(len(hotels))

		// for _, info := range hotels {
		// 	fmt.Println("Title:", info.DisplayName.Text)
		// }

		if len(hotels) > 1 {
			hotel_location := models.Hotel_Locations {
				LocationName: location,
			}

			dbIns.Create(&hotel_location)
			// for _, info := range hotels {
			// 	fmt.Println("Title:", info.DisplayName.Text)
			// }
		}
	}

	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
