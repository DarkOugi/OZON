package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	conn *pgxpool.Pool
}

// type Fake struct {
// }
//
//	func (f *Fake) GetValues(ctx context.Context, date string) ([]*entity.DailyValueSQL, error) {
//		return []*entity.DailyValueSQL{
//			{
//				Day: date,
//			},
//		}, nil
//	}
func NewPostgresDB(ctx context.Context, host, port, username, password, dbname string) (*PostgresDB, error) {
	strConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbname)

	conn, err := pgxpool.New(ctx, strConn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	return &PostgresDB{conn: conn}, nil
}

func (db *PostgresDB) Close() {
	db.conn.Close()
}
