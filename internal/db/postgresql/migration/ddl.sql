select 'create database "simple-twitter"'
where not exists (select from pg_database where datname = 'simple-twitter')\gexec

create extension if not exists "uuid-ossp";

create table users (
    id           uuid primary key default uuid_generate_v4(),
    display_name varchar(255) not null,
    username     varchar(100) not null unique,
    born_date    date,
    address      text,
    passwd varchar(20) not null
);

create table posts (
    id         uuid primary key default uuid_generate_v4(),
    content    text not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    user_id    uuid not null,
    constraint fk_user foreign key(user_id) references users(id) on delete cascade
);

create index idx_posts_user_id on posts(user_id);
create index idx_users_username on users(username);
