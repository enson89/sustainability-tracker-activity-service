.PHONY: build run test migrate

build:
	go build -o sustainability-tracker-activity-service ./cmd

run:
	go run ./cmd/main.go

test:
	go test ./...

# Run migrations using the migrate/migrate Docker image.
# Ensure DATABASE_URL is set in your environment, for example:
# export DATABASE_URL="postgres://user:password@localhost:5432/sustainability?sslmode=disable"
migrate:
	docker run --rm \
	  -v $(PWD)/internal/migration:/migrations \
	  migrate/migrate:v4.15.2 -path=/migrations -database="$(DATABASE_URL)" up