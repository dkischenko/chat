CREATE schema lets_go_chat;
SET search_path to lets_go_chat, public;
create table IF NOT EXISTS users
(
    id            uuid default gen_random_uuid() PRIMARY KEY,
    username      varchar(100) not null unique,
    password_hash varchar(100) not null
);

CREATE TABLE IF NOT EXISTS messages
(
    id         serial PRIMARY KEY,
    text       varchar not null,
    u_from     uuid    not null references users (id),
    created_at integer default null
);

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS key varchar default '';

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_online bool default false;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS last_online integer default 0;