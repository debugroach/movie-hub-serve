package db

import (
	"context"
	"testing"

	"github.com/debugroach/movie-hub-serve/util"
	"github.com/stretchr/testify/require"
)

func TestGetMovie(t *testing.T) {
	args := CreateMovieParams{
		ID:           util.GenerateRandomInt(1, 100),
		Title:        util.GenerateRandomString(6),
		BackdropPath: util.GenerateRandomString(10),
		PosterPath:   util.GenerateRandomString(10),
		VoteAverage:  float64(util.GenerateRandomInt(1, 10)),
	}
	result, err := queries.CreateMovie(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	movie, err := queries.GetMovie(context.Background(), args.ID)
	require.NoError(t, err)
	require.NotEmpty(t, movie)
	require.Equal(t, args.ID, movie.ID)
	require.Equal(t, args.Title, movie.Title)
	require.Equal(t, args.BackdropPath, movie.BackdropPath)
	require.Equal(t, args.PosterPath, movie.PosterPath)
	require.Equal(t, args.VoteAverage, movie.VoteAverage)
}
