package api

import (
	"fmt"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/token"
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store          db.Store
	router         *gin.Engine
	tokenGenerator token.Generator
	config         util.Config
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenGenerator, err := token.NewJWTGenerator(config.TokenKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token generator: %w", err)
	}

	server := &Server{
		store:          store,
		tokenGenerator: tokenGenerator,
		config:         config,
	}

	server.setupRoute()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return server, nil
}

func (server *Server) setupRoute() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenGenerator))

	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccount)
	authRouter.PUT("/accounts/:id", server.updateAccount)
	authRouter.DELETE("/accounts/:id", server.deleteAccount)

	authRouter.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error) *gin.H {
	return &gin.H{
		"error": err.Error(),
	}
}
