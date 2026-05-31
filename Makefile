.PHONY: create_migration migration_up migration_down run

include .env
export

create_migration:
	@migrate create -ext sql -dir ./migrations/ -seq $(SQL_NAME)

migration_up:
	@migrate -path ./migrations/ -database $(DATABASE_URL) -verbose up

migration_down:
	@migrate -path ./migrations/ -database $(DATABASE_URL) -verbose down

run:
	@go run ./cmd/main.go
