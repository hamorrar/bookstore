create table if not exists orders (
    order_id serial unique primary key,
    order_user_id int not null,
    order_status varchar(256) not null,
    order_total_price int not null,
    foreign key (order_user_id) references users(user_id) on delete cascade
);