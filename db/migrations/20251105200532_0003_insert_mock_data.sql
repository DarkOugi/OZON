-- +goose Up
-- +goose StatementBegin
INSERT INTO ValueInfo (Name, CharCode)
VALUES ('Российский рубль', 'RUB'),
       ('Доллар США', 'USD'),
       ('Евро', 'EUR');
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO DailyValue (ValueInfo, Nominal, Value, VunitRate, Day)
SELECT v.NumCode, 1, (100 + i), (100 + i), (DATE '2024-01-01' + i)
FROM generate_series(0, 365) i
         JOIN ValueInfo v ON v.CharCode = 'RUB';

INSERT INTO DailyValue (ValueInfo, Nominal, Value, VunitRate, Day)
SELECT v.NumCode, 1, (90 + i * 0.05), (90 + i * 0.05), (DATE '2024-01-01' + i)
FROM generate_series(0, 365) i
         JOIN ValueInfo v ON v.CharCode = 'USD';

INSERT INTO DailyValue (ValueInfo, Nominal, Value, VunitRate, Day)
SELECT v.NumCode, 1, (100 + i * 0.03), (100 + i * 0.03), (DATE '2024-01-01' + i)
FROM generate_series(0, 365) i
         JOIN ValueInfo v ON v.CharCode = 'EUR';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM DailyValue;
DELETE
FROM ValueInfo
WHERE CharCode IN ('RUB', 'USD', 'EUR');
-- +goose StatementEnd
