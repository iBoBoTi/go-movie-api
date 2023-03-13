package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/go-movie-api/cache"
	"github.com/iBoBoTi/go-movie-api/internal/api/response"
	"github.com/iBoBoTi/go-movie-api/internal/domain"
	"github.com/iBoBoTi/go-movie-api/swapi"
	"github.com/redis/go-redis/v9"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type CharacterHandler interface {
	GetCharactersByMovie(c *gin.Context)
}

type characterHandler struct {
	Cache       cache.GoMovieCache
	SwapiClient *swapi.SwapiClient
}

func NewCharacterHandler(cache cache.GoMovieCache, swapiClient *swapi.SwapiClient) CharacterHandler {
	return &characterHandler{
		Cache:       cache,
		SwapiClient: swapiClient,
	}
}

func (m *characterHandler) GetCharactersByMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "", nil, fmt.Errorf("invalid id value"))
		return
	}

	var swapiMovie domain.SwapiMovie

	cachedMovie, err := m.Cache.Get(fmt.Sprintf("movie-%v", movieID), &swapiMovie)
	if err == redis.Nil && cachedMovie == nil {
		responseStatusCode, err := m.SwapiClient.Get(fmt.Sprintf("films/%v/", movieID), &swapiMovie)
		if err != nil {
			log.Printf("error making swapi client call %#v", err)

			var errorStr string
			SetErrorString(responseStatusCode, &errorStr)

			response.JSON(c, responseStatusCode, "", nil, fmt.Errorf(errorStr))
			return
		}

		if err := m.Cache.Set(fmt.Sprintf("movie-%v", movieID), &swapiMovie); err != nil {
			log.Printf("error setting movie-%v in cache: %v", movieID, err)
		}
	}

	if cachedMovie != nil {
		swapiMovie, _ = cachedMovie.(domain.SwapiMovie)
	}

	//get characters from movie from
	var swapiMovieCharacters []domain.MovieCharacter
	for _, v := range swapiMovie.CharacterURLs {
		var movieCharacter domain.MovieCharacter

		url := strings.Replace(v, "https://swapi.dev/api/", "", 1)

		cachedMovieCharacter, err := m.Cache.Get(fmt.Sprintf("character-%v", url), &movieCharacter)
		if err == redis.Nil && cachedMovieCharacter == nil {
			_, _ = m.SwapiClient.Get(url, &movieCharacter)

			if err := m.Cache.Set(fmt.Sprintf("character-%v", url), &movieCharacter); err != nil {
				log.Printf("error setting character-%v in cache: %v", url, err)
			}
		}

		if cachedMovieCharacter != nil {
			movieCharacter, _ = cachedMovieCharacter.(domain.MovieCharacter)
		}

		swapiMovieCharacters = append(swapiMovieCharacters, movieCharacter)
	}

	// Sortby and filterby
	sortBy := c.Query("sort_by")
	filterByGender := c.Query("gender")
	if sortBy != "" {
		SortCharacterList(sortBy, swapiMovieCharacters)
	}

	if filterByGender != "" {
		filteredCharacters := FilterCharacterList(filterByGender, swapiMovieCharacters)
		swapiMovieCharacters = filteredCharacters
	}

	//get total number of height in CM/ Feet
	heightsInFeet, heightsInCM := GetTotalHeightOfCharacter(swapiMovieCharacters)

	characterResponse := domain.CharacterListResponse{
		Characters:                    swapiMovieCharacters,
		CharactersCount:               len(swapiMovieCharacters),
		TotalHeightOfCharactersInFeet: heightsInFeet,
		TotalHeightOfCharactersInCM:   heightsInCM,
	}

	response.JSON(c, http.StatusOK, "movie characters retrieved successfully", characterResponse, nil)
}

func SortCharacterList(sortBy string, characters []domain.MovieCharacter) {
	switch sortBy {
	case "name.asc":
		sort.Slice(characters, func(i, j int) bool {
			return characters[i].Name < characters[j].Name
		})
	case "name.desc":
		sort.Slice(characters, func(i, j int) bool {
			return characters[i].Name > characters[j].Name
		})
	case "gender.asc":
		sort.Slice(characters, func(i, j int) bool {
			return characters[i].Gender < characters[j].Gender
		})
	case "gender.desc":
		sort.Slice(characters, func(i, j int) bool {
			return characters[i].Gender > characters[j].Gender
		})
	case "height.asc":
		sort.Slice(characters, func(i, j int) bool {
			heightI, _ := strconv.ParseFloat(characters[i].Height, 64)
			heightJ, _ := strconv.ParseFloat(characters[j].Height, 64)
			return heightI < heightJ
		})
	case "height.desc":
		sort.Slice(characters, func(i, j int) bool {
			heightI, _ := strconv.ParseFloat(characters[i].Height, 64)
			heightJ, _ := strconv.ParseFloat(characters[j].Height, 64)
			return heightI > heightJ
		})

	}
}

func FilterCharacterList(filterByGender string, characters []domain.MovieCharacter) []domain.MovieCharacter {
	filteredCharacters := make([]domain.MovieCharacter, 0)
	for _, v := range characters {
		if strings.ToLower(filterByGender) == strings.ToLower(v.Gender) {
			filteredCharacters = append(filteredCharacters, v)
		}
	}
	return filteredCharacters
}

func GetTotalHeightOfCharacter(characters []domain.MovieCharacter) (string, string) {
	var totalHeight float64
	for _, v := range characters {
		height, _ := strconv.ParseFloat(v.Height, 64)
		totalHeight += height
	}
	totalHeightInFeet := totalHeight * 0.0328084

	ft, fraction := math.Modf(totalHeightInFeet)

	inch := fraction * 12

	return fmt.Sprintf("%v ft and %.2f inches", ft, inch), fmt.Sprintf("%v cm", totalHeight)
}
