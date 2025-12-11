.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: startdb
startdb:
	docker run --name task-postgres -e POSTGRES_PASSWORD=12345 -p 5432:5432 -d postgres
