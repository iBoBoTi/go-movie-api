package comment

import (
	"fmt"
	"github.com/iBoBoTi/go-movie-api/errors"
	"github.com/iBoBoTi/go-movie-api/internal/cache"
	"github.com/iBoBoTi/go-movie-api/pkg/comment/types"
	"github.com/iBoBoTi/go-movie-api/pkg/models"
	"github.com/iBoBoTi/go-movie-api/swapi"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

type Service interface {
	GetComments(movieID int) ([]types.Comment, *errors.Error)
	AddComment(comment *types.Comment, movieID int) (*types.Comment, *errors.Error)
	GetCountForMovie(title string) (int, error)
}

type service struct {
	CommentRepository Repository
	Cache             cache.GoMovieCache
}

func NewService(cache cache.GoMovieCache, repo Repository) Service {
	return &service{
		Cache:             cache,
		CommentRepository: repo,
	}
}

func (s *service) AddComment(comment *types.Comment, movieID int) (*types.Comment, *errors.Error) {
	var swapiMovie models.SwapiMovie
	cachedMovie, err := s.Cache.Get(fmt.Sprintf("movie-%v", movieID), &swapiMovie)
	if err == redis.Nil && cachedMovie == nil {

		err := swapi.DefaultClient.Get(fmt.Sprintf("films/%v/", movieID), &swapiMovie)
		if err != nil {
			log.Printf("error making swapi client call %#v", err)

			return nil, errors.New("internal server error", http.StatusInternalServerError)
		}
		if err := s.Cache.Set(fmt.Sprintf("movie-%v", movieID), &swapiMovie); err != nil {
			log.Printf("error setting movie-%v in cache: %v", movieID, err)
		}
	}

	if cachedMovie != nil {
		movie, _ := cachedMovie.(*models.SwapiMovie)
		swapiMovie = *movie
	}

	comment.MovieTitle = swapiMovie.Title
	comment.MovieID = movieID

	createdComment, err := s.CommentRepository.AddComment(comment)
	if err != nil {
		log.Printf("error adding Comment: %v", err)
		return nil, errors.New("internal server error", http.StatusInternalServerError)
	}

	return createdComment, nil
}

func (s *service) GetComments(movieID int) ([]types.Comment, *errors.Error) {
	var swapiMovie models.SwapiMovie
	// Cache  call based off the id
	_, err := s.Cache.Get(fmt.Sprintf("movie-%v", movieID), &swapiMovie)
	if err == redis.Nil {
		err := swapi.DefaultClient.Get(fmt.Sprintf("films/%v/", movieID), &swapiMovie)
		if err != nil {
			log.Printf("error making swapi client call %#v", err)

			return nil, errors.New("internal server error", http.StatusInternalServerError)
		}
		if err := s.Cache.Set(fmt.Sprintf("movie-%v", movieID), &swapiMovie); err != nil {
			log.Printf("error setting movie-%v in cache: %v", movieID, err)
		}
	}

	comments, err := s.CommentRepository.GetCommentsByMovieID(movieID)
	if err != nil {
		log.Printf("error getting comments: %v", err)
		return nil, errors.New("internal server error", http.StatusInternalServerError)
	}
	return comments, nil
}

func (s *service) GetCountForMovie(title string) (int, error) {
	return s.CommentRepository.GetCountForMovie(title)
}

func GetRealIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-IP")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarder-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
