package swapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var DefaultClient *Client

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Client struct {
	baseURL    string
	HTTPClient HTTPClient
}

func newClient() *Client {
	return &Client{
		baseURL: "https://swapi.dev/api/",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func InitClient() {
	DefaultClient = newClient()
}

func (s *Client) Get(url string, response interface{}) *Error {
	url = fmt.Sprintf("%v%v", s.baseURL, url)
	res, err := s.HTTPClient.Get(url)
	if err != nil {
		return &Error{
			Message:     "internal server error",
			StatusCode:  http.StatusInternalServerError,
			ActualError: fmt.Errorf("error sending request %+v", err),
		}
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return &Error{
				Message:     "",
				StatusCode:  res.StatusCode,
				ActualError: fmt.Errorf("error sending request %+v", err),
			}
		}
		return &Error{
			Message:    string(b),
			StatusCode: res.StatusCode,
		}
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return &Error{
			Message:     "internal server error",
			StatusCode:  http.StatusInternalServerError,
			ActualError: fmt.Errorf("error decoding response %+v", err),
		}
	}

	return nil
}
