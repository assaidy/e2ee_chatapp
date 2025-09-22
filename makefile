include .env

GOOSE_ENV=GOOSE_DRIVER="postgres" GOOSE_DBSTRING="$(PG_URL)" GOOSE_MIGRATION_DIR="./db/migrations/"

all: build

run: build 
	@./bin/app

build:
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -o ./bin/app main.go

test:
	@go test -v ./...

compose-up:
	@docker-compose up

compose-down:
	@docker-compose down

goose-up:
	@$(GOOSE_ENV) goose up

goose-down:
	@$(GOOSE_ENV) goose down

goose-reset:
	@$(GOOSE_ENV) goose reset

goose-migration:
	@if [ -z "$(name)" ]; then echo "ERROR: 'name' variable is required." && exit 1; fi
	@$(GOOSE_ENV) goose create -s $(name) sql

sqlc:
	@sqlc generate

pg:
	@psql $(PG_URL)

vk:
	@valkey-cli -p $(VALKEY_PORT)
