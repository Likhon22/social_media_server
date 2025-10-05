include .envrc

MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migration migrate-up migrate-down migrate-status migrate-force

# Create a new migration file
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

# Apply all up migrations
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

# Show the current migration version
migrate-status:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) version

# Roll back migrations
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down

# Force set the migration version
migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(VERSION)


# Show current migration version
migrate-version:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) version

.PHONY: seed
seed:
	@bash -c 'source .envrc && go run cmd/migrate/seed/main.go'
