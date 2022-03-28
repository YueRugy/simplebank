test:
	go test -v -cover ./...
sqlc:
	sqlc generate
server:
	go run main.go