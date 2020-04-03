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


go
select *
from proxy
where deleted_at is not null


exec uspDeleteBadProxy 0.1, 7

create or alter view proxy_stat_view as
select count(p.proxy_id) proxy_count,
       count(deleted_at) proxy_deleted,
       s.stat_count,
       avg_conn_time,
       avg_conn_status
from proxy p,
     (select count(s.proxy_id)        stat_count,
             avg(s.conn_time)         avg_conn_time,
             avg(s.conn_status * 1.0) avg_conn_status
      from stat s) as s
group by s.stat_count, avg_conn_time, avg_conn_status

--          join stat s on p.proxy_id = s.proxy_id
-- group by p.proxy_id

select max(created_at), v, d
from country,
     (select version v, dirty d from schema_migrations) as t
group by v, d
go

with t as (
    select version v, dirty d
    from schema_migrations
)
select user_id, user_name, created_at, v, d
from proxy_user,
     t



select *
from proxy
where deleted_at is not null