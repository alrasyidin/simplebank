package gapi

import (
	"fmt"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/pb"
	"github.com/alrasyidin/simplebank-go/token"
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/alrasyidin/simplebank-go/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	tokenGenerator  token.Generator
	config          util.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, config util.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenGenerator, err := token.NewJWTGenerator(config.TokenKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token generator: %w", err)
	}

	server := &Server{
		store:           store,
		tokenGenerator:  tokenGenerator,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
