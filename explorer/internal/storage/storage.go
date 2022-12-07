package storage

import (
	"context"
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
}

type storage struct {
	config *config.Storage

	pool *pgxpool.Pool
}

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
	query := "SELECT s.city_name, c.name from subscribtion s JOIN countries c ON s.country_code=c.code;"
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		city := domain.City{}
		err = rows.Scan(&city.Name, city.CountryName)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}

func (s *storage) GetCountryCode(countryName string) (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	var cCode string

	query := "SELECT code FROM countries WHERE name='$1' ORDER BY code LIMIT 1;"
	row, err := s.pool.Query(ctx, query, countryName)
	if err != nil {
		return "", err
	}
	defer row.Close()

	if err := row.Scan(&cCode); err != nil {
		return "", err
	}

	return cCode, nil
}

func (s *storage) AddCity(city domain.City) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	query := "INSERT INTO subscription(city_name, country_code, latitude, longitude) VALUES ('$1', '$2', $3, $4);"
	if _, err := s.pool.Exec(ctx, query, city.Name, city.CountryCode, city.Latitude, city.Longitude); err != nil {
		return err
	}

	return nil
}

func (s *storage) GetTemperatureStat(cityName, countryName string) (float64, float64, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	var avcC, avgF float64

	query := "SELECT avg_temperature('$1', '$2');"
	row, err := s.pool.Query(ctx, query, countryName)
	if err != nil {
		return 0, 0, err
	}
	defer row.Close()

	if err := row.Scan(&avcC, &avgF); err != nil {
		return 0, 0, err
	}

	return avcC, avgF, nil

}

func (s *storage) DeleteCity(cityName, countryName string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	query := "DELETE FROM subscription WHERE city_name='$1' AND country_code=(SELECT code from countries WHERE name='$2');"
	if _, err := s.pool.Exec(ctx, query, cityName, countryName); err != nil {
		return err
	}

	return nil
}

func (s *storage) AddTemperatureForCities(cities domain.Cities) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	for _, city := range cities {
		query := "INSERT INTO temperature (city_name, country_code, celsius, fahrenheit) VALUES ('$1', '$2', '$3', '$4');"
		if _, err := tx.Exec(ctx, query, city.CountryName, city.CountryCode, city.TempC, city.TempF); err != nil {
			return err
		}
	}

	return nil
}
