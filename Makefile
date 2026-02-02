# Load environment variables if you have a .env file
include .env

# Create a new migration file (Usage: make migration name=create_users_table)
migration:
	migrate create -ext sql -dir migrations -seq $(name)

# Run migrations up
migrate-up:
	migrate -path migrations -database "$(DB_URI)" up

# Run migrations down (rollback 1 step)
migrate-down:
	migrate -path migrations -database "$(DB_URI)" down 1

# Force a specific version (if things break)
migrate-force:
	migrate -path migrations -database "$(DB_URI)" force $(version)

# Start everything
dev-up:
	docker-compose up -d

# Shut everything down and wipe volumes (start fresh)
dev-down:
	docker-compose down -v

# Run your new Seeder
seed:
	go run cmd/seed/main.go

# View Logs
logs:
	docker-compose logs -f

.PHONY: migration migrate-up migrate-down migrate-force dev-up dev-down migrate-up seed