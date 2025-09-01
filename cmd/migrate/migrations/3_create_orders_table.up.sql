create table if not exists orders (
    order_id serial primary key,
    order_user_id serial not null,
    order_status varchar(256) not null,
    order_total_price int not null,
    constraint fk_user foreign key (order_user_id) references users(user_id)
);