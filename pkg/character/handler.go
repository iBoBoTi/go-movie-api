package character

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/go-movie-api/pkg/render"
	"net/http"
	"strconv"
)

type Handler interface {
	GetCharactersByMovie(c *gin.Context)
}

type handler struct {
	CharacterService Service
}

func NewHandler(characterService Service) Handler {
	return &handler{
		CharacterService: characterService,
	}
}

func (h *handler) GetCharactersByMovie(c *gin.Context) {
	sortBy := c.Query("sort_by")
	filterByGender := c.Query("gender")
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		render.JSON(c, http.StatusBadRequest, "", nil, fmt.Errorf("invalid id value"))
		return
	}

	characterResponse, er := h.CharacterService.List(movieID, sortBy, filterByGender)
	if err != nil {
		render.JSON(c, er.Status, "", nil, er)
		return
	}

	render.JSON(c, http.StatusOK, "movie characters retrieved successfully", characterResponse, nil)
}
