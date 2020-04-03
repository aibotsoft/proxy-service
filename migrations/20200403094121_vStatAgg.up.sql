create or alter view vStatAgg as
select top 1000 s.proxy_id,
                p.proxy_addr,
                p.deleted_at,
                count(s.proxy_id)        stat_count,
                avg(s.conn_time)         avg_conn_time,
                avg(s.conn_status * 1.0) avg_conn_status
from stat s
         join proxy p on s.proxy_id = p.proxy_id
group by s.proxy_id, p.proxy_addr, p.deleted_at
order by stat_count desc