CREATE OR REPLACE FUNCTION expire_old_rows_from_temperature() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  DELETE FROM temperature WHERE timestamp < NOW() - INTERVAL $CNTDAYARCHIVE;
  RETURN NEW;
END;
$$;

CREATE TRIGGER delete_old_rows_trigger
    AFTER INSERT ON temperature
    EXECUTE PROCEDURE expire_old_rows_from_temperature();