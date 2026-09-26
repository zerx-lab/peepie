-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
    ADD COLUMN retry_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN last_error  TEXT    NOT NULL DEFAULT '';

ALTER TABLE subtasks
    ADD COLUMN retry_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN last_error  TEXT    NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE subtasks
    DROP COLUMN retry_count,
    DROP COLUMN last_error;

ALTER TABLE tasks
    DROP COLUMN retry_count,
    DROP COLUMN last_error;
-- +goose StatementEnd
