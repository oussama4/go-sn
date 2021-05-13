-- +goose Up
-- SQL in this section is executed when the migration is applied.

create table activities (
    id serial primary key,
    type varchar(20) not null,
    actor int not null references users(id) on delete cascade,
    content bytea not null,
    created timestamptz default now() not null
);

create index activity_actor_idx on activities(actor);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

drop table activities;