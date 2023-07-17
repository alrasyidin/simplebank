postgres:
	docker run --name postgres12 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=root -p 5432:5432 -d postgres:12-alpine
createdb: 
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb: 
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migrations -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migrations -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migrations -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migrations -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	air
mock:
	mockgen -package mockdb -destination ./db/mock/store.go github.com/alrasyidin/simplebank-go/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test mock