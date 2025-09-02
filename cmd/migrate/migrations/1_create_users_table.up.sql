create table if not exists users (
    user_id serial unique primary key,
    user_email varchar(256) unique not null,
    user_password varchar(256) not null,
    user_role varchar(12) not null
);