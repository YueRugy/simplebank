package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"simplebank/util"
	"testing"
)

var (
	testQuery *Queries
	testDB    *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalln("load config fail", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect db", err)
	}
	testQuery = New(testDB)
	os.Exit(m.Run())
}
