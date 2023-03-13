package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/iBoBoTi/go-movie-api/internal/domain"
	"github.com/iBoBoTi/go-movie-api/mocks"
	"github.com/iBoBoTi/go-movie-api/swapi"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetCommentByMovie(t *testing.T) {

	swapiMovieJSON := `{
		"title": "A New Hope",
		"opening_crawl": "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
		"release_date": "1977-05-25", 
		"characters": [
		"https://swapi.dev/api/people/1/", 
		"https://swapi.dev/api/people/2/", 
		"https://swapi.dev/api/people/3/"
		],
		"url": "https://swapi.dev/api/films/1/"
	}`

	hclient := &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
		require.Equal(t, "https://swapi.dev/api/films/1/", req.URL.String())
		require.Equal(t, http.MethodGet, req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(swapiMovieJSON)),
		}
	})}

	var sMovieResponse domain.SwapiMovie
	swapiMovieResponse := domain.SwapiMovie{
		Title:        "A New Hope",
		OpeningCrawl: "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
		CharacterURLs: []string{
			"https://swapi.dev/api/people/1/",
			"https://swapi.dev/api/people/2/",
			"https://swapi.dev/api/people/3/",
		},
		ReleaseDate: "1977-05-25",
		URL:         "https://swapi.dev/api/films/1/",
	}

	ct, _ := time.Parse("2023-03-10T00:57:48.186695Z", "2023-03-10T00:57:48.186695Z")

	commentsResponse := []domain.Comment{domain.Comment{
		ID:         2,
		MovieTitle: "A New Hope",
		MovieID:    1,
		Author:     "172.21.0.1:61380",
		Content:    "Super Dope Movie",
		CreatedAt:  ct,
	},
	}
	testCases := []struct {
		name         string
		urlParam     int
		responseData interface{}
		responseCode int
		mocks        func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int)
	}{
		{
			name:         "StatusOK Case",
			responseCode: http.StatusOK,
			urlParam:     1,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int) {
				cache.EXPECT().Get(fmt.Sprintf("movie-%v", urlParam), &sMovieResponse).Times(1).Return(nil, redis.Nil)
				cache.EXPECT().Set(fmt.Sprintf("movie-%v", urlParam), &swapiMovieResponse).Times(1).Return(nil)
				commentService.EXPECT().GetCommentsByMovieByID(urlParam).Times(1).Return(commentsResponse, nil)
			},
		},
		{
			name:         "StatusOK Case 2",
			responseCode: http.StatusOK,
			urlParam:     1,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int) {
				cache.EXPECT().Get(fmt.Sprintf("movie-%v", urlParam), &sMovieResponse).Times(1).Return(swapiMovieResponse, nil)
				commentService.EXPECT().GetCommentsByMovieByID(urlParam).Times(1).Return(commentsResponse, nil)
			},
		},
		{
			name:         "StatusInternalServerError Case",
			responseCode: http.StatusInternalServerError,
			urlParam:     1,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int) {
				cache.EXPECT().Get(fmt.Sprintf("movie-%v", urlParam), &sMovieResponse).Times(1).Return(swapiMovieResponse, nil)
				commentService.EXPECT().GetCommentsByMovieByID(urlParam).Times(1).Return(nil, fmt.Errorf("internal server error"))
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := mocks.NewMockGoMovieCache(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)

	testSwapiClient := swapi.NewSwapiClient()
	testSwapiClient.HTTPClient = hclient

	testCommentHandler := &commentHandler{
		CommentService: mockCommentService,
		Cache:          mockCache,
		SwapiClient:    testSwapiClient,
	}

	router := gin.Default()
	router.GET("/api/v1/movie/:movie_id/comments", testCommentHandler.GetCommentsByMovie)

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			c.mocks(mockCommentService, mockCache, c.urlParam)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/movie/%v/comments", c.urlParam), nil)
			require.Nil(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			require.Equal(t, c.responseCode, recorder.Code)
		})
	}
}

