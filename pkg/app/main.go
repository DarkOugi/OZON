package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/DarkOugi/OZON/pkg/db"
	proto "github.com/DarkOugi/OZON/pkg/grpc/pb"
	"github.com/DarkOugi/OZON/pkg/server"
	"github.com/DarkOugi/OZON/pkg/service"
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"net"
)

func RunServer(pSQL *db.PostgresDB, typeSV int, port string, log zerolog.Logger) {
	svXml := service.NewService(pSQL, typeSV)

	srXml := server.NewXmlServer(svXml)

	rXml := router.New()
	rXml.GET("/scripts/XML_daily.asp", srXml.GetDailyValueXml)

	fs := &fasthttp.FS{
		Root:               "",
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: false,
	}
	handler := fs.NewRequestHandler()

	rXml.GET("/docs/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
		handler(ctx)
	})

	go func() {
		if errServer := fasthttp.ListenAndServe(fmt.Sprintf(":%s", port), rXml.Handler); errServer != nil {
			log.Fatal().Err(errServer).Msg("server critical error")
		}
	}()
}

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	msg := "log interceptor"
	if err != nil {
		log.Info().Str("full Method", info.FullMethod).Any("request", req).Err(err).Msg(msg)
	} else {
		log.Info().Str("full method", info.FullMethod).Any("request", req).Msg(msg)
	}
	return resp, err
}

func RunServerGrpc(host, port string, pSQL *db.PostgresDB, typeSV int) {
	grpcServe := service.NewService(pSQL, typeSV)
	netAddr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", netAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen TCP")
	}
	opts := make([]grpc.ServerOption, 0, 2)
	interceptors := []grpc.UnaryServerInterceptor{
		LogInterceptor,
	}
	opts = append(opts, grpc.ChainUnaryInterceptor(interceptors...))
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterMockDailyValueServiceServer(grpcServer, server.NewProtoServer(grpcServe))
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil && errors.Is(err, grpc.ErrServerStopped) {
			log.Fatal().Err(err).Msg("Error while server working: %v")
		}
	}()
}
