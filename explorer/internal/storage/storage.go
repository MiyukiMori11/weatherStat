// Package storage represents connection and query methods for database
package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MiyukiMori11/weatherstat/explorer/internal/config"
	"github.com/MiyukiMori11/weatherstat/explorer/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	GetSubscribedCities() (domain.Cities, error)
	GetCountryCode(countryName string) (string, error)
	AddCity(domain.City) error
	GetTemperatureStat(cityName, countryName string) (float64, float64, error)
	DeleteCity(cityName, countryName string) error
	AddTemperatureForCities(domain.Cities) error
	GetCoordinates(cityName, countryCode string) (float64, float64, error)
}

type storage struct {
	config *config.Storage

	pool *pgxpool.Pool
}

// New is storage constructor
func New(ctx context.Context, config *config.Storage) (Storage, error) {
	pool, err := pgxpool.New(context.Background(), config.URL())
	if err != nil {
		return nil, err
	}

	go func(ctx context.Context) {
		ctx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()

		<-ctx.Done()
		pool.Close()
	}(ctx)

	return &storage{
		config: config,
		pool:   pool,
	}, nil
}

func (s *storage) GetSubscribedCities() (domain.Cities, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	cities := make(domain.Cities, 0)

	//language=sql
	query := `SELECT s.city_name, c.name from subscription s JOIN countries c ON s.country_code=c.code;`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("can't exec query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		city := domain.City{}
		err = rows.Scan(&city.Name, &city.CountryName)
		if err != nil {
			return nil, fmt.Errorf("can't scan result: %w", err)
		}
		cities = append(cities, &city)
	}

	return cities, nil
}

func (s *storage) GetCountryCode(countryName string) (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	var cCode string

	query := `SELECT code FROM countries WHERE name LIKE $1 ORDER BY code LIMIT 1;`
	row := s.pool.QueryRow(ctx, query, "%"+countryName+"%")

	if err := row.Scan(&cCode); err != nil {
		return "", err
	}

	return cCode, nil
}

func (s *storage) AddCity(city domain.City) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	query := `INSERT INTO subscription(city_name, country_code, latitude, longitude) VALUES ($1, $2, $3, $4);`
	if _, err := s.pool.Exec(ctx, query, city.Name, city.CountryCode, city.Latitude, city.Longitude); err != nil {
		return err
	}

	return nil
}

func (s *storage) GetTemperatureStat(cityName, countryName string) (float64, float64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	var avcC, avgF float64

	query := `SELECT * FROM avg_temperature($1::text, $2::text);`
	row := s.pool.QueryRow(ctx, query, cityName, "%"+countryName+"%")

	if err := row.Scan(&avcC, &avgF); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, err
		}
	}

	return avcC, avgF, nil

}

func (s *storage) DeleteCity(cityName, countryName string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	query := `DELETE FROM subscription WHERE city_name=$1 AND country_code=(SELECT code from countries WHERE name like $2);`
	if _, err := s.pool.Exec(ctx, query, cityName, "%"+countryName+"%"); err != nil {
		return err
	}

	return nil
}

func (s *storage) AddTemperatureForCities(cities domain.Cities) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("can't begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	for _, city := range cities {
		query := `INSERT INTO temperature (city_name, country_code, celsius, fahrenheit) VALUES ($1, $2, $3, $4);`
		if _, err := tx.Exec(ctx, query, city.Name, city.CountryCode, city.TempC, city.TempF); err != nil {
			return fmt.Errorf("can't exec query: %w", err)
		}
	}

	return nil
}

func (s *storage) GetCoordinates(cityName, countryCode string) (lat, lon float64, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	query := `SELECT latitude, longitude FROM subscription WHERE city_name=$1 AND country_code=$2 ORDER BY city_name LIMIT 1;`
	row := s.pool.QueryRow(ctx, query, cityName, countryCode)

	if err := row.Scan(&lat, &lon); err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}
