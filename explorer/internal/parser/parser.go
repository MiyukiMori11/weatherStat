// Package parser represents methods for parse temperature information from external api
package parser

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/MiyukiMori11/weatherstat/explorer/internal/config"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/domain"
)

type Parser interface {
	Run(ctx context.Context)
}

type Storage interface {
	GetSubscribedCities() (domain.Cities, error)
	GetCountryCode(countryName string) (string, error)
	AddTemperatureForCities(domain.Cities) error
	GetCoordinates(cityName, countryCode string) (float64, float64, error)
}

type Client interface {
	GetTemperature(Latitude, Longitude float64) (float64, error)
}

type parser struct {
	config *config.Parser
	logger *zap.Logger

	storage Storage
	client  Client
}

var (
	errChan = make(chan error)
	done    = make(chan struct{})
)

// New is parser constructor
func New(cfg *config.Parser, logger *zap.Logger, storage Storage, client Client) Parser {
	return &parser{
		config:  cfg,
		logger:  logger,
		storage: storage,
		client:  client,
	}
}

// Run runs the parser to update the temperature
func (p *parser) Run(ctx context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	ticker := time.NewTicker(p.config.Timeout())
	go p.parse()

loop:
	for {
		select {
		case <-ctx.Done():
			<-done
			break loop
		case <-ticker.C:
			go p.parse()
		case err := <-errChan:
			p.logger.Error("error on temperature parsing", zap.Error(err))
		case <-done:
		}
	}

}

func (p *parser) parse() {
	defer func() { done <- struct{}{} }()

	p.logger.Info("start parsing")
	defer p.logger.Info("end parsing")

	cities, err := p.storage.GetSubscribedCities()
	if err != nil {
		errChan <- err
		return
	}
	wg := &sync.WaitGroup{}

	wg.Add(len(cities))
	for _, city := range cities {
		go func(c *domain.City) {
			defer wg.Done()
			var err error
			c.CountryCode, err = p.storage.GetCountryCode(c.CountryName)
			if err != nil {
				p.logger.Warn("can't get country code",
					zap.Error(err), zap.String("country", c.CountryCode), zap.String("city", c.Name))
				return
			}

			c.Latitude, c.Longitude, err = p.storage.GetCoordinates(c.Name, c.CountryCode)
			if err != nil {
				p.logger.Warn("can't get coordinates",
					zap.Error(err), zap.String("country", c.CountryCode), zap.String("city", c.Name))
				return
			}

			c.TempK, err = p.client.GetTemperature(c.Latitude, c.Longitude)
			if err != nil {
				p.logger.Warn("can't get temperature",
					zap.Error(err), zap.String("country", c.CountryCode), zap.String("city", c.Name))
				return
			}
			c.TempC = c.TempK - 273.15
			c.TempF = 1.8*(c.TempK-273) + 32
		}(city)
	}

	wg.Wait()

	if err := p.storage.AddTemperatureForCities(cities); err != nil {
		errChan <- err
		return
	}

}
