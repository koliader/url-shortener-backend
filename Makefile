createdb:
	 docker exec -it pg12 createdb --username=root --owner=root urls
dropdb:
	 docker exec -it pg12 dropdb urls
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/urls?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/urls?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: createdb dropdb migrateup migratedown sqlc test server