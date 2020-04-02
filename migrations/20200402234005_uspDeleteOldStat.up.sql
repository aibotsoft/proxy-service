create or alter proc uspDeleteOldStat @minCheckCount int = 20 as
-- задача удалять старые проверки если их больше 20
with t as (
    select proxy_id, min(created_at) created_at, count(proxy_id) check_count
    from stat
    group by proxy_id
    having count(proxy_id) > @minCheckCount
)
delete stat
output deleted.proxy_id, deleted.created_at
from stat
         inner join t on stat.created_at = t.created_at