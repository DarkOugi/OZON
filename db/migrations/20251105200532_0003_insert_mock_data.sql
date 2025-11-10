-- +goose Up
-- +goose StatementBegin
INSERT INTO ValueInfo (ValueId, NumCode, Name, CharCode)
VALUES
    ('R00000', 643, 'Российский рубль', 'RUB'),
    ('R01235', 840, 'Доллар США', 'USD'),
    ('R01239', 978, 'Евро', 'EUR')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO DailyValue (ValueInfo, Nominal, Value, Day)
SELECT
    vi.NumCode,
    1 AS Nominal,
    CASE
        WHEN vi.CharCode = 'RUB' THEN 1.0
        WHEN vi.CharCode = 'USD' THEN (70 + random() * 50)   -- 70–120
        WHEN vi.CharCode = 'EUR' THEN (75 + random() * 55)   -- 75–130
        END AS Value,
    d::date AS Day
FROM generate_series('2024-01-01'::date, '2024-12-31', '1 day') d
         JOIN ValueInfo vi ON vi.CharCode IN ('USD', 'EUR', 'RUB')
WHERE EXTRACT(DOW FROM d) NOT IN (0, 6) -- исключаем выходные
ON CONFLICT (ValueInfo, Day) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM DailyValue;
DELETE
FROM ValueInfo
WHERE CharCode IN ('RUB', 'USD', 'EUR');
-- +goose StatementEnd
