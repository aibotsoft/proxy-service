ALTER DATABASE proxy_service SET timezone TO 'UTC';

create table if not exists country
(
    created_at   timestamp with time zone not null default now(),
    country_id   smallserial              not null,
    country_name varchar(30)              not null default 'Unknown',
    country_code varchar(2)               not null default 'NA',
    constraint country_pk primary key (country_id),
    constraint country_code_ui unique (country_code)
);

create table if not exists proxy
(
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone null,
    proxy_id   serial                   not null,
    proxy_ip   inet                     not null,
    proxy_port int                      not null,
    country_id smallint                 not null default 1,
    constraint proxy_pk primary key (proxy_id),
    constraint proxy_ui unique (proxy_ip, proxy_port),
    constraint proxy_country_fk foreign key (country_id) references country on delete cascade
);

create table if not exists proxy_user
(
    created_at timestamp with time zone not null default now(),
    user_id    smallserial              not null,
    user_name  varchar(30)              not null default 'checker',
    constraint proxy_user_pk primary key (user_id)
);

create table if not exists stat
(
    created_at  timestamp with time zone not null default now(),
    stat_id     serial                   not null,
    proxy_id    int                      not null,
    conn_time   int                      not null,
    conn_status bool                     not null,
    user_id     smallint                 not null default 1,

    constraint stat_pk primary key (stat_id),
    constraint stat_proxy_fk foreign key (proxy_id) references proxy on delete cascade,
    constraint stat_user_fk foreign key (user_id) references proxy_user on delete cascade,
);
create index stat_proxy_id_index on stat (proxy_id);
create index stat_created_at_index on stat (created_at);

insert into country (country_id, country_code, country_name)
values (1, default, default);
insert into proxy_user (user_id, user_name)
values (1, default);

create or replace view proxy_stat_view as
select p.proxy_id,
       proxy_ip,
       proxy_port,
       country_id,
       s.created_at,
       stat_id,
       conn_time,
       conn_status,
       user_id
from proxy p
         left join stat s using (proxy_id)
where s.stat_id is null
  and deleted_at is null;
---- create above / drop below ----
drop table proxy;
drop table country;
drop table proxy_user;
drop table stat;
