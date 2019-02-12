package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

type Weather struct {
	gorm.Model
	Airport     Airport `gorm:"foreignkey:AirportId"`
	AirportId   int
	Timestamp   int
	Ceiling     int
	Visibility  int
	Wind        int
	Temperature int
	Sky         int
}

type ADDSRoot struct {
	XMLName xml.Name `xml:"response"`
	Data    ADDSData `xml:"data"`
}

type ADDSData struct {
	XMLName xml.Name `xml:"data"`
	METAR  METAR  `xml:"METAR"`
}

type METAR struct {
	XMLName xml.Name `xml:"METAR"`
	ICAO        string  `xml:"station_id"`
	Observation string  `xml:"observation_time"`
	Temperature float32 `xml:"temp_c"`
	Wind        int     `xml:"wind_speed_kt"`
	Visibility  float32     `xml:"visibility_statute_mi"`
	Category    string  `xml:"flight_category"`
}

func GetAirportWx(c *gin.Context) {
	var wx []Weather
	iata := c.Param("iata")

	airport := GetAirport(iata)

	db.Where("airport_id = ?", airport.ID).Find(&wx)
	c.JSON(http.StatusOK, wx)
}

func worker(id int, jobs <-chan string, results <-chan string) {
	for j := range jobs {
		url := fmt.Sprintf("https://www.aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&hoursBeforeNow=3&mostRecentForEachStation=true&stationString=K%s", j)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var result ADDSRoot
		err = xml.Unmarshal(body, &result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result.Data.METAR)
	}
}

func UpdateWxData(c *gin.Context) {
	var airports []Airport
	db.Find(&airports)

	jobs := make(chan string, len(airports))
	results := make(chan string, len(airports))

	for i := 0; i < 10; i++ {
		go worker(i, jobs, results)
	}

	for i := 0; i < len(airports); i++ {
		jobs <- airports[i].Code
	}
	close(jobs)

	for i := 0; i < 3; i++ {
		<-results
	}
}
