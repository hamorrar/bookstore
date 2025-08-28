create table users (
    user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email varchar(256) unique not null,
    user_password_hash varchar(256) not null,
    user_role varchar(6) not null,
    user_created_at timestamptz CURRENT_TIMESTAMP,
    user_updated_at timestamptz
);

insert into users (user_email, user_password_hash, user_role) values ("testuser1@gmail.com", "testuser1password", "customer")
insert into users (user_email, user_password_hash, user_role) values ("testuser2@gmail.com", "testuser2password", "customer")
insert into users (user_email, user_password_hash, user_role) values ("testuser3@gmail.com", "testuser3password", "admin")