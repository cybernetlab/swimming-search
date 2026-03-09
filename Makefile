run:
	go run ./cmd/app

up:
	docker compose up --build --force-recreate

down:
	docker compose down

# .PHONY: test
# test:
# 	go test -v -cover ./...

# integration-test:
# 	go test -count=1 -v -tags=integration ./test/integration

# generate:
# 	go generate ./...

# mockery-install:
# 	go install github.com/vektra/mockery/v3@v3.2.5