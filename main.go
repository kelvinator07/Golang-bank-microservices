package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kelvinator07/golang-bank-microservices/api"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
	"github.com/kelvinator07/golang-bank-microservices/util"

	_ "github.com/lib/pq"
)

var dbSource string = fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s?sslmode=disable",
	util.ViperEnvVariable("POSTGRES_USER"),
	util.ViperEnvVariable("POSTGRES_PASSWORD"),
	util.ViperEnvVariable("POSTGRES_DATABASE"))

const (
	dbDriver      = "postgres"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	fmt.Print(util.ViperEnvVariable(""))
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
