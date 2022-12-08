package handler_test

import (
	"errors"
	"testing"

	"github.com/MiyukiMori11/weatherstat/explorer/internal/client"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/config"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/domain"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/models"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/server/handler"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/server/operations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go.uber.org/zap"
)

func TestGetCities(t *testing.T) {
	storage := handler.NewMockStorage(t)
	h := handler.New(
		zap.L(),
		storage,
		client.New(&config.Client{}, zap.L()),
	)

	t.Run("positive case", func(t *testing.T) {
		storage.On("GetSubscribedCities").Once().Return(domain.Cities{
			&domain.City{},
		}, nil)
		r := h.GetCities(operations.GetCitiesParams{})

		assert.IsType(t, &operations.GetCitiesOK{}, r)
		assert.NotEmpty(t, r)
	})

	t.Run("positive case with empty slice", func(t *testing.T) {
		storage.On("GetSubscribedCities").Once().Return(domain.Cities{}, nil)
		r := h.GetCities(operations.GetCitiesParams{})

		assert.IsType(t, &operations.GetCitiesOK{}, r)
		assert.NotEmpty(t, r)
	})

	t.Run("error on sql query", func(t *testing.T) {
		storage.On("GetSubscribedCities").Once().Return(nil, errors.New("test"))
		r := h.GetCities(operations.GetCitiesParams{})

		assert.IsType(t, &operations.GetCitiesBadRequest{}, r)
		assert.Empty(t, r)
	})

}

func TestPostCities(t *testing.T) {
	client := handler.NewMockClient(t)
	storage := handler.NewMockStorage(t)
	h := handler.New(
		zap.L(),
		storage,
		client,
	)

	t.Run("positive case", func(t *testing.T) {
		storage.On("GetCountryCode", mock.Anything).Once().Return("TS", nil)
		client.On("GetCoordinates", mock.Anything, mock.Anything).Once().Return(1.1, 2.2, nil)
		storage.On("AddCity", mock.AnythingOfType("domain.City")).Once().Return(nil)
		r := h.PostCities(operations.PostCitiesParams{
			City: &models.CityBodySchema{},
		})

		assert.IsType(t, &operations.PostCitiesOK{}, r)
		assert.Empty(t, r)
	})

	t.Run("error on country code", func(t *testing.T) {
		storage.On("GetCountryCode", mock.Anything).Once().Return("", errors.New("test"))
		r := h.PostCities(operations.PostCitiesParams{
			City: &models.CityBodySchema{},
		})

		assert.IsType(t, &operations.PostCitiesBadRequest{}, r)
		assert.Empty(t, r)
	})

	t.Run("error in getting coordinates", func(t *testing.T) {
		storage.On("GetCountryCode", mock.Anything).Once().Return("TS", nil)
		client.On("GetCoordinates", mock.Anything, mock.Anything).Once().Return(0.0, 0.0, errors.New("test"))
		r := h.PostCities(operations.PostCitiesParams{
			City: &models.CityBodySchema{},
		})

		assert.IsType(t, &operations.PostCitiesBadRequest{}, r)
		assert.Empty(t, r)
	})

	t.Run("empty result in getting coordinates", func(t *testing.T) {
		storage.On("GetCountryCode", mock.Anything).Once().Return("TS", nil)
		client.On("GetCoordinates", mock.Anything, mock.Anything).Once().Return(0.0, 0.0, domain.ErrNotFound)
		r := h.PostCities(operations.PostCitiesParams{
			City: &models.CityBodySchema{},
		})

		assert.IsType(t, &operations.PostCitiesNotFound{}, r)
		assert.Empty(t, r)
	})

	t.Run("error on adding city", func(t *testing.T) {
		storage.On("GetCountryCode", mock.Anything).Once().Return("TS", nil)
		client.On("GetCoordinates", mock.Anything, mock.Anything).Once().Return(1.0, 1.0, nil)
		storage.On("AddCity", mock.AnythingOfType("domain.City")).Once().Return(errors.New("test"))
		r := h.PostCities(operations.PostCitiesParams{
			City: &models.CityBodySchema{},
		})

		assert.IsType(t, &operations.PostCitiesBadRequest{}, r)
		assert.Empty(t, r)
	})

}

func TestDeleteCities(t *testing.T) {
	storage := handler.NewMockStorage(t)
	h := handler.New(
		zap.L(),
		storage,
		client.New(&config.Client{}, zap.L()),
	)

	t.Run("positive case", func(t *testing.T) {
		storage.On("DeleteCity", mock.Anything, mock.Anything).Once().Return(nil)
		r := h.DeleteCities(operations.DeleteCitiesParams{})

		assert.IsType(t, &operations.DeleteCitiesOK{}, r)
		assert.Empty(t, r)
	})

	t.Run("error on sql query", func(t *testing.T) {
		storage.On("DeleteCity", mock.Anything, mock.Anything).Once().Return(errors.New("test"))
		r := h.DeleteCities(operations.DeleteCitiesParams{})

		assert.IsType(t, &operations.DeleteCitiesBadRequest{}, r)
		assert.Empty(t, r)
	})

}

func TestGetTemp(t *testing.T) {
	storage := handler.NewMockStorage(t)
	h := handler.New(
		zap.L(),
		storage,
		client.New(&config.Client{}, zap.L()),
	)

	t.Run("positive case", func(t *testing.T) {
		storage.On("GetTemperatureStat", mock.Anything, mock.Anything).Once().Return(1.0, 0.0, nil)
		r := h.GetTemp(operations.GetTempParams{})

		assert.IsType(t, &operations.GetTempOK{}, r)
		assert.NotEmpty(t, r)
	})

	t.Run("temp info not found", func(t *testing.T) {
		storage.On("GetTemperatureStat", mock.Anything, mock.Anything).Once().Return(0.0, 0.0, nil)
		r := h.GetTemp(operations.GetTempParams{})

		assert.IsType(t, &operations.GetTempNotFound{}, r)
		assert.Empty(t, r)
	})

	t.Run("error on sql query", func(t *testing.T) {
		storage.On("GetTemperatureStat", mock.Anything, mock.Anything).Once().Return(0.0, 0.0, errors.New("test"))
		r := h.GetTemp(operations.GetTempParams{})

		assert.IsType(t, &operations.GetTempBadRequest{}, r)
		assert.Empty(t, r)
	})

}
