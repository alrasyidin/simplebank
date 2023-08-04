package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/alrasyidin/simplebank-go/util"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var testStore Store
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	poolConn, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testStore = NewStore(poolConn)

	os.Exit(m.Run())
}
