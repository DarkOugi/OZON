-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ValueInfo
(
    ValueId CHAR(6) PRIMARY KEY CHECK ( ValueId ~ '^R[0-9]{5}[A-Z]*$' ),
    NumCode  SMALLINT UNIQUE,
    Name     VARCHAR NOT NULL,
    CharCode CHAR(3) NOT NULL UNIQUE CHECK ( CharCode ~ '^[A-Z]{3}$' )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ValueInfo;
-- +goose StatementEnd
