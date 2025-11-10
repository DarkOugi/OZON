package main

import (
	"context"
	"errors"
	"flag"
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
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUsername = "ozon"
	dbPassword = "0000"
	dbName     = "ozondb"

	serverPortXml              = "8080"
	serverPortXmlSlow          = "8081"
	serverPortXmlNotAvailable  = "8082"
	serverPortXmlWithoutValute = "8083"

	serverPortGrpc              = "9080"
	serverPortGrpcSlow          = "9081"
	serverPortGrpcNotAvailable  = "9082"
	serverPortGrpcWithoutValute = "9083"
)

func RunServer(pSQL *db.PostgresDB, typeSV int, port string) {
	svXml := service.NewService(pSQL, typeSV)

	srXml := server.NewXmlServer(svXml)

	rXml := router.New()
	rXml.GET("/scripts/XML_daily.asp", srXml.GetDailyValueXml)

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

func main() {
	var err error
	var pSQL *db.PostgresDB

	ctx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.Stamp,
	}).Level(zerolog.DebugLevel)

	flag.StringVar(&dbHost, "dbHost", dbHost, "dbHost pgx connect")
	flag.StringVar(&dbPort, "dbPort", dbPort, "dbPort pgx connect")
	flag.StringVar(&dbUsername, "dbUsername", dbUsername, "dbUsername pgx connect")
	flag.StringVar(&dbPassword, "dbPassword", dbPassword, "dbPassword pgx connect")
	flag.StringVar(&dbName, "dbName", dbName, "dbName pgx connect")

	flag.StringVar(&serverPortXml, "serverPortXml", serverPortXml, "serverXml run in this port")
	flag.StringVar(&serverPortGrpc, "serverPortGrpc", serverPortGrpc, "serverGrpc run in this port")
	flag.Parse()

	if pSQL == nil {
		pSQL, err = db.NewPostgresDB(ctx, dbHost, dbPort, dbUsername, dbPassword, dbName)
		if err != nil {
			log.Error().Err(err).Msg("don't create connect with db")
			return
		}
	}
	defer func() {
		pSQL.Close()
	}()
	//pSQL := db.Fake{}
	RunServer(pSQL, 0, serverPortXml)
	RunServer(pSQL, 1, serverPortXmlSlow)
	RunServer(pSQL, 2, serverPortXmlNotAvailable)
	RunServer(pSQL, 3, serverPortXmlWithoutValute)

	RunServerGrpc(dbHost, serverPortGrpc, pSQL, 0)
	//svXml := service.NewService(pSQL, 0)
	//srXml := server.NewXmlServer(svXml)
	//
	//svXmlSlow := service.NewService(pSQL, 1)
	//srXmlSlow := server.NewXmlServer(svXmlSlow)
	//
	//svXmlNotAvailable := service.NewService(pSQL, 2)
	//srXmlNotAvailable := server.NewXmlServer(svXmlNotAvailable)
	//
	//svXmlWithoutValute := service.NewService(pSQL, 3)
	//srXmlWithoutValute := server.NewXmlServer(svXmlWithoutValute)
	//
	//rXml := router.New()
	//rXml.GET("/scripts/XML_daily.asp", srXml.GetDailyValueXml)
	//
	//rXmlSlow := router.New()
	//rXmlSlow.GET("/scripts/XML_daily.asp", srXmlSlow.GetDailyValueXml)
	//
	//rXmlNotAvailable := router.New()
	//rXmlNotAvailable.GET("/scripts/XML_daily.asp", srXmlNotAvailable.GetDailyValueXml)
	//
	//rXmlWithoutValute := router.New()
	//rXmlWithoutValute.GET("/scripts/XML_daily.asp", srXmlWithoutValute.GetDailyValueXml)
	//
	//go func() {
	//	if errServer := fasthttp.ListenAndServe(fmt.Sprintf(":%s", serverPortXml), rXml.Handler); errServer != nil {
	//		log.Fatal().Err(errServer).Msg("server critical error")
	//	}
	//}()
	//go func() {
	//	if errServer := fasthttp.ListenAndServe(fmt.Sprintf(":%s", serverPortXmlSlow), rXmlSlow.Handler); errServer != nil {
	//		log.Fatal().Err(errServer).Msg("server critical error")
	//	}
	//}()
	//go func() {
	//	if errServer := fasthttp.ListenAndServe(
	//		fmt.Sprintf(":%s", serverPortXmlNotAvailable), rXmlNotAvailable.Handler); errServer != nil {
	//		log.Fatal().Err(errServer).Msg("server critical error")
	//	}
	//}()
	//
	//go func() {
	//	if errServer := fasthttp.ListenAndServe(
	//		fmt.Sprintf(":%s", serverPortXmlWithoutValute), rXmlWithoutValute.Handler); errServer != nil {
	//		log.Fatal().Err(errServer).Msg("server critical error")
	//	}
	//}()

	log.Info().Msg("SERVER SUCCESS START")
	<-ctx.Done()
}
