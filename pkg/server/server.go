package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/DarkOugi/OZON/pkg/grpc/pb"
	"github.com/DarkOugi/OZON/pkg/helpers"
	"github.com/DarkOugi/OZON/pkg/service"
	"github.com/DarkOugi/OZON/pkg/xml"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (sr *ProtoServer) GetDailyValue(ctx context.Context, data *pb.RequestDate) (*pb.ResponseDailyValues, error) {
	dv, err := sr.sv.GetDailyValue(ctx, data.Date)
	if err != nil {
		if errors.Is(err, service.ErrorServerNotAvailable) {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("service unavailavle: %w", err))
		} else if errors.Is(err, service.ErrorUnknownType) {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("err with server type: %w", err))
		}
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("err with get data: %w", err))

	}
	return helpers.ConvertSqlDvToResponseMock(service.Name, dv), nil
}

func (sr *XmlServer) GetDailyValueXml(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/xml; charset=windows-1251")
	// для swagger
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")

	date := string(ctx.QueryArgs().Peek("date_req"))

	dv, err := sr.sv.GetDailyValue(ctx, date)

	if err != nil {
		if errors.Is(err, service.ErrorServerNotAvailable) {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
			//return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("service unavailavle: %w", err))
		} else if errors.Is(err, service.ErrorUnknownType) {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		} else if errors.Is(err, service.ErrorDate) {
			ctx.SetStatusCode(fasthttp.StatusOK)

			dataErr, errXml := xml.CreateErrorXml()
			if errXml != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			}
			ctx.SetBody(dataErr)
			return
		}
		log.Err(err).Msg("Unknown Error")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return

	}

	buildXml, errXml := xml.CreateXML(service.Name, date, dv)
	if errXml != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(buildXml)
	return
}
