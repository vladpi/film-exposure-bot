-- +goose Up
-- +goose StatementBegin
CREATE TABLE films (
    id INTEGER PRIMARY KEY,
    name TEXT,
    iso INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE films;
-- +goose StatementEnd
