package character

import (
	"fmt"
	"github.com/iBoBoTi/go-movie-api/errors"
	"github.com/iBoBoTi/go-movie-api/internal/cache"
	"github.com/iBoBoTi/go-movie-api/pkg/character/types"
	"github.com/iBoBoTi/go-movie-api/pkg/models"
	"github.com/iBoBoTi/go-movie-api/swapi"
	"github.com/redis/go-redis/v9"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Service interface {
	List(movieID int, sortBy, filterByGender string) (*types.CharacterListResponse, *errors.Error)
}

type service struct {
	Cache cache.GoMovieCache
}

func NewService(cache cache.GoMovieCache) Service {
	return &service{
		Cache: cache,
	}
}

func (s *service) List(movieID int, sortBy, filterByGender string) (*types.CharacterListResponse, *errors.Error) {

	var swapiMovie models.SwapiMovie

	cachedMovie, err := s.Cache.Get(fmt.Sprintf("movie-%v", movieID), &swapiMovie)
	if err == redis.Nil && cachedMovie == nil {
		err := swapi.DefaultClient.Get(fmt.Sprintf("films/%v/", movieID), &swapiMovie)
		if err != nil {
			log.Printf("error making swapi client call %#v", err)

			return nil, errors.New("internal server error", http.StatusInternalServerError)
		}

		if err := s.Cache.Set(fmt.Sprintf("movie-%v", movieID), &swapiMovie); err != nil {
			log.Printf("error setting movie-%v in cache: %v", movieID, err)
		}
	}

	if cachedMovie != nil {
		movie, _ := cachedMovie.(*models.SwapiMovie)
		swapiMovie = *movie
	}

	//get characters from movie from
	var swapiMovieCharacters []types.Character
	for _, v := range swapiMovie.CharacterURLs {
		var movieCharacter types.Character

		url := strings.Replace(v, "https://swapi.dev/api/", "", 1)

		cachedMovieCharacter, err := s.Cache.Get(fmt.Sprintf("Character-%v", url), &movieCharacter)
		if err == redis.Nil && cachedMovieCharacter == nil {
			_ = swapi.DefaultClient.Get(url, &movieCharacter)

			if err := s.Cache.Set(fmt.Sprintf("Character-%v", url), &movieCharacter); err != nil {
				log.Printf("error setting Character-%v in cache: %v", url, err)
			}
		}

		if cachedMovieCharacter != nil {
			character, _ := cachedMovieCharacter.(*types.Character)
			movieCharacter = *character
		}

		swapiMovieCharacters = append(swapiMovieCharacters, movieCharacter)
	}

	// Sortby and filterby

	if sortBy != "" {
		SortCharacterList(sortBy, swapiMovieCharacters)
	}

	if filterByGender != "" {
		filteredCharacters := FilterCharacterList(filterByGender, swapiMovieCharacters)
		swapiMovieCharacters = filteredCharacters
	}

	//get total number of height in CM/ Feet
	heightsInFeet, heightsInCM := GetTotalHeightOfCharacter(swapiMovieCharacters)

	characterResponse := &types.CharacterListResponse{
		Characters:                    swapiMovieCharacters,
		CharactersCount:               len(swapiMovieCharacters),
		TotalHeightOfCharactersInFeet: heightsInFeet,
		TotalHeightOfCharactersInCM:   heightsInCM,
	}
	return characterResponse, nil
}

func SortCharacterList(sortBy string, characters []types.Character) {
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

func FilterCharacterList(filterByGender string, characters []types.Character) []types.Character {
	filteredCharacters := make([]types.Character, 0)
	for _, v := range characters {
		if strings.ToLower(filterByGender) == strings.ToLower(v.Gender) {
			filteredCharacters = append(filteredCharacters, v)
		}
	}
	return filteredCharacters
}

func GetTotalHeightOfCharacter(characters []types.Character) (string, string) {
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
