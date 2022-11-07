gen:
		protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   internal/pb/user_service.proto
clean:
		rm ./internal/pb/*.go
server:
		go run cmd/server/main.go -port 8080
client:
		go run cmd/client/main.go -address 0.0.0.0:8080
test:
		go test -cover -race ./...
up:
		docker-compose up -d
down:
		docker-compose down
logs:
		docker logs hexxdev-clickhouse-1

