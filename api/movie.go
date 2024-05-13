package api

import (
	"fmt"
	"math"
	"sort"

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
	// 获取为用户 User1 推荐的电影
	recommendedMovies := recommendMovies(userRatings, req.Username, 20)
	fmt.Println("Recommended Movies for", req.Username)

	for _, movie := range recommendedMovies {
		fmt.Println("MovieID:", movie.MovieID, "Predicted Score:", movie.Score)
	}
	ctx.JSON(200, gin.H{"movies": recommendedMovies})
}

// 计算两个用户之间的皮尔逊相关系数
func pearsonCorrelation(user1, user2 map[int]int) float64 {
	var sum1, sum2, sum1Sq, sum2Sq, pSum float64
	var n float64

	for item, rating := range user1 {
		if rating2, ok := user2[item]; ok {
			n++
			sum1 += float64(rating)
			sum2 += float64(rating2)
			sum1Sq += float64(rating) * float64(rating)
			sum2Sq += float64(rating2) * float64(rating2)
			pSum += float64(rating) * float64(rating2)
		}
	}

	if n == 0 {
		return 0
	}

	num := pSum - (sum1 * sum2 / n)
	fmt.Println(pSum, sum1, sum2, n)
	fmt.Println(num)
	den := math.Sqrt((sum1Sq - math.Pow(sum1, 2)/n) * (sum2Sq - math.Pow(sum2, 2)/n))
	if den == 0 {
		fmt.Println("den = 0")
		return 0
	}
	return num / den
}

// 推荐电影
func recommendMovies(userRatings UserRatings, targetUser string, n int) []MovieScore {
	// 存储未评分电影及其预测评分
	scores := make(map[int]float64)
	simSums := make(map[int]float64)

	// 遍历每个用户计算与目标用户的相似度
	for otherUser, ratings := range userRatings {
		if otherUser == targetUser {
			continue
		}
		sim := pearsonCorrelation(userRatings[targetUser], ratings)
		fmt.Println(targetUser, otherUser, sim)
		if sim <= 0 { // 只考虑正相关性
			continue
		}

		for movieID, rating := range ratings {
			if _, ok := userRatings[targetUser][movieID]; !ok { // 目标用户未评分的电影
				scores[movieID] += float64(rating) * sim // 累加加权评分
				simSums[movieID] += sim                  // 累加相似度
			}
		}
	}

	// 计算加权平均评分
	var movieScores []MovieScore
	for movieID, scoreSum := range scores {
		if simSums[movieID] > 0 { // 避免除以零
			movieScores = append(movieScores, MovieScore{
				MovieID: movieID,
				Score:   scoreSum / simSums[movieID],
			})
		}
	}

	// 按评分降序排序
	sort.Slice(movieScores, func(i, j int) bool {
		return movieScores[i].Score > movieScores[j].Score
	})

	// 返回评分最高的n部电影
	if len(movieScores) > n {
		movieScores = movieScores[:n]
	}
	return movieScores
}
