
create table if not exists t2
(
    created_at   timestamp with time zone not null default now()
);

---- create above / drop below ----

drop table t2;