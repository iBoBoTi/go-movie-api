package swapi

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/iBoBoTi/go-movie-api/swapi/mocks"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

var integrationClient *Client

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("couldn't load env vars: %v", err)
	}
	if os.Getenv("IntegrationTest") == "true" {
		integrationClient = newClient()
	}
	os.Exit(m.Run())
}

func TestSwapiClient_Get(t *testing.T) {

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

	type args struct {
		url      string
		response interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *moviePayload
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Success",
			args: args{
				url: "films/1/",
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(swapiMovieJSON)),
				},
			},
			want: &moviePayload{
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
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockHTTPClient(ctrl)

	s := Client{
		baseURL:    "https://swapi.dev/api/",
		HTTPClient: mockClient,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.EXPECT().Get(fmt.Sprintf("%v%v", s.baseURL, tt.args.url)).Times(1).Return(tt.args.response, nil)
			result := &moviePayload{}
			err := s.Get(tt.args.url, result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("got = %v, want %v", result, tt.want)
				return
			}

		})
	}
	t.Run("Integration Test", func(t *testing.T) {
		if integrationClient == nil {
			log.Println("skipping integration...")
			return
		}
		result := &moviePayload{}
		err := integrationClient.Get("films/1/", result)
		if err != nil {
			t.Errorf("Get() error = %v", err)
			return
		}
		if reflect.TypeOf(*result) != reflect.TypeOf(moviePayload{}) {
			t.Errorf("expect result types to be %v but got %v", reflect.TypeOf(moviePayload{}), reflect.TypeOf(*result))
			return
		}

	})

}
