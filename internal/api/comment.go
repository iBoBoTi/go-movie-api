package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/go-movie-api/cache"
	"github.com/iBoBoTi/go-movie-api/internal/api/response"
	"github.com/iBoBoTi/go-movie-api/internal/domain"
	"github.com/iBoBoTi/go-movie-api/swapi"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"

	"github.com/iBoBoTi/go-movie-api/internal/usecase"
	"log"
	"net/http"
)

type CommentHandler interface {
	GetCommentsByMovie(c *gin.Context)
	AddCommentToMovie(c *gin.Context)
}

type commentHandler struct {
	CommentService usecase.CommentService
	Cache          cache.GoMovieCache
	SwapiClient    *swapi.SwapiClient
}

func NewCommentHandler(commentService usecase.CommentService, cache cache.GoMovieCache, swapiClient *swapi.SwapiClient) CommentHandler {
	return &commentHandler{
		CommentService: commentService,
		Cache:          cache,
		SwapiClient:    swapiClient,
	}
}

func (m *commentHandler) GetCommentsByMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "", nil, fmt.Errorf("invalid id value"))
		return
	}

	var swapiMovie domain.SwapiMovie
	// Cache  call based off the id
	_, err = m.Cache.Get(fmt.Sprintf("movie-%v", movieID), &swapiMovie)
	if err == redis.Nil {
		responseStatusCode, err := m.SwapiClient.Get(fmt.Sprintf("films/%v/", movieID), &swapiMovie)
		if err != nil {
			log.Printf("error making swapi client call %#v", err)

			var errorStr string
			SetErrorString(responseStatusCode, &errorStr)

			response.JSON(c, responseStatusCode, "", nil, fmt.Errorf(errorStr))
			return
		}
		if err := m.Cache.Set(fmt.Sprintf("movie-%v", movieID), &swapiMovie); err != nil {
			log.Printf("error setting movie-%v in cache: %v", movieID, err)
		}
	}

	comments, err := m.CommentService.GetCommentsByMovieByID(movieID)
	if err != nil {
		log.Printf("error getting comments: %v", err)
		response.JSON(c, http.StatusInternalServerError, "", nil, fmt.Errorf("internal server error"))
		return
	}

	response.JSON(c, http.StatusOK, "comments retrieved successfully", comments, nil)
}

func (m *commentHandler) AddCommentToMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "", nil, fmt.Errorf("invalid id value"))
		return
	}

	var swapiMovie domain.SwapiMovie
	cachedMovie, err := m.Cache.Get(fmt.Sprintf("movie-%v", movieID), &swapiMovie)
	if err == redis.Nil && cachedMovie == nil {

		responseStatusCode, err := m.SwapiClient.Get(fmt.Sprintf("films/%v/", movieID), &swapiMovie)
		if err != nil {
			log.Printf("error making swapi client call %#v", err)

			var errorStr string
			SetErrorString(responseStatusCode, &errorStr)

			response.JSON(c, responseStatusCode, "", nil, fmt.Errorf(errorStr))
			return
		}
		if err := m.Cache.Set(fmt.Sprintf("movie-%v", movieID), &swapiMovie); err != nil {
			log.Printf("error setting movie-%v in cache: %v", movieID, err)
		}
	}

	if cachedMovie != nil {
		swapiMovie, _ = cachedMovie.(domain.SwapiMovie)
	}

	var comment domain.Comment
	if err := decode(c, &comment); err != nil {
		if strings.Contains(err.Error(), "Content is invalid") {
			err = fmt.Errorf("content exceeds 500 characters or is empty")
		}
		response.JSON(c, http.StatusBadRequest, "", nil, err)
		return
	}

	comment.Author = GetRealIP(c.Request)
	comment.MovieTitle = swapiMovie.Title
	comment.MovieID = movieID

	createdComment, err := m.CommentService.AddComment(&comment)
	if err != nil {
		log.Printf("error adding comment: %v", err)
		response.JSON(c, http.StatusInternalServerError, "", nil, err)
		return
	}

	response.JSON(c, http.StatusCreated, "comment added successfully", createdComment, nil)
}

func SetErrorString(resCode int, errorStr *string) {
	switch {
	case resCode == http.StatusBadRequest:
		*errorStr = "bad request"
		break
	case resCode == http.StatusNotFound:
		*errorStr = "record not found"
		break
	case resCode == http.StatusInternalServerError:
		*errorStr = "internal server error"
		break
	}
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
