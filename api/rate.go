package api

import (
	"database/sql"

	"github.com/debugroach/movie-hub-serve/db"
	"github.com/gin-gonic/gin"
)

// rateRequest represents the request body for rating a movie.
type rateRequest struct {
	Username     string  `json:"username" binding:"required"`
	MovieID      int     `json:"movieID" binding:"required"`
	Rating       int     `json:"rating" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	PosterPath   string  `json:"posterPath" binding:"required"`
	BackdropPath string  `json:"backdropPath" binding:"required"`
	VoteAverage  float64 `json:"voteAverage" binding:"required"`
}

// createMovie creates a new movie if it doesn't exist in the database.
func (s *Server) createMovie(ctx *gin.Context, req rateRequest) error {
	_, err := s.GetMovie(ctx, req.MovieID)

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

	_, err = s.CreateMovie(ctx, arg)
	return err
}

// rate handles the rating of a movie.
func (s *Server) rate(ctx *gin.Context) {
	var req rateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	if err := s.createMovie(ctx, req); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	arg := db.CreateRatingParams{
		Username: req.Username,
		MovieID:  req.MovieID,
		Rating:   req.Rating,
	}

	_, err := s.GetRating(ctx, db.GetRatingParams{
		Username: arg.Username, MovieID: arg.MovieID,
	})

	if err == nil {
		if _, err = s.UpdateRating(ctx, db.UpdateRatingParams{
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

	if _, err := s.CreateRating(ctx, arg); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{
		"hasError": false,
		"message":  "Rating created successfully",
	})
}
