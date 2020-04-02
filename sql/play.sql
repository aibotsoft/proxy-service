declare @code varchar(2) = '66';
declare @var int;
select @var = country_id
from country
where country_code = @code
select @var


if @@rowcount = 0
    insert into country(country_name, country_code)
    output inserted.country_id
    values ('UNKNOWN', @code)


--     select @@IDENTITY

create or alter proc uspTest as
begin
    set nocount on;
    declare @code varchar(2) = '66';
    select @code
end
go

exec uspTest