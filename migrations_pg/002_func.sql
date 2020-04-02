create or replace function get_next_proxy_for_check(min_interval_minutes int, max_limit int)
    returns table
            (
                proxy_id   int,
                proxy_ip   inet,
                proxy_port int
            )
as
$$
select proxy.proxy_id,
       proxy_ip,
       proxy_port
from proxy
         left join (
    select max(created_at) as last_check, proxy_id
    from stat
    group by proxy_id
) t2 on proxy.proxy_id = t2.proxy_id
where now() - coalesce(last_check, '2000.01.01'::timestamp) > make_interval(mins => min_interval_minutes)
order by last_check is not null, last_check
limit max_limit;
$$ language sql;

create or replace function get_or_create_proxy_country(in name varchar(30), in code varchar(2), out id smallint)
as
$$
begin
    select country_id into id from country where country_code = code;
    if not found then
        insert into country (country_name, country_code)
        values (name, code)
        returning country_id into id;
    end if;
end
$$ language plpgsql;

create or replace function get_or_create_proxy_item(in ip inet, in port int, in name varchar(30), in code varchar(2),
                                                    out id int)
as
$$
begin
    select p.proxy_id into id from proxy p where p.proxy_ip = ip and p.proxy_port = port;
    if not found then
        insert into proxy (proxy_ip, proxy_port, country_id)
        values (ip, port, (select get_or_create_proxy_country(name, code)))
        returning proxy_id into id;
    end if;
end
$$ language plpgsql;

---- create above / drop below ----
drop function get_next_proxy_for_check(min_interval_minutes int, max_limit int);
drop function get_or_create_proxy_country(in name varchar(30), in code varchar(2), out id smallint);
drop function get_or_create_proxy_item(in ip inet, in port int, in name varchar(30), in code varchar(2),
    out id int);
