package helpers

import (
	"errors"
	"github.com/DarkOugi/OZON/pkg/entity"
	"github.com/DarkOugi/OZON/pkg/grpc/pb"
	"time"
)

var errBadDateFormat = errors.New("Error Parse Date")
var (
	DateInputFormat  = "02/01/2006"
	DateOutputFormat = "02.01.2006"
	DateSqlFormat    = "2006-01-02"
)

type TimeHelper struct {
	inputTime time.Time
}

func NewTimeHelper(date string) (*TimeHelper, error) {
	if date == "" {
		return &TimeHelper{
			inputTime: time.Now(),
		}, nil
	}
	t, errParse := time.Parse(DateInputFormat, date)
	if errParse != nil {
		return nil, errBadDateFormat
	} else {
		return &TimeHelper{
			inputTime: t,
		}, nil
	}

}
func (th *TimeHelper) GetDayInMonth() int {
	newDate := time.Date(th.inputTime.Year(), th.inputTime.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	return newDate.Day()
}

func (th *TimeHelper) ConvertToSqlDateFormat() string {
	return th.inputTime.Format(DateSqlFormat)
}

func (th *TimeHelper) FutureDate() string {
	newDate := th.inputTime.AddDate(0, 0, -th.GetDayInMonth())

	return newDate.Format(DateOutputFormat)
}

// 0 eq
// 1 one > two
// 2 one < two
func (th *TimeHelper) CompareDate(two string) int {
	tOne := th.inputTime
	//tOne, _ := time.Parse(DateSqlFormat, one)
	tTwo, _ := time.Parse(DateSqlFormat, two)

	if tOne.Equal(tTwo) {
		return 0
	} else if status := tOne.After(tTwo); status {
		return 2
	} else {
		return 1
	}
}

func ConvertSqlDvToResponseMock(name string, dv []*entity.DailyValueSQL) *pb.ResponseDailyValues {
	value := []*pb.ResponseDailyValues_Value{}
	var date string

	for _, el := range dv {
		if date == "" {
			date = el.Day
		}
		value = append(value, &pb.ResponseDailyValues_Value{
			ID: el.ValuteId,
			MetaValue: &pb.ResponseDailyValues_Value_MetaValue{
				NumCode:   el.NumCode,
				CharCode:  el.CharCode,
				Nominal:   int64(el.Nominal),
				Name:      el.Name,
				Value:     el.Value,
				VunitRate: el.VunitRate,
			},
		})
	}

	return &pb.ResponseDailyValues{
		Valute: value,
		Name:   name,
		Date:   date,
	}
}
