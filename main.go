package main

import (
	"context"
	"flag"
	"github.com/DarkOugi/OZON/pkg/app"
	"github.com/DarkOugi/OZON/pkg/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	serverHost = "localhost"

	serverPortXml              = "8080"
	serverPortXmlSlow          = "8081"
	serverPortXmlNotAvailable  = "8082"
	serverPortXmlWithoutValute = "8083"

	serverPortGrpc              = "4080"
	serverPortGrpcSlow          = "4081"
	serverPortGrpcNotAvailable  = "4082"
	serverPortGrpcWithoutValute = "4083"
)

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
	flag.StringVar(&serverPortXmlSlow, "serverPortXmlSlow",
		serverPortXmlSlow, "serverPortXmlSlow run in this port")
	flag.StringVar(&serverPortXmlNotAvailable, "serverPortXmlNotAvailable",
		serverPortXmlNotAvailable, "serverPortXmlNotAvailable run in this port")
	flag.StringVar(&serverPortXmlWithoutValute, "serverPortXmlWithoutValute",
		serverPortXmlWithoutValute, "serverPortXmlWithoutValute run in this port")

	flag.StringVar(&serverPortGrpc, "serverPortGrpc", serverPortGrpc, "serverGrpc run in this port")
	flag.StringVar(&serverPortGrpcSlow, "serverPortGrpcSlow",
		serverPortGrpcSlow, "serverPortGrpcSlow run in this port")
	flag.StringVar(&serverPortGrpcNotAvailable, "serverPortGrpcNotAvailable",
		serverPortGrpcNotAvailable, "serverPortGrpcNotAvailable run in this port")
	flag.StringVar(&serverPortGrpcWithoutValute, "serverPortGrpcWithoutValute",
		serverPortGrpcWithoutValute, "serverPortGrpcWithoutValute run in this port")
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

	app.RunServer(pSQL, 0, serverPortXml, log.Logger)
	app.RunServer(pSQL, 1, serverPortXmlSlow, log.Logger)
	app.RunServer(pSQL, 2, serverPortXmlNotAvailable, log.Logger)
	app.RunServer(pSQL, 3, serverPortXmlWithoutValute, log.Logger)

	app.RunServerGrpc(serverHost, serverPortGrpc, pSQL, 0)
	app.RunServerGrpc(serverHost, serverPortGrpcSlow, pSQL, 1)
	app.RunServerGrpc(serverHost, serverPortGrpcNotAvailable, pSQL, 2)
	app.RunServerGrpc(serverHost, serverPortGrpcWithoutValute, pSQL, 3)

	log.Info().Msg("SERVER SUCCESS START")
	<-ctx.Done()
}
