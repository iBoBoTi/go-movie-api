package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/go-movie-api/cache"
	"github.com/iBoBoTi/go-movie-api/internal/api/response"
	"github.com/iBoBoTi/go-movie-api/internal/domain"
	"github.com/iBoBoTi/go-movie-api/internal/usecase"
	"github.com/iBoBoTi/go-movie-api/swapi"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"sort"
)

type MovieHandler interface {
	GetMovieList(c *gin.Context)
}

type movieHandler struct {
	CommentService usecase.CommentService
	Cache          cache.GoMovieCache
	SwapiClient    *swapi.SwapiClient
}

func NewMovieHandler(commentService usecase.CommentService, cache cache.GoMovieCache, swapiClient *swapi.SwapiClient) MovieHandler {
	return &movieHandler{
		CommentService: commentService,
		Cache:          cache,
		SwapiClient:    swapiClient,
	}
}

func (m *movieHandler) GetMovieList(c *gin.Context) {
	var swapiMoviesResponse domain.SwapiMovieListResponse

	cachedMovies, err := m.Cache.Get("movies", &swapiMoviesResponse)
	if err == redis.Nil && cachedMovies == nil {
		responseStatusCode, err := m.SwapiClient.Get("films/", &swapiMoviesResponse)
		if err != nil {
			log.Printf("error making swapi call %#v", err)

			var errorStr string
			SetErrorString(responseStatusCode, &errorStr)

			response.JSON(c, responseStatusCode, "", nil, fmt.Errorf(errorStr))
			return
		}
		if err := m.Cache.Set("movies", &swapiMoviesResponse); err != nil {
			log.Printf("error setting movielist in cache: %v", err)
		}
	}

	if cachedMovies != nil {
		movies, _ := cachedMovies.(*domain.SwapiMovieListResponse)
		swapiMoviesResponse = *movies
	}

	sort.Slice(swapiMoviesResponse.Results, func(i, j int) bool {
		return swapiMoviesResponse.Results[i].ReleaseDate < swapiMoviesResponse.Results[j].ReleaseDate
	})

	moviesResponse := make([]domain.Movie, 0)
	for _, v := range swapiMoviesResponse.Results {
		var movie domain.Movie
		movie.Title = v.Title
		movie.ReleaseDate = v.ReleaseDate
		movie.OpeningCrawl = v.OpeningCrawl
		count, err := m.CommentService.GetCommentCountForMovie(movie.Title)
		if err != nil {
			log.Printf("error getting comments: %v", err)
			response.JSON(c, http.StatusInternalServerError, "", nil, fmt.Errorf("internal server error"))
			return
		}
		movie.CommentsCount = count
		moviesResponse = append(moviesResponse, movie)
	}

	response.JSON(c, http.StatusOK, "movie list retrieved successfully", moviesResponse, nil)
}
