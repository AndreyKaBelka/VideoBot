GOOSE_VERSION := v3.24.1
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

MIGRATIONS_DIR := migrations

-include configs/.env
export

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= videobot
DB_PASSWORD ?= videobot
DB_NAME ?= videobot

DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: migrate-up migrate-down migrate-down-to migrate-redo migrate-reset migrate-status migrate-version migrate-create

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-down-to:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down-to $(VERSION)

migrate-redo:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" redo

migrate-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" reset

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

migrate-version:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" version

migrate-create:
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(NAME) sql
