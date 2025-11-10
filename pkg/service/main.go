package service

import (
	"context"
	"github.com/DarkOugi/OZON/pkg/entity"
)

type Repository interface {
	GetValues(ctx context.Context, date string, dayInMonth int) ([]*entity.DailyValueSQL, error)
	GetMostPastAndLatestDate(ctx context.Context) (string, string, error)
}

//	typeSV {
//		0 : base
//		1 : slow
//		2 : not available
//	}
type Service struct {
	rep    Repository
	typeSV int
}

func NewService(rep Repository, typeSV int) *Service {
	return &Service{
		rep:    rep,
		typeSV: typeSV,
	}
}
