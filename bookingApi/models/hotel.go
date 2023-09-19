package models

import (
	"github.com/lib/pq"
)

type Hotel_Locations struct {
	LocationID   int 	`gorm:"primaryKey;autoIncrement:true"`
    LocationName string 	  
}

type Hotel_Lists struct {
	HotelID 		string 			`gorm:"primaryKey"`
    HotelName 		string 			`gorm:"unique"`
	HotelCity 		string
	HotelImageUrl 	string 
	HotelPrice 		string 
	LocationID   	int 			
	Hotel_Locations Hotel_Locations `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Hotel_Details struct {
	HotelID 	 	 string  		`gorm:"primaryKey"`
	HotelReviewCount int 	 	
	HotelRating 	 int 	 	
	HotelNoOfBed   	 int      	
	HotelAmenities   pq.StringArray `gorm:"type:text[]"`
	HotelDescription string   	
	HotelImageUrls   pq.StringArray `gorm:"type:text[]"`
}

type HotelData struct {
	Data []struct {
		BasicPropertyData struct {
			Location struct {
				City 	string `json:"city"`
			} `json:"location"`
			Photos struct {
				Main struct {
					HighResUrl struct {
						RelativeUrl string `json:"relativeUrl"`
					} `json:"highResUrl"`
				} `json:"main"`
			} `json:"photos"`
			Reviews struct {
				ReviewsCount int `json:"reviewsCount"`
			} `json:"reviews"`
			StarRating struct {
				Value int `json:"value"`
			} `json:"starRating"`
		} `json:"basicPropertyData"`
		DisplayName struct {
			Text string `json:"text"`
		} `json:"displayName"`
		IdDetail string `json:"idDetail"`
		MatchingUnitConfigurations struct {
			CommonConfiguration struct {
				NbAllBeds int `json:"nbAllBeds"`
			} `json:"commonConfiguration"`
		} `json:"matchingUnitConfigurations"`
		PriceDisplayInfoIrene struct {
			DisplayPrice struct {
				AmountPerStay struct {
					AmountRounded string `json:"amountRounded"`
				} `json:"amountPerStay"`
			} `json:"displayPrice"`
		} `json:"priceDisplayInfoIrene"`
	} `json:"data"`
}

type HotelDetails struct {
	Data struct {
		GenericFacilityHighlight []struct {
			Title string `json:"title"`
		} `json:"GenericFacilityHighlight"`
		HotelTranslation []struct {
			Description string `json:"description"`
		} `json:"HotelTranslation"`
	} `json:"data"`
}

type HotelPhotos struct {
	Data struct {
		Photos []struct {
			PhotoUri string `json:"photoUri"`
		} `json:"photos"`
	} `json:"data"`
}
