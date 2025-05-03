-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE chirps(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
     user_id UUID NOT NULL,
    body TEXT NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

);
-- +goose Down
DROP TABLE chirps;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
