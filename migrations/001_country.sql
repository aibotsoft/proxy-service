create table if not exists country
(
    created_at   timestamp with time zone not null default now(),
    country_id   smallserial              not null,
    country_name varchar(30)              not null default 'Unknown',
    country_code varchar(2)               not null default 'NA',
    constraint country_pk primary key (country_id),
    constraint country_name_ui unique (country_name),
    constraint country_code_ui unique (country_code)

);
create table if not exists proxy
(
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    proxy_id   serial                   not null,
    proxy_ip   inet                     not null,
    proxy_port int                      not null,
    country_id smallint                 not null,
    constraint proxy_pk primary key (proxy_id),
    constraint proxy_ui unique (proxy_ip, proxy_port),
    constraint proxy_country_fk foreign key (country_id) references public.country on delete cascade
);

---- create above / drop below ----
drop table country;
drop table proxy;