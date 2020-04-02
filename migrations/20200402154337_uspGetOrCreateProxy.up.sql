create or alter proc uspGetOrCreateProxy @proxy_addr varchar(30), @country_id smallint as
begin
    select proxy_id from proxy where proxy_addr = @proxy_addr
    if @@rowcount = 0
        insert into proxy (proxy_addr, country_id)
        output inserted.proxy_id
        values (@proxy_addr, @country_id)
end