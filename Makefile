run:
	@go run server/routes.go server/template.go server/server.go

postgresinit:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root go-chat

dropdb:
	docker exec -it postgres15 dropdb go-chat

createmigration:
	migrate create -ext sql -dir db/migrations add_users_table

migrateup:
	migrate -path db/migrations/ -database "postgresql://root:password@localhost:5432/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations/ -database "postgresql://root:password@localhost:5432/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit postgres createdb dropdb createmigration migrateup migratedown