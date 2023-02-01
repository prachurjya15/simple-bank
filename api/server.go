package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/prachurjya15/simple-bank/db/sqlc"
	"github.com/prachurjya15/simple-bank/token"
	"github.com/prachurjya15/simple-bank/util"
)

type Server struct {
	config     util.Config
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.TokenMaker
}

func NewServer(store *db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize token maker %s", err)
	}
	server := Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	r.POST("/login", server.LoginUser)
	r.POST("/users", server.CreateUser)
	r.GET("/users/:username", server.GetUser)

	authGroup := r.Group("/").Use(AuthMiddleware(server.tokenMaker))
	authGroup.POST("/accounts", server.CreateAccount)
	authGroup.GET("/accounts/:id", server.GetAccount)
	authGroup.GET("/accounts", server.listAccounts)

	authGroup.POST("/transfers", server.CreateTransfer)

	server.router = r
	return &server, nil
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
