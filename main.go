package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/alrasyidin/simplebank-go/api"
	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/docs"
	"github.com/alrasyidin/simplebank-go/gapi"
	"github.com/alrasyidin/simplebank-go/pb"
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	store := db.NewStore(conn)
	go runGatewayServer(store, config)
	runGRPCServer(store, config)
	// runGinServer(store, config)
}

func runGinServer(store db.Store, config util.Config) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal().Msg("cannot start the server")
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}

func runGRPCServer(store db.Store, config util.Config) {
	logger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	grpcServer := grpc.NewServer(logger)
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}
}

func runGatewayServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}
	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register handler server:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	swaggerHandler := http.FileServer(http.FS(docs.StaticSwagger))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener:")
	}

	log.Info().Msgf("Start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Msg("cannot start HTTP gateway server:")
	}
}
