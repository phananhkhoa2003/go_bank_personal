package db

import (
	"context"
	"log"
	"os"
	"simple_bank/util"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	// Override with test database if available from environment
	if testDBSource := os.Getenv("TEST_DB_SOURCE"); testDBSource != "" {
		config.DBSource = testDBSource
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
