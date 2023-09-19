package controllers

import (
	"bookingApi/db"
	"bookingApi/models"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"io"
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

var CheckIn string
var CheckOut string
var rapidApiKey string
var rapidApiHost string
var hotelDescription string

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	dbIns := db.Db
	rapidApiKey, _ = beego.AppConfig.String("rapidapikey")
	rapidApiHost, _ = beego.AppConfig.String("rapidapihost")

	location := strings.ToLower(c.GetString("location"))
	checkIn := c.GetString("t-start")
	CheckIn = checkIn
	checkOut := c.GetString("t-end")
	CheckOut = checkOut
	page := c.GetString("page")

	if location == "" || checkIn == "" || checkOut == "" || page == "" {
		c.Data["Error"] = "Please fill all the required fields"
		fmt.Println("Please fill all the required fields")
	} else {
		url := "https://booking-com13.p.rapidapi.com/stays/properties/list-v2" +
		"?location=" + location + 
		"&checkin_date=" + checkIn + 
		"&checkout_date=" + checkOut + 
		"&language_code=en-us&currency_code=USD" +
		"&page=" + page

		// fmt.Println(url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.Data["Error"] = "Error creating request"
			return
		}

		req.Header.Add("X-RapidAPI-Key", rapidApiKey)
		req.Header.Add("X-RapidAPI-Host", rapidApiHost)

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

		if len(hotels) > 1 {
			var existingLocation models.Hotel_Locations
			var existingHotel models.Hotel_Lists

			newLocation := models.Hotel_Locations {
				LocationName: location,
			}

			err := dbIns.Where("Location_name = ?", location).First(&existingLocation).Error
			if err != nil {
				if err := dbIns.Create(&newLocation).Error; err != nil {
					fmt.Println("Error creating location:", err)
				} else {
					fmt.Println("Location added to the database")
				}
			} else {
				fmt.Println("Location already exists to the database")
			}
			
			for _, hotel := range hotels {
				hotelDetailsInfo, err := GetHotelDetails(hotel.IdDetail)
				if err != nil {
					c.Data["Error"] = "Error fetching hotel details: " + err.Error()
					return
				}
				hotelDetails := hotelDetailsInfo.Data
				c.Data["HotelDetails"] = hotelDetails

				amenities := []string{}
				for _, amenity := range hotelDetails.GenericFacilityHighlight {
					amenities = append(amenities, amenity.Title)
				}
				hotelAmenities := pq.StringArray(amenities)

				hotelDetailsPhotos, err := GetHotelPhotos(hotel.IdDetail)
				if err != nil {
					c.Data["Error"] = "Error fetching hotel details: " + err.Error()
					return
				}
				hotelPhotos := hotelDetailsPhotos.Data.Photos
				c.Data["HotelPhotos"] = hotelPhotos

				photos := []string{}
				for _, photo := range hotelPhotos {
					photos = append(photos, "https://cf.bstatic.com" + photo.PhotoUri)
				}
				hotelImageUrls := pq.StringArray(photos)

				if len(hotelDetails.HotelTranslation) > 0 {
					hotelDescription = hotelDetails.HotelTranslation[0].Description
				}

				newHotel := models.Hotel_Lists {
					HotelID: hotel.IdDetail,
					HotelName: hotel.DisplayName.Text,
					HotelCity: hotel.BasicPropertyData.Location.City,
					HotelImageUrl: "https://cf.bstatic.com" + hotel.BasicPropertyData.Photos.Main.HighResUrl.RelativeUrl,
					HotelPrice: hotel.PriceDisplayInfoIrene.DisplayPrice.AmountPerStay.AmountRounded,
					LocationID: existingLocation.LocationID,
				}

				fmt.Println(hotel.IdDetail)
				newHotelDetails := models.Hotel_Details {
					HotelID: hotel.IdDetail,
					HotelReviewCount: hotel.BasicPropertyData.Reviews.ReviewsCount,
					HotelRating: hotel.BasicPropertyData.StarRating.Value,
					HotelNoOfBed: hotel.MatchingUnitConfigurations.CommonConfiguration.NbAllBeds,
					HotelAmenities: hotelAmenities,
					HotelDescription: hotelDescription,
					HotelImageUrls: hotelImageUrls,
				}

				if err := dbIns.Where("Hotel_name = ?", newHotel.HotelName).First(&existingHotel).Error; err != nil {
					if err := dbIns.Create(&newHotel).Error; err != nil {
						fmt.Println("Error creating hotel list:", err)
					} else {
						fmt.Println("Hotel list added to the database")

						if err := dbIns.Create(&newHotelDetails).Error; err != nil {
							fmt.Println("Error creating hotel details:", err)
						} else {
							fmt.Println("Hotel details added to the database")
						}
					}
				} else {
					fmt.Println("Hotel list already exists to the database")
				}
			}
		}
	}

	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func GetHotelDetails(id string) (models.HotelDetails, error) {
    url := "https://booking-com13.p.rapidapi.com/stays/properties/detail" +
        "?id_detail=" + id +
        "&checkin_date=" + CheckIn +
        "&checkout_date=" + CheckOut +
        "&language_code=en-us&currency_code=USD"

	// fmt.Println(url)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return models.HotelDetails{}, err
    }

    req.Header.Add("X-RapidAPI-Key", rapidApiKey)
    req.Header.Add("X-RapidAPI-Host", rapidApiHost)

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return models.HotelDetails{}, err
    }

    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return models.HotelDetails{}, err
    }

    var hotelDetails models.HotelDetails
    if err = json.Unmarshal(body, &hotelDetails); err != nil {
        return models.HotelDetails{}, err
    }

    return hotelDetails, nil
}

func GetHotelPhotos(id string) (models.HotelPhotos, error) {
	url := "https://booking-com13.p.rapidapi.com/stays/properties/detail/photos" +
		"?id_detail=" + id +
		"&language_code=en-us"

	fmt.Println(url)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return models.HotelPhotos{}, err
    }

    req.Header.Add("X-RapidAPI-Key", rapidApiKey)
    req.Header.Add("X-RapidAPI-Host", rapidApiHost)

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return models.HotelPhotos{}, err
    }

    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return models.HotelPhotos{}, err
    }

    var hotelPhotos models.HotelPhotos
    if err = json.Unmarshal(body, &hotelPhotos); err != nil {
        return models.HotelPhotos{}, err
    }

    return hotelPhotos, nil
}
