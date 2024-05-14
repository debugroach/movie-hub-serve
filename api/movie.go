package api

import (
	"fmt"
	"math"
	"sync"

	db "github.com/debugroach/video-hub-serve/db/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// 使用 map 来模拟用户和电影评分的数据结构
type UserRatings map[string]map[int]int

// MovieScore 用于存储电影评分和电影ID
type MovieScore struct {
	MovieID int
	Score   float64
}

type recommendRequest struct {
	Username string `json:"username" binding:"required"`
}

func (server *Server) recommend(ctx *gin.Context) {
	req := recommendRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, errorResponse(err.Error()))
		return
	}
	// 读取数据
	userRatings := make(UserRatings)

	ratings, _ := server.ListRatings(ctx)
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
			movie, err := server.GetMovie(ctx, movieID)
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
func findMostSimilarUser(targetUser string, ratings UserRatings) (string, float64) {
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

func recommendMoviesForUser(userID string, ratings UserRatings) []int {
	// 假设只取一个最相似的用户进行推荐，实际应用中可以考虑多个相似用户
	mostSimilarUser, _ := findMostSimilarUser(userID, ratings)
	similarUserRatings := ratings[mostSimilarUser]

	recommendations := []int{}
	for movieID, rating := range similarUserRatings {
		if _, exists := ratings[userID][movieID]; !exists && rating > 0 {
			recommendations = append(recommendations, movieID)
		}
	}

	return recommendations
}
