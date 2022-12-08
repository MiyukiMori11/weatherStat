// Package handler package with custom handlers
package handler

//go:generate mockery --dir=. --name=Storage --filename=storage_mock.go --output=. --inpackage

import (
	"errors"

	"github.com/MiyukiMori11/weatherstat/explorer/internal/domain"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/models"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/server/operations"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

type Handler interface {
	GetCities(operations.GetCitiesParams) middleware.Responder
	PostCities(operations.PostCitiesParams) middleware.Responder
	DeleteCities(operations.DeleteCitiesParams) middleware.Responder
	GetTemp(operations.GetTempParams) middleware.Responder
}

type handler struct {
	logger  *zap.Logger
	storage Storage
	client  Client
}

type Storage interface {
	GetSubscribedCities() (domain.Cities, error)
	GetCountryCode(countryName string) (string, error)
	AddCity(domain.City) error
	GetTemperatureStat(cityName, countryName string) (float64, float64, error)
	DeleteCity(cityName, countryName string) error
}

type Client interface {
	GetCoordinates(city, countryCode string) (float64, float64, error)
}

// New is handler constructor
func New(logger *zap.Logger, storage Storage, client Client) Handler {
	return &handler{
		logger:  logger,
		storage: storage,
		client:  client,
	}
}

// GetCities handles GET /cities request
func (h *handler) GetCities(gcp operations.GetCitiesParams) middleware.Responder {
	var citiesCount int64
	response := &models.CitiesStat{}

	cities, err := h.storage.GetSubscribedCities()
	if err != nil {
		h.logger.Error("can't get subscribed countries", zap.Error(err))
		return operations.NewGetCitiesBadRequest()
	}

	citiesCount = int64(len(cities))
	response.CitiesCount = &citiesCount

	for _, city := range cities {
		response.Cities = append(response.Cities, &models.CityInfo{
			CityName:    models.CityName(city.Name),
			CountryName: models.CountryName(city.CountryName),
		})
	}

	return operations.NewGetCitiesOK().WithPayload(response)
}

// PostCities handles POST /cities request
func (h *handler) PostCities(pcp operations.PostCitiesParams) middleware.Responder {
	var err error
	city := domain.City{
		Name:        string(pcp.City.Name),
		CountryName: string(pcp.City.Country),
	}

	city.CountryCode, err = h.storage.GetCountryCode(string(pcp.City.Country))
	if err != nil {
		h.logger.Error("can't get country code", zap.Error(err), zap.String("country", city.CountryName), zap.String("city", city.Name))
		return operations.NewPostCitiesBadRequest()
	}

	city.Latitude, city.Longitude, err = h.client.GetCoordinates(city.Name, city.CountryCode)
	if err != nil {
		h.logger.Error("can't get coordinates", zap.Error(err), zap.String("country", city.CountryName), zap.String("city", city.Name))
		if errors.Is(err, domain.ErrNotFound) {
			return operations.NewPostCitiesNotFound()
		}
		return operations.NewPostCitiesBadRequest()
	}

	if err = h.storage.AddCity(city); err != nil {
		h.logger.Error("can't add new city", zap.Error(err), zap.String("country", city.CountryName), zap.String("city", city.Name))
		return operations.NewPostCitiesBadRequest()
	}

	return operations.NewPostCitiesOK()

}

// DeleteCities handles DELETE /cities request
func (h *handler) DeleteCities(dcp operations.DeleteCitiesParams) middleware.Responder {
	if err := h.storage.DeleteCity(
		dcp.CityName,
		dcp.CountryName,
	); err != nil {
		h.logger.Error("can't delete city", zap.Error(err), zap.String("country", dcp.CountryName), zap.String("city", dcp.CityName))
		return operations.NewDeleteCitiesBadRequest()
	}

	return operations.NewDeleteCitiesOK()
}

// GetTemp handles GET /temp request
func (h *handler) GetTemp(gtp operations.GetTempParams) middleware.Responder {
	var err error
	response := &models.CityTemp{
		Name: gtp.CityName,
		Avgc: 0,
		Avgf: 0,
	}

	if response.Avgc, response.Avgf, err = h.storage.GetTemperatureStat(gtp.CityName, gtp.CountryName); err != nil {
		h.logger.Error("can't get temperature statistics", zap.Error(err), zap.String("country", gtp.CountryName), zap.String("city", gtp.CityName))
		return operations.NewGetTempBadRequest()
	}

	if response.Avgc == response.Avgf && response.Avgc == 0 {
		h.logger.Debug("temp stat for country not found", zap.Error(err), zap.String("country", gtp.CountryName), zap.String("city", gtp.CityName))
		return operations.NewGetTempNotFound()
	}

	return operations.NewGetTempOK().WithPayload(response)
}
