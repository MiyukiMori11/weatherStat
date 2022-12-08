CREATE TABLE subscription(
    city_name varchar(25)  NOT NULL,
    country_code varchar(5) NOT NULL,
    latitude numeric,
    longitude numeric
);
CREATE UNIQUE INDEX idx_subscribed_cities ON subscription(city_name, country_code);