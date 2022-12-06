package domain

type Country struct {
	Name        string
	CountryCode string
	Latitude    float64
	Longitude   float64
	TempK       float64
	TempC       float64
	TempF       float64
}

type Countries []Country
