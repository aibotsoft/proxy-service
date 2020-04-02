create or alter proc uspDeleteBadProxy @minSuccessRate real = 0.1, @minCheckCount int = 20 as
--     Задача пометить удаленными прокси с плохой статистикой
begin
    with t as (
        select s.proxy_id, avg(conn_time) avgTime, avg(conn_status * 1.0) successRate, count(s.proxy_id) checkCount
        from stat s
                 inner join proxy p on s.proxy_id = p.proxy_id
        where p.deleted_at is null
        group by s.proxy_id
        having avg(conn_status * 1.0) < @minSuccessRate
           and count(s.proxy_id) >= @minCheckCount
    )
    update proxy
    set deleted_at = sysdatetimeoffset()
    output inserted.proxy_id, inserted.deleted_at
    from t
    where proxy.proxy_id = t.proxy_id
end