package api

import (
	"log"
	"net/http"

	"github.com/debugroach/movie-hub-serve/config"
	"github.com/debugroach/movie-hub-serve/db"
	"github.com/debugroach/movie-hub-serve/token"
	"github.com/gin-gonic/gin"
)

// Server represents the API server.
type Server struct {
	*db.Queries
	token.Maker
	router *gin.Engine
}

// NewServer creates a new instance of the API server.
func NewServer(q *db.Queries) *Server {
	m, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker", err)
	}

	s := &Server{
		Queries: q,
		Maker:   m,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures the server's routes.
func (s *Server) setupRoutes() {
	r := gin.Default()
	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")                                                              // Allow requests from all domains
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")                       // Allowed methods
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization") // Allowed headers
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	})

	r.POST("/login", s.login)
	r.POST("/rate", s.rate)
	r.POST("/recommend", s.recommend)
	s.router = r
}

// Start starts the server on the specified address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// errorResponse returns a JSON response with an error message.
func errorResponse(err string) gin.H {
	return gin.H{
		"hasError": "true",
		"message":  err,
	}
}
