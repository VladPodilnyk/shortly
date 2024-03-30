# syntax=docker/dockerfile:1

FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

# copy all src files
COPY . ./

RUN go build -ldflags="-s" -o /docker-app-shortener ./cmd/api

EXPOSE 4000

CMD [ "/docker-app-shortener" ]