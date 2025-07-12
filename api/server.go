package api

import (
	db "simle_bank/db/sqlc"
	"simle_bank/token"
	"simle_bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config util.Config
	db db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		panic(err)
	}
	server := &Server{
		config: config,
		db: store,
		tokenMaker: tokenMaker,
		router: gin.Default(),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {

	server.router.POST("/users", server.createUser)
	server.router.GET("/users/:username", server.getUser)
	server.router.POST("/users/login", server.loginUser)
	server.router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := server.router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/transfers", server.createTransfer)
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}