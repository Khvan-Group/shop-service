create table if not exists t_categories
(
    code varchar(255) primary key not null,
    name varchar(255)             not null
);

create table if not exists t_items
(
    code        varchar(255) primary key not null,
    name        varchar(255)             not null,
    description text,
    count       int                      not null default 0,
    price       int                      not null default 0,
    category    varchar(255)             not null
);

alter table if exists t_items
    add constraint fk_items_categories
        foreign key (category) references t_categories (code)
            on delete cascade;