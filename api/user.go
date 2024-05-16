package api

import (
	"database/sql"

	"github.com/debugroach/movie-hub-serve/db"
	"github.com/debugroach/movie-hub-serve/util"
	"github.com/gin-gonic/gin"
)

// loginRequest represents the request body for the login endpoint.
type loginRequest struct {
	Username string `json:"username" binding:"required"` // Username is the username of the user.
	Password string `json:"password" binding:"required"` // Password is the password of the user.
}

// login handles the login endpoint.
func (s *Server) login(ctx *gin.Context) {
	var req loginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	user, err := s.GetUser(ctx, req.Username)

	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(200, gin.H{"message": err.Error()})
			return
		}

		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(200, gin.H{"message": err.Error()})
			return
		}

		arg := db.CreateUserParams{
			Username: req.Username,
			Password: hashedPassword,
		}

		if _, err := s.CreateUser(ctx, arg); err != nil {
			ctx.JSON(200, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "User created successfully"})
		return
	}

	if util.CheckPassword(req.Password, user.Password) == nil {
		ctx.JSON(200, gin.H{"message": "Login successful"})
		return
	}
	ctx.JSON(200, errorResponse("Invalid username or password"))

}
