package helpers

import (
	"github.com/DarkOugi/OZON/pkg/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTimeHelper(t *testing.T) {
	t.Run("Get correct day in month", func(t *testing.T) {
		_, err := NewTimeHelper("25/12/2000")
		assert.Nil(t, err)
	})

	t.Run("Get un correct day in month", func(t *testing.T) {
		_, err := NewTimeHelper("25-12-2000")
		assert.Error(t, errBadDateFormat, err)
	})
}
func TestTimeHelper_GetDayInMonth(t *testing.T) {
	th, _ := NewTimeHelper("25/12/2000")
	t.Run("Get count day in month", func(t *testing.T) {
		day := th.GetDayInMonth()
		assert.Equal(t, 31, day)
	})
}

func TestTimeHelper_ConvertToSqlDateFormat(t *testing.T) {
	th, _ := NewTimeHelper("25/12/2000")
	t.Run("Get count day in month", func(t *testing.T) {
		date := th.ConvertToSqlDateFormat()
		assert.Equal(t, "2000-12-25", date)
	})
}

func TestHelpers_FutureDate(t *testing.T) {
	th, _ := NewTimeHelper("25/12/2000")
	t.Run("Correct date", func(t *testing.T) {
		res := th.FutureDate()
		assert.Equal(t, "24.11.2000", res)
	})
}

func TestTimeHelper_CompareDate(t *testing.T) {
	th, _ := NewTimeHelper("25/12/2000")
	t.Run("eq date", func(t *testing.T) {
		res := th.CompareDate("2000-12-25")
		assert.Equal(t, 0, res)
	})
	t.Run("after", func(t *testing.T) {
		res := th.CompareDate("2000-12-24")
		assert.Equal(t, 2, res)
	})
	t.Run("after", func(t *testing.T) {
		res := th.CompareDate("2000-12-26")
		assert.Equal(t, 1, res)
	})
}

func TestHelpers_ConvertSqlDvToResponseMock(t *testing.T) {
	t.Run("Correct date", func(t *testing.T) {
		res := ConvertSqlDvToResponseMock("test", []*entity.DailyValueSQL{
			{
				ValuteId:  "R01235",
				NumCode:   "840",
				CharCode:  "USD",
				Nominal:   1,
				Name:      "Доллар США",
				Value:     "30,9436",
				VunitRate: "30,9436",
				Day:       "02.03.2002",
			},
			{
				ValuteId:  "R01239",
				NumCode:   "978",
				CharCode:  "EUR",
				Nominal:   1,
				Name:      "Евро",
				Value:     "26,8343",
				VunitRate: "26,8343",
				Day:       "02.03.2002",
			},
			{
				ValuteId:  "R01350",
				NumCode:   "124",
				CharCode:  "CAD",
				Nominal:   1,
				Name:      "Канадский доллар",
				Value:     "19,3240",
				VunitRate: "19,3240",
				Day:       "02.03.2002",
			},
		})
		assert.Equal(t, len(res.Valute), 3)
	})

}
