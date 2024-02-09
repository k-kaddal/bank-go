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
- [x] Transaction Isolation Level
- [x] build_test actions workflow
- [x] RESTful HTTP API; createAccount, getAccount, listAccounts
- [x] load config uing viper
- [] MockDb for testing using STUBS https://github.com/golang/mock

## Improvements:

- [] add APIs: for update and delete account.
