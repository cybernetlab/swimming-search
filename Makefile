COVERAGE=test/coverage

run:
	go run ./cmd/app

up:
	docker compose up --build --force-recreate

down:
	docker compose down

.PHONY: test
test:
	make unit-test

test-coverage:
	go tool cover -html=${COVERAGE}/unit

unit-test:
	go test -cover -coverprofile=${COVERAGE}/unit.out ./internal/domain/... ./internal/usecase/...

# integration-test:
# 	go test -count=1 -v -tags=integration ./test/integration

generate:
	go generate ./...

mockery-install:
	go install github.com/vektra/mockery/v3@v3.7.0