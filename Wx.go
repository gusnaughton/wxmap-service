package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Weather struct {
	gorm.Model
	Airport		Airport `gorm:"foreignkey:AirportId"`
	AirportId	int
	Timestamp   int
	Ceiling     int
	Visibility  int
	Wind        int
	Temperature int
	Sky         int
}

func GetAirportWx(c *gin.Context) {
	var wx []Weather
	iata := c.Param("iata")

	airport := GetAirport(iata)

	db.Where("airport_id = ?", airport.ID).Find(&wx)
	c.JSON(http.StatusOK, wx)
}