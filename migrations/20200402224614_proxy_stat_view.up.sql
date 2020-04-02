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