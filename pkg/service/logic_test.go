package service

import (
	"context"
	"github.com/DarkOugi/OZON/pkg/entity"
	"github.com/DarkOugi/OZON/pkg/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestRep struct {
}

func (r *TestRep) GetValues(ctx context.Context, date string, dayInMonth int) ([]*entity.DailyValueSQL, error) {
	switch ctx.Value("behavior") {
	case 0:
		return []*entity.DailyValueSQL{{Day: date}}, nil
	case 1:
		return nil, nil
	default:
		return nil, ErrorUnknownType
	}
}
func (r *TestRep) GetMostPastAndLatestDate(ctx context.Context) (string, string, error) {
	return "2024-01-01", "2024-12-12", nil
}
func TestService_GetDailyValue(t *testing.T) {
	svWork := NewService(&TestRep{}, 0)
	svSlowWork := NewService(&TestRep{}, 1)
	svNotWork := NewService(&TestRep{}, 2)

	baseContext := context.WithValue(context.Background(), "behavior", 0)
	nilContext := context.WithValue(context.Background(), "behavior", 1)
	errContext := context.WithValue(context.Background(), "behavior", 2)

	t.Run("Stable Work", func(t *testing.T) {
		res, err := svWork.GetDailyValue(baseContext, "24/12/2000")

		assert.Nil(t, err)
		assert.Equal(t, "24.12.2000", res[0].Day)
	})

	t.Run("Stable Work not correct date", func(t *testing.T) {
		_, err := svWork.GetDailyValue(baseContext, "24/13/2000")
		assert.ErrorIs(t, ErrorDate, err)
	})

	t.Run("Stable Work date in past with no data", func(t *testing.T) {
		_, err := svWork.GetDailyValue(nilContext, "24/12/2000")
		assert.ErrorIs(t, ErrorDate, err)
	})

	t.Run("Stable Work date in future", func(t *testing.T) {
		dateF := time.Now().Add(25 * time.Hour)
		_, err := svWork.GetDailyValue(nilContext, dateF.Format(helpers.DateInputFormat))
		assert.Nil(t, err)
	})

	t.Run("Stable Work errSql", func(t *testing.T) {
		dateF := time.Now().Add(25 * time.Hour)
		_, err := svWork.GetDailyValue(errContext, dateF.Format(helpers.DateInputFormat))
		assert.Error(t, err)
	})

	t.Run("Slow Work", func(t *testing.T) {
		res, err := svSlowWork.GetDailyValue(baseContext, "24/12/2000")
		assert.Nil(t, err)
		assert.Equal(t, "24.12.2000", res[0].Day)
	})

	t.Run("Not Work", func(t *testing.T) {
		_, err := svNotWork.GetDailyValue(baseContext, "24/12/2000")
		assert.ErrorIs(t, ErrorServerNotAvailable, err)

	})

}
