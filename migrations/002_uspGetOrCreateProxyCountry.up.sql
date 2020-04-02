create or alter proc uspGetOrCreateProxyCountry @country_name varchar(30), @country_code varchar(2) as
begin
    select country_id from country where country_code = @country_code
    if @@rowcount = 0
        insert into country(country_name, country_code)
        output inserted.country_id
        values (@country_name, @country_code)
end





