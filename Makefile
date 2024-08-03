sinclude .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

run-dev: \
	run-infra \
	run-go

stop-dev: \
	stop-infra

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@latest
	GOBIN=$(LOCAL_BIN) go install github.com/joho/godotenv/cmd/godotenv@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

update-modules:
	GOPRIVATE=$(GOPRIVATE) go mod tidy

run-go:
	$(LOCAL_BIN)/godotenv -f .env go run cmd/main.go

run-infra:
	docker compose up -d --build

stop-go:
	pkill -f 'go run cmd/main.go'

stop-infra:
	docker compose down --remove-orphans

build:
	@read -p "Enter build name: " build_name; \
	$(LOCAL_BIN)/godotenv -f .env go build -o $$build_name cmd/main.go

mock-generate:
	go generate stats-back-minio/internal/services
	go generate stats-back-minio/internal/s3

test:
	go clean -testcache
	go test ./.../tests -covermode count -coverpkg=imp-api/internal/services/job/...

docs-generate:
	$(LOCAL_BIN)/swag init -g cmd/main.go

migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

migration-create:
	@read -p "Enter migration name: " migration_name; \
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) create $$migration_name sql