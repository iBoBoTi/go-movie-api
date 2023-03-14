package comment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/go-movie-api/pkg"
	"github.com/iBoBoTi/go-movie-api/pkg/comment/types"
	"github.com/iBoBoTi/go-movie-api/pkg/render"
	"strconv"
	"strings"

	"net/http"
)

type Handler interface {
	GetComments(c *gin.Context)
	AddComment(c *gin.Context)
}

type handler struct {
	CommentService Service
}

func NewHandler(commentService Service) Handler {
	return &handler{
		CommentService: commentService,
	}
}

func (h *handler) GetComments(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		render.JSON(c, http.StatusBadRequest, "", nil, fmt.Errorf("invalid id value"))
		return
	}

	comments, errr := h.CommentService.GetComments(movieID)
	if errr != nil {
		render.JSON(c, errr.Status, "", nil, errr)
		return
	}

	render.JSON(c, http.StatusOK, "comments retrieved successfully", comments, nil)
}

func (h *handler) AddComment(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		render.JSON(c, http.StatusBadRequest, "", nil, fmt.Errorf("invalid id value"))
		return
	}

	var comment types.Comment
	if err := pkg.Decode(c, &comment); err != nil {
		if strings.Contains(err.Error(), "Content is invalid") {
			err = fmt.Errorf("content exceeds 500 characters or is empty")
		}
		render.JSON(c, http.StatusBadRequest, "", nil, err)
		return
	}
	comment.Author = GetRealIP(c.Request)

	createdComment, errr := h.CommentService.AddComment(&comment, movieID)
	if errr != nil {
		render.JSON(c, errr.Status, "", nil, errr)
		return
	}

	render.JSON(c, http.StatusCreated, "Comment added successfully", createdComment, nil)
}
