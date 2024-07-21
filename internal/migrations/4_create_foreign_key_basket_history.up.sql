alter table if exists t_basket_history
    add constraint fk_basket_history_items
        foreign key (item_code) references t_items (code)
            on delete cascade;