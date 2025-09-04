MIGRATE = migrate
MIGRATION_DIR = ./db/migrations
DB_URL = "postgres://user_rw:user_rw@localhost:5432/expensemanagement?sslmode=disable"

.PHONY: test

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	$(MIGRATE) -path $(MIGRATION_DIR) -database $(DB_URL) up

migrate-down:
	$(MIGRATE) -path $(MIGRATION_DIR) -database $(DB_URL) down

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	$(MIGRATE) create -ext sql -dir $(MIGRATION_DIR) $(name)

env:
	cp env.sample .env

run-api:
	go run cmd/api/main.go

run-consumer:
	go run cmd/consumer/main.go

test:
	go test -v ./...
