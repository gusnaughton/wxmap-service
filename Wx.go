package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type METAR struct {
	gorm.Model `json:"-"`
	Airport     Airport `gorm:"foreignkey:AirportId" json:"-"`
	AirportId   uint `json:"-"`
	Timestamp   int64 `json:"-"`
	Ceiling     int `json:"C"`
	Visibility  int `json:"V"`
	Wind        int `json:"W"`
	Temperature int `json:"T"`
	Sky         int `json:"S"`
}

type tinyMETAR struct {
	C []int
	V []int
	W []int
	T []int
	S []int
}

type miniMETAR struct {
	C [][]int
	V [][]int
	W [][]int
	T [][]int
	S [][]int
}

func GetAirportWx(c *gin.Context) {
	var wx []METAR
	var codes []string
	final := miniMETAR{}

	iata := c.Param("iata")
	codes = strings.Split(iata, ",")

	for _, code := range codes {
		airport := GetAirport(code)
		db.Limit(10).Where("airport_id = ?", airport.ID).Find(&wx)

		res := tinyMETAR{}
		for _, metar := range wx {
			res.C = append(res.C, metar.Ceiling)
			res.V = append(res.V, metar.Visibility)
			res.W = append(res.W, metar.Wind)
			res.T = append(res.T, metar.Temperature)
			res.S = append(res.S, metar.Sky)
		}
		final.C = append(final.C, res.C)
		final.V = append(final.V, res.V)
		final.W = append(final.W, res.W)
		final.T = append(final.T, res.T)
		final.S = append(final.S, res.S)
	}


	c.JSON(http.StatusOK, final)
}

func scraper(id int, jobs <-chan []string, results chan *METAR) {
	for row := range jobs {
		// IATA code
		airport := GetAirport(row[1])
		if airport.ID != 0 {
			// Observation Time
			fmt.Println(row[2])
			timestamp, err := time.Parse(time.RFC3339Nano, row[2])
			if err != nil {
				log.Fatal(err)
			}

			// Ceiling 24
			ceiling, _ := strconv.ParseInt(row[24], 0, 64)
			// Visibility 10
			visibility, _ := strconv.ParseInt(row[10], 0, 64)
			// Wind 8
			wind, _ := strconv.ParseInt(row[8], 0, 64)
			// Temperature 5
			temperature, _ := strconv.ParseInt(row[5], 0, 64)

			db.Create(&METAR{
				Airport:     airport,
				AirportId:   airport.ID,
				Timestamp:   timestamp.Unix(),
				Ceiling:     CeilingParser(ceiling),
				Visibility:  VisibilityParser(visibility),
				Wind:        WindParser(wind),
				Temperature: TemperatureParser(temperature),
				Sky:         SkyParser(row[30]),
			})


		}
	}
}

const (
	Purple int = iota + 1
	Red
	Blue
	Green
)

func CeilingParser(ceil int64) int {
	if ceil < 5 {
		return Purple
	} else if ceil < 10 {
		return Red
	} else if ceil < 30 {
		return Blue
	} else if ceil >= 30 {
		return Green
	} else {
		return Purple
	}

}

func VisibilityParser(vis int64) int {
	if vis < 1 {
		return Purple
	} else if vis < 3 {
		return Red
	} else if vis < 5 {
		return Blue
	} else if vis >= 5 {
		return Green
	} else {
		return Purple
	}

}

func WindParser(wind int64) int {
	if wind > 30 {
		return Purple
	} else if wind > 20 {
		return Red
	} else if wind > 10 {
		return Blue
	} else if wind <= 10 {
		return Green
	} else {
		return Purple
	}
}

func TemperatureParser(temp int64) int {
	if temp <= 0 {
		return Blue
	} else if temp < 30 {
		return Green
	} else if temp >= 30 {
		return Red
	} else {
		return Red
	}
}

func SkyParser(sky string) int {
	switch sky {
	case "VFR":
		return Green
	case "MVFR":
		return Blue
	case "IFR":
		return Red
	case "LIFR":
		return Purple
	default:
		return Purple
	}

}

func UpdateWxData(c *gin.Context) {
	var airports []Airport
	db.Find(&airports)

	// get wx data from adds
	resp, err := http.Get("https://aviationweather.gov/adds/dataserver_current/current/metars.cache.csv")
	if err != nil {
		log.Fatal(err)
	}

	// parse into a string to trim first 60 characters
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan []string, 5000)
	results := make(chan *METAR, 5000)

	for i := 0; i < 10; i++ {
		go scraper(i, jobs, results)
	}

	// parse the csv
	r := csv.NewReader(bytes.NewReader(body[659:]))
	for {
		row, err := r.Read()
		// eof case
		if err == io.EOF {
			break
		}
		// error case
		if err != nil {
			log.Fatal(err)
		}
		jobs <- row
	}
	close(jobs)

	c.JSON(http.StatusOK, "ok")
}
