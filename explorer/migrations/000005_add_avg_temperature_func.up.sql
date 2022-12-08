CREATE FUNCTION avg_temperature(city_name text, country_name text) RETURNS table(avg_c numeric, avg_f numeric) AS $$
    SELECT AVG(celsius), AVG(fahrenheit) FROM temperature WHERE city_name=$1 AND country_code=(SELECT code FROM countries WHERE name like $2) GROUP BY city_name, country_code;
$$ LANGUAGE SQL;