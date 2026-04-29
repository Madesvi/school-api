include .env

# Migration commands
create_migration:
	migrate create -ext=sql -dir=internal/database/migrations -seq init

migrate_up:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up $(STEP)

migrate_down:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down $(STEP)

migrate_force:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" force 1

css:
	npx @tailwindcss/cli -i ./views/css/input.css -o ./public/output.css --content "./views/**/*.templ" --watch

templ:
	templ generate --proxy=http://localhost:3000 --watch

.PHONY: create_migration migrate_up migrate_down migrate_force css
