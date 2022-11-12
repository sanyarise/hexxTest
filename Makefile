dbmock:
	mockgen -source=internal/usecases/userrepo/userrepo.go -destination=internal/usecases/userrepo/userrepo_mock.go -package=userrepo UserStore

qmock:
	mockgen -source=internal/usecases/qrepo/qrepo.go -destination=internal/usecases/qrepo/qrepo_mock.go -package=qrepo Queue

cashmock:
	mockgen -source=internal/usecases/cashrepo/cashrepo.go -destination=internal/usecases/cashrepo/cashrepo_mock.go -package=cashrepo Cash

up:
	docker-compose up -d

down:
	docker-compose down

client:
	go run ./cmd/client/client.go

slog:
	docker logs server

clog:
	docker logs logs_client
