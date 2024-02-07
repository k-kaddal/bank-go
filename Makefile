migrateup:
	migrate -path database/migration -database "postgresql://root:password@localhost:5432/bank_db?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migration -database "postgresql://root:password@localhost:5432/bank_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: migrateup migratedown sqlc test
