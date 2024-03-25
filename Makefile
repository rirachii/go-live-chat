export PORT=:8080

run:
	PORT=${PORT} go run cmd/main.go

postgresinit:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgresup:
	docker start postgres15

postgres:
	docker exec -it postgres15 psql

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root go-chat

dropdb:
	docker exec -it postgres15 dropdb go-chat

migrateup:
	migrate -path service/db/migrations/ -database "postgresql://root:password@localhost:5432/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path service/db/migrations/ -database "postgresql://root:password@localhost:5432/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit postgres createdb dropdb createmigration migrateup migratedown