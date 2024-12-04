include .envrc

## run/api: run the cmd/api application

.PHONY: run/api
run/api:
	@echo 'Running Application...'
	@go run ./cmd/api \
	-port=4000 \
	-env=development \
	-limiter-burst=5 \
	-limiter-rps=2 \
	-limiter-enabled=false \
	-db-dsn=${TEST3_DB_DSN} \
	-smtp-host=${SMTP_HOST} \
	-smtp-port=${SMTP_PORT} \
	-smtp-username=${SMTP_USERNAME} \
	-smtp-password=${SMTP_PASSWORD} \
	-smtp-sender=${SMTP_SENDER} \
	-limiter-rps=3 \
	-limiter-burst=5 \
	-limiter-enabled=false \
	-cors-trusted-origin="http://localhost:9000 http://localhost:9001"
## @go run ./cmd/api/ -port=4000 -env=production -db-dsn=${COMMENTS_DB_DSN}

## db/psql: connect to the database using psql (terminal)
.PHONY: db/psql
db/psql: 
	psql ${TEST3_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY db/migrations/new:
	@echo 'creating migration fles for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running UP migrations...'
	migrate -path ./migrations -database ${TEST3_DB_DSN} up

.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Running DOWN migrations...'
	migrate -path ./migrations -database ${TEST3_DB_DSN} down


.PHONY: db/migrations/version
db/migrations/version:
	@echo 'Checking current database migration version.....'
	migrate -path ./migrations -database ${TEST3_DB_DSN} version

.PHONY: db/migrations/force
db/migrations/force:
	@echo 'Chenging current database migration version to ${version}'
	migrate -path ./migrations -database ${TEST3_DB_DSN} force ${version}

.PHONY: db/migrations/goto
db/migrations/goto:
	@echo 'ROlling back to version: ${version}'
	migrate -path ./migrations -database ${TEST3_DB_DSN} goto ${version}

.PHONY: db/migrations/down
db/migrations/step_down:
	@echo 'Reverting migrations ${num} steps back'
	migrate -path ./migrations -database ${TEST3_DB_DSN} down ${num}
