run:
	go run cmd/main.go --config=./config/local.yaml
migrate_up:
	goose -dir db/migrations postgres "postgresql://admin:1234@localhost:5432/db?sslmode=disable" up