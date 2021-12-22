package main

import (
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := gin.Default()
	router.GET("/movies", listMovies)
	err := router.Run("127.0.0.1:5000")
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}

type Movie struct {
	Film                     string
	Genre                    string
	LeadStudio               string
	AudienceScorePercentage  string
	Profitability            string
	RottenTomatoesPercentage string
	WorldwideGross           string
	Year                     string
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
		}
		movies = append(movies, movie)
	}
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
	PageSize int `form:"page_size" binding:"required,min=5,max=10"`
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
