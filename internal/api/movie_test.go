package api

import (
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
)

func TestGetMovieList(t *testing.T) {

	swapiMoviesJSON := `{
		"count": 2, 
    	"next": null, 
    	"previous": null,
		"results":[
				{
            		"title": "A New Hope",
            		"opening_crawl": "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
            		"release_date": "1977-05-25", 
            		"characters": [
                		"https://swapi.dev/api/people/1/", 
                		"https://swapi.dev/api/people/2/", 
                		"https://swapi.dev/api/people/3/"
            		],
            		"url": "https://swapi.dev/api/films/1/"
        		}, 
        		{
            		"title": "The Empire Strikes Back",
            		"opening_crawl": "It is a dark time for the\r\nRebellion. Although the Death\r\nStar....",
            		"release_date": "1980-05-17", 
            		"characters": [
                		"https://swapi.dev/api/people/1/", 
                		"https://swapi.dev/api/people/2/", 
                		"https://swapi.dev/api/people/3/"
            		],
            		"url": "https://swapi.dev/api/films/2/"
        		}
			]
		}`

	hclient := &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
		require.Equal(t, "https://swapi.dev/api/films/", req.URL.String())
		require.Equal(t, http.MethodGet, req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(swapiMoviesJSON)),
		}
	})}

	var sMoviesResponse domain.SwapiMovieListResponse
	swapiMoviesResponse := domain.SwapiMovieListResponse{
		Count: 2,
		Results: []domain.SwapiMovie{
			domain.SwapiMovie{
				Title:        "A New Hope",
				OpeningCrawl: "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
				CharacterURLs: []string{
					"https://swapi.dev/api/people/1/",
					"https://swapi.dev/api/people/2/",
					"https://swapi.dev/api/people/3/",
				},
				ReleaseDate: "1977-05-25",
				URL:         "https://swapi.dev/api/films/1/",
			},
			domain.SwapiMovie{
				Title:        "The Empire Strikes Back",
				OpeningCrawl: "It is a dark time for the\r\nRebellion. Although the Death\r\nStar....",
				CharacterURLs: []string{
					"https://swapi.dev/api/people/1/",
					"https://swapi.dev/api/people/2/",
					"https://swapi.dev/api/people/3/",
				},
				ReleaseDate: "1980-05-17",
				URL:         "https://swapi.dev/api/films/2/",
			},
		},
	}

	movieResponse := []domain.Movie{
		domain.Movie{
			Title:         "A New Hope",
			OpeningCrawl:  "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
			ReleaseDate:   "1977-05-25",
			CommentsCount: 0,
		},
		domain.Movie{
			Title:         "The Empire Strikes Back",
			OpeningCrawl:  "It is a dark time for the\r\nRebellion. Although the Death\r\nStar....",
			ReleaseDate:   "1980-05-17",
			CommentsCount: 0,
		},
	}
	testCases := []struct {
		name         string
		request      interface{}
		responseData interface{}
		responseCode int
		mocks        func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, count int)
	}{
		{
			name:         "StatusOK Case",
			request:      nil,
			responseCode: http.StatusOK,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, count int) {
				cache.EXPECT().Get("movies", &sMoviesResponse).Times(1).Return(nil, redis.Nil)
				cache.EXPECT().Set("movies", &swapiMoviesResponse).Times(1).Return(nil)
				commentService.EXPECT().GetCommentCountForMovie(gomock.Any()).Times(len(swapiMoviesResponse.Results)).Return(count, nil)
			},
		},
		{
			name:         "StatusOK Case 2",
			request:      nil,
			responseCode: http.StatusOK,
			mocks: func(commentService *mocks.MockCommentService, cache *mocks.MockGoMovieCache, count int) {
				cache.EXPECT().Get("movies", &sMoviesResponse).Times(1).Return(swapiMoviesResponse, nil)
				commentService.EXPECT().GetCommentCountForMovie(gomock.Any()).Times(len(swapiMoviesResponse.Results)).Return(count, nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := mocks.NewMockGoMovieCache(ctrl)
	mockCommentService := mocks.NewMockCommentService(ctrl)

	testSwapiClient := swapi.NewSwapiClient()
	testSwapiClient.HTTPClient = hclient

	testMovieHandler := &movieHandler{
		CommentService: mockCommentService,
		Cache:          mockCache,
		SwapiClient:    testSwapiClient,
	}

	router := gin.Default()
	router.GET("/api/v1/movies", testMovieHandler.GetMovieList)

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			c.mocks(mockCommentService, mockCache, movieResponse[0].CommentsCount)

			req, err := http.NewRequest(http.MethodGet, "/api/v1/movies", nil)
			require.Nil(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			require.Equal(t, c.responseCode, recorder.Code)
		})
	}

}
