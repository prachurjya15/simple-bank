package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/prachurjya15/simple-bank/util"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Error Loading Configs")
	}
	testDB, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalf("error connecting to DB: %#v \n", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
