package movie

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/go-movie-api/pkg/render"

	"net/http"
)

type Handler interface {
	GetMovieList(c *gin.Context)
}

type handler struct {
	MovieService Service
}

func NewHandler(movieService Service) Handler {
	return &handler{
		MovieService: movieService,
	}
}

func (h *handler) GetMovieList(c *gin.Context) {
	moviesResponse, err := h.MovieService.List()
	if err != nil {
		render.JSON(c, err.Status, "", nil, err)
		return
	}

	render.JSON(c, http.StatusOK, "Movie list retrieved successfully", moviesResponse, nil)
}
