package api

import (
	"log"
	"net/http"

	"github.com/debugroach/video-hub-serve/config"
	db "github.com/debugroach/video-hub-serve/db/sqlc"
	"github.com/debugroach/video-hub-serve/token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	db.Store
	token.Maker
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker", err)
	}

	server := &Server{
		Store: store,
		Maker: maker,
	}

	server.setupRoutes()
	return server
}

func (server *Server) setupRoutes() {
	r := gin.Default()
	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")                                                              // 允许所有域名的请求
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")                       // 允许的方法
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization") // 允许的头部
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	})

	r.POST("/login", server.login)
	r.POST("/rate", server.rate)
	r.POST("/recommend", server.recommend)
	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err string) gin.H {
	return gin.H{
		"hasError": "true",
		"message":  err,
	}
}
