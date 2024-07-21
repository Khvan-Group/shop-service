create table if not exists t_basket_history
(
    username  varchar(255)          not null,
    item_code varchar(255)          not null,
    count     int                   not null,
    payed_at  timestamp             not null default now()
);
