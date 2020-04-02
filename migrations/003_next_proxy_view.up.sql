create or alter view next_proxy_view as
select top 100 proxy.proxy_id,
               proxy_addr,
               max(s.created_at) last
from proxy
         left join stat s on proxy.proxy_id = s.proxy_id
where deleted_at is null
group by proxy.proxy_id, proxy_addr
having datediff(minute, isnull(max(s.created_at), '2019'), sysdatetimeoffset()) > 30
order by last
