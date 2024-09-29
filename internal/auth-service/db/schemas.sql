CREATE TABLE IF NOT EXISTS Users
(
    Id       serial primary key,
    Username varchar(255)   not null unique,
    Email    varchar(255)   not null unique,
    Password varchar(255)   not null,
    Balance  numeric(10, 2) not null default 0 check ( Balance >= 0),
    Role     text[]         default '{"user"}',
    CreateAt timestamp      not null
);

create index if not exists Users_email_idx on Users using btree (lower(Email));
create index if not exists Users_username_idx on Users using btree (lower(username));

--drop table Users;