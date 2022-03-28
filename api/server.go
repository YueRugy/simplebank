package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"
)

type Server struct {
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(store *db.Store, config util.Config) (*Server, error) {
	server := &Server{store: store, config: config}
	maker, err := token.NewPasetoMaker(config.TokenSymKey)
	if err != nil {
		return nil, err
	}
	server.tokenMaker = maker
	if valid, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = valid.RegisterValidation("currency", validCurrency)
	}
	//设置路由
	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	authRouterGroup := router.Group("/",authMiddleWare(s.tokenMaker))
	authRouterGroup.POST("/account", s.createAccount)
	authRouterGroup.GET("/account/:id", s.getAccount)
	authRouterGroup.GET("/account", s.listAccount)
	authRouterGroup.POST("/transfer", s.createTransfer)
	s.router = router
}

func responseError(err error) gin.H {
	return gin.H{"err": err.Error()}
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
