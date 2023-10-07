-- +goose Up
-- +goose StatementBegin
CREATE TABLE films (
    id INTEGER PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    iso INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE films;
-- +goose StatementEnd
