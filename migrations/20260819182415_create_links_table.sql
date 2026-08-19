-- +goose Up
-- +goose StatementBegin
CREATE TABLE Links
(
    id        UUID PRIMARY KEY NOT NULL,
    link      VARCHAR(200)     NOT NULL,
    link_type smallint         NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Links;
-- +goose StatementEnd