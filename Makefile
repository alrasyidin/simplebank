DB_URL=postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=root -p 5432:5432 -d postgres:12-alpine
createdb: 
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb: 
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migrations -database "${DB_URL}" -verbose up
migrateup1:
	migrate -path db/migrations -database "${DB_URL}" -verbose up 1
migratedown:
	migrate -path db/migrations -database "${DB_URL}" -verbose down
migratedown1:
	migrate -path db/migrations -database "${DB_URL}" -verbose down 1
new_migration:
	migrate create -ext sql -dir db/migrations -seq ${name}
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	air
mock:
	mockgen -package mockdb -destination ./db/mock/store.go github.com/alrasyidin/simplebank-go/db/sqlc Store
dbdocs:
	dbdocs build docs/db.dbml
proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
evans:
	evans --host localhost --port 8082 -r repl
redis:
	docker run --name redis -p 6378:6379 -d redis:7.0.12-alpine

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test mock server dbdocs proto redis new_migration