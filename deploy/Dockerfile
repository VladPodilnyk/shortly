# syntax=docker/dockerfile:1

FROM golang:1.22-alpine

RUN apt-get update && apt-get install -y build-essential make

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

# copy all src files
COPY . ./

RUN make build

EXPOSE 4000

CMD [ "/docker-app-shortener" ]