.PHONY: run migrate

include .env

run:
	go run cmd/main.go --config=.env

migrate:
	goose -dir db/migrations postgres "postgresql://$(POSTGRES_USERNAME):$(POSTRGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

test:
	go test -v ./... 
