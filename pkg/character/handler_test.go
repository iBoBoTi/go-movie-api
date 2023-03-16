package character

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/iBoBoTi/go-movie-api/errors"
	character_mocks "github.com/iBoBoTi/go-movie-api/pkg/character/mocks"
	"github.com/iBoBoTi/go-movie-api/pkg/character/types"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestGetCharacterByMovies(t *testing.T) {
//
//	swapiMovieJSON := `{
//		"title": "A New Hope",
//		"opening_crawl": "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
//		"release_date": "1977-05-25",
//		"characters": [
//		"https://swapi.dev/api/people/1/",
//		"https://swapi.dev/api/people/2/",
//		"https://swapi.dev/api/people/3/"
//		],
//		"url": "https://swapi.dev/api/films/1/"
//	}`
//
//	swapiCharacterJSON := `
//	"name": "Luke Skywalker",
//	"height": "172",
//	"gender": "male",
//	"films": [
//	  "https://swapi.dev/api/films/1/",
//	  "https://swapi.dev/api/films/2/"
//	],
//	"url": "https://swapi.dev/api/people/1/"
//	`
//
//	hclient := &http.Client{Transport: api.RoundTripFunc(func(req *http.Request) *http.Response {
//
//		require.Equal(t, http.MethodGet, req.Method)
//
//		if req.URL.String() == "https://swapi.dev/api/films/" {
//			return &http.Response{
//				StatusCode: http.StatusOK,
//				Body:       io.NopCloser(strings.NewReader(swapiMovieJSON)),
//			}
//		} else {
//			return &http.Response{
//				StatusCode: http.StatusOK,
//				Body:       io.NopCloser(strings.NewReader(swapiCharacterJSON)),
//			}
//		}
//	})}
//
//	//hclient := &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
//	//	require.Equal(t, "https://swapi.dev/api/films/1/", req.URL.String())
//	//	require.Equal(t, http.MethodGet, req.Method)
//	//
//	//	return &http.Response{
//	//		StatusCode: http.StatusOK,
//	//		Body:       io.NopCloser(strings.NewReader(swapiMovieJSON)),
//	//	}
//	//})}
//
//	var sMovieResponse domain.SwapiMovie
//	swapiMovieResponse := domain.SwapiMovie{
//		Title:        "A New Hope",
//		OpeningCrawl: "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
//		CharacterURLs: []string{
//			"https://swapi.dev/api/people/1/",
//			"https://swapi.dev/api/people/2/",
//			"https://swapi.dev/api/people/3/",
//		},
//		ReleaseDate: "1977-05-25",
//		URL:         "https://swapi.dev/api/films/1/",
//	}
//
//	var sCharacter domain.MovieCharacter
//	swapiCharacter := domain.MovieCharacter{
//		Name:   "Luke Skywalker",
//		Height: "72",
//		Gender: "male",
//		Films: []string{
//			"https://swapi.dev/api/films/1/",
//			"https://swapi.dev/api/films/2/",
//		},
//		Url: "https://swapi.dev/api/people/1/",
//	}
//
//	characterResponse := domain.CharacterListResponse{
//		Characters:                    []domain.MovieCharacter{swapiCharacter},
//		CharactersCount:               1,
//		TotalHeightOfCharactersInCM:   "",
//		TotalHeightOfCharactersInFeet: "",
//	}
//	testCases := []struct {
//		name         string
//		movieID      int
//		responseData interface{}
//		responseCode int
//		mocks        func(cache *mocks.MockGoMovieCache, movieID int)
//	}{
//		{
//			name:         "StatusOK Case",
//			movieID:      1,
//			responseCode: http.StatusOK,
//			mocks: func(cache *mocks.MockGoMovieCache, movieID int) {
//				cache.EXPECT().Get(fmt.Sprintf("movie-%v", movieID), &sMovieResponse).Times(1).Return(nil, redis.Nil)
//				cache.EXPECT().Set(fmt.Sprintf("movie-%v", movieID), &swapiMovieResponse).Times(1).Return(nil)
//				cache.EXPECT().Get(fmt.Sprintf("Character-%v", 1), &sCharacter).Times(1).Return(nil, redis.Nil)
//				cache.EXPECT().Set(fmt.Sprintf("Character-%v", 1), &characterResponse).Times(1).Return(nil)
//			},
//		},
//	}
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockCache := mocks.NewMockGoMovieCache(ctrl)
//
//	testSwapiClient := swapi.NewClient()
//	testSwapiClient.HTTPClient = hclient
//
//	testCharacterHandler := &characterHandler{
//		Cache:       mockCache,
//		SwapiClient: testSwapiClient,
//	}
//
//	router := gin.Default()
//	router.GET("/movie/:movie_id/characters", testCharacterHandler.GetCharactersByMovie)
//
//	for _, c := range testCases {
//		t.Run(c.name, func(t *testing.T) {
//			c.mocks(mockCache, c.movieID)
//
//			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/movie/%v/characters", c.movieID), nil)
//			require.Nil(t, err)
//
//			recorder := httptest.NewRecorder()
//			router.ServeHTTP(recorder, req)
//
//			require.Equal(t, c.responseCode, recorder.Code)
//		})
//	}
//}

func Test_handler_GetCharactersByMovie(t *testing.T) {

	characters := &types.CharacterListResponse{
		Characters: []types.Character{
			types.Character{
				Name:   "Luke Skywalker",
				Height: "72",
				Gender: "male",
				Films: []string{
					"https://swapi.dev/api/films/1/",
					"https://swapi.dev/api/films/2/",
				},
				Url: "https://swapi.dev/api/people/1/",
			},
		},
		CharactersCount:               1,
		TotalHeightOfCharactersInCM:   "",
		TotalHeightOfCharactersInFeet: "",
	}

	tests := []struct {
		name         string
		sortBy       string
		filterBy     string
		urlParam     interface{}
		responseCode int
		mocks        func(characterService *character_mocks.MockService, urlParam interface{}, filterBy, sortBy string)
	}{
		{
			name:         "status_ok",
			sortBy:       "gender.asc",
			filterBy:     "male",
			urlParam:     1,
			responseCode: http.StatusOK,
			mocks: func(characterService *character_mocks.MockService, urlParam interface{}, filterBy, sortBy string) {
				characterService.EXPECT().List(urlParam, sortBy, filterBy).Times(1).Return(characters, nil)
			},
		},
		{
			name:         "bad_request",
			sortBy:       "gender.asc",
			urlParam:     "a",
			filterBy:     "male",
			responseCode: http.StatusBadRequest,
			mocks: func(characterService *character_mocks.MockService, urlParam interface{}, filterBy, sortBy string) {

			},
		},
		{
			name:         "internal_server_error",
			sortBy:       "gender.asc",
			filterBy:     "male",
			urlParam:     1,
			responseCode: http.StatusInternalServerError,
			mocks: func(characterService *character_mocks.MockService, urlParam interface{}, filterBy, sortBy string) {
				characterService.EXPECT().List(urlParam, sortBy, filterBy).Times(1).Return(nil, errors.New("internal server error", http.StatusInternalServerError))
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCharacterService := character_mocks.NewMockService(ctrl)
	testCharacterHandler := &handler{
		CharacterService: mockCharacterService,
	}

	router := gin.Default()
	router.GET("/api/v1/movie/:movie_id/characters", testCharacterHandler.GetCharactersByMovie)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(mockCharacterService, tt.urlParam, tt.filterBy, tt.sortBy)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/movie/%v/characters?gender=%v&sort_by=%v", tt.urlParam, tt.filterBy, tt.sortBy), nil)
			require.Nil(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			require.Equal(t, tt.responseCode, recorder.Code)
		})
	}
}
