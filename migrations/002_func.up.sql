create or alter proc uspGetOrCreateProxyCountry @country_name varchar(30), @country_code varchar(2) as
begin
    select country_id from country where country_code = @country_code
    if @@rowcount = 0
        insert into country(country_name, country_code)
        output inserted.country_id
        values (@country_name, @country_code)
end
-- go
-- create or alter proc uspGetOrCreateProxy @proxy_addr varchar(30), @country_id smallint as
-- begin
--     select proxy_id from proxy where proxy_addr = @proxy_addr
--     if @@rowcount = 0
--         insert into proxy (proxy_addr, country_id)
--         output inserted.proxy_id
--         values (@proxy_addr, @country_id)
-- end
--
-- create or alter proc uspGetNextProxy @returnCount int = null, @minCheckInterval int = null as
-- begin
--     select top (isnull(@returnCount, 100)) proxy.proxy_id,
--                                         proxy_addr,
--                                         max(s.created_at) last
--     from proxy
--              left join stat s on proxy.proxy_id = s.proxy_id
--     where deleted_at is null
--     group by proxy.proxy_id, proxy_addr
--     having datediff(minute, isnull(max(s.created_at), '2019'), sysdatetimeoffset()) > isnull(@minCheckInterval, 30)
--     order by last
-- end
--
-- create or alter proc uspGetBestProxy @returnCount int = null, @minSuccessRate int = null, @minCheckCount int = null as
-- begin
--     with top_stat as (
--         select top 50 proxy_id, avg(conn_time) avgTime, avg(conn_status * 1.0) successRate, count(proxy_id) checkCount
--         from stat
--         group by proxy_id
--         having count(proxy_id) > isnull(@minCheckCount, 0)
--            and avg(conn_status * 1.) > isnull(@minSuccessRate, 0.3)
--         order by avgTime
--     )
--     select top(isnull(@returnCount, 1)) t.proxy_id,
--                                          proxy_addr,
--                                          avgTime,
--                                          successRate,
--                                          checkCount
--     from top_stat t
--              inner join proxy p on p.proxy_id = t.proxy_id
--     where p.deleted_at is null
-- end
