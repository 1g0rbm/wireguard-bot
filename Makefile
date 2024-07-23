sinclude .env

LOCAL_BIN:=$(CURDIR)/bin

run-dev: \
	run-infra \
	run-go

stop-dev: \
	stop-infra

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@latest
	GOBIN=$(LOCAL_BIN) go install github.com/joho/godotenv/cmd/godotenv@latest
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest

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