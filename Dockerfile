# docker run --net=host -it --rm server:v1

FROM golang:1.23.1-alpine3.19 as base
FROM base as dev
WORKDIR /app
COPY . .
EXPOSE 8000
RUN apk update
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN apk update && apk add --no-cache postgresql-client
RUN go mod download
RUN go build -o api cmd/main.go
RUN apk add make 