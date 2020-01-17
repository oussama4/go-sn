-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table users (
    id serial primary key,
    user_name varchar(50) unique not null ,
    email varchar(50) unique not null ,
    password text not null,
    bio text,
    avatar text default 'img/avatars/default.jpeg' not null,
    created_at timestamptz default now() not null
);

create table connections (
    user_one int not null references users(id) on delete cascade,
    user_two int not null references users(id) on delete cascade,
    status  smallint not null,
    primary key (user_one, user_two)
);

CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  data BYTEA NOT NULL,
  expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table connections;
drop table sessions;
drop table users;