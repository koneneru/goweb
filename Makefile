include .envrc

# =============================================================================================== #
# HELPERS
# =============================================================================================== #

## help: print this help information
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /' | sort

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [Y/N] ' && read ans && [ $${ans:-N} = y ]

# =============================================================================================== #
# DEVELOPMENT
# =============================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api  -db-dsn=${GOWEB_DB_DSN}

## db/psql: connect to database using psql
.PHONY: db/psql
db/psql:
	psql ${GOWEB_DB_DSN}

## db/migrations/new name=$1: create a new database migrations
.PHONY: db/migrations/new
db/migrations/new:
	@echo Creating migration files for ${name}...
	migrate create -seq -ext=sql -dir=migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo Running up migrations...
	migrate -path=migrations -database=${GOWEB_DB_DSN} up

# =============================================================================================== #
# QUALITY CONTROL
# =============================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## Vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api/goweb_api.exe ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api/goweb_api.exe ./cmd/api