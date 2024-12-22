.PHONY: run migrate

ENV ?= local
ENV_FILE = $(ENV).env

include $(ENV_FILE)

run:
	go run cmd/main.go --config=$(ENV_FILE)

migrate:
	goose -dir db/migrations postgres "postgresql://$(POSTGRES_USERNAME):$(POSTRGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

test:
	go test -v ./... 
