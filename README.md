# bank-go

- [x] DB Schema : [https://dbdiagram.io/d/Bank-Go-65c363f7ac844320aea4730f]
- [x] DB migrate :
      `brew install golang-migrate`
      `migrate create -ext sql -dir db/migration -seq add_users`
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
- [x] MockDb for testing using STUBS https://github.com/golang/mock
- [x] AccountApi unit test
- [x] createTransfer API
- [x] createTransfer API unit test
- [x] add a user table && its migration
- [x] DB CRUD for the user table + Unit tests
- [x] fix unit tests and api of accounts to adobt the user constrain keys
- [x] Hash password bcrypt

## Improvements:

- [] add APIs: for update and delete account.
