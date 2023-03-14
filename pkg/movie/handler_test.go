package movie

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/iBoBoTi/go-movie-api/errors"
	movie_mocks "github.com/iBoBoTi/go-movie-api/pkg/movie/mocks"
	"github.com/iBoBoTi/go-movie-api/pkg/movie/types"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_GetMovieList(t *testing.T) {

	movieResponse := []types.Movie{
		types.Movie{
			Title:         "A New Hope",
			OpeningCrawl:  "It is a period of civil war.\r\nRebel spaceships, striking\r\nfrom a hidden...",
			ReleaseDate:   "1977-05-25",
			CommentsCount: 0,
		},
		types.Movie{
			Title:         "The Empire Strikes Back",
			OpeningCrawl:  "It is a dark time for the\r\nRebellion. Although the Death\r\nStar....",
			ReleaseDate:   "1980-05-17",
			CommentsCount: 0,
		},
	}
	tests := []struct {
		name         string
		responseCode int
		want         *errors.Error
		want2        []types.Movie
	}{
		{
			name:         "status_ok",
			responseCode: http.StatusOK,
			want:         nil,
			want2:        movieResponse,
		},
		{
			name:         "status_ok",
			responseCode: http.StatusInternalServerError,
			want:         errors.New("internal server error", http.StatusInternalServerError),
			want2:        nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieService := movie_mocks.NewMockService(ctrl)
	testMovieHandler := &handler{
		MovieService: mockMovieService,
	}

	router := gin.Default()
	router.GET("/api/v1/movies", testMovieHandler.GetMovieList)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMovieService.EXPECT().List().Times(1).Return(tt.want2, tt.want)
			req, err := http.NewRequest(http.MethodGet, "/api/v1/movies", nil)
			require.Nil(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			require.Equal(t, tt.responseCode, recorder.Code)
		})
	}
}
