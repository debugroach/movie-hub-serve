package db

import (
	"context"
	"testing"

	"github.com/debugroach/video-hub-serve/util"
	"github.com/stretchr/testify/require"
)

func TestGetMovie(t *testing.T) {
	args := CreateMovieParams{
		ID:           util.GenerateRandomInt(1, 100),
		Title:        "title",
		BackdropPath: "1",
		PosterPath:   "2",
		VoteAverage:  1,
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
