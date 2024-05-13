package db

import (
	"context"
	"testing"

	"github.com/debugroach/video-hub-serve/util"
	"github.com/stretchr/testify/require"
)

func TestCreateRating(t *testing.T) {
	user := createRandomUser(t)
	rating := CreateRatingParams{
		Username: user.Username,
		MovieID:  util.GenerateRandomInt(1, 100),
		Rating:   util.GenerateRandomInt(1, 10),
	}
	result, err := queries.CreateRating(context.Background(), rating)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	ratings, err := queries.GetRatings(context.Background(), rating.Username)
	require.NoError(t, err)
	require.Equal(t, ratings[0].Username, rating.Username)
	require.Equal(t, ratings[0].MovieID, rating.MovieID)
	require.Equal(t, ratings[0].Rating, rating.Rating)
	require.NotZero(t, ratings[0].CreatedAt)
}
