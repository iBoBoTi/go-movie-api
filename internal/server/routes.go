package server

import (
	"fmt"
	"github.com/iBoBoTi/go-movie-api/internal/cache"
	repo "github.com/iBoBoTi/go-movie-api/internal/database"
	"github.com/iBoBoTi/go-movie-api/pkg/character"
	"github.com/iBoBoTi/go-movie-api/pkg/comment"
	"github.com/iBoBoTi/go-movie-api/pkg/movie"
	"github.com/iBoBoTi/go-movie-api/swapi"

	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) defineRoutes(router *gin.Engine) {
	apirouter := router.Group("/api/v1")

	db, err := repo.ConnectPostgres(s.Config)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	swapi.InitClient()
	goMovieCache := cache.NewGoMovieCache(cache.NewRedisClient(s.Config))

	commentRepository := comment.NewRespository(db)
	commentService := comment.NewService(goMovieCache, commentRepository)
	movieService := movie.NewService(commentService, goMovieCache)

	// movie
	movieHandler := movie.NewHandler(movieService)
	apirouter.GET("/movies", movieHandler.GetMovieList)

	// comment
	commentHandler := comment.NewHandler(commentService)
	apirouter.POST("/movie/:movie_id/comments", commentHandler.AddComment)
	apirouter.GET("/movie/:movie_id/comments", commentHandler.GetComments)

	// character
	characterService := character.NewService(goMovieCache)
	characterHandler := character.NewHandler(characterService)
	apirouter.GET("/movie/:movie_id/characters", characterHandler.GetCharactersByMovie)
}

func (s *Server) setupRouter() *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "test" {
		r := gin.New()
		s.defineRoutes(r)
		return r
	}

	r := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	// setup cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.defineRoutes(r)

	return r
}
