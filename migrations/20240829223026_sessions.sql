-- +goose Up
-- +goose StatementBegin
--------------------------- SESSIONS ------------------------------------
SELECT 'up SQL query';
create table if not exists sessions (
    id UUID primary key default gen_random_uuid(),
    user_id int not null,
    expired_at timestamp not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
