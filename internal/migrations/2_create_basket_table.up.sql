create table if not exists t_basket
(
    username  varchar(255)          not null,
    item_code varchar(255)          not null,
    count     int                   not null default 1
);

alter table if exists t_basket
    add constraint fk_basket_items
        foreign key (item_code) references t_items (code)
            on delete cascade;

create unique index udx_basket_username_items_status on t_basket (username, item_code);