create or alter proc uspGetBestProxy @returnCount int = null, @minSuccessRate int = null, @minCheckCount int = null as
begin
    set nocount on;
    with top_stat as (
        select top 50 proxy_id, avg(conn_time) avgTime, avg(conn_status * 1.0) successRate, count(proxy_id) checkCount
        from stat
        group by proxy_id
        having count(proxy_id) > isnull(@minCheckCount, 0)
           and avg(conn_status * 1.) > isnull(@minSuccessRate, 0.3)
        order by avgTime
    )
    select top(isnull(@returnCount, 1)) t.proxy_id,
                                        proxy_addr,
                                        avgTime,
                                        successRate,
                                        checkCount
    from top_stat t
             inner join proxy p on p.proxy_id = t.proxy_id
    where p.deleted_at is null
end