package xml

import (
	"bytes"
	"github.com/DarkOugi/OZON/pkg/entity"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"testing"
	"time"
)

func check1251(date []byte) bool {
	_, err := io.ReadAll(transform.NewReader(bytes.NewReader(date), charmap.Windows1251.NewDecoder()))
	if err != nil {
		return false
	}
	return true
}
func TestBuilder_createXML(t *testing.T) {

	name := "Foreign Currency Market"
	date := time.Now().Format("02.01.2006")
	t.Run("Correct params", func(t *testing.T) {
		values := []*entity.DailyValueSQL{
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
		}
		res, err := CreateXML(name, date, values)

		assert.Nil(t, err)
		assert.True(t, check1251(res))
	})
	t.Run("Empty data", func(t *testing.T) {
		res, err := CreateXML(name, date, nil)

		assert.Nil(t, err)
		assert.True(t, check1251(res))
	})
	t.Run("Correct params", func(t *testing.T) {
		res, err := CreateErrorXml()

		assert.Nil(t, err)
		assert.True(t, check1251(res))
	})
}
