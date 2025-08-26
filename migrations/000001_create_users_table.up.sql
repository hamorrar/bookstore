create table users (
    user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email varchar(256) unique not null,
    user_password_hash varchar(256) not null,
    user_role varchar(6) not null,
    user_created_at timestamptz CURRENT_TIMESTAMP,
    user_updated_at timestamptz
);