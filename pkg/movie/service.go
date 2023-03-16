package movie

import (
	"github.com/iBoBoTi/go-movie-api/errors"
	"github.com/iBoBoTi/go-movie-api/internal/cache"
	"github.com/iBoBoTi/go-movie-api/pkg/comment"
	"github.com/iBoBoTi/go-movie-api/pkg/models"
	"github.com/iBoBoTi/go-movie-api/pkg/movie/types"
	"github.com/iBoBoTi/go-movie-api/swapi"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"sort"
)

type Service interface {
	List() ([]types.Movie, *errors.Error)
}

type movieService struct {
	CommentService comment.Service
	Cache          cache.GoMovieCache
}

func NewService(commentService comment.Service, cache cache.GoMovieCache) Service {
	return &movieService{
		CommentService: commentService,
		Cache:          cache,
	}
}

func (m movieService) List() ([]types.Movie, *errors.Error) {
	var swapiMoviesResponse models.SwapiMovieListResponse

	cachedMovies, err := m.Cache.Get("movies", &swapiMoviesResponse)
	if err == redis.Nil && cachedMovies == nil {
		if err := swapi.DefaultClient.Get("films/", &swapiMoviesResponse); err != nil {
			log.Println(err.ActualError)
			return nil, errors.New(err.Message, err.StatusCode)
		}

		if err := m.Cache.Set("movies", &swapiMoviesResponse); err != nil {
			log.Printf("error setting movielist in cache: %v", err)
		}
	}

	// watch out for this
	if cachedMovies != nil {
		movies, _ := cachedMovies.(*models.SwapiMovieListResponse)
		swapiMoviesResponse = *movies
	}

	sort.Slice(swapiMoviesResponse.Results, func(i, j int) bool {
		return swapiMoviesResponse.Results[i].ReleaseDate < swapiMoviesResponse.Results[j].ReleaseDate
	})

	moviesResponse := make([]types.Movie, 0)
	for _, v := range swapiMoviesResponse.Results {
		var movie types.Movie
		movie.Title = v.Title
		movie.ReleaseDate = v.ReleaseDate
		movie.OpeningCrawl = v.OpeningCrawl
		count, err := m.CommentService.GetCountForMovie(movie.Title)
		if err != nil {
			log.Printf("error getting comments: %v", err)
			return nil, errors.New("internal server error", http.StatusInternalServerError)
		}
		movie.CommentsCount = count
		moviesResponse = append(moviesResponse, movie)
	}

	return moviesResponse, nil
}
