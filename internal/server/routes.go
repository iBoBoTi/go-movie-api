package server

import (
	"fmt"
	cache2 "github.com/iBoBoTi/go-movie-api/cache"
	"github.com/iBoBoTi/go-movie-api/internal/api"
	repo "github.com/iBoBoTi/go-movie-api/internal/respository"
	"github.com/iBoBoTi/go-movie-api/internal/usecase"
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
	swapiClient := swapi.NewSwapiClient()

	goMovieCache := cache2.NewGoMovieCache(cache2.NewRedisClient(s.Config))

	commentRepository := repo.NewCommentRespository(db)
	commentService := usecase.NewCommentService(commentRepository)

	// movie
	movieHandler := api.NewMovieHandler(commentService, goMovieCache, swapiClient)
	apirouter.GET("/movies", movieHandler.GetMovieList)

	// comment
	commentHandler := api.NewCommentHandler(commentService, goMovieCache, swapiClient)
	apirouter.POST("/movie/:movie_id/comments", commentHandler.AddCommentToMovie)
	apirouter.GET("/movie/:movie_id/comments", commentHandler.GetCommentsByMovie)

	// character
	characterHandler := api.NewCharacterHandler(goMovieCache, swapiClient)
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
