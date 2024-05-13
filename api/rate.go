package api

import (
	"database/sql"

	db "github.com/debugroach/video-hub-serve/db/sqlc"
	"github.com/gin-gonic/gin"
)

type rateRequest struct {
	Username     string  `json:"username" bingding:"required"`
	MovieID      int     `json:"movieID" binding:"required"`
	Rating       int     `json:"rating" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	PosterPath   string  `json:"posterPath" binding:"required"`
	BackdropPath string  `json:"backdropPath" binding:"required"`
	VoteAverage  float64 `json:"voteAverage" binding:"required"`
}

func (server *Server) createMovie(ctx *gin.Context, req rateRequest) error {
	_, err := server.GetMovie(ctx, req.MovieID)
	if err == nil {
		return nil
	}

	if err != sql.ErrNoRows {
		return err
	}

	arg := db.CreateMovieParams{
		ID:           req.MovieID,
		Title:        req.Title,
		PosterPath:   req.PosterPath,
		BackdropPath: req.BackdropPath,
		VoteAverage:  req.VoteAverage,
	}
	if _, err := server.CreateMovie(ctx, arg); err != nil {
		return err
	}

	return nil
}

func (server *Server) rate(ctx *gin.Context) {
	var req rateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	if err := server.createMovie(ctx, req); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	arg := db.CreateRatingParams{
		Username: req.Username,
		MovieID:  req.MovieID,
		Rating:   req.Rating,
	}

	_, err := server.GetRating(ctx, db.GetRatingParams{
		Username: arg.Username, MovieID: arg.MovieID,
	})

	if err == nil {
		if _, err = server.UpdateRating(ctx, db.UpdateRatingParams{
			Username: arg.Username, MovieID: arg.MovieID, Rating: arg.Rating,
		}); err != nil {
			ctx.JSON(200, errorResponse(err.Error()))
		}

		ctx.JSON(200, gin.H{
			"hasError": false,
			"message":  "Updated rating successfully",
		})
		return
	}

	if err != sql.ErrNoRows {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	if _, err := server.CreateRating(ctx, arg); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{
		"hasError": false,
		"message":  "Rating created successfully",
	})
}
