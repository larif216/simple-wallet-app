.PHONY: migrate

MODULES	 = $(shell cd module && \ls -d */)
TMP_DIR	:= $(shell mktemp -d)
UNAME		:= $(shell uname)

# PostgreSQL Default Settings
export POSTGRES_USER ?= postgres
export POSTGRES_PASS ?= postgres
export POSTGRES_DB_NAME ?= wallet
export POSTGRES_DB_HOST ?= 127.0.0.1
export POSTGRES_DB_PORT ?= 5432

# General Commands
bin:
	@mkdir -p bin

tool-migrate: bin
ifeq ($(UNAME), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else ifeq ($(UNAME), Darwin)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else
	@echo "Your OS is not supported."
endif

migrate: tool-migrate
	@$(foreach module, $(MODULES), cp module/$(module)/db/migrations/postgresql/*.sql $(TMP_DIR) 2>/dev/null;)
	@bin/migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASS)@$(POSTGRES_DB_HOST):$(POSTGRES_DB_PORT)/$(POSTGRES_DB_NAME)?sslmode=disable" \
	-path $(TMP_DIR) \
	$(MIGRATE_ARGS)
