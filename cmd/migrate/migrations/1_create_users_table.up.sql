create table if not exists users (
    user_id serial primary key,
    user_email varchar(256) unique not null,
    user_password_hash varchar(256) not null,
    user_role varchar(6) not null
);