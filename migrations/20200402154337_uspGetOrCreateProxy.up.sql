create or alter proc uspGetOrCreateProxy @proxy_addr varchar(30), @country_id smallint as
begin
    set nocount on
    declare @proxy_id int

    select @proxy_id = proxy_id from proxy where proxy_addr = @proxy_addr
    if @@rowcount = 0
        insert into proxy (proxy_addr, country_id)
        output inserted.proxy_id
        values (@proxy_addr, @country_id)
    else
        select @proxy_id
end