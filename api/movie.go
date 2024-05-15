package api

import (
	"fmt"
	"math"
	"sync"

	"github.com/debugroach/movie-hub-serve/db"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type userRating map[string]map[int]int

type recommendRequest struct {
	Username string `json:"username" binding:"required"`
}

// recommend is a handler function that recommends movies for a given user.
func (s *Server) recommend(ctx *gin.Context) {
	req := recommendRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}

	userRatings := make(userRating)

	ratings, _ := s.ListRatings(ctx)
	for _, r := range ratings {
		if userRatings[r.Username] == nil {
			userRatings[r.Username] = make(map[int]int)
		}
		userRatings[r.Username][r.MovieID] = r.Rating
		fmt.Println(r.Username, r.MovieID, userRatings[r.Username][r.MovieID])
	}

	recommendedMovies := recommendMoviesForUser(req.Username, userRatings)
	movies := make([]db.Movie, 0, len(recommendedMovies))

	var mu sync.Mutex
	var wg sync.WaitGroup
	var getMovieError error
	for _, movieID := range recommendedMovies {
		wg.Add(1)
		go func(movieID int) {
			defer wg.Done()
			fmt.Println(movieID)
			movie, err := s.GetMovie(ctx, movieID)
			mu.Lock()
			if err != nil {
				getMovieError = err
				return
			}
			movies = append(movies, movie)
			mu.Unlock()
		}(movieID)
	}
	wg.Wait()
	if getMovieError != nil {
		ctx.JSON(200, errorResponse(getMovieError.Error()))
		return
	}

	ctx.JSON(200, gin.H{"movies": movies})
}

// pearsonCorrelation calculates the Pearson correlation coefficient between two users.
func pearsonCorrelation(user1, user2 map[int]int) float64 {
	var sum1, sum2, sum1Sq, sum2Sq, pSum float64
	commonItems := make(map[int]bool)

	for itemID := range user1 {
		if rating, exists := user2[itemID]; exists {
			commonItems[itemID] = true
			sum1 += float64(user1[itemID])
			sum2 += float64(rating)
			sum1Sq += float64(user1[itemID]) * float64(user1[itemID])
			sum2Sq += float64(rating) * float64(rating)
			pSum += float64(user1[itemID]) * float64(rating)
		}
	}

	n := float64(len(commonItems))
	if n == 0 {
		return 0
	}

	numerator := pSum - (sum1 * sum2 / n)
	denominator := math.Sqrt((sum1Sq - sum1*sum1/n) * (sum2Sq - sum2*sum2/n))

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// findMostSimilarUser finds the user with the highest similarity to the target user.
func findMostSimilarUser(targetUser string, ratings userRating) (string, float64) {
	targetRatings := ratings[targetUser]
	maxSimilarity := 0.0
	var mostSimilar string

	for user, ratings := range ratings {
		if user != targetUser {
			similarity := pearsonCorrelation(targetRatings, ratings)
			if similarity > maxSimilarity {
				maxSimilarity = similarity
				mostSimilar = user
			}
		}
	}

	return mostSimilar, maxSimilarity
}

// recommendMoviesForUser recommends movies for a given user based on the ratings of similar users.
func recommendMoviesForUser(userID string, ratings userRating) []int {
	mostSimilarUser, _ := findMostSimilarUser(userID, ratings)
	similaruserRatings := ratings[mostSimilarUser]

	recommendations := []int{}
	for movieID, rating := range similaruserRatings {
		if _, exists := ratings[userID][movieID]; !exists && rating > 0 {
			recommendations = append(recommendations, movieID)
		}
	}

	return recommendations
}
