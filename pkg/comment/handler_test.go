package comment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/iBoBoTi/go-movie-api/errors"
	comment_mocks "github.com/iBoBoTi/go-movie-api/pkg/comment/mocks"
	"github.com/iBoBoTi/go-movie-api/pkg/comment/types"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_GetComments1(t *testing.T) {

	commentsResponse := []types.Comment{types.Comment{
		MovieTitle: "A New Hope",
		MovieID:    1,
		Author:     "172.21.0.1:61380",
		Content:    "Super Dope Movie",
	},
	}

	tests := []struct {
		name         string
		urlParam     interface{}
		responseCode int
		mocks        func(commentService *comment_mocks.MockService, urlParam interface{})
	}{
		{
			name:         "status_ok",
			urlParam:     1,
			responseCode: http.StatusOK,
			mocks: func(commentService *comment_mocks.MockService, urlParam interface{}) {
				commentService.EXPECT().GetComments(urlParam).Times(1).Return(commentsResponse, nil)
			},
		},
		{
			name:         "bad_request",
			urlParam:     "a",
			responseCode: http.StatusBadRequest,
			mocks: func(commentService *comment_mocks.MockService, urlParam interface{}) {

			},
		},
		{
			name:         "internal_server_error",
			urlParam:     1,
			responseCode: http.StatusInternalServerError,
			mocks: func(commentService *comment_mocks.MockService, urlParam interface{}) {
				commentService.EXPECT().GetComments(urlParam).Times(1).Return(nil, errors.New("internal server error", http.StatusInternalServerError))
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := comment_mocks.NewMockService(ctrl)
	testCommentHandler := &handler{
		CommentService: mockCommentService,
	}

	router := gin.Default()
	router.GET("/api/v1/movie/:movie_id/comments", testCommentHandler.GetComments)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(mockCommentService, tt.urlParam)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/movie/%v/comments", tt.urlParam), nil)
			require.Nil(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			require.Equal(t, tt.responseCode, recorder.Code)
		})
	}
}

func Test_handler_AddComment(t *testing.T) {

	createdComment := &types.Comment{
		MovieTitle: "A New Hope",
		MovieID:    1,
		Content:    "Super Dope Movie",
	}

	tests := []struct {
		name         string
		urlParam     interface{}
		requestBody  interface{}
		responseCode int
		mocks        func(commentService *comment_mocks.MockService, reqBody interface{}, urlParam interface{})
	}{
		{
			name:     "status_created",
			urlParam: 1,
			requestBody: &types.Comment{
				Content: "Awesome movie",
			},
			responseCode: http.StatusCreated,
			mocks: func(commentService *comment_mocks.MockService, reqBody interface{}, urlParam interface{}) {
				commentService.EXPECT().AddComment(reqBody, urlParam).Times(1).Return(createdComment, nil)
			},
		},
		{
			name:         "bad_request",
			urlParam:     "A",
			responseCode: http.StatusBadRequest,
			mocks: func(commentService *comment_mocks.MockService, reqBody interface{}, urlParam interface{}) {

			},
		},
		{
			name:         "Content Exceeds 500 Characters Case",
			urlParam:     1,
			responseCode: http.StatusBadRequest,
			requestBody: &types.Comment{
				Content: `Nam quis nulla. Integer malesuada. In in enim a arcu imperdiet malesuada. Sed vel lectus. 
				Donec odio urna, tempus molestie, porttitor ut, iaculis quis, sem. Phasellus rhoncus. Aenean id metus 
				id velit ullamcorper pulvinar. Vestibulum fermentum tortor id mi. Pellentesque ipsum. Nulla non arcu 
				lacinia neque faucibus fringilla. Nulla non lectus sed nisl molestie malesuada. Proin in tellus sit amet 
				nibh dignissim sagittis. Vivamus luctus egestas leo. Maecenas sollicitudin. Nullam rhoncus aliquam metu`,
			},
			mocks: func(commentService *comment_mocks.MockService, reqBody interface{}, urlParam interface{}) {

			},
		},
		{
			name:     "internal server error",
			urlParam: 1,
			requestBody: &types.Comment{
				Content: "Awesome movie",
			},
			responseCode: http.StatusInternalServerError,
			mocks: func(commentService *comment_mocks.MockService, reqBody interface{}, urlParam interface{}) {
				commentService.EXPECT().AddComment(reqBody, urlParam).Times(1).Return(nil, errors.New("internal server error", http.StatusInternalServerError))
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := comment_mocks.NewMockService(ctrl)
	testcommentHandler := &handler{
		CommentService: mockCommentService,
	}

	router := gin.Default()
	router.POST("/api/v1/movie/:movie_id/comments", testcommentHandler.AddComment)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(mockCommentService, tt.requestBody, tt.urlParam)

			data, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/movie/%v/comments", tt.urlParam), bytes.NewReader(data))
			require.Nil(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			require.Equal(t, tt.responseCode, recorder.Code)
		})
	}
}
