package domain

import "errors"

type City struct {
	Name        string
	CountryCode string
	CountryName string
	Latitude    float64
	Longitude   float64
	TempK       float64
	TempC       float64
	TempF       float64
}

type Cities []City

var ErrNotFound = errors.New("not found")
