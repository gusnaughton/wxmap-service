package main

import (
	"fmt"
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

	fmt.Println(code[0])
	code = code[1:]

	fmt.Println("GetAirport" + code)
	db.Where("code = ?", code).First(&airport)
	fmt.Println("Airport" + code + " ID: " + string(airport.ID))
	return airport
}