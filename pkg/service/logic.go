package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/DarkOugi/OZON/pkg/entity"
	"github.com/DarkOugi/OZON/pkg/helpers"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	ErrorServerNotAvailable = errors.New("service not available")
	ErrorUnknownType        = errors.New("this type not support")
	ErrorDate               = errors.New("date not good")
	//ErrorInParameters       = errors.New("date didn't correct")
)

var Name = "Foreign Currency Market"

func (sv *Service) GetDailyValue(ctx context.Context, date string) ([]*entity.DailyValueSQL, error) {
	switch sv.typeSV {
	case 0:
		return sv.GetDailyValueWork(ctx, date)
	case 1:
		time.Sleep(10 * time.Second)
		return sv.GetDailyValueWork(ctx, date)
	case 2:
		return nil, ErrorServerNotAvailable
	case 3:
		return nil, nil
	default:
		return nil, ErrorUnknownType
	}
}

//Давайте немного подумаем над логикой
//Если это выходной мы ищем котировку за последний рабочий день - пятница например
//Я не хочу держать логику определения даты выдачи на беке, это должно быть в бд
//по типу - если есть сущность за этот день, мы её выдаем
//если нет - выдаем раньше
//в самом цб - выдается предыдущая дата за период не больше 33 дней
//то есть - если сегодня 5 ноября а я спрошу 8 декабря (+ 33 дня) - то мне дадут ответ за 5 ноября
//но если позже, то уже пустую сущность
//Я в данный момент не совсем понимаю логику даты в хэдере xml т.к. при запросе например данных за 2045 год
//она будет выдавать рандомную дату 2045(в моем понимании она такова)

//Давайте подытожим логику:
//	Если есть валидный запрос:
//		данные есть в бд - отдаем валидный ответ
//		данных нет в бд:
//			дата больше текущей - отдаем пустой массив
//			(так как я не понимаю суть даты в этом массиве в таких случаях - будет такая же как и запрашивали)
//			дата меньше самой ранней в бд - отдаем ошибку Error in parameters
//	Если запрос не валидный:
//		отдаем также Error in parameters

func (sv *Service) GetDailyValueWork(ctx context.Context, date string) ([]*entity.DailyValueSQL, error) {
	dateHelper, err := helpers.NewTimeHelper(date)
	if err != nil {
		log.Err(err).Msg("Error convert date")
		return nil, ErrorDate
	}
	dv, errSQL := sv.rep.GetValues(ctx, dateHelper.ConvertToSqlDateFormat(), dateHelper.GetDayInMonth())
	if errSQL != nil {
		return nil, fmt.Errorf("can't get DailyValue: %w", errSQL)
	}
	if dv != nil {
		return dv, nil
	}
	past, latest, errPL := sv.rep.GetMostPastAndLatestDate(ctx)
	if errPL != nil {
		return nil, fmt.Errorf("can't get GetMostPastAndLatestDate: %w", errPL)
	}
	pastCompare := dateHelper.CompareDate(past)
	latestCompare := dateHelper.CompareDate(latest)
	fmt.Println(past)
	fmt.Println(latest)
	if pastCompare == 0 || latestCompare == 0 || latestCompare == 2 || pastCompare == 2 {
		log.Info().Msg("THIS")
		return []*entity.DailyValueSQL{
			{Day: dateHelper.FutureDate()},
		}, nil
	} else {
		return nil, ErrorDate
	}

}
