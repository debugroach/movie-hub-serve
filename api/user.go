package api

import (
	"database/sql"

	db "github.com/debugroach/video-hub-serve/db/sqlc"
	"github.com/debugroach/video-hub-serve/util"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func loginResponse(hasError bool, err string) gin.H {
	return gin.H{
		"hasError": hasError,
		"message":  err,
	}
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, loginResponse(true, err.Error()))
		return
	}

	user, err := server.GetUser(ctx, req.Username)

	if err == nil {
		if util.CheckPassword(req.Password, user.Password) == nil {
			ctx.JSON(200, loginResponse(false, "Login successful"))
			return
		}
		ctx.JSON(200, loginResponse(true, "Invalid password"))
		return
	}

	if err != sql.ErrNoRows {
		ctx.JSON(200, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(200, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
	}

	if _, err := server.CreateUser(ctx, arg); err != nil {
		ctx.JSON(200, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User created successfully"})
}
