package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kelvinator07/golang-bank-microservices/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	config, err := util.LoadConfig("../../app-test.env")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
