postgres:
	docker run --name postgres -p 5555:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root technodom

dropdb:
	docker exec -it postgres dropdb technodom

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/technodom?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/technodom?sslmode=disable" -verbose down

run:
	go run cmd/main.go

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown