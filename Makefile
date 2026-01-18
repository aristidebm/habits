.PHONY: run format migrate-up migrate-down migrate-create build

build:
	# @mkdir ./build/
	@go build -o ./build ./...

run: build
	@./build/cmd

format:
	@go fmt ./...

test:
	@go test ./...

migrate-up:
	goose -dir migrations sqlite3 ./habits.db up

migrate-down:
	goose -dir migrations sqlite3 ./habits.db down

migrate-create:
	@read -p "Enter migration name: " name && \
	goose -dir migrations create "$$name" sql
