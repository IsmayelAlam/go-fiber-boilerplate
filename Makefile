# Include variables from the .envrc file
include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run the main application
.PHONY: run
run:
	air

.PHONY: docs
docs:
	swag fmt
	swag init -g cmd/docs/main.go --output temp/
	go run cmd/docs/main.go

## module/new module='module_name': create a new module
.PHONY: module/new
module/new:
	@echo 'Creating new module...'
	@if [ -z "${name}" ]; then \
		echo "Error: 'name' is required. Usage: make module/new name=module_name"; \
		exit 1; \
	fi
	@mkdir -p internal/modules/${name}/{migrations,queries}
	@echo 'package ${name}' > internal/modules/${name}/${name}.go
	@echo 'package ${name}' > internal/modules/${name}/controllers.go
	@echo 'package ${name}' > internal/modules/${name}/routes.go
	@echo 'package ${name}' > internal/modules/${name}/validator.go
	@echo 'package ${name}' > internal/modules/${name}/utils.go
	@echo '-- add new queries here' > internal/modules/${name}/queries/query.sql
	@echo '{"version":"2","sql":[{"engine":"postgresql","schema":"../migrations/*.sql","queries":"../queries/*.sql","gen":{"go":{"package":"${name}Services","out":"../services","emit_json_tags":true}}}]}' > \
	internal/modules/${name}/migrations/sqlc.json

## sqlc/generate module='module_name': use sqlc to generate type-safe database code
.PHONY: sqlc/generate
sqlc/generate:
	@if [ -z "$(module)" ]; then \
		echo "Error: 'module' is required. Usage: make sqlc/generate module=user"; \
		exit 1; \
	fi

	@echo "Generating SQLC for module: $(module)"
	@sqlc generate -f internal/modules/$(module)/migrations/sqlc.json

# ==================================================================================== #
# Database Migrations
# ==================================================================================== #
## db/migrations/new name='file_name' module='module_name': create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@if [ -z "${name}" ]; then \
		echo "Error: 'name' is required. Usage: make db/migrations/new name=file_name module=module_name"; \
		exit 1; \
	fi
	@if [ -z "${module}" ]; then \
		echo "Error: 'module' is required. Usage: make db/migrations/new name=file_name module=module_name"; \
		exit 1; \
	fi
	@echo 'Creating migration files for ${name} in ${module} module...'
	goose create -dir ./internal/modules/${module}/migrations ${module}_${name} sql 

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	@mkdir -p migrations_flat
	@find internal/modules/**/migrations -name "*.sql" -exec cp {} migrations_flat/ \;
	@goose up -dir ./migrations_flat up
	@rm -rf migrations_flat

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	@mkdir -p migrations_flat
	@find internal/modules/**/migrations -name "*.sql" -exec cp {} migrations_flat/ \;
	@goose down -dir ./migrations_flat down
	@rm -rf migrations_flat
