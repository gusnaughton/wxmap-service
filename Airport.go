package main

import (
	"github.com/jinzhu/gorm"
)

type Airport struct {
	gorm.Model
	Code     string
	Name     string
	Location string
}

func GetAirport(iata string)Airport {
	var airport Airport
	db.Where("name = ?", iata).First(&airport)
	return airport
}