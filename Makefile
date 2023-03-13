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
	mockgen -destination=mocks/comment_service_mock.go -package=mocks github.com/iBoBoTi/go-movie-api/internal/usecase CommentService
	mockgen -destination=mocks/cache_mock.go -package=mocks github.com/iBoBoTi/go-movie-api/cache GoMovieCache