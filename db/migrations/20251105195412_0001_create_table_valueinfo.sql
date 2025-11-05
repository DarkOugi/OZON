-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ValueInfo
(
    NumCode  SERIAL PRIMARY KEY,
    Name     VARCHAR NOT NULL,
    CharCode VARCHAR(3) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ValueInfo;
-- +goose StatementEnd
