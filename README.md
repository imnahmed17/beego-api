# Beego Booking App API

The Beego Booking App API is a way to store data from RESTful web service to PostgeSQL database for hotel rooms.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)

## Prerequisites

- Beego 
- Gorm
- Docker
- PostgreSQL / PGAdmin

## Installation

**1. Clone this repository:**
```bash
git clone https://github.com/imnahmed17/beego-api.git
cd beego-api
```

**2. PostgreSQL setup:**
```bash
cd pg
docker-compose up -d
```
To find containers name and ports:
```bash
docker ps -a
```
To stop docker containers:
```bash
docker-compose down
```
Now open your browser and type `http://localhost:8000` into the url to access PostgreSQL. For login give Email: admin@user.com and Password: adminuser. After that, right click on `Server > Register > Server` in left sidemenu. In general tab give type any name you want as Server name. In connection tab give Host name: postgres, Username: myuser, Password: mypassword and save it. Then nevigate your server name and create a database `bookingApp` right clicking on your server name. For writing database queries open `Query Tool`. To open `Query Tool` right click on your database name. Now write the below query to create a table using `Query Tool`:
```bash
CREATE TABLE Hotel_Locations (
    Location_id SERIAL PRIMARY KEY,
    Location_name VARCHAR (70) NOT NULL
);

CREATE TABLE Hotel_Lists (  
    Hotel_id VARCHAR(100) PRIMARY KEY,  
    Hotel_name TEXT NOT NULL,
	Hotel_city TEXT NOT NULL,
	Hotel_image_url TEXT NOT NULL,
	Hotel_price NUMERIC (10, 2) NOT NULL,
	Location_id INT REFERENCES Hotel_Locations (Location_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE Hotel_Details (
	Hotel_id VARCHAR (100) PRIMARY KEY,
	Hotel_review_count INT NOT NULL,
	Hotel_rating INT NOT NULL,
	Hotel_no_of_bed INT NOT NULL,
	Hotel_amenities TEXT[] NOT NULL,
	Hotel_description TEXT NOT NULL,
	Hotel_image_urls TEXT[] NOT NULL,
	Hotel_property_type VARCHAR (20)
);
```

**3. Project setup:**

Create a `conf` folder inside `bookingApi` folder as well as create a `app.conf` file inside `conf` folder (`bookingApi/conf/app.conf`). Now, add these below lines to `app.conf` file:
```bash
appname = bookingApi
httpport = 8080
runmode = dev

dbdriver = postgres
dbhost = localhost
dbport = 5432
dbuser = myuser
dbpass = mypassword
dbname = bookingApp
dbsslmode = disable
dbtimezone = your_preferred_time_zone

rapidapikey = your_api_key
rapidapihost = your_api_host
```
To get RapidApi key and host browse [Booking.com API | RapidAPI](https://rapidapi.com/ntd119/api/booking-com13?fbclid=IwAR2aC91bQeRddPSQZ7szn93Ck7hMdmRUwpZ9EBHQf-RPps0lua_Qe3jLd8I).

**4. Run this application:**
```bash
cd ..
cd bookingApi
go mod tidy
bee run
```
