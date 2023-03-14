include .env
export

dev:
	go run cmd/main.go

up:
	docker compose up --build

migrate-create:
	migrate create -ext sql -dir database/migration/ -seq go-movie_init_schema

migrate-up:
	migrate -path database/migration/ -database $(DATABASE_URL) -verbose up

migrate-down:
	migrate -path database/migration/ -database $(DATABASE_URL) -verbose down

migrate-clean:
	migrate -path database/migration/ -database $(DATABASE_URL) force  1

generate-mock:
	mockgen -destination=./pkg/comment/mocks/service_mock.go -package=comment_mocks github.com/iBoBoTi/go-movie-api/pkg/comment Service
	mockgen -destination=./pkg/comment/mocks/repository_mock.go -package=comment_mocks github.com/iBoBoTi/go-movie-api/pkg/comment Repository
	mockgen -destination=./pkg/movie/mocks/service_mock.go -package=movie_mocks github.com/iBoBoTi/go-movie-api/pkg/movie Service
	mockgen -destination=./pkg/character/mocks/service_mock.go -package=character_mocks github.com/iBoBoTi/go-movie-api/pkg/character Service
	mockgen -destination=./internal/cache/mocks/cache_mock.go -package=cache_mocks github.com/iBoBoTi/go-movie-api/internal/cache GoMovieCache