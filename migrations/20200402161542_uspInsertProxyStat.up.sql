create or alter proc uspInsertProxyStat @proxy_id int, @conn_time int, @conn_status bit as
begin
    set nocount on
    insert into stat (proxy_id, conn_time, conn_status)
    output inserted.created_at
    values (@proxy_id, @conn_time, @conn_status)
end
-- go
--
-- exec uspInsertProxyStat 1, 1, 1
