package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"log"
	"os"
)

var db *gorm.DB

type Fields struct {
	AirportInfo []Airport `json:"airports"`
}

func importAirports()Fields {
	airports, err := ioutil.ReadFile("airports.json")
	if err != nil {
		log.Fatal(err)
	}

	var data Fields
	err = json.Unmarshal(airports, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "wxmap.db")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Airport{})
	db.AutoMigrate(&Weather{})
	var airportItem Airport
	db.Last(&airportItem, 1)
	if airportItem.ID == 0 {
		fmt.Println(importAirports())
		for _, airport := range importAirports().AirportInfo {
			db.Create(&airport)
		}
	}

}



func main() {
	port := os.Getenv("WXPORT")

	if port == "" {
		log.Fatal("$WXPORT must be set")
	}
	defer db.Close()

	router := gin.New()
	router.Use(gin.Logger())
	apiRouter := router.Group("/")
	{
		apiRouter.GET("/wx/:iata", GetAirportWx)
	}
	addr := fmt.Sprintf(":%s", port)

	router.Run(addr)
}