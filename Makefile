createdb:
	 docker exec -it pg12 createdb --username=root --owner=root url
dropdb:
	 docker exec -it pg12 dropdb url
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/url?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/url?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: createdb dropdb migrateup migratedown sqlc test server