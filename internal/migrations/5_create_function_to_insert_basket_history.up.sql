create or replace function insert_to_basket_history(input_username text)
returns void as $$
    begin
        insert into t_basket_history
        select username, item_code, count, now() from t_basket where username = input_username;
    end;
$$ language plpgsql;