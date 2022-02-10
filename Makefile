postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root beatstore
dropdb:
	docker exec -it postgres14 dropdb beatstore
migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/beatstore?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/beatstore?sslmode=disable" -verbose down
sqlc:
	sqlc generate
mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/danglebary/beatstore-backend-go/db/sqlc Store
test:
	go test -v -cover ./...
server:
	go run main.go
tidy:
	go mod tidy

.PHONY: postgres createdb dropdb migrateup migratedown sqlc mockgen test server tidy