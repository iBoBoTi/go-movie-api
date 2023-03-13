package main

import (
	"fmt"
	"github.com/iBoBoTi/go-movie-api/config"
	"github.com/iBoBoTi/go-movie-api/internal/server"
	"log"
	"net/http"
	"time"
)

// TODO: set cors
func main() {
	fmt.Println("Movie Api")
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	s := &server.Server{
		Config: conf,
	}
	s.Start()
}
