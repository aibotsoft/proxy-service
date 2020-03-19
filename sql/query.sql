select *
from proxy;

create or replace function find_country_by_code(code character varying)
    returns smallint
    language plpgsql
as
$$
DECLARE
    country_id smallint;
BEGIN
    SELECT c.country_id
    INTO country_id
    FROM proxy_service.public.country c
    where country_code = code;
    RETURN country_id;
END;
$$;

select find_country_by_code('NA');