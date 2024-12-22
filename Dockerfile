# docker run --net=host -it --rm server:v1

FROM golang:1.23.1-alpine3.19 as base
FROM base as dev
WORKDIR /app
COPY . .
EXPOSE 44044
RUN apk update
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN apk add --no-cache postgresql-client
RUN apk add make 