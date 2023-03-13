package swapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SwapiClient struct {
	baseURL    string
	HTTPClient *http.Client
}

func NewSwapiClient() *SwapiClient {
	return &SwapiClient{
		baseURL: "https://swapi.dev/api/",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *SwapiClient) Get(url string, response interface{}) (int, error) {
	url = fmt.Sprintf("%v%v", s.baseURL, url)
	res, err := s.HTTPClient.Get(url)
	if err != nil {

		return res.StatusCode, fmt.Errorf("error sending request %+v", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return http.StatusBadRequest, fmt.Errorf("error marshalling response: %s", err)
	}

	if res.StatusCode >= http.StatusBadRequest {
		return res.StatusCode, fmt.Errorf("response with status code: %v message: %v", res.StatusCode, response)
	}
	return res.StatusCode, nil
}
