package gapi

import (
	"fmt"

	db "github.com/k-kaddal/bank-go/db/sqlc"
	"github.com/k-kaddal/bank-go/pb"
	"github.com/k-kaddal/bank-go/token"
	"github.com/k-kaddal/bank-go/util"
)

type Server struct {
	pb.UnimplementedBankGoServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	var tokenMaker token.Maker
	var err error

	/*
		todo: it is better to switch the authentication mechanism according to the config.TokenType 
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


	return server, nil
}