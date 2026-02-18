ROOT := $(CURDIR)
SERVER_DIR := $(ROOT)/apps/server
MIGRATIONS_DIR := $(SERVER_DIR)/migrations
BACKUP_DIR := $(ROOT)/backups
DB_URL ?= $(shell grep '^DATABASE_URL=' $(SERVER_DIR)/.env 2>/dev/null | cut -d '=' -f2-)
MIGRATE := cd $(SERVER_DIR) && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
MIGRATE_FLAGS := -path $(MIGRATIONS_DIR) -database "$(DB_URL)"

.PHONY: deps swagger run db-up db-down db-dump db-restore db-restore-reset migrate-up migrate-down migrate-force migrate-create

deps:
	cd $(SERVER_DIR) && go get -u ./... && go mod tidy

swagger:
	cd $(SERVER_DIR) && go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -d cmd/api,internal/handlers,internal/catalog -o docs --parseInternal

run:
	cd $(SERVER_DIR) && go run ./cmd/api

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-dump:
	@mkdir -p $(BACKUP_DIR)
	@FILE="$(BACKUP_DIR)/naengrijeongdon_$$(date +%Y%m%d_%H%M%S).sql"; \
	docker exec naengrijeongdon-postgres pg_dump -U naengrijeongdon -d naengrijeongdon > "$$FILE"; \
	echo "Database backup saved: $$FILE"

db-restore:
	@FILE_PATH="$(FILE)"; \
	if [ -z "$$FILE_PATH" ]; then \
		set -- backups/*.sql; \
		if [ "$$1" = "backups/*.sql" ]; then \
			echo "No backup files found in backups/"; \
			exit 1; \
		fi; \
		echo "Select backup file:"; \
		i=1; \
		for f in $$@; do echo "  $$i) $$f"; i=$$((i+1)); done; \
		printf "Enter number: "; \
		read n; \
		i=1; \
		for f in $$@; do \
			if [ "$$i" = "$$n" ]; then FILE_PATH="$$f"; break; fi; \
			i=$$((i+1)); \
		done; \
	fi; \
	if [ -z "$$FILE_PATH" ] || [ ! -f "$$FILE_PATH" ]; then \
		echo "Invalid FILE. Example: make db-restore FILE=backups/naengrijeongdon_20260218_120000.sql"; \
		exit 1; \
	fi; \
	cat "$$FILE_PATH" | docker exec -i naengrijeongdon-postgres psql -U naengrijeongdon -d naengrijeongdon; \
	echo "Database restored from: $$FILE_PATH"

db-restore-reset:
	@FILE_PATH="$(FILE)"; \
	if [ -z "$$FILE_PATH" ]; then \
		set -- backups/*.sql; \
		if [ "$$1" = "backups/*.sql" ]; then \
			echo "No backup files found in backups/"; \
			exit 1; \
		fi; \
		echo "Select backup file:"; \
		i=1; \
		for f in $$@; do echo "  $$i) $$f"; i=$$((i+1)); done; \
		printf "Enter number: "; \
		read n; \
		i=1; \
		for f in $$@; do \
			if [ "$$i" = "$$n" ]; then FILE_PATH="$$f"; break; fi; \
			i=$$((i+1)); \
		done; \
	fi; \
	if [ -z "$$FILE_PATH" ] || [ ! -f "$$FILE_PATH" ]; then \
		echo "Invalid FILE. Example: make db-restore-reset FILE=backups/naengrijeongdon_20260218_120000.sql"; \
		exit 1; \
	fi; \
	docker exec naengrijeongdon-postgres psql -U naengrijeongdon -d naengrijeongdon -c "drop schema if exists public cascade; create schema public;"; \
	cat "$$FILE_PATH" | docker exec -i naengrijeongdon-postgres psql -U naengrijeongdon -d naengrijeongdon; \
	echo "Database reset + restored from: $$FILE_PATH"

migrate-up:
	@if [ -z "$(DB_URL)" ]; then echo "DB_URL is empty. Create $(SERVER_DIR)/.env first."; exit 1; fi
	$(MIGRATE) $(MIGRATE_FLAGS) up

migrate-down:
	@if [ -z "$(DB_URL)" ]; then echo "DB_URL is empty. Create $(SERVER_DIR)/.env first."; exit 1; fi
	$(MIGRATE) $(MIGRATE_FLAGS) down 1

migrate-force:
	@if [ -z "$(DB_URL)" ]; then echo "DB_URL is empty. Create $(SERVER_DIR)/.env first."; exit 1; fi
	@if [ -z "$(VERSION)" ]; then echo "Set VERSION, e.g. make migrate-force VERSION=1"; exit 1; fi
	$(MIGRATE) $(MIGRATE_FLAGS) force $(VERSION)

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Set NAME, e.g. make migrate-create NAME=add_items_index"; exit 1; fi
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
