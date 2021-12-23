package main

import (
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/movies", listMovies)
	err := router.Run("127.0.0.1:5000")
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}

type Movie struct {
	Film                     string `json:"film"`
	Genre                    string `json:"genre"`
	LeadStudio               string `json:"lead_studio"`
	AudienceScorePercentage  string `json:"audience_score_percentage"`
	Profitability            string `json:"profitability"`
	RottenTomatoesPercentage string `json:"rotten_tomatoes_percentage"`
	WorldwideGross           string `json:"worldwide_gross"`
	Year                     string `json:"year"`
	Size                     int    `json:"size"`
}

type ListMoviesParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func ListMovies(ctx context.Context, arg ListMoviesParams) ([]Movie, error) {
	file, err := os.Open("movies.csv")
	if err != nil {
		return nil, err
	}
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}
	movies := []Movie{}
	size := len(lines) - 1
	for _, line := range lines {
		movie := Movie{
			Film:                     line[0],
			Genre:                    line[1],
			LeadStudio:               line[2],
			AudienceScorePercentage:  line[3],
			Profitability:            line[4],
			RottenTomatoesPercentage: line[5],
			WorldwideGross:           line[6],
			Year:                     line[7],
			Size:                     size,
		}
		movies = append(movies, movie)
	}
	movies = movies[1:]
	return paginateMovies(movies, arg.Offset, arg.Limit), nil
}

func paginateMovies(movies []Movie, offset int, limit int) []Movie {
	if offset > len(movies) {
		offset = len(movies)
	}
	end := offset + limit
	if end > len(movies) {
		end = len(movies)
	}
	return movies[offset:end]
}

type listMoviesRequest struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=10"`
}

func listMovies(ctx *gin.Context) {
	var req listMoviesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := ListMoviesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	products, err := ListMovies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, products)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
