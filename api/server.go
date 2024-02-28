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

	/*
		todo: it is bettern to switch the authentication mechanism according to the config.TokenType 
		todo: the pending issue here is that the test would fail. 
		todo: env vars need to be passed in tests
	*/

	// switch config.Token_Type {
	// case "JWT":
	// 	tokenMaker, err = token.NewJWTMaker(config.TokenSymmetricKey)
	// case "PASETO":
	// 	tokenMaker, err = token.NewPasetoMaker(config.TokenSymmetricKey)
	// }
	tokenMaker, err = token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create a token maker error: %w", err)
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
	router.POST("/token/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	
	authRoutes.POST("/transfers", server.createTransfer)
	
	// todo : create a post router "/entry" to deposit or withdraw to an account's balance
	authRoutes.POST("/entries", server.createEntry)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}