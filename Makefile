include app.env

postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:14-alpine

postgres-stop:
	docker stop postgres14 && docker rm postgres14

createdb:
	docker exec -it postgres14 createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DATABASE}

dropdb:
	docker exec -it postgres14 dropdb ${POSTGRES_DATABASE}

migrateup:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DATABASE}?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DATABASE}?sslmode=disable" -verbose down

addmigration:
	migrate create -ext sql -dir db/migration -seq ${name}

sqlc:
	sqlc generate

format-check:
	gofmt -l .

format-lint:
	gofmt -s -w .

test:
	go test -v -cover -short ./...

test-package: 
	go test -v ./${pn}/...

server:
	go run main.go

start-docker:
	docker compose down && docker rmi golang-bank-microservices-api 2> /dev/null || true && docker compose up

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/kelvinator07/golang-bank-microservices/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/kelvinator07/golang-bank-microservices/worker TaskDistributor

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
 
.PHONY: potgres postgres-stop createdb dropdb migrateup migratedown addmigration sqlc format-check format-lint test test-package server start-docker mockgen redis
