-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS DailyValue
(
    id        SERIAL PRIMARY KEY,
    ValueInfo INTEGER REFERENCES ValueInfo (NumCode) NOT NULL,
    Nominal   INTEGER CHECK ( Nominal >= 0 ),
    Value     DOUBLE PRECISION CHECK ( Value > 0 ),
    VunitRate DOUBLE PRECISION CHECK ( VunitRate > 0 ),
    Day       DATE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS DailyValue;
-- +goose StatementEnd
