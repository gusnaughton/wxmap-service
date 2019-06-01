package main

import (
	"github.com/jinzhu/gorm"
	"strings"
)

type Airport struct {
	gorm.Model
	Code     string
	Name     string
	Location string
}

func GetAirport(iata string)Airport {
	var airport Airport
	code := strings.ToUpper(iata)

	if code[0] == 'K' {
		code = code[1:]
		db.Where("code = ?", code).First(&airport)
	}
	return airport
}