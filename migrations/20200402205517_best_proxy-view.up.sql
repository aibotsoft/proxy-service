create or alter view best_proxy_view as
with top_stat as (
    select top 200 proxy_id, avg(conn_time) avgTime, avg(conn_status * 1.0) successRate, count(proxy_id) checkCount
    from stat
    group by proxy_id
    having count(proxy_id) > 2
       and avg(conn_status * 1.) > 0.3
    order by avgTime
)
select top (200) t.proxy_id,
                 proxy_addr,
                 avgTime,
                 successRate,
                 checkCount
from top_stat t
         inner join proxy p on p.proxy_id = t.proxy_id
where p.deleted_at is null
order by avgTime