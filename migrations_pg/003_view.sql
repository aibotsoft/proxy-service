-- На всякий случай аналог этой функции в виде view
create or replace view get_next_proxy_for_check as
select proxy.proxy_id,
       proxy_ip,
       proxy_port,
       last_check
from proxy
         left join (
    select max(created_at) as last_check, proxy_id
    from stat
    group by proxy_id
) t2 on proxy.proxy_id = t2.proxy_id
where now() - coalesce(last_check, '2000.01.01'::timestamp) > make_interval(mins => 60)
order by last_check is not null, last_check
limit 100;


---- create above / drop below ----
drop view if exists get_next_proxy_for_check;