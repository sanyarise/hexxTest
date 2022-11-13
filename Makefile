.PHONY: test
test:
	go test ./... -v

.PHONY: up
up: test
	docker-compose up -d

.PHONY: run
run: test up
	go run ./cmd/client/client.go