# syntax=docker/dockerfile:1

FROM alpine:latest as tailwindbuild
WORKDIR /app
RUN apk update && apk add curl
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
RUN chmod +x tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 tailwindcss

FROM golang:1.22-alpine as gobuild
WORKDIR /app
RUN apk update && apk add --no-cache make bash
COPY go.mod go.sum ./
COPY . ./
COPY --from=tailwindbuild /app/tailwindcss ./
RUN go mod download
RUN make build

FROM alpine:latest as final
WORKDIR /app
COPY --from=gobuild /app/bin/shortly ./

EXPOSE 4000

CMD [ "./shortly"]