-- +goose Up
-- +goose StatementBegin
--------------------------- USERS ------------------------------------
create table if not exists users (
                                     id serial primary key,
                                     username varchar(255) ,
                                     first_name varchar(255),
                                     last_name varchar(255),
                                     role smallint not null,
                                     public_key varchar(44) not null unique,
                                     private_key varchar(44) not null unique,
                                     state varchar(20) not null default 'disabled',
                                     created_at timestamp not null default now(),
                                     updated_at timestamp
);

--------------------------- SERVERS ------------------------------------
create table if not exists servers (
                                       id serial primary key,
                                       name varchar(255) not null unique,
                                       address varchar(21) not null unique,
                                       public_key varchar(44) not null unique,
                                       private_key varchar(44) not null unique,
                                       created_at timestamp not null default now(),
                                       updated_at timestamp
);

--------------------------- USERS2SERVERS ------------------------------------
create table if not exists users2servers (
                                             user_id int not null,
                                             server_id int not null,
                                             address varchar(21) not null,
                                             primary key (user_id, server_id),
                                             foreign key (user_id) references users (id) on delete cascade,
                                             foreign key (server_id) references servers (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users2servers;
drop table if exists servers;
drop table if exists users;
-- +goose StatementEnd