func TestAddCommentToMovie(t *testing.T) {
	swapiMovieJSON := `{
		"title": "A New Hope",
		"opening_crawl": "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
		"release_date": "1977-05-25", 
		"characters": [
		"https://swapi.dev/api/people/1/", 
		"https://swapi.dev/api/people/2/", 
		"https://swapi.dev/api/people/3/"
		],
		"url": "https://swapi.dev/api/films/1/"
	}`

	hclient := &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
		require.Equal(t, "https://swapi.dev/api/films/1/", req.URL.String())
		require.Equal(t, http.MethodGet, req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(swapiMovieJSON)),
		}
	})}

	var sMovieResponse domain.SwapiMovie //get cache param

	//set cache param
	swapiMovieResponse := domain.SwapiMovie{
		Title:        "A New Hope",
		OpeningCrawl: "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
		CharacterURLs: []string{
			"https://swapi.dev/api/people/1/",
			"https://swapi.dev/api/people/2/",
			"https://swapi.dev/api/people/3/",
		},
		ReleaseDate: "1977-05-25",
		URL:         "https://swapi.dev/api/films/1/",
	}

	commentResponse := domain.Comment{
		MovieTitle: "A New Hope",
		MovieID:    1,
		Content:    "Super Dope Movie",
	}
	testCases := []struct {
		name         string
		urlParam     int
		request      *domain.Comment
		responseData interface{}
		responseCode int
		mocks        func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int)
	}{
		{
			name:         "StatusCreated Case",
			responseCode: http.StatusCreated,
			request:      &domain.Comment{Content: "Super Dope Movie"},
			urlParam:     1,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int) {
				cache.EXPECT().Get(fmt.Sprintf("movie-%v", urlParam), &sMovieResponse).Times(1).Return(nil, redis.Nil)
				cache.EXPECT().Set(fmt.Sprintf("movie-%v", urlParam), &swapiMovieResponse).Times(1).Return(nil)
				commentService.EXPECT().AddComment(&commentResponse).Times(1).Return(&commentResponse, nil)
			},
		},
		{
			name:         "StatusCreated Case 2",
			responseCode: http.StatusCreated,
			request:      &domain.Comment{Content: "Super Dope Movie"},
			urlParam:     1,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int) {
				cache.EXPECT().Get(fmt.Sprintf("movie-%v", urlParam), &sMovieResponse).Times(1).Return(swapiMovieResponse, nil)
				commentService.EXPECT().AddComment(&commentResponse).Times(1).Return(&commentResponse, nil)
			},
		},
		{
			name:         "StatusInternalServerError Case",
			responseCode: http.StatusInternalServerError,
			request:      &domain.Comment{Content: "Super Dope Movie"},
			urlParam:     1,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int) {
				cache.EXPECT().Get(fmt.Sprintf("movie-%v", urlParam), &sMovieResponse).Times(1).Return(nil, redis.Nil)
				cache.EXPECT().Set(fmt.Sprintf("movie-%v", urlParam), &swapiMovieResponse).Times(1).Return(nil)
				commentService.EXPECT().AddComment(&commentResponse).Times(1).Return(nil, fmt.Errorf("internal server error"))
			},
		},
		{
			name:         "Content Exceeds 500 Characters Case",
			responseCode: http.StatusBadRequest,
			request: &domain.Comment{
				Content: `Nam quis nulla. Integer malesuada. In in enim a arcu imperdiet malesuada. Sed vel lectus. 
				Donec odio urna, tempus molestie, porttitor ut, iaculis quis, sem. Phasellus rhoncus. Aenean id metus 
				id velit ullamcorper pulvinar. Vestibulum fermentum tortor id mi. Pellentesque ipsum. Nulla non arcu 
				lacinia neque faucibus fringilla. Nulla non lectus sed nisl molestie malesuada. Proin in tellus sit amet 
				nibh dignissim sagittis. Vivamus luctus egestas leo. Maecenas sollicitudin. Nullam rhoncus aliquam metu`,
			},
			urlParam: 1,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, urlParam int) {
				cache.EXPECT().Get(fmt.Sprintf("movie-%v", urlParam), &sMovieResponse).Times(1).Return(nil, redis.Nil)
				cache.EXPECT().Set(fmt.Sprintf("movie-%v", urlParam), &swapiMovieResponse).Times(1).Return(nil)
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := mocks.NewMockGoMovieCache(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)

	testSwapiClient := swapi.NewSwapiClient()
	testSwapiClient.HTTPClient = hclient

	testCommentHandler := &commentHandler{
		CommentService: mockCommentService,
		Cache:          mockCache,
		SwapiClient:    testSwapiClient,
	}

	router := gin.Default()
	router.POST("/api/v1/movie/:movie_id/comments", testCommentHandler.AddCommentToMovie)

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			c.mocks(mockCommentService, mockCache, c.urlParam)

			data, err := json.Marshal(c.request)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/movie/%v/comments", c.urlParam), bytes.NewReader(data))
			require.Nil(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			require.Equal(t, c.responseCode, recorder.Code)
		})
	}
}
