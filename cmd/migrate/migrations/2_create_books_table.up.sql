create table if not exists books (
    book_id serial primary key,
    book_title varchar(256) not null,
    book_author varchar(256) not null,
    book_price int not null
);