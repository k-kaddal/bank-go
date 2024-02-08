# bank-go

- [x] DB Schema : [https://dbdiagram.io/d/Bank-Go-65c363f7ac844320aea4730f]
- [x] DB migrate :
      `brew install golang-migrate`
      `migrate create -ext sql -dir database/migration -seq init_schema`
      `make migrateup` : to migrate the database up
      `make migratedown` : to migrate the database down
- [x] Generate CRUD using sql:
      `brew install kyleconroy/sqlc/sqlc`
      `sqlc init`
      `make sqlc`
- [x] DB CRUD Unit testing
- [x] DB Transaction (store)
- [x] Resolving DB deadlock
- [] Transaction Isolation Level
