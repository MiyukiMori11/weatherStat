CREATE TABLE temperature(
    city_name varchar(25)  NOT NULL,
    country_code varchar(5) NOT NULL,
    celsius numeric(4,0) NOT NULL,
    fahrenheit numeric(4,0) NOT NULL,
    expires timestamp NOT NULL DEFAULT NOW()
);