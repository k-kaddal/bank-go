migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank_db?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank_db?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/bank_db?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/k-kaddal/bank-go/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: migrateup migratedown sqlc test server mock migrateup1 migratedown1 proto
