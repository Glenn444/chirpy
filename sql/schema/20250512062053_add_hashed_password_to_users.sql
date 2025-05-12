-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE users ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE users DROP COLUMN hashed_password;