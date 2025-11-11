package service

import (
	"context"
	"github.com/DarkOugi/OZON/pkg/entity"
	"github.com/stretchr/testify/assert"
	"testing"
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
	svEmptyXml := NewService(&TestRep{}, 3)
	svUnknownType := NewService(&TestRep{}, 4)

	baseContext := context.WithValue(context.Background(), "behavior", 0)
	nilContext := context.WithValue(context.Background(), "behavior", 1)
	errContext := context.WithValue(context.Background(), "behavior", 2)
	t.Run("Not Work", func(t *testing.T) {
		_, err := svNotWork.GetDailyValue(baseContext, "24/12/2000")
		assert.ErrorIs(t, ErrorServerNotAvailable, err)

	})

	t.Run("Empty XML", func(t *testing.T) {
		res, err := svEmptyXml.GetDailyValue(baseContext, "24/12/2000")
		assert.Nil(t, res)
		assert.Nil(t, err)
	})

	t.Run("UnknownType", func(t *testing.T) {
		res, err := svUnknownType.GetDailyValue(baseContext, "24/12/2000")
		assert.ErrorIs(t, ErrorUnknownType, err)
		assert.Nil(t, res)
	})

	t.Run("Stable Work", func(t *testing.T) {
		res, err := svWork.GetDailyValue(baseContext, "24/12/2024")

		assert.Nil(t, err)
		assert.Equal(t, "2024-12-24", res[0].Day)
	})

	t.Run("Slow Work", func(t *testing.T) {
		res, err := svSlowWork.GetDailyValue(baseContext, "24/12/2024")

		assert.Nil(t, err)
		assert.Equal(t, "2024-12-24", res[0].Day)
	})

	t.Run("Stable Work Input > date in base", func(t *testing.T) {
		res, err := svWork.GetDailyValue(nilContext, "24/12/2025")

		assert.Nil(t, err)
		assert.Equal(t, "23.11.2025", res[0].Day)
	})

	t.Run("Stable Work Input < date in base", func(t *testing.T) {
		_, err := svWork.GetDailyValue(nilContext, "24/12/2020")

		assert.ErrorIs(t, ErrorDate, err)
	})

	t.Run("Stable Work err GetValue", func(t *testing.T) {
		_, err := svWork.GetDailyValue(errContext, "24/12/2020")

		assert.Error(t, err)
	})

	t.Run("Stable Work not correct date", func(t *testing.T) {
		_, err := svWork.GetDailyValue(errContext, "24-12-2020")

		assert.ErrorIs(t, ErrorDate, err)
	})

}
