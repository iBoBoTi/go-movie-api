# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
COPY . ./
RUN go get ./...
RUN go mod download


RUN go build -o /go-movie-api

EXPOSE $PORT

CMD [ "/go-movie-api" ]