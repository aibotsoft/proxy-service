use proxy_service

create table country
(
    country_id   smallint identity,
    country_name varchar(30)       not null default 'Unknown',
    country_code varchar(2)        not null default 'NA',
    created_at   datetimeoffset(2) not null default SYSDATETIMEOFFSET(),
    constraint country_pk primary key nonclustered (country_id),
    constraint country_code_ui unique (country_code)
)

insert into country (country_name, country_code)
values (default, default)

create table proxy
(
    proxy_id   int identity      not null,
    proxy_addr varchar(30)       not null,
    country_id smallint          not null default 1,
    created_at datetimeoffset(2) not null default sysdatetimeoffset(),
    updated_at datetimeoffset(2) not null default sysdatetimeoffset(),
    deleted_at datetimeoffset(2) null,
    constraint proxy_pk primary key nonclustered (proxy_id),
    constraint proxy_ui unique (proxy_addr),
    constraint proxy_country_fk foreign key (country_id) references country on delete cascade
)

insert into proxy (proxy_addr, country_id, deleted_at) values ('0.0.0.0:80', 1, sysdatetimeoffset())

create table proxy_user
(
    user_id    smallint identity not null,
    user_name  varchar(30)       not null default 'checker',
    created_at datetimeoffset(2) not null default sysdatetimeoffset(),
    constraint proxy_user_pk primary key nonclustered (user_id)
)

insert into proxy_user (user_name)
values (default)

create table stat
(
    created_at  datetimeoffset not null default sysdatetimeoffset(),
    proxy_id    int            not null,
    conn_time   int            not null,
    conn_status bit            not null default 0,
    user_id     smallint       not null default 1,
    constraint stat_pk primary key nonclustered (created_at),
    constraint stat_proxy_fk foreign key (proxy_id) references proxy on delete cascade,
    constraint stat_user_fk foreign key (user_id) references proxy_user on delete cascade,
)