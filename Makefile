migrateup:
	migrate -path database/migration -database "postgresql://root:password@localhost:5432/bank_db?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migration -database "postgresql://root:password@localhost:5432/bank_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/k-kaddal/bank-go/db/sqlc Store

.PHONY: migrateup migratedown sqlc test server mock
