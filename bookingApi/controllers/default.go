package controllers

import (
	"bookingApi/db"
	"bookingApi/models"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

var dbIns *gorm.DB
var rapidApiKey string
var rapidApiHost string
var CheckIn string
var CheckOut string
var hotelDescription string
var newLocation models.Hotel_Locations
var existingLocation models.Hotel_Locations
var existingHotel models.Hotel_Lists
var hotels []models.HotelData

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	dbIns = db.Db
	rapidApiKey, _ = beego.AppConfig.String("rapidapikey")
	rapidApiHost, _ = beego.AppConfig.String("rapidapihost")

	location := strings.ToLower(c.GetString("location"))
	checkIn := c.GetString("t-start")
	CheckIn = checkIn
	checkOut := c.GetString("t-end")
	CheckOut = checkOut
	page := c.GetString("page")

	// check whether the input fields are empty or not
	if location == "" || checkIn == "" || checkOut == "" || page == "" {
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
			fmt.Println("Error creating request")
			return
		}

		req.Header.Add("X-RapidAPI-Key", rapidApiKey)
		req.Header.Add("X-RapidAPI-Host", rapidApiHost)

		hotelDataChan := make(chan []models.HotelData)

		go func() {
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error making the request")
				hotelDataChan <- []models.HotelData{}
				return
			}
			fmt.Println(res)
			
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading the response")
				hotelDataChan <- []models.HotelData{}
				return
			}
			
			var allHotels struct {
				Data []models.HotelData `json:"data"`
			}

			if err = json.Unmarshal(body, &allHotels); err != nil {
				fmt.Println("Error parsing JSON response")
				hotelDataChan <- []models.HotelData{}
				return
			}

			hotelDataChan <- allHotels.Data
		}()

		extractedData := <- hotelDataChan
		hotels =  extractedData

		if len(hotels) > 1 {
			newLocation = models.Hotel_Locations {
				LocationName: location,
			}

			// check whether the entered location is already exists to the hotel_locations table or not
			if err := dbIns.Where("Location_name = ?", location).First(&existingLocation).Error; err != nil {
				if err := dbIns.Create(&newLocation).Error; err != nil {
					fmt.Println("Error creating new location:", err)
				} else {
					fmt.Println("Location added to the database")
					InsertHotelData(hotels, newLocation)
				}
			} else {
				fmt.Println("Location already exists to the database")
				InsertHotelData(hotels, existingLocation)
			}
		}
	}

	c.TplName = "index.tpl"
}

func InsertHotelData (hotels []models.HotelData, location models.Hotel_Locations) {
	for _, hotel := range hotels {
		fmt.Println(hotel.IdDetail)

		newPrice, err := strconv.ParseFloat(hotel.PriceDisplayInfoIrene.DisplayPrice.AmountPerStay.AmountRounded[1:], 64)
		if err != nil {
			fmt.Println("Error type convertion: ", err.Error())
			return
		}

		newHotel := models.Hotel_Lists {
			HotelID: hotel.IdDetail,
			HotelName: hotel.DisplayName.Text,
			HotelCity: hotel.BasicPropertyData.Location.City,
			HotelImageUrl: "https://cf.bstatic.com" + hotel.BasicPropertyData.Photos.Main.HighResUrl.RelativeUrl,
			HotelPrice: newPrice,
			LocationID: location.LocationID,
		}

		// fetch hotel details information
		hotelDetailsInfo, err := GetHotelDetails(hotel.IdDetail)
		if err != nil {
			fmt.Println("Error fetching hotel details: ", err.Error())
			return
		}
		hotelDetails := hotelDetailsInfo.Data

		// fetch hotel details photo
		hotelDetailsPhotos, err := GetHotelPhotos(hotel.IdDetail)
		if err != nil {
			fmt.Println("Error fetching hotel details: ", err.Error())
			return
		}
		hotelPhotos := hotelDetailsPhotos.Data.Photos

		amenities := []string{}
		for _, amenity := range hotelDetails.GenericFacilityHighlight {
			amenities = append(amenities, amenity.Title)
		}

		if len(hotelDetails.HotelTranslation) > 0 {
			hotelDescription = hotelDetails.HotelTranslation[0].Description
		} else {
			hotelDescription = ""
		}

		photos := []string{}
		for _, photo := range hotelPhotos {
			photos = append(photos, "https://cf.bstatic.com" + photo.PhotoUri)
		}

		propertyTypes := []string{"Apartment", "Cottage", "Hotel", "House", "Kabin", "Resort", "Villa"}
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(propertyTypes))
		
		newHotelDetails := models.Hotel_Details {
			HotelID: hotel.IdDetail,
			HotelReviewCount: hotel.BasicPropertyData.Reviews.ReviewsCount,
			HotelRating: hotel.BasicPropertyData.StarRating.Value,
			HotelNoOfBed: hotel.MatchingUnitConfigurations.CommonConfiguration.NbAllBeds,
			HotelAmenities: pq.StringArray(amenities),
			HotelDescription: hotelDescription,
			HotelImageUrls: pq.StringArray(photos),
			HotelPropertyType: propertyTypes[randomIndex],
		}

		// check whether the entered hotel id is already exists to the hotel_lists table or not
		if err := dbIns.Where("Hotel_id = ?", newHotel.HotelID).First(&existingHotel).Error; err != nil {
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
	// fmt.Println(url)

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
