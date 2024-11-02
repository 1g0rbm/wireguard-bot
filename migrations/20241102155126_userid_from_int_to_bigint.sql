-- +goose Up
-- +goose StatementBegin
alter table users
    alter column id type bigint using id::bigint;

alter table users2servers
    alter column user_id type bigint using user_id::bigint;

alter table sessions
    alter column user_id type bigint using user_id::bigint;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    alter column id type int;

alter table users2servers
    alter column user_id type int;

alter table sessions
    alter column user_id type int;
-- +goose StatementEnd
