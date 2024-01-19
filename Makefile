include .env

postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DATABASE}

dropdb:
	docker exec -it postgres14 dropdb ${POSTGRES_DATABASE}

migrateup:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DATABASE}?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DATABASE}?sslmode=disable" -verbose down

sqlc:
	sqlc generate

format-check:
	gofmt -l .

format-lint:
	gofmt -s -w .

test:
	go test -v -cover ./...

test-package: 
	go test -v ./${pn}/...

server:
	go run main.go

start-docker:
	docker compose down && docker rmi golang-bank-microservices-api 2> /dev/null || true && docker compose up

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/kelvinator07/golang-bank-microservices/db/sqlc Store
 
.PHONY: potgres createdb dropdb migrateup migratedown sqlc format-check format-lint test test-package server start-docker mockgen
