-- +goose Up
-- +goose StatementBegin
ALTER TABLE links
    ADD COLUMN chat_id int;
-- +goose StatementEnd

-- +goose Down
SELECT 'down SQL query';
