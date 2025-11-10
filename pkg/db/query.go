package db

import (
	"context"
	"github.com/DarkOugi/OZON/pkg/entity"
)

func (db *PostgresDB) GetValues(ctx context.Context, date string, dayInMonth int) ([]*entity.DailyValueSQL, error) {
	selectSQL := "SELECT\n" +
		"ValueId,\n" +
		"LPAD(vi.NumCode::text,3,'0'),\n" +
		"vi.CharCode,\n" +
		"dv.Nominal,\n" +
		"vi.Name,\n" +
		"REPLACE(TO_CHAR(dv.Value, 'FM999999999.0000'), '.', ','),\n" +
		"REPLACE(TO_CHAR(dv.Nominal * dv.value, 'FM999999999.0000'), '.', ',') AS VunitRate,\n" +
		"TO_CHAR(dv.Day, 'DD.MM.YYYY')\n" +
		"FROM DailyValue AS dv\n" +
		"JOIN ValueInfo as vi ON dv.ValueInfo = vi.NumCode\n" +
		"WHERE dv.Day = (SELECT dv.Day FROM DailyValue AS dv\n" +
		"WHERE dv.Day BETWEEN ($1::date - MAKE_INTERVAL(DAYS => $2)) AND $3::date\n" +
		"ORDER BY dv.Day DESC\n" +
		"LIMIT 1)"
	ValueRows, errQuery := db.conn.Query(ctx, selectSQL, date, dayInMonth, date)
	if errQuery != nil {
		return nil, errQuery
	}

	var dv []*entity.DailyValueSQL
	for ValueRows.Next() {
		var vid, nc, cc, name, val, vunit, day string
		var nom int
		errScan := ValueRows.Scan(&vid, &nc, &cc, &nom, &name, &val, &vunit, &day)
		if errScan != nil {
			return make([]*entity.DailyValueSQL, 0), errScan
		}
		dv = append(dv, &entity.DailyValueSQL{
			ValuteId:  vid,
			NumCode:   nc,
			CharCode:  cc,
			Nominal:   nom,
			Name:      name,
			Value:     val,
			VunitRate: vunit,
			Day:       day,
		})
	}

	return dv, nil
}

func (db *PostgresDB) GetMostPastAndLatestDate(ctx context.Context) (string, string, error) {
	selectPast := "SELECT\n" +
		"TO_CHAR(MIN(dv.Day), 'YYYY-MM-DD') AS past,\n" +
		"TO_CHAR(MAX(dv.Day), 'YYYY-MM-DD') AS latest\n" +
		"from  DailyValue as dv"

	var past, latest string

	err := db.conn.QueryRow(ctx, selectPast).Scan(&past, &latest)

	if err != nil {
		return "", "", err
	}

	return past, latest, nil
}
