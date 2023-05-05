ifneq (,$(wildcard ./.env))
  include .env
  export
endif

migrateup:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose up

migrateup1:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose up 1

migratedown:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose down

migratedown1:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose down 1

migrateforce:
	migrate -path db/migration -database "${DB_SOURCE}" force 1

sqlc:
	sqlc generate

test-migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/next-bazaar?sslmode=disable" -verbose up

test:
	go test -v -cover -shuffle=on ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ot07/next-bazaar/db/sqlc Store

swag:
	swag init --pd

seed:
	go run ./scripts/seed.go

.PHONY: migrateup migratedown mmigrateup1 migratedown1 igrateforce sqlc test-migrateup test server mock
