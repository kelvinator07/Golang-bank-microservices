package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/hibiken/asynq"
	"github.com/kelvinator07/golang-bank-microservices/api"
	db "github.com/kelvinator07/golang-bank-microservices/db/sqlc"
	"github.com/kelvinator07/golang-bank-microservices/util"
	"github.com/kelvinator07/golang-bank-microservices/worker"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("app.env")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	// run db migration
	runDBMigration(config.MigrationURL, config.DBSource)

	// load test data
	// loadTestData(conn)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)

	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("Cannot connect to config: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Println("Starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("Failed to start task processor: ", err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("Cannot create new migration instance: ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}

	log.Println("db migrated succesfully")
}

func loadTestData(db *sql.DB) {
	// Read file
	file, err := os.ReadFile("./testdata.sql")
	if err != nil {
		log.Fatal("read file error: ", err.Error())
	}

	// Execute file
	_, err = db.Exec(string(file))
	if err != nil {
		log.Fatal("execute file error: ", err.Error())
	}

	log.Println("test data loaded succesfully")
}
