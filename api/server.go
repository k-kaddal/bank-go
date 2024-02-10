package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/k-kaddal/bank-go/db/sqlc"
	"github.com/k-kaddal/bank-go/token"
	"github.com/k-kaddal/bank-go/util"
)

type Server struct {
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	var tokenMaker token.Maker
	var err error

	switch config.Token_type {
	case "JWT":
		tokenMaker, err = token.NewJWTMaker(config.Token_Symmetric_Key)
	case "PASETO":
		tokenMaker, err = token.NewPasetoMaker(config.Token_Symmetric_Key)
	}

	if err != nil {
		return nil, fmt.Errorf("cannot create a paseto error: %w", err)
	}

	server := &Server{
		config: config,
		store:store,
		tokenMaker: tokenMaker,
	}

	if v, ok :=binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	
	router.POST("/transfers", server.createTransfer)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}