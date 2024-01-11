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

test:
	go test -v -cover ./...
 
.PHONY: potgres createdb dropdb migrateup migratedown sqlc test
